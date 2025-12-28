// Package jtbd provides a high-performance test execution engine for JTBD tests.
package jtbd

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ExecutionMode defines how tests are executed.
type ExecutionMode string

const (
	// ExecutionModeSequential runs tests one at a time.
	ExecutionModeSequential ExecutionMode = "sequential"
	// ExecutionModeParallel runs independent tests concurrently.
	ExecutionModeParallel ExecutionMode = "parallel"
	// ExecutionModeFailFast stops on first failure.
	ExecutionModeFailFast ExecutionMode = "fail-fast"
	// ExecutionModeComprehensive runs all tests regardless of failures.
	ExecutionModeComprehensive ExecutionMode = "comprehensive"
)

// TestStatus represents the outcome of a test execution.
type TestStatus string

const (
	TestStatusPending   TestStatus = "pending"
	TestStatusRunning   TestStatus = "running"
	TestStatusPassed    TestStatus = "passed"
	TestStatusFailed    TestStatus = "failed"
	TestStatusSkipped   TestStatus = "skipped"
	TestStatusRetrying  TestStatus = "retrying"
)

// Test represents a single test with lifecycle hooks.
type Test struct {
	ID           string
	Name         string
	Description  string
	Dependencies []string
	Priority     int
	Timeout      time.Duration
	MaxRetries   int

	// Lifecycle hooks
	Setup    func(ctx context.Context) error
	Execute  func(ctx context.Context) error // Required
	Teardown func(ctx context.Context) error
}

// ExecutionResult contains the outcome of a test execution.
type ExecutionResult struct {
	TestID       string        `json:"test_id"`
	Status       TestStatus    `json:"status"`
	Error        error         `json:"-"`
	ErrorMessage string        `json:"error_message,omitempty"`
	Duration     time.Duration `json:"duration"`
	RetryCount   int           `json:"retry_count"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Output       string        `json:"output,omitempty"`
	SkipReason   string        `json:"skip_reason,omitempty"`
}

// RunConfig configures test execution.
type RunConfig struct {
	Mode           ExecutionMode
	MaxWorkers     int
	GlobalTimeout  time.Duration
	TestTimeout    time.Duration
	EnableRetry    bool
	IsolateTests   bool
}

// DefaultRunConfig returns default configuration.
func DefaultRunConfig() *RunConfig {
	return &RunConfig{
		Mode:          ExecutionModeParallel,
		MaxWorkers:    10,
		GlobalTimeout: 30 * time.Minute,
		TestTimeout:   5 * time.Minute,
		EnableRetry:   true,
		IsolateTests:  true,
	}
}

// ExecutionEngine orchestrates test execution.
type ExecutionEngine struct {
	tests    []*Test
	config   *RunConfig
	plan     *ExecutionPlan
	workChan chan *Test
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc

	// Metrics (atomic)
	totalTests    atomic.Int32
	passedTests   atomic.Int32
	failedTests   atomic.Int32
	skippedTests  atomic.Int32
	retryAttempts atomic.Int32

	// Results
	results   []*ExecutionResult
	resultsMu sync.Mutex

	// Shared state
	mu              sync.RWMutex
	completedTests  map[string]bool
	failedTestsList map[string]bool
}

// ExecutionPlan determines test execution order based on dependencies.
type ExecutionPlan struct {
	mu           sync.RWMutex
	tests        []*Test
	dependencies map[string][]string
	completed    map[string]bool
	failed       map[string]bool
}

// TestMetrics provides execution statistics.
type TestMetrics struct {
	Total    int32
	Passed   int32
	Failed   int32
	Skipped  int32
	Retries  int32
}

// NewExecutionEngine creates a new test execution engine.
func NewExecutionEngine(tests []*Test, config *RunConfig) (*ExecutionEngine, error) {
	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests provided")
	}
	if config == nil {
		config = DefaultRunConfig()
	}

	// Validate and cap max workers
	if config.MaxWorkers < 1 {
		config.MaxWorkers = 1
	}
	if config.MaxWorkers > 100 {
		config.MaxWorkers = 100 // Safety cap
	}

	plan, err := NewExecutionPlan(tests)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution plan: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.GlobalTimeout)

	ee := &ExecutionEngine{
		tests:           tests,
		config:          config,
		plan:            plan,
		workChan:        make(chan *Test, len(tests)),
		ctx:             ctx,
		cancel:          cancel,
		results:         make([]*ExecutionResult, 0, len(tests)),
		completedTests:  make(map[string]bool),
		failedTestsList: make(map[string]bool),
	}

	ee.totalTests.Store(int32(len(tests)))

	return ee, nil
}

// Run executes all tests according to the configuration.
func (ee *ExecutionEngine) Run() ([]*ExecutionResult, error) {
	defer ee.cancel()

	switch ee.config.Mode {
	case ExecutionModeSequential:
		return ee.runSequential()
	case ExecutionModeParallel:
		return ee.runParallel()
	case ExecutionModeFailFast:
		return ee.runFailFast()
	case ExecutionModeComprehensive:
		return ee.runComprehensive()
	default:
		return nil, fmt.Errorf("unknown execution mode: %s", ee.config.Mode)
	}
}

// runSequential executes tests one at a time.
func (ee *ExecutionEngine) runSequential() ([]*ExecutionResult, error) {
	ordered, err := ee.plan.GetExecutionOrder()
	if err != nil {
		return nil, err
	}

	for _, test := range ordered {
		if !ee.shouldRunTest(test) {
			ee.skipTest(test, "dependencies failed")
			continue
		}

		result := ee.executeTest(ee.ctx, test)
		ee.recordResult(result)

		if result.Status == TestStatusFailed {
			ee.markTestFailed(test.ID)
		} else if result.Status == TestStatusPassed {
			ee.markTestCompleted(test.ID)
		}
	}

	return ee.results, nil
}

// runParallel executes independent tests concurrently.
func (ee *ExecutionEngine) runParallel() ([]*ExecutionResult, error) {
	// Start worker pool
	for i := 0; i < ee.config.MaxWorkers; i++ {
		ee.wg.Add(1)
		go ee.worker(i)
	}

	// Dispatcher goroutine
	go func() {
		ee.dispatchTests()
		close(ee.workChan) // Signal workers to finish
	}()

	// Wait for all workers to complete
	ee.wg.Wait()

	return ee.results, nil
}

// runFailFast executes tests and stops on first failure.
func (ee *ExecutionEngine) runFailFast() ([]*ExecutionResult, error) {
	ordered, err := ee.plan.GetExecutionOrder()
	if err != nil {
		return nil, err
	}

	for _, test := range ordered {
		select {
		case <-ee.ctx.Done():
			return ee.results, ee.ctx.Err()
		default:
		}

		if !ee.shouldRunTest(test) {
			ee.skipTest(test, "dependencies failed")
			continue
		}

		result := ee.executeTest(ee.ctx, test)
		ee.recordResult(result)

		if result.Status == TestStatusFailed {
			ee.markTestFailed(test.ID)
			return ee.results, fmt.Errorf("test failed: %s", test.ID)
		}

		ee.markTestCompleted(test.ID)
	}

	return ee.results, nil
}

// runComprehensive executes all tests regardless of failures.
func (ee *ExecutionEngine) runComprehensive() ([]*ExecutionResult, error) {
	return ee.runParallel() // Same as parallel but don't stop on failure
}

// worker processes tests from the work channel.
func (ee *ExecutionEngine) worker(id int) {
	defer ee.wg.Done()

	for test := range ee.workChan {
		select {
		case <-ee.ctx.Done():
			ee.skipTest(test, "context canceled")
			return
		default:
		}

		if !ee.shouldRunTest(test) {
			ee.skipTest(test, "dependencies not met")
			continue
		}

		result := ee.executeTest(ee.ctx, test)
		ee.recordResult(result)

		if result.Status == TestStatusPassed {
			ee.markTestCompleted(test.ID)
			ee.plan.MarkCompleted(test.ID)
		} else if result.Status == TestStatusFailed {
			ee.markTestFailed(test.ID)
			ee.plan.MarkFailed(test.ID)
		}
	}
}

// dispatchTests sends tests to workers as dependencies are satisfied.
func (ee *ExecutionEngine) dispatchTests() {
	dispatched := make(map[string]bool)

	for {
		select {
		case <-ee.ctx.Done():
			return
		default:
		}

		// Find tests ready to run
		ready := ee.plan.GetReadyTests()
		if len(ready) == 0 {
			// Check if all tests dispatched
			if len(dispatched) == len(ee.tests) {
				return
			}
			time.Sleep(50 * time.Millisecond)
			continue
		}

		for _, test := range ready {
			if !dispatched[test.ID] {
				dispatched[test.ID] = true
				ee.workChan <- test
			}
		}
	}
}

// executeTest runs a single test with retry logic.
func (ee *ExecutionEngine) executeTest(ctx context.Context, test *Test) *ExecutionResult {
	result := &ExecutionResult{
		TestID:    test.ID,
		StartTime: time.Now(),
	}

	maxAttempts := 1
	if ee.config.EnableRetry && test.MaxRetries > 0 {
		maxAttempts = test.MaxRetries + 1
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			result.RetryCount = attempt
			ee.retryAttempts.Add(1)
			// Exponential backoff with jitter
			baseDelay := 100 * time.Millisecond
			backoff := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
			jitter := time.Duration(rand.Float64()*float64(baseDelay)*2 - float64(baseDelay))
			time.Sleep(backoff + jitter)
		}

		err := ee.runTestLifecycle(ctx, test)
		if err == nil {
			result.Status = TestStatusPassed
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}

		lastErr = err
	}

	result.Status = TestStatusFailed
	result.Error = lastErr
	result.ErrorMessage = lastErr.Error()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	return result
}

// runTestLifecycle executes setup, execute, and teardown.
func (ee *ExecutionEngine) runTestLifecycle(ctx context.Context, test *Test) error {
	// Create test-specific context with timeout
	timeout := ee.config.TestTimeout
	if test.Timeout > 0 {
		timeout = test.Timeout
	}
	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Setup
	if test.Setup != nil {
		if err := test.Setup(testCtx); err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
	}

	// Teardown (always run, even on failure)
	if test.Teardown != nil {
		defer func() {
			if err := test.Teardown(context.Background()); err != nil {
				// Log but don't fail test
			}
		}()
	}

	// Execute
	if test.Execute == nil {
		return fmt.Errorf("test has no Execute function")
	}

	if err := test.Execute(testCtx); err != nil {
		return fmt.Errorf("execute failed: %w", err)
	}

	return nil
}

// shouldRunTest checks if test dependencies are satisfied.
func (ee *ExecutionEngine) shouldRunTest(test *Test) bool {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	for _, depID := range test.Dependencies {
		if ee.failedTestsList[depID] {
			return false
		}
		if !ee.completedTests[depID] {
			return false
		}
	}
	return true
}

// recordResult adds a result to the results list.
func (ee *ExecutionEngine) recordResult(result *ExecutionResult) {
	ee.resultsMu.Lock()
	defer ee.resultsMu.Unlock()

	ee.results = append(ee.results, result)

	switch result.Status {
	case TestStatusPassed:
		ee.passedTests.Add(1)
	case TestStatusFailed:
		ee.failedTests.Add(1)
	case TestStatusSkipped:
		ee.skippedTests.Add(1)
	}
}

// skipTest marks a test as skipped.
func (ee *ExecutionEngine) skipTest(test *Test, reason string) {
	result := &ExecutionResult{
		TestID:     test.ID,
		Status:     TestStatusSkipped,
		SkipReason: reason,
		StartTime:  time.Now(),
		EndTime:    time.Now(),
	}
	ee.recordResult(result)
}

// markTestCompleted marks a test as completed.
func (ee *ExecutionEngine) markTestCompleted(testID string) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	ee.completedTests[testID] = true
}

// markTestFailed marks a test as failed.
func (ee *ExecutionEngine) markTestFailed(testID string) {
	ee.mu.Lock()
	defer ee.mu.Unlock()
	ee.failedTestsList[testID] = true
}

// Stop gracefully stops the execution engine.
func (ee *ExecutionEngine) Stop(timeout time.Duration) error {
	done := make(chan struct{})

	go func() {
		ee.cancel()
		ee.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}

// GetMetrics returns current execution metrics.
func (ee *ExecutionEngine) GetMetrics() TestMetrics {
	return TestMetrics{
		Total:   ee.totalTests.Load(),
		Passed:  ee.passedTests.Load(),
		Failed:  ee.failedTests.Load(),
		Skipped: ee.skippedTests.Load(),
		Retries: ee.retryAttempts.Load(),
	}
}

// String returns a string representation of test metrics.
func (tm TestMetrics) String() string {
	return fmt.Sprintf("Tests: %d total, %d passed, %d failed, %d skipped (retries: %d)",
		tm.Total, tm.Passed, tm.Failed, tm.Skipped, tm.Retries)
}

// NewExecutionPlan creates an execution plan with dependency resolution.
func NewExecutionPlan(tests []*Test) (*ExecutionPlan, error) {
	plan := &ExecutionPlan{
		tests:        tests,
		dependencies: make(map[string][]string),
		completed:    make(map[string]bool),
		failed:       make(map[string]bool),
	}

	// Build dependency graph
	for _, test := range tests {
		plan.dependencies[test.ID] = test.Dependencies
	}

	// Validate no circular dependencies
	if err := plan.detectCircularDependencies(); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetExecutionOrder returns tests in dependency-respecting order (topological sort).
func (ep *ExecutionPlan) GetExecutionOrder() ([]*Test, error) {
	ep.mu.RLock()
	defer ep.mu.RUnlock()

	// Build in-degree map
	inDegree := make(map[string]int)
	for _, test := range ep.tests {
		if _, exists := inDegree[test.ID]; !exists {
			inDegree[test.ID] = 0
		}
		for _, dep := range test.Dependencies {
			inDegree[dep]++
		}
	}

	// Topological sort using Kahn's algorithm
	var ordered []*Test
	queue := make([]*Test, 0)

	// Find tests with no dependencies
	for _, test := range ep.tests {
		if len(test.Dependencies) == 0 {
			queue = append(queue, test)
		}
	}

	testMap := make(map[string]*Test)
	for _, test := range ep.tests {
		testMap[test.ID] = test
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		ordered = append(ordered, current)

		// Find tests that depend on current
		for _, test := range ep.tests {
			for _, dep := range test.Dependencies {
				if dep == current.ID {
					inDegree[test.ID]--
					if inDegree[test.ID] == 0 {
						queue = append(queue, test)
					}
				}
			}
		}
	}

	if len(ordered) != len(ep.tests) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return ordered, nil
}

// GetReadyTests returns tests whose dependencies are satisfied.
func (ep *ExecutionPlan) GetReadyTests() []*Test {
	ep.mu.RLock()
	defer ep.mu.RUnlock()

	var ready []*Test
	for _, test := range ep.tests {
		if ep.completed[test.ID] || ep.failed[test.ID] {
			continue
		}

		allDepsSatisfied := true
		for _, depID := range test.Dependencies {
			if !ep.completed[depID] || ep.failed[depID] {
				allDepsSatisfied = false
				break
			}
		}

		if allDepsSatisfied {
			ready = append(ready, test)
		}
	}

	return ready
}

// MarkCompleted marks a test as completed.
func (ep *ExecutionPlan) MarkCompleted(testID string) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	ep.completed[testID] = true
}

// MarkFailed marks a test as failed.
func (ep *ExecutionPlan) MarkFailed(testID string) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	ep.failed[testID] = true
}

// detectCircularDependencies checks for circular dependencies using DFS.
func (ep *ExecutionPlan) detectCircularDependencies() error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, test := range ep.tests {
		if !visited[test.ID] {
			if ep.hasCycle(test.ID, visited, recStack) {
				return fmt.Errorf("circular dependency detected involving test: %s", test.ID)
			}
		}
	}

	return nil
}

// hasCycle performs DFS to detect cycles.
func (ep *ExecutionPlan) hasCycle(testID string, visited, recStack map[string]bool) bool {
	visited[testID] = true
	recStack[testID] = true

	for _, depID := range ep.dependencies[testID] {
		if !visited[depID] {
			if ep.hasCycle(depID, visited, recStack) {
				return true
			}
		} else if recStack[depID] {
			return true
		}
	}

	recStack[testID] = false
	return false
}
