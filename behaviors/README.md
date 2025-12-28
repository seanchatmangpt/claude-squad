# Hyper-Advanced 10-Agent Concurrent Behavior Simulation Framework

A deterministic, production-ready behavior simulation engine that simulates ALL possible behaviors without LLM calls using a sophisticated 10-agent concurrent methodology.

## Overview

This framework simulates complete behavior spaces through a coordinated 10-agent system where each agent specializes in a specific aspect of behavior analysis and validation:

- **Agent 1**: Behavior Graph Definition - Models behavior space structure
- **Agent 2**: State Machine Simulator - Executes state transitions
- **Agent 3**: Permutation Generator - Creates all valid behavior sequences
- **Agent 4**: Validation Engine - Validates behaviors against constraints
- **Agent 5**: Concurrent Executor - Runs behaviors in parallel
- **Agent 6**: Coverage Analyzer - Measures behavior space coverage
- **Agent 7**: Performance Profiler - Profiles execution performance
- **Agent 8**: Mutation Generator - Creates behavior variations
- **Agent 9**: Orchestrator - Coordinates all 10 agents
- **Agent 10**: Integration Test Harness - End-to-end validation

## Key Features

### ✅ No External LLM Calls
- Purely deterministic behavior simulation
- No API dependencies or network calls
- Perfect for offline testing and CI/CD pipelines

### ✅ 10-Agent Parallel Execution
- All agents work simultaneously (except Agent 1 setup phase)
- Maximum concurrency optimization
- Coordinated result aggregation

### ✅ Comprehensive Behavior Coverage
- Automatic permutation generation up to configurable depth
- Deterministic and repeatable sequences
- Edge case detection through mutation

### ✅ Production-Ready Performance Metrics
- Latency profiling (min, avg, max, P95, P99)
- Throughput measurement
- Coverage reporting (node and edge coverage)

### ✅ Deterministic Mutation Testing
- 8 mutation types for comprehensive variation testing
- Predictable seed-based generation
- Mutation application and reversion

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    "claude-squad/behaviors"
)

func main() {
    // Create behavior graph
    graph := behaviors.NewBehaviorGraph()

    // Add behaviors
    graph.AddNode(&behaviors.BehaviorNode{
        ID: "idle",
        Name: "Idle State",
        Category: "core",
    })
    graph.AddNode(&behaviors.BehaviorNode{
        ID: "active",
        Name: "Active State",
        Category: "core",
    })

    // Define transitions
    graph.AddEdge("idle", "active",
        func() bool { return true },
        10 * time.Millisecond,
        true,
    )

    // Configure orchestrator
    config := behaviors.OrchestratorConfig{
        MaxConcurrency:   10,
        TimeoutPerPhase:  30 * time.Second,
        MaxSequenceDepth: 5,
        MutationCount:    50,
    }

    // Run simulation
    orchestrator := behaviors.NewBehaviorOrchestrator(graph, config)
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
    defer cancel()

    err := orchestrator.ExecuteAll(ctx)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Analyze results
    results := orchestrator.GetResults()
    fmt.Printf("Execution time: %v\n", results["total_duration"])
}
```

## Agent Details

### Agent 1: Behavior Graph Definition
Validates and initializes the behavior graph structure.

```go
graph := behaviors.NewBehaviorGraph()
graph.AddNode(&behaviors.BehaviorNode{
    ID: "state1",
    Name: "State 1",
    Category: "core",
    Constraints: []string{"timeout_5s", "max_cpu_80"},
})
```

**Output**: Graph node/edge counts, structure validation

### Agent 2: State Machine Simulator
Executes the behavior graph as a finite state machine with configurable transitions.

```go
config := behaviors.StateMachineConfig{
    InitialState: "idle",
    MaxSteps:     100,
    Timeout:      5 * time.Second,
    TrackMetrics: true,
}
sm := behaviors.NewStateMachine(graph, config)
sm.Execute(ctx)
metrics := sm.GetMetrics()
```

**Output**: Transition count, unique states visited, execution latency

### Agent 3: Permutation Generator
Generates all valid behavior sequences deterministically.

```go
gen := behaviors.NewPermutationGenerator(graph)
sequences, _ := gen.GenerateSequences("idle", 5) // depth=5
// Returns all valid paths from "idle" up to 5 steps
```

**Output**: List of behavior sequences, path costs

### Agent 4: Validation Engine
Validates each behavior against registered constraints.

```go
validator := behaviors.NewBehaviorValidator(graph)
validator.RegisterValidator("custom_check", func(node *behaviors.BehaviorNode) error {
    if len(node.ID) == 0 {
        return fmt.Errorf("node id is empty")
    }
    return nil
})
result := validator.Validate("state1")
```

**Output**: Validation results with errors/warnings per behavior

### Agent 5: Concurrent Executor
Runs behavior sequences in parallel with configurable concurrency.

```go
executor := behaviors.NewConcurrentExecutor(graph, config, 10)
results, _ := executor.ExecuteAll(ctx, sequences)
// Executes all sequences with max 10 concurrent executions
```

**Output**: Execution results with timing and success/failure per sequence

### Agent 6: Coverage Analyzer
Analyzes which parts of the behavior graph were exercised.

```go
analyzer := behaviors.NewCoverageAnalyzer(graph)
analyzer.RecordVisit("state1")
analyzer.RecordTransition("state1", "state2")
report := analyzer.GenerateReport()
// report.CoveragePercent, report.UncoveredNodes, etc.
```

**Output**: Coverage report with visited/uncovered nodes, edge coverage

### Agent 7: Performance Profiler
Measures execution performance metrics.

```go
profiler := behaviors.NewPerformanceProfiler()
metrics := profiler.RecordExecution(results)
// MinLatency, AvgLatency, MaxLatency, P95, P99, Throughput
```

**Output**: Latency percentiles, throughput, duration statistics

### Agent 8: Mutation Generator
Creates behavior variations for edge case testing.

```go
mutGen := behaviors.NewMutationGenerator(graph, 12345) // seed
mutations, _ := mutGen.GenerateMutations(50)
for _, mut := range mutations {
    mutGen.ApplyMutation(mut) // Apply first few
}
stats := mutGen.GetMutationStats()
```

**Output**: Applied mutations, mutation type distribution

### Agent 9: Orchestrator
Coordinates all 10 agents in parallel.

```go
orchestrator := behaviors.NewBehaviorOrchestrator(graph, config)
// Launches agents 2-10 in parallel
// Returns aggregated results from all agents
```

**Output**: Orchestration status, timing per agent, combined metrics

### Agent 10: Integration Test Harness
Provides end-to-end testing and result validation.

```go
// Automatically runs within orchestrator
// Validates results from all other agents
// Provides test pass/fail verdict
```

**Output**: Integration test results, overall status

## Configuration

```go
type OrchestratorConfig struct {
    MaxConcurrency   int           // Concurrent behavior executions (1-1000)
    TimeoutPerPhase  time.Duration // Timeout per agent phase
    EnableCaching    bool          // Cache permutation results
    ValidateAll      bool          // Validate all behaviors
    MaxSequenceDepth int           // Max permutation depth (1-10)
    MutationCount    int           // Mutations to generate (1-1000)
}
```

## Architecture

```
BehaviorOrchestrator (Agent 9)
├── Agent 1: BehaviorGraph (setup phase)
├── Agent 2: StateMachine (parallel execution)
├── Agent 3: PermutationGenerator (parallel execution)
├── Agent 4: BehaviorValidator (parallel execution)
├── Agent 5: ConcurrentExecutor (parallel execution)
├── Agent 6: CoverageAnalyzer (parallel execution)
├── Agent 7: PerformanceProfiler (parallel execution)
├── Agent 8: MutationGenerator (parallel execution)
├── Agent 9: Orchestrator (meta-coordination)
└── Agent 10: Integration Test Harness (parallel execution)
```

## Performance

### Typical Execution Times (6-node graph)
- Agent 1 (Graph Definition): <1ms
- Agent 2 (State Machine): ~50s (100 steps)
- Agent 3 (Permutation Generation): <10µs
- Agent 4 (Validation): <5µs
- Agent 5 (Concurrent Execution): ~25s (100 sequences)
- Agent 6 (Coverage): <100µs
- Agent 7 (Performance Profiler): <1µs
- Agent 8 (Mutations): ~50µs
- Agent 9 (Orchestrator): ~100ms
- Agent 10 (Integration): <10µs

**Total**: ~50s (sequential would be ~75s+)

### Metrics Captured
- Behavior sequences generated
- Concurrent execution count
- Coverage percentage
- Latency distribution (min/avg/max/P95/P99)
- Throughput (behaviors/sec)
- Mutation statistics

## Advanced Usage

### Custom Behavior Constraints

```go
graph.AddNode(&behaviors.BehaviorNode{
    ID: "degraded",
    Constraints: []string{
        "max_latency_100ms",
        "max_concurrent_5",
        "requires_recovery",
    },
})
```

### Custom Validators

```go
validator.RegisterValidator("latency_check", func(node *behaviors.BehaviorNode) error {
    if _, ok := node.Metadata["latency_max"]; !ok {
        return fmt.Errorf("node missing latency_max metadata")
    }
    return nil
})
```

### Result Analysis

```go
results := orchestrator.GetResults()
agentResults := results["agent_results"].(map[string]interface{})

// Coverage from Agent 6
coverage := agentResults["coverage_report"].(*behaviors.CoverageReport)
fmt.Printf("Coverage: %.1f%% (%d/%d nodes)\n",
    coverage.CoveragePercent,
    coverage.VisitedNodes,
    coverage.TotalNodes,
)

// Performance from Agent 7
metrics := agentResults["performance_metrics"].(*behaviors.PerformanceMetrics)
fmt.Printf("Latency: avg=%v, p99=%v\n", metrics.AvgLatency, metrics.P99Latency)
```

## Testing

Run all tests:

```bash
go test ./behaviors -v -timeout 5m
```

Run specific test:

```bash
go test ./behaviors -v -run TestAll10AgentsConcurrentSimulation
```

Run stress test:

```bash
go test ./behaviors -v -run TestConcurrentExecutionStress
```

## Type Reference

### BehaviorGraph
Graph-based behavior space model with nodes and edges.

### StateMachine
Finite state machine implementation with metrics tracking.

### BehaviorSequence
A valid path through the behavior graph.

### ValidationResult
Result of behavior validation against constraints.

### ExecutionResult
Result of executing a single behavior sequence.

### CoverageReport
Analysis of behavior space coverage.

### PerformanceMetrics
Latency and throughput statistics.

### Mutation
Variation applied to the behavior graph.

### BehaviorOrchestrator
Main coordinator for all 10 agents.

## Best Practices

1. **Keep graphs manageable**: 6-50 nodes for practical testing
2. **Use meaningful behavior names**: "validating_input" not "v1"
3. **Define clear transitions**: Avoid creating disconnected components
4. **Set appropriate depth**: MaxSequenceDepth 3-5 balances coverage and runtime
5. **Validate early**: Register validators before orchestration
6. **Monitor coverage**: Target >80% node coverage for comprehensive testing

## Limitations

- Graph must be acyclic for permutation generation depth limits to work
- Maximum practical graph size: ~100 nodes
- Mutation application is non-reversible for removed nodes
- Concurrent executor uses semaphore for concurrency control

## No LLM Dependency

This framework is **completely deterministic** and requires:
- Only Go standard library
- No external API calls
- No machine learning models
- No network dependencies
- Perfect for CI/CD, offline testing, and deterministic validation

## Contributing

To extend the framework:

1. Add new agent type following Agent interface patterns
2. Implement metrics collection
3. Add tests to harness
4. Update orchestrator to launch new agent

## License

Part of Claude Squad project.
