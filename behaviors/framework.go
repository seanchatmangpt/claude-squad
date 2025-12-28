// Package behaviors provides a comprehensive 10-agent framework for simulating
// all possible behaviors deterministically without LLM calls.
package behaviors

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// AGENT 1: Behavior Graph Definition
// ============================================================================

// BehaviorNode represents a distinct behavior or state in the system
type BehaviorNode struct {
	ID          string
	Name        string
	Description string
	Category    string
	Constraints []string
	Metadata    map[string]interface{}
}

// BehaviorEdge represents a transition from one behavior to another
type BehaviorEdge struct {
	From       string
	To         string
	Condition  func() bool
	Weight     int // For prioritization
	Latency    time.Duration
	Deterministic bool
}

// BehaviorGraph models the complete behavior space as a directed graph
type BehaviorGraph struct {
	mu    sync.RWMutex
	Nodes map[string]*BehaviorNode
	Edges map[string][]*BehaviorEdge
}

// NewBehaviorGraph creates an empty behavior graph
func NewBehaviorGraph() *BehaviorGraph {
	return &BehaviorGraph{
		Nodes: make(map[string]*BehaviorNode),
		Edges: make(map[string][]*BehaviorEdge),
	}
}

// AddNode adds a behavior node to the graph
func (bg *BehaviorGraph) AddNode(node *BehaviorNode) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	if _, exists := bg.Nodes[node.ID]; exists {
		return fmt.Errorf("node %s already exists", node.ID)
	}
	bg.Nodes[node.ID] = node
	bg.Edges[node.ID] = []*BehaviorEdge{}
	return nil
}

// AddEdge adds a transition from one behavior to another
func (bg *BehaviorGraph) AddEdge(from, to string, cond func() bool, latency time.Duration, deterministic bool) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	if _, exists := bg.Nodes[from]; !exists {
		return fmt.Errorf("source node %s does not exist", from)
	}
	if _, exists := bg.Nodes[to]; !exists {
		return fmt.Errorf("target node %s does not exist", to)
	}

	edge := &BehaviorEdge{
		From:          from,
		To:            to,
		Condition:     cond,
		Weight:        1,
		Latency:       latency,
		Deterministic: deterministic,
	}
	bg.Edges[from] = append(bg.Edges[from], edge)
	return nil
}

// GetSuccessors returns all valid next behaviors from a given node
func (bg *BehaviorGraph) GetSuccessors(nodeID string) ([]*BehaviorEdge, error) {
	bg.mu.RLock()
	defer bg.mu.RUnlock()

	if _, exists := bg.Nodes[nodeID]; !exists {
		return nil, fmt.Errorf("node %s does not exist", nodeID)
	}

	successors := bg.Edges[nodeID]
	valid := make([]*BehaviorEdge, 0)
	for _, edge := range successors {
		if edge.Condition == nil || edge.Condition() {
			valid = append(valid, edge)
		}
	}
	return valid, nil
}

// ============================================================================
// AGENT 2: State Machine Simulator
// ============================================================================

// StateMachineConfig defines the configuration for state machine execution
type StateMachineConfig struct {
	InitialState string
	MaxSteps     int
	Timeout      time.Duration
	TrackMetrics bool
}

// StateTransition represents a single state change
type StateTransition struct {
	From      string
	To        string
	Timestamp time.Time
	Latency   time.Duration
}

// StateMachine executes behavior transitions according to the graph
type StateMachine struct {
	mu           sync.RWMutex
	graph        *BehaviorGraph
	current      string
	transitions  []StateTransition
	visited      map[string]int
	config       StateMachineConfig
	startTime    time.Time
}

// NewStateMachine creates a new state machine for the behavior graph
func NewStateMachine(bg *BehaviorGraph, config StateMachineConfig) *StateMachine {
	return &StateMachine{
		graph:       bg,
		current:     config.InitialState,
		transitions: make([]StateTransition, 0),
		visited:     make(map[string]int),
		config:      config,
		startTime:   time.Now(),
	}
}

// Execute runs the state machine for the configured duration
func (sm *StateMachine) Execute(ctx context.Context) error {
	sm.mu.Lock()
	sm.visited[sm.current]++
	sm.mu.Unlock()

	steps := 0
	for steps < sm.config.MaxSteps {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		successors, err := sm.graph.GetSuccessors(sm.current)
		if err != nil {
			return err
		}

		if len(successors) == 0 {
			break // Dead end state
		}

		// Choose next transition (deterministically, first valid edge)
		edge := successors[0]
		startTime := time.Now()

		// Simulate latency
		if edge.Latency > 0 {
			time.Sleep(edge.Latency)
		}

		latency := time.Since(startTime)
		sm.mu.Lock()
		sm.transitions = append(sm.transitions, StateTransition{
			From:      sm.current,
			To:        edge.To,
			Timestamp: time.Now(),
			Latency:   latency,
		})
		sm.current = edge.To
		sm.visited[sm.current]++
		sm.mu.Unlock()

		steps++
	}

	return nil
}

// GetMetrics returns execution metrics
func (sm *StateMachine) GetMetrics() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	totalLatency := time.Duration(0)
	for _, t := range sm.transitions {
		totalLatency += t.Latency
	}

	avgLatency := time.Duration(0)
	if len(sm.transitions) > 0 {
		avgLatency = totalLatency / time.Duration(len(sm.transitions))
	}

	return map[string]interface{}{
		"current_state":     sm.current,
		"transitions_count": len(sm.transitions),
		"unique_states":     len(sm.visited),
		"visited_states":    sm.visited,
		"total_latency":     totalLatency,
		"avg_latency":       avgLatency,
		"execution_time":    time.Since(sm.startTime),
	}
}

// ============================================================================
// AGENT 3: Permutation Generator
// ============================================================================

// BehaviorSequence represents a sequence of behaviors
type BehaviorSequence struct {
	Path      []string
	Cost      int
	Valid     bool
	Coverage  float64
	Timestamp time.Time
}

// PermutationGenerator generates all valid behavior sequences
type PermutationGenerator struct {
	mu    sync.Mutex
	graph *BehaviorGraph
	cache map[string][]*BehaviorSequence
}

// NewPermutationGenerator creates a new permutation generator
func NewPermutationGenerator(bg *BehaviorGraph) *PermutationGenerator {
	return &PermutationGenerator{
		graph: bg,
		cache: make(map[string][]*BehaviorSequence),
	}
}

// GenerateSequences generates all valid behavior sequences up to maxDepth
func (pg *PermutationGenerator) GenerateSequences(startNode string, maxDepth int) ([]*BehaviorSequence, error) {
	pg.mu.Lock()
	if cached, ok := pg.cache[startNode]; ok && len(cached) > 0 {
		pg.mu.Unlock()
		return cached, nil
	}
	pg.mu.Unlock()

	sequences := make([]*BehaviorSequence, 0)
	visited := make(map[string]bool)

	pg.generateSequencesRecursive(startNode, []string{startNode}, maxDepth, &sequences, visited)

	pg.mu.Lock()
	pg.cache[startNode] = sequences
	pg.mu.Unlock()

	return sequences, nil
}

func (pg *PermutationGenerator) generateSequencesRecursive(current string, path []string, depth int, sequences *[]*BehaviorSequence, visited map[string]bool) {
	if depth == 0 {
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		*sequences = append(*sequences, &BehaviorSequence{
			Path:      pathCopy,
			Cost:      len(pathCopy),
			Valid:     true,
			Timestamp: time.Now(),
		})
		return
	}

	successors, err := pg.graph.GetSuccessors(current)
	if err != nil || len(successors) == 0 {
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		*sequences = append(*sequences, &BehaviorSequence{
			Path:      pathCopy,
			Cost:      len(pathCopy),
			Valid:     true,
			Timestamp: time.Now(),
		})
		return
	}

	for _, edge := range successors {
		newPath := append(path, edge.To)
		pg.generateSequencesRecursive(edge.To, newPath, depth-1, sequences, visited)
	}
}

// ============================================================================
// AGENT 4: Validation Engine
// ============================================================================

// ValidationResult represents the result of behavior validation
type ValidationResult struct {
	BehaviorID string
	Valid      bool
	Errors     []string
	Warnings   []string
	Timestamp  time.Time
}

// BehaviorValidator validates behaviors against constraints
type BehaviorValidator struct {
	mu         sync.RWMutex
	graph      *BehaviorGraph
	validators map[string]func(*BehaviorNode) error
}

// NewBehaviorValidator creates a new validator
func NewBehaviorValidator(bg *BehaviorGraph) *BehaviorValidator {
	return &BehaviorValidator{
		graph:      bg,
		validators: make(map[string]func(*BehaviorNode) error),
	}
}

// RegisterValidator registers a custom validation function
func (bv *BehaviorValidator) RegisterValidator(name string, fn func(*BehaviorNode) error) {
	bv.mu.Lock()
	defer bv.mu.Unlock()
	bv.validators[name] = fn
}

// Validate validates a behavior node
func (bv *BehaviorValidator) Validate(nodeID string) *ValidationResult {
	bv.mu.RLock()
	defer bv.mu.RUnlock()

	result := &ValidationResult{
		BehaviorID: nodeID,
		Valid:      true,
		Errors:     make([]string, 0),
		Warnings:   make([]string, 0),
		Timestamp:  time.Now(),
	}

	node, exists := bv.graph.Nodes[nodeID]
	if !exists {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("node %s does not exist", nodeID))
		return result
	}

	// Run registered validators
	for name, validator := range bv.validators {
		if err := validator(node); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", name, err))
		}
	}

	return result
}

// ============================================================================
// AGENT 5: Concurrent Executor
// ============================================================================

// ExecutionResult represents the result of behavior execution
type ExecutionResult struct {
	BehaviorID      string
	Success         bool
	Duration        time.Duration
	Error           error
	Metrics         map[string]interface{}
	StateTransitions []StateTransition
}

// ConcurrentExecutor runs multiple state machines in parallel
type ConcurrentExecutor struct {
	mu           sync.RWMutex
	graph        *BehaviorGraph
	config       StateMachineConfig
	results      []*ExecutionResult
	maxConcurrency int
}

// NewConcurrentExecutor creates a new concurrent executor
func NewConcurrentExecutor(bg *BehaviorGraph, config StateMachineConfig, maxConcurrency int) *ConcurrentExecutor {
	return &ConcurrentExecutor{
		graph:          bg,
		config:         config,
		results:        make([]*ExecutionResult, 0),
		maxConcurrency: maxConcurrency,
	}
}

// ExecuteAll runs all behavior sequences concurrently
func (ce *ConcurrentExecutor) ExecuteAll(ctx context.Context, sequences []*BehaviorSequence) ([]*ExecutionResult, error) {
	var wg sync.WaitGroup
	resultsChan := make(chan *ExecutionResult, len(sequences))
	semaphore := make(chan struct{}, ce.maxConcurrency)

	for _, seq := range sequences {
		wg.Add(1)
		go func(sequence *BehaviorSequence) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			startTime := time.Now()
			sm := NewStateMachine(ce.graph, ce.config)
			err := sm.Execute(ctx)

			result := &ExecutionResult{
				BehaviorID:       fmt.Sprintf("%v", sequence.Path),
				Success:          err == nil,
				Duration:         time.Since(startTime),
				Error:            err,
				Metrics:          sm.GetMetrics(),
			}
			resultsChan <- result
		}(seq)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	results := make([]*ExecutionResult, 0)
	for result := range resultsChan {
		results = append(results, result)
	}

	ce.mu.Lock()
	ce.results = append(ce.results, results...)
	ce.mu.Unlock()

	return results, nil
}

// ============================================================================
// AGENT 6: Coverage Analyzer
// ============================================================================

// CoverageReport represents coverage analysis results
type CoverageReport struct {
	TotalNodes        int
	VisitedNodes      int
	CoveragePercent   float64
	UncoveredNodes    []string
	EdgeCoverage      map[string]int
	SequenceCoverage  float64
	Timestamp         time.Time
}

// CoverageAnalyzer analyzes behavior space coverage
type CoverageAnalyzer struct {
	mu            sync.RWMutex
	graph         *BehaviorGraph
	visitedNodes  map[string]int
	edgeCoverage  map[string]map[string]int
}

// NewCoverageAnalyzer creates a new coverage analyzer
func NewCoverageAnalyzer(bg *BehaviorGraph) *CoverageAnalyzer {
	return &CoverageAnalyzer{
		graph:        bg,
		visitedNodes: make(map[string]int),
		edgeCoverage: make(map[string]map[string]int),
	}
}

// RecordVisit records a node visit
func (ca *CoverageAnalyzer) RecordVisit(nodeID string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	ca.visitedNodes[nodeID]++
}

// RecordTransition records an edge traversal
func (ca *CoverageAnalyzer) RecordTransition(from, to string) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if ca.edgeCoverage[from] == nil {
		ca.edgeCoverage[from] = make(map[string]int)
	}
	ca.edgeCoverage[from][to]++
}

// GenerateReport generates a coverage report
func (ca *CoverageAnalyzer) GenerateReport() *CoverageReport {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	report := &CoverageReport{
		TotalNodes:     len(ca.graph.Nodes),
		VisitedNodes:   len(ca.visitedNodes),
		UncoveredNodes: make([]string, 0),
		EdgeCoverage:   make(map[string]int),
		Timestamp:      time.Now(),
	}

	for nodeID := range ca.graph.Nodes {
		if _, visited := ca.visitedNodes[nodeID]; !visited {
			report.UncoveredNodes = append(report.UncoveredNodes, nodeID)
		}
	}

	report.CoveragePercent = float64(report.VisitedNodes) / float64(report.TotalNodes) * 100

	totalEdges := 0
	totalTransitions := 0
	for from, transitions := range ca.edgeCoverage {
		for _, count := range transitions {
			totalEdges++
			totalTransitions += count
		}
		for to, count := range transitions {
			report.EdgeCoverage[fmt.Sprintf("%s->%s", from, to)] = count
		}
	}

	if totalEdges > 0 {
		report.SequenceCoverage = float64(totalTransitions) / float64(totalEdges)
	}

	return report
}

// ============================================================================
// AGENT 7: Performance Profiler
// ============================================================================

// PerformanceMetrics represents performance measurements
type PerformanceMetrics struct {
	MinLatency       time.Duration
	MaxLatency       time.Duration
	AvgLatency       time.Duration
	P95Latency       time.Duration
	P99Latency       time.Duration
	Throughput       float64 // behaviors per second
	TotalDuration    time.Duration
	MemoryUsage      uint64
	GoroutineCount   int
	Timestamp        time.Time
}

// PerformanceProfiler measures execution performance
type PerformanceProfiler struct {
	mu      sync.RWMutex
	metrics []*PerformanceMetrics
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler() *PerformanceProfiler {
	return &PerformanceProfiler{
		metrics: make([]*PerformanceMetrics, 0),
	}
}

// RecordExecution records execution metrics
func (pp *PerformanceProfiler) RecordExecution(results []*ExecutionResult) *PerformanceMetrics {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if len(results) == 0 {
		return &PerformanceMetrics{Timestamp: time.Now()}
	}

	latencies := make([]time.Duration, 0)
	totalDuration := time.Duration(0)

	for _, result := range results {
		latencies = append(latencies, result.Duration)
		totalDuration += result.Duration
	}

	// Sort latencies for percentile calculation
	for i := 0; i < len(latencies); i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[j] < latencies[i] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}

	metrics := &PerformanceMetrics{
		MinLatency:    latencies[0],
		MaxLatency:    latencies[len(latencies)-1],
		AvgLatency:    totalDuration / time.Duration(len(results)),
		P95Latency:    latencies[int(float64(len(latencies))*0.95)],
		P99Latency:    latencies[int(float64(len(latencies))*0.99)],
		Throughput:    float64(len(results)) / totalDuration.Seconds(),
		TotalDuration: totalDuration,
		Timestamp:     time.Now(),
	}

	pp.metrics = append(pp.metrics, metrics)
	return metrics
}

// GetAverageMetrics returns average metrics across all recordings
func (pp *PerformanceProfiler) GetAverageMetrics() *PerformanceMetrics {
	pp.mu.RLock()
	defer pp.mu.RUnlock()

	if len(pp.metrics) == 0 {
		return &PerformanceMetrics{Timestamp: time.Now()}
	}

	totalMin := time.Duration(0)
	totalMax := time.Duration(0)
	totalAvg := time.Duration(0)
	totalThroughput := 0.0

	for _, m := range pp.metrics {
		totalMin += m.MinLatency
		totalMax += m.MaxLatency
		totalAvg += m.AvgLatency
		totalThroughput += m.Throughput
	}

	count := time.Duration(len(pp.metrics))
	return &PerformanceMetrics{
		MinLatency:    totalMin / count,
		MaxLatency:    totalMax / count,
		AvgLatency:    totalAvg / count,
		Throughput:    totalThroughput / float64(len(pp.metrics)),
		Timestamp:     time.Now(),
	}
}
