# JTBD Testing Framework for Fortune 5 Companies

A comprehensive Jobs-to-be-Done (JTBD) testing framework implemented in pure Go, designed for testing products and services at Fortune 5 companies: Walmart, Amazon, Apple, CVS Health, and UnitedHealth Group.

## Overview

This framework implements the Jobs-to-be-Done theory pioneered by Clayton Christensen, enabling product teams to test features based on what customers are trying to accomplish rather than demographic profiles or feature lists.

### Key Principles

1. **Customer Progress**: Focus on the progress customers want to make, not their characteristics
2. **Circumstance Matters**: The context (when, where, why) determines which solution customers "hire"
3. **Three Dimensions**: Every job has functional, emotional, and social aspects
4. **Measurable Outcomes**: Success is defined by specific, measurable outcomes

## Features

- ✅ **Pure Go Implementation**: No LLM dependencies, fully testable
- ✅ **Concurrent-Safe**: All operations protected with appropriate mutexes
- ✅ **Industry Agnostic**: Works across retail, e-commerce, technology, and healthcare
- ✅ **Comprehensive JTBD Model**: Includes all three job dimensions
- ✅ **Progress Indicators**: Track leading, concurrent, and lagging metrics
- ✅ **Fluent Builder API**: Easy job definition with method chaining
- ✅ **Test Framework**: Built-in test execution and result tracking

## Installation

```bash
go get claude-squad/jtbd
```

## Core Concepts

### Job Structure

Every job has three dimensions:

1. **Functional**: The practical task (e.g., "Get groceries into my house")
2. **Emotional**: How they want to feel (e.g., "Feel confident my family is fed")
3. **Social**: How they want to be perceived (e.g., "Be seen as a reliable provider")

### Circumstances

Circumstances describe when and why a job arises:

- **Temporal**: Time-related constraints
- **Spatial**: Location-related factors
- **Situational**: Broader context or triggers
- **Social**: Social context and audience

### Outcomes

Measurable results that indicate job completion:

- **Speed**: Time to complete
- **Quality**: How well it's done
- **Cost**: Resource efficiency
- **Experience**: Customer satisfaction

### Progress Indicators

Three types of metrics:

- **Leading**: Predict future success (e.g., "task started")
- **Concurrent**: Track during execution (e.g., "50% complete")
- **Lagging**: Measure final outcomes (e.g., "quality score")

## Quick Start

### 1. Define a Job

```go
package main

import (
    "claude-squad/jtbd"
)

func main() {
    // Build a job definition for Walmart monthly shopping
    job, err := jtbd.NewJobBuilder("walmart-monthly-pantry", "Stock pantry for the month").
        WithDescription("Purchase groceries to feed family for one month").
        WithFunctional("Get enough groceries into my house to feed my family for a month").
        WithEmotional("Feel confident that my family won't run out of essential food").
        WithSocial("Be seen as a reliable and organized provider by my family").
        WithIndustry("retail").
        WithCompany("walmart").
        AddCircumstance(&jtbd.Circumstance{
            Type:        jtbd.CircumstanceTypeTemporal,
            Description: "End of month before next paycheck",
            Constraints: map[string]interface{}{
                "days_until_paycheck": 3,
                "current_pantry_level": "20%",
            },
            Intensity: 0.8,
        }).
        AddOutcome(&jtbd.Outcome{
            Type:        jtbd.OutcomeTypeSpeed,
            Description: "Complete shopping in under 90 minutes",
            Metric:      "total_shopping_time_minutes",
            Target:      90.0,
            Unit:        "minutes",
            Priority:    2,
            Direction:   "minimize",
            Threshold:   120.0,
        }).
        AddOutcome(&jtbd.Outcome{
            Type:        jtbd.OutcomeTypeCost,
            Description: "Stay within budget",
            Metric:      "total_cost_dollars",
            Target:      450.00,
            Unit:        "dollars",
            Priority:    1,
            Direction:   "minimize",
            Threshold:   500.00,
        }).
        Build()

    if err != nil {
        panic(err)
    }
}
```

### 2. Register the Job

```go
registry := jtbd.NewJobRegistry()
err := registry.RegisterJob(job)
if err != nil {
    panic(err)
}
```

### 3. Create a Test

```go
test := jtbd.NewSimpleJobTest(
    "outcome_validation",
    "Validates that priority 1 outcomes meet thresholds",
    func(ctx context.Context, job *jtbd.Job) (*jtbd.TestResult, error) {
        // Your test logic here
        // In production, measure actual system performance

        result := &jtbd.TestResult{
            TestName:             "outcome_validation",
            JobID:                job.ID,
            Success:              true,
            Score:                0.95,
            Message:              "All critical outcomes met",
            ProgressMeasurements: make(map[string]float64),
            OutcomeResults:       make(map[string]*jtbd.OutcomeResult),
        }

        return result, nil
    },
)
```

### 4. Execute Tests

```go
executor := jtbd.NewTestExecutor(registry)
executor.RegisterTest(test)

ctx := context.Background()
result, err := executor.ExecuteTest(ctx, "outcome_validation", job.ID)
if err != nil {
    panic(err)
}

fmt.Printf("Test: %s - Success: %v, Score: %.2f\n",
    result.TestName, result.Success, result.Score)
```

## Fortune 5 Examples

The framework includes complete example implementations for each Fortune 5 company:

### Walmart - Monthly Pantry Stocking

```go
job, err := jtbd.ExampleWalmartPantryStocking()
// Complete job definition with:
// - Functional: Get groceries for a month
// - Emotional: Feel confident family is fed
// - Social: Be seen as reliable provider
// - Outcomes: Speed, cost, quality
```

### Amazon - Quick Gift Finding

```go
job, err := jtbd.ExampleAmazonGiftFinding()
// Complete job definition for time-pressured gift shopping
```

### Apple - Family Connection

```go
job, err := jtbd.ExampleAppleFamilyConnection()
// Job definition for staying connected with distant family
```

### CVS Health - Prescription Management

```go
job, err := jtbd.ExampleCVSPrescriptionManagement()
// Managing multiple prescriptions without errors
```

### UnitedHealth - Coverage Understanding

```go
job, err := jtbd.ExampleUnitedHealthCoverageUnderstanding()
// Understanding medical coverage before procedures
```

## Architecture

### Core Types

```
Job
├── ID, Name, Description
├── Functional, Emotional, Social dimensions
├── Circumstances []Circumstance
├── Outcomes []Outcome
└── Indicators []ProgressIndicator

Circumstance
├── Type (Temporal, Spatial, Situational, Social)
├── Description
├── Constraints map[string]interface{}
├── Triggers []string
└── Intensity float64

Outcome
├── Type (Speed, Quality, Cost, Experience)
├── Description, Metric, Unit
├── Target, Threshold values
├── Priority, Direction
└── Metadata

TestResult
├── TestName, JobID
├── Success bool, Score float64
├── ProgressMeasurements map[string]float64
├── OutcomeResults map[string]*OutcomeResult
└── ExecutionTime, Timestamp
```

### Interfaces

```go
// ProgressIndicator measures progress toward job completion
type ProgressIndicator interface {
    Measure(ctx context.Context, job *Job) (float64, error)
    GetName() string
    GetType() IndicatorType
    IsComplete(ctx context.Context, job *Job) (bool, error)
}

// JobTest validates job fulfillment
type JobTest interface {
    Execute(ctx context.Context, job *Job) (*TestResult, error)
    GetTestName() string
    GetDescription() string
    Validate() error
}
```

## Performance

Benchmark results on Intel Xeon @ 2.60GHz:

```
BenchmarkJobRegistry_RegisterJob-16     1000000    2541 ns/op    495 B/op    6 allocs/op
BenchmarkJobRegistry_GetJob-16         10595208     116 ns/op     13 B/op    1 allocs/op
BenchmarkTestExecutor_ExecuteTest-16    3453210     378 ns/op    168 B/op    1 allocs/op
```

- **Job Registration**: ~2.5 µs per job
- **Job Retrieval**: ~117 ns per lookup
- **Test Execution**: ~378 ns overhead (excludes actual test logic)

## Concurrency Safety

All registry operations are protected with `sync.RWMutex`:

- **Read operations** (GetJob, ListJobs): Use `RLock()` for concurrent reads
- **Write operations** (RegisterJob, RemoveJob): Use `Lock()` for exclusive access
- **Test execution**: Results are safely appended to executor

Example concurrent usage:

```go
var wg sync.WaitGroup
registry := jtbd.NewJobRegistry()

// Safe concurrent registration
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        job := &jtbd.Job{
            ID:   fmt.Sprintf("job-%d", id),
            Name: fmt.Sprintf("Job %d", id),
        }
        registry.RegisterJob(job)
    }(i)
}

wg.Wait()
```

## Testing

Run tests:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

Test coverage:

```bash
go test -cover
```

## Design Patterns

### Builder Pattern

Jobs use the builder pattern for clean, readable construction:

```go
job, err := jtbd.NewJobBuilder("id", "name").
    WithFunctional("functional job").
    WithEmotional("emotional job").
    WithSocial("social job").
    AddOutcome(outcome).
    Build()
```

### Registry Pattern

Centralized job management:

```go
registry := jtbd.NewJobRegistry()
registry.RegisterJob(job)
job, err := registry.GetJob("job-id")
jobs := registry.ListJobsByIndustry("retail")
```

### Strategy Pattern

Progress indicators and tests implement interfaces for flexibility:

```go
type CustomProgressIndicator struct { /* ... */ }

func (cpi *CustomProgressIndicator) Measure(ctx context.Context, job *Job) (float64, error) {
    // Custom measurement logic
}
```

## Error Handling

The framework uses structured errors:

```go
type JTBDError struct {
    Code    string
    Message string
    Cause   error
}
```

Common error codes:

- `invalid_job`: Job validation failed
- `invalid_test`: Test validation failed
- `job_not_found`: Job ID not in registry
- `test_not_found`: Test name not registered
- `test_failed`: Test execution error
- `invalid_input`: Invalid parameter
- `internal_error`: Internal framework error

## Best Practices

### 1. Define Jobs from Customer Perspective

❌ **Bad**: "Enable users to click the checkout button"
✅ **Good**: "Get my purchase completed quickly before I change my mind"

### 2. Include All Three Dimensions

```go
// Don't just define functional
job.Functional = "Purchase item"

// Include emotional and social too
job.Emotional = "Feel confident I made the right choice"
job.Social = "Be seen as a smart shopper by my family"
```

### 3. Use Specific, Measurable Outcomes

❌ **Bad**: "Make checkout better"
✅ **Good**: "Complete checkout in under 60 seconds with zero errors"

```go
&jtbd.Outcome{
    Type:        jtbd.OutcomeTypeSpeed,
    Description: "Complete checkout quickly",
    Metric:      "checkout_time_seconds",
    Target:      60.0,
    Unit:        "seconds",
    Threshold:   90.0,
    Direction:   "minimize",
}
```

### 4. Model Real Circumstances

Include the actual constraints customers face:

```go
&jtbd.Circumstance{
    Type:        jtbd.CircumstanceTypeTemporal,
    Description: "Shopping during lunch break",
    Constraints: map[string]interface{}{
        "time_available_minutes": 30,
        "must_return_to_work":    true,
    },
    Intensity: 0.9,
}
```

## Contributing

When adding new functionality:

1. Follow existing code patterns
2. Add comprehensive godoc comments
3. Include unit tests
4. Ensure concurrent safety
5. Add benchmark tests for performance-critical code
6. Update this README

## License

See LICENSE.md in the project root.

## References

- Christensen, Clayton M. "The Innovator's Dilemma" (1997)
- Christensen, Clayton M. "Competing Against Luck" (2016)
- Ulwick, Anthony W. "Jobs to be Done: Theory to Practice" (2016)

## Support

For questions or issues, please file an issue in the project repository.

---

**Framework Version**: 1.0.0
**Go Version**: 1.16+
**Status**: Production Ready ✅
