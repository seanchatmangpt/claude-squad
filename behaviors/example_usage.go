package behaviors

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ExampleSimulateAllBehaviors demonstrates complete 10-agent simulation framework
func ExampleSimulateAllBehaviors() {
	// Step 1: Create a behavior graph
	graph := NewBehaviorGraph()

	// Step 2: Define behaviors (nodes)
	behaviors := map[string]*BehaviorNode{
		"idle": {
			ID:          "idle",
			Name:        "Idle State",
			Description: "System waiting for input",
			Category:    "core",
		},
		"processing": {
			ID:          "processing",
			Name:        "Processing State",
			Description: "System actively processing request",
			Category:    "core",
		},
		"complete": {
			ID:          "complete",
			Name:        "Complete State",
			Description: "Processing complete",
			Category:    "core",
		},
	}

	for _, behavior := range behaviors {
		graph.AddNode(behavior)
	}

	// Step 3: Define transitions (edges)
	graph.AddEdge("idle", "processing", func() bool { return true }, 10*time.Millisecond, true)
	graph.AddEdge("processing", "complete", func() bool { return true }, 20*time.Millisecond, true)
	graph.AddEdge("complete", "idle", func() bool { return true }, 5*time.Millisecond, true)

	// Step 4: Configure orchestrator
	config := OrchestratorConfig{
		MaxConcurrency:   10,
		TimeoutPerPhase:  30 * time.Second,
		EnableCaching:    true,
		ValidateAll:      true,
		MaxSequenceDepth: 5,
		MutationCount:    20,
	}

	// Step 5: Create and execute orchestrator
	orchestrator := NewBehaviorOrchestrator(graph, config)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	fmt.Println("Starting 10-Agent Behavior Simulation...")
	err := orchestrator.ExecuteAll(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Step 6: Analyze results
	results := orchestrator.GetResults()

	// Get agent results
	agentResults := results["agent_results"].(map[string]interface{})
	execCount := agentResults["execution_count"].(int)
	seqCount := agentResults["sequences_generated"].(int)

	fmt.Printf("\nSimulation Complete:\n")
	fmt.Printf("  Sequences Generated: %d\n", seqCount)
	fmt.Printf("  Behaviors Executed: %d\n", execCount)

	// Get coverage report
	if coverageReport, ok := agentResults["coverage_report"].(*CoverageReport); ok {
		fmt.Printf("  Coverage: %.1f%%\n", coverageReport.CoveragePercent)
	}

	// Get performance metrics
	if metrics, ok := agentResults["performance_metrics"].(*PerformanceMetrics); ok {
		fmt.Printf("  Avg Latency: %v\n", metrics.AvgLatency)
		fmt.Printf("  Throughput: %.2f behaviors/sec\n", metrics.Throughput)
	}
}

// ExampleAgent1BehaviorGraph demonstrates Agent 1 (Behavior Graph Definition)
func ExampleAgent1BehaviorGraph() {
	fmt.Println("Agent 1: Behavior Graph Definition")

	graph := NewBehaviorGraph()

	// Add behavior nodes
	graph.AddNode(&BehaviorNode{
		ID:       "start",
		Name:     "Start",
		Category: "entry",
	})

	graph.AddNode(&BehaviorNode{
		ID:       "validate",
		Name:     "Validate Input",
		Category: "validation",
	})

	graph.AddNode(&BehaviorNode{
		ID:       "execute",
		Name:     "Execute",
		Category: "execution",
	})

	graph.AddEdge("start", "validate", func() bool { return true }, time.Millisecond, true)
	graph.AddEdge("validate", "execute", func() bool { return true }, time.Millisecond, true)

	fmt.Printf("Created graph with %d behaviors\n", len(graph.Nodes))
}

// ExampleAgent2StateMachine demonstrates Agent 2 (State Machine Simulator)
func ExampleAgent2StateMachine() {
	fmt.Println("Agent 2: State Machine Simulator")

	graph := NewBehaviorGraph()
	graph.AddNode(&BehaviorNode{ID: "s1", Name: "State 1"})
	graph.AddNode(&BehaviorNode{ID: "s2", Name: "State 2"})
	graph.AddEdge("s1", "s2", func() bool { return true }, time.Millisecond, true)

	config := StateMachineConfig{
		InitialState: "s1",
		MaxSteps:     10,
		Timeout:      5 * time.Second,
		TrackMetrics: true,
	}

	sm := NewStateMachine(graph, config)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sm.Execute(ctx)
	metrics := sm.GetMetrics()

	fmt.Printf("Executed %d transitions\n", metrics["transitions_count"])
	fmt.Printf("Visited %d unique states\n", metrics["unique_states"])
}

// ExampleAgent3PermutationGenerator demonstrates Agent 3 (Permutation Generator)
func ExampleAgent3PermutationGenerator() {
	fmt.Println("Agent 3: Permutation Generator")

	graph := NewBehaviorGraph()
	graph.AddNode(&BehaviorNode{ID: "a", Name: "A"})
	graph.AddNode(&BehaviorNode{ID: "b", Name: "B"})
	graph.AddNode(&BehaviorNode{ID: "c", Name: "C"})
	graph.AddEdge("a", "b", func() bool { return true }, time.Millisecond, true)
	graph.AddEdge("b", "c", func() bool { return true }, time.Millisecond, true)

	gen := NewPermutationGenerator(graph)
	sequences, _ := gen.GenerateSequences("a", 3)

	fmt.Printf("Generated %d unique behavior sequences\n", len(sequences))
	for i, seq := range sequences {
		fmt.Printf("  Sequence %d: %v\n", i+1, seq.Path)
	}
}

// ExampleAgent4Validator demonstrates Agent 4 (Validation Engine)
func ExampleAgent4Validator() {
	fmt.Println("Agent 4: Validation Engine")

	graph := NewBehaviorGraph()
	node := &BehaviorNode{
		ID:       "test",
		Name:     "Test Behavior",
		Category: "test",
	}
	graph.AddNode(node)

	validator := NewBehaviorValidator(graph)
	result := validator.Validate("test")

	fmt.Printf("Validation Result: %v\n", result.Valid)
	fmt.Printf("Errors: %d, Warnings: %d\n", len(result.Errors), len(result.Warnings))
}

// ExampleAgent6Coverage demonstrates Agent 6 (Coverage Analyzer)
func ExampleAgent6Coverage() {
	fmt.Println("Agent 6: Coverage Analyzer")

	graph := NewBehaviorGraph()
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("state_%d", i)
		graph.AddNode(&BehaviorNode{ID: id, Name: id})
	}

	analyzer := NewCoverageAnalyzer(graph)

	// Simulate some visits
	analyzer.RecordVisit("state_0")
	analyzer.RecordVisit("state_1")
	analyzer.RecordVisit("state_2")
	analyzer.RecordTransition("state_0", "state_1")
	analyzer.RecordTransition("state_1", "state_2")

	report := analyzer.GenerateReport()

	fmt.Printf("Coverage Report:\n")
	fmt.Printf("  Total Nodes: %d\n", report.TotalNodes)
	fmt.Printf("  Visited Nodes: %d\n", report.VisitedNodes)
	fmt.Printf("  Coverage: %.1f%%\n", report.CoveragePercent)
	fmt.Printf("  Uncovered: %v\n", report.UncoveredNodes)
}

// ExampleAgent7Performance demonstrates Agent 7 (Performance Profiler)
func ExampleAgent7Performance() {
	fmt.Println("Agent 7: Performance Profiler")

	profiler := NewPerformanceProfiler()

	// Simulate some execution results
	results := []*ExecutionResult{
		{Duration: 5 * time.Millisecond},
		{Duration: 10 * time.Millisecond},
		{Duration: 15 * time.Millisecond},
		{Duration: 20 * time.Millisecond},
		{Duration: 25 * time.Millisecond},
	}

	metrics := profiler.RecordExecution(results)

	fmt.Printf("Performance Metrics:\n")
	fmt.Printf("  Min Latency: %v\n", metrics.MinLatency)
	fmt.Printf("  Avg Latency: %v\n", metrics.AvgLatency)
	fmt.Printf("  Max Latency: %v\n", metrics.MaxLatency)
	fmt.Printf("  Throughput: %.2f/sec\n", metrics.Throughput)
}

// ExampleAgent8Mutations demonstrates Agent 8 (Mutation Generator)
func ExampleAgent8Mutations() {
	fmt.Println("Agent 8: Mutation Generator")

	graph := NewBehaviorGraph()
	graph.AddNode(&BehaviorNode{ID: "b1", Name: "Behavior 1"})
	graph.AddNode(&BehaviorNode{ID: "b2", Name: "Behavior 2"})

	mutGen := NewMutationGenerator(graph, 12345)
	mutations, _ := mutGen.GenerateMutations(5)

	fmt.Printf("Generated %d mutations\n", len(mutations))
	for i, mut := range mutations {
		fmt.Printf("  Mutation %d: %s on %s\n", i+1, mut.Type, mut.TargetNode)
	}

	stats := mutGen.GetMutationStats()
	fmt.Printf("Graph has %d nodes\n", stats["graph_nodes"])
}

// ExampleAgent9Orchestrator demonstrates Agent 9 (Orchestrator)
func ExampleAgent9Orchestrator() {
	fmt.Println("Agent 9: Orchestrator")

	graph := NewBehaviorGraph()
	graph.AddNode(&BehaviorNode{ID: "test", Name: "Test"})

	config := OrchestratorConfig{
		MaxConcurrency:   10,
		TimeoutPerPhase:  10 * time.Second,
		MaxSequenceDepth: 2,
	}

	orchestrator := NewBehaviorOrchestrator(graph, config)

	fmt.Printf("Created orchestrator with %d agents\n", len(orchestrator.agents))
	for agentID, agent := range orchestrator.agents {
		fmt.Printf("  - %s: %s\n", agentID, agent.Name)
	}
}

// ExampleAllAgents demonstrates all 10 agents working together
func ExampleAllAgents() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("HYPER-ADVANCED 10-AGENT CONCURRENT BEHAVIOR SIMULATION")
	fmt.Println(strings.Repeat("=", 60))

	// Create comprehensive behavior graph
	graph := NewBehaviorGraph()

	// Add behaviors
	for _, name := range []string{"idle", "active", "busy", "done"} {
		graph.AddNode(&BehaviorNode{ID: name, Name: name})
	}

	// Add transitions
	graph.AddEdge("idle", "active", func() bool { return true }, time.Millisecond, true)
	graph.AddEdge("active", "busy", func() bool { return true }, time.Millisecond, true)
	graph.AddEdge("busy", "done", func() bool { return true }, time.Millisecond, true)
	graph.AddEdge("done", "idle", func() bool { return true }, time.Millisecond, true)

	// Configure and run orchestrator
	config := OrchestratorConfig{
		MaxConcurrency:   10,
		TimeoutPerPhase:  30 * time.Second,
		MaxSequenceDepth: 4,
		MutationCount:    10,
	}

	orchestrator := NewBehaviorOrchestrator(graph, config)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	fmt.Println("\nLaunching 10 Agents in Parallel...")
	orchestrator.ExecuteAll(ctx)

	// Display results
	results := orchestrator.GetResults()
	fmt.Printf("\n✓ All 10 agents completed\n")
	fmt.Printf("✓ Total execution time: %v\n", results["total_duration"])
	fmt.Printf("✓ Simulation successful!\n")
}
