// Package jtbd provides assertions for JTBD (Jobs-to-be-Done) test validation.
package jtbd

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AssertionResult represents the outcome of a single assertion.
type AssertionResult struct {
	Pass     bool        `json:"pass"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
	Message  string      `json:"message"`
	Diff     string      `json:"diff,omitempty"`
}

// AssertionChain allows fluent chaining of multiple assertions.
type AssertionChain struct {
	mu           sync.RWMutex
	results      []AssertionResult
	errors       []error
	failOnError  bool
}

// AssertionReport aggregates multiple assertion results with statistics.
type AssertionReport struct {
	mu            sync.RWMutex
	TotalTests    int               `json:"total_tests"`
	PassedTests   int               `json:"passed_tests"`
	FailedTests   int               `json:"failed_tests"`
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	Duration      time.Duration     `json:"duration"`
	Results       []AssertionResult `json:"results"`
	Errors        []string          `json:"errors"`
}

// AssertionConstraint defines a constraint to validate against.
type AssertionConstraint struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"` // max, min, equals, range, contains
	Value  interface{} `json:"value"`
	Min    interface{} `json:"min,omitempty"`
	Max    interface{} `json:"max,omitempty"`
	Strict bool        `json:"strict"` // Strict comparison vs fuzzy
}

// Expectations defines what constitutes success for a job.
type Expectations struct {
	FunctionalCriteria []string               `json:"functional_criteria"`
	EmotionalCriteria  []string               `json:"emotional_criteria"`
	SocialCriteria     []string               `json:"social_criteria"`
	Metrics            map[string]interface{} `json:"metrics"`
	MaxDuration        time.Duration          `json:"max_duration"`
}

// ProgressSnapshot represents a point-in-time progress measurement.
type ProgressSnapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

// Money represents a monetary value with currency.
type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Result represents the outcome of a job execution.
type Result struct {
	JobID     string                 `json:"job_id"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data"`
	Errors    []error                `json:"-"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
}

// ProgressTracker tracks progress over time with checkpoints.
type ProgressTracker struct {
	mu          sync.RWMutex
	snapshots   map[string]ProgressSnapshot
	checkpoints map[string]time.Time
}

// NewProgressTracker creates a new progress tracker.
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		snapshots:   make(map[string]ProgressSnapshot),
		checkpoints: make(map[string]time.Time),
	}
}

// RecordProgress records a progress snapshot.
func (pt *ProgressTracker) RecordProgress(name string, values map[string]interface{}) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.snapshots[name] = ProgressSnapshot{
		Timestamp: time.Now(),
		Values:    values,
	}
}

// RecordCheckpoint records a time-based checkpoint.
func (pt *ProgressTracker) RecordCheckpoint(name string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.checkpoints[name] = time.Now()
}

// GetProgress retrieves a progress snapshot.
func (pt *ProgressTracker) GetProgress(name string) (ProgressSnapshot, bool) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	snapshot, ok := pt.snapshots[name]
	return snapshot, ok
}

// GetCheckpoint retrieves a checkpoint time.
func (pt *ProgressTracker) GetCheckpoint(name string) (time.Time, bool) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	t, ok := pt.checkpoints[name]
	return t, ok
}

// AllIndicators returns all recorded progress indicators.
func (pt *ProgressTracker) AllIndicators() map[string]ProgressSnapshot {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	result := make(map[string]ProgressSnapshot, len(pt.snapshots))
	for k, v := range pt.snapshots {
		result[k] = v
	}
	return result
}

// AssertJobCompleted validates that a job was completed successfully.
func AssertJobCompleted(ctx context.Context, job *Job) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}
	if job.ID == "" {
		return fmt.Errorf("job ID is empty")
	}
	if job.Name == "" {
		return fmt.Errorf("job name is empty")
	}
	if len(job.Outcomes) == 0 {
		return fmt.Errorf("job has no outcomes defined")
	}
	return nil
}

// AssertProgressMade compares two progress snapshots and validates progress.
func AssertProgressMade(before, after ProgressSnapshot) error {
	if after.Timestamp.Before(before.Timestamp) {
		return fmt.Errorf("'after' snapshot (%v) is before 'before' snapshot (%v)",
			after.Timestamp, before.Timestamp)
	}

	// Compare values
	for key, beforeVal := range before.Values {
		afterVal, exists := after.Values[key]
		if !exists {
			return fmt.Errorf("metric '%s' missing in 'after' snapshot", key)
		}

		// Try numeric comparison
		beforeNum, beforeOK := toFloat64(beforeVal)
		afterNum, afterOK := toFloat64(afterVal)
		if beforeOK && afterOK {
			if afterNum <= beforeNum {
				return fmt.Errorf("no progress made for '%s': before=%.2f, after=%.2f",
					key, beforeNum, afterNum)
			}
		}
	}

	return nil
}

// AssertWithinConstraints validates results against constraints.
func AssertWithinConstraints(result Result, constraints []AssertionConstraint) error {
	for _, constraint := range constraints {
		value, exists := result.Data[constraint.Name]
		if !exists {
			return fmt.Errorf("constraint '%s' not found in result", constraint.Name)
		}

		switch constraint.Type {
		case "max":
			num, ok := toFloat64(value)
			maxNum, maxOK := toFloat64(constraint.Value)
			if !ok || !maxOK {
				return fmt.Errorf("cannot compare non-numeric values for max constraint")
			}
			if num > maxNum {
				return fmt.Errorf("'%s' exceeds max: %.2f > %.2f", constraint.Name, num, maxNum)
			}

		case "min":
			num, ok := toFloat64(value)
			minNum, minOK := toFloat64(constraint.Value)
			if !ok || !minOK {
				return fmt.Errorf("cannot compare non-numeric values for min constraint")
			}
			if num < minNum {
				return fmt.Errorf("'%s' below min: %.2f < %.2f", constraint.Name, num, minNum)
			}

		case "equals":
			if value != constraint.Value {
				return fmt.Errorf("'%s' does not equal expected: got %v, want %v",
					constraint.Name, value, constraint.Value)
			}

		case "range":
			num, ok := toFloat64(value)
			minNum, minOK := toFloat64(constraint.Min)
			maxNum, maxOK := toFloat64(constraint.Max)
			if !ok || !minOK || !maxOK {
				return fmt.Errorf("cannot perform range check on non-numeric values")
			}
			if num < minNum || num > maxNum {
				return fmt.Errorf("'%s' out of range: %.2f not in [%.2f, %.2f]",
					constraint.Name, num, minNum, maxNum)
			}

		case "contains":
			strValue, ok := value.(string)
			strConstraint, cOK := constraint.Value.(string)
			if !ok || !cOK {
				return fmt.Errorf("'contains' constraint requires string values")
			}
			if !stringContains(strValue, strConstraint) {
				return fmt.Errorf("'%s' does not contain '%s'", constraint.Name, strConstraint)
			}

		default:
			return fmt.Errorf("unknown constraint type: %s", constraint.Type)
		}
	}
	return nil
}

// AssertSatisfaction validates job satisfaction against expectations.
func AssertSatisfaction(ctx context.Context, job *Job, expectations Expectations) error {
	// Validate functional criteria
	if job.Functional == "" && len(expectations.FunctionalCriteria) > 0 {
		return fmt.Errorf("job has no functional dimension but expectations require it")
	}

	// Validate emotional criteria
	if job.Emotional == "" && len(expectations.EmotionalCriteria) > 0 {
		return fmt.Errorf("job has no emotional dimension but expectations require it")
	}

	// Validate social criteria
	if job.Social == "" && len(expectations.SocialCriteria) > 0 {
		return fmt.Errorf("job has no social dimension but expectations require it")
	}

	// Validate metrics
	for metricName := range expectations.Metrics {
		found := false
		for _, outcome := range job.Outcomes {
			if outcome.Metric == metricName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected metric '%s' not found in job outcomes", metricName)
		}
	}

	return nil
}

// AssertTimeCompliance validates execution time is within limits.
func AssertTimeCompliance(duration, limit time.Duration) error {
	if duration > limit {
		pct := float64(duration) / float64(limit) * 100.0
		return fmt.Errorf("execution time exceeded limit: %v > %v (%.1f%%)",
			duration, limit, pct)
	}
	return nil
}

// AssertCostCompliance validates spending is within budget.
func AssertCostCompliance(spent, budget Money) error {
	if spent.Currency != budget.Currency {
		return fmt.Errorf("currency mismatch: spent=%s, budget=%s",
			spent.Currency, budget.Currency)
	}
	if spent.Amount > budget.Amount {
		pct := (spent.Amount / budget.Amount) * 100.0
		return fmt.Errorf("spending exceeded budget: %.2f > %.2f (%.1f%%)",
			spent.Amount, budget.Amount, pct)
	}
	return nil
}

// NewAssertionChain creates a new assertion chain.
func NewAssertionChain() *AssertionChain {
	return &AssertionChain{
		results:     make([]AssertionResult, 0),
		errors:      make([]error, 0),
		failOnError: false,
	}
}

// Add adds an assertion result to the chain.
func (ac *AssertionChain) Add(result AssertionResult) *AssertionChain {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.results = append(ac.results, result)
	return ac
}

// AddError adds an error to the chain.
func (ac *AssertionChain) AddError(err error) *AssertionChain {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if err != nil {
		ac.errors = append(ac.errors, err)
	}
	return ac
}

// WithFailOnError sets whether to fail on first error.
func (ac *AssertionChain) WithFailOnError(fail bool) *AssertionChain {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.failOnError = fail
	return ac
}

// Results returns all assertion results.
func (ac *AssertionChain) Results() []AssertionResult {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	results := make([]AssertionResult, len(ac.results))
	copy(results, ac.results)
	return results
}

// Errors returns all accumulated errors.
func (ac *AssertionChain) Errors() []error {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	errors := make([]error, len(ac.errors))
	copy(errors, ac.errors)
	return errors
}

// IsValid returns true if all assertions passed.
func (ac *AssertionChain) IsValid() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	if len(ac.errors) > 0 {
		return false
	}

	for _, result := range ac.results {
		if !result.Pass {
			return false
		}
	}

	return true
}

// String returns a string representation of the assertion chain.
func (ac *AssertionChain) String() string {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	passed := 0
	for _, result := range ac.results {
		if result.Pass {
			passed++
		}
	}

	return fmt.Sprintf("AssertionChain: %d/%d passed, %d errors",
		passed, len(ac.results), len(ac.errors))
}

// NewAssertionReport creates a new assertion report.
func NewAssertionReport() *AssertionReport {
	return &AssertionReport{
		StartTime: time.Now(),
		Results:   make([]AssertionResult, 0),
		Errors:    make([]string, 0),
	}
}

// AddResult adds an assertion result to the report.
func (ar *AssertionReport) AddResult(result AssertionResult) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.Results = append(ar.Results, result)
	ar.TotalTests++
	if result.Pass {
		ar.PassedTests++
	} else {
		ar.FailedTests++
	}
}

// AddError adds an error to the report.
func (ar *AssertionReport) AddError(err error) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if err != nil {
		ar.Errors = append(ar.Errors, err.Error())
	}
}

// Complete marks the report as complete and calculates duration.
func (ar *AssertionReport) Complete() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.EndTime = time.Now()
	ar.Duration = ar.EndTime.Sub(ar.StartTime)
}

// IsSuccessful returns true if all tests passed.
func (ar *AssertionReport) IsSuccessful() bool {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	return ar.FailedTests == 0 && len(ar.Errors) == 0
}

// PassRate returns the percentage of tests that passed.
func (ar *AssertionReport) PassRate() float64 {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	if ar.TotalTests == 0 {
		return 0.0
	}
	return float64(ar.PassedTests) / float64(ar.TotalTests) * 100.0
}

// Summary returns a summary string of the report.
func (ar *AssertionReport) Summary() string {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	return fmt.Sprintf("Tests: %d total, %d passed, %d failed (%.1f%% pass rate) in %v",
		ar.TotalTests, ar.PassedTests, ar.FailedTests, ar.PassRate(), ar.Duration)
}

// Helper functions

// toFloat64 converts various numeric types to float64.
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}

// stringContains checks if a string contains a substring (simple implementation).
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

// findSubstring searches for substring in s.
func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
