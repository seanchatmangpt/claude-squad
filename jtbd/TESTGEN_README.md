# JTBD Test Case Generator

## Overview

A deterministic, template-based test case generator for Fortune 5 Jobs-to-be-Done (JTBD) testing. This implementation requires **no LLM dependencies** and generates comprehensive test scenarios across five major industries.

## Features

- **Deterministic Generation**: Template-based approach with no randomness
- **Fortune 5 Industry Coverage**: 
  - Retail (Walmart)
  - E-commerce (Amazon)
  - Technology (Apple)
  - Healthcare/Pharmacy (CVS Health)
  - Insurance (UnitedHealth Group)

- **Test Case Types**:
  - Happy path scenarios
  - Edge cases (time constraints, budget limits, urgent needs)
  - Failure scenarios
  - Multi-step workflows
  - Competing jobs with trade-offs

- **JTBD Framework Integration**: Converts test cases to framework Job objects with all three dimensions:
  - Functional (practical task)
  - Emotional (desired feelings)
  - Social (how to be perceived)

## Core Structures

### TestCase
```go
type TestCase struct {
    ID               string
    Industry         string
    JobSpec          TestJobSpec
    CircumstanceSpec TestCircumstanceSpec
    OutcomeSpec      TestOutcomeSpec
    Constraints      []Constraint
    CompetingJobs    []TestJobSpec
    TradeOffs        []string
    IsEdgeCase       bool
    IsHappyPath      bool
    MultiStep        bool
}
```

### TestJobSpec
```go
type TestJobSpec struct {
    Name        string
    Description string
    Category    string
    Steps       []string
    Priority    string
    Functional  string  // JTBD functional dimension
    Emotional   string  // JTBD emotional dimension
    Social      string  // JTBD social dimension
}
```

## Usage Examples

### Basic Generation

```go
gen := jtbd.NewTestCaseGenerator()

options := jtbd.TestGenerationOptions{
    IncludeHappyPath: true,
    IncludeEdgeCases: true,
}

// Generate test cases for Walmart
walmartCases := gen.GenerateTestCases("retail", options)

for _, tc := range walmartCases {
    fmt.Printf("Test: %s - %s\n", tc.ID, tc.JobSpec.Name)
    fmt.Printf("  Functional: %s\n", tc.JobSpec.Functional)
    fmt.Printf("  Emotional: %s\n", tc.JobSpec.Emotional)
    fmt.Printf("  Social: %s\n", tc.JobSpec.Social)
}
```

### Generate All Industries

```go
allOptions := jtbd.TestGenerationOptions{
    IncludeHappyPath:   true,
    IncludeEdgeCases:   true,
    IncludeFailures:    true,
    IncludeMultiStep:   true,
    IncludeCompeting:   true,
    CombinatorialLevel: 2,  // Increases variation
}

allCases := gen.GenerateAllTestCases(allOptions)

for industry, cases := range allCases {
    fmt.Printf("%s: %d test cases\n", industry, len(cases))
}
```

### Convert to Framework Job

```go
testCase := walmartCases[0]
job := testCase.ToJob()

// Now you have a framework Job with:
// - job.Functional, job.Emotional, job.Social (JTBD dimensions)
// - job.Circumstances (converted from CircumstanceSpec)
// - job.Outcomes (converted from OutcomeSpec)
// - job.Metadata (includes test case info)
```

## Industry-Specific Examples

### Walmart - Weekly Grocery Shopping
- **Functional**: Get enough groceries to feed family for a week
- **Emotional**: Feel confident family has enough food
- **Social**: Be seen as responsible provider
- **Edge Case**: 30-minute time limit, $150 budget, family of 4

### Amazon - Gift Shopping
- **Functional**: Find perfect gift within budget and timeframe
- **Emotional**: Feel confident recipient will love it
- **Social**: Be seen as thoughtful and generous
- **Edge Case**: Next-day delivery required, $50 budget, birthday tomorrow

### Apple - Device Setup for Elderly
- **Functional**: Get device working with all data transferred
- **Emotional**: Feel confident they can use it without frustration
- **Social**: Be seen as helpful and patient by family
- **Constraints**: Must preserve all photos, minimize tech confusion

### CVS Health - Prescription Refill
- **Functional**: Get medication refilled without running out
- **Emotional**: Feel secure about health management
- **Social**: Maintain independence in managing health
- **Edge Case**: Urgent refill, running out of medication

### UnitedHealth - Find Provider
- **Functional**: Find qualified provider who accepts insurance
- **Emotional**: Feel confident about avoiding unexpected costs
- **Social**: Make responsible healthcare decisions for family
- **Constraints**: In-network only, deductible not met

## Test Generation Options

```go
type TestGenerationOptions struct {
    IncludeHappyPath    bool  // Standard success scenarios
    IncludeEdgeCases    bool  // Time pressure, budget constraints
    IncludeFailures     bool  // Scenarios that fail
    IncludeMultiStep    bool  // Complex workflows
    IncludeCompeting    bool  // Jobs with competing priorities
    CombinatorialLevel  int   // 0=none, 1=low, 2=medium, 3=high
    MaxCasesPerCategory int   // Limit per category
}
```

## Combinatorial Explosion

The generator supports combinatorial test case generation:

```go
options := jtbd.TestGenerationOptions{
    IncludeHappyPath:   true,
    CombinatorialLevel: 2,  // Creates 2 variations per base case
}
```

- **Level 0**: No variations (base cases only)
- **Level 1**: 2 variations per case
- **Level 2**: 4 variations per case
- **Level 3**: 6 variations per case

## Files

- `/home/user/claude-squad/jtbd/testgen.go` - Main implementation (582 lines)
- `/home/user/claude-squad/jtbd/testgen_test.go` - Comprehensive tests
- `/home/user/claude-squad/jtbd/testgen_example.go` - Usage examples

## Running Tests

```bash
# All tests
go test ./jtbd/ -v

# Generator-specific tests
go test ./jtbd/ -run TestGenerate -v

# Run example
go run /home/user/claude-squad/jtbd/testgen_example.go
```

## Implementation Highlights

1. **No LLM Dependencies**: Pure Go, template-based generation
2. **Deterministic**: Same inputs always produce same test cases
3. **Industry-Specific**: Real Fortune 5 job patterns
4. **Framework Integration**: Seamless conversion to JTBD framework Jobs
5. **Extensible**: Easy to add new industries or job patterns
6. **Type-Safe**: Leverages Go's type system for correctness
7. **Well-Tested**: Comprehensive test coverage

## Test Results

```
=== RUN   TestNewTestCaseGenerator
--- PASS: TestNewTestCaseGenerator (0.00s)
=== RUN   TestGenerateRetailTestCases
    Generated 3 retail test cases
--- PASS: TestGenerateRetailTestCases (0.00s)
=== RUN   TestGenerateAllIndustries
    Total test cases generated: 10
--- PASS: TestGenerateAllIndustries (0.00s)
=== RUN   TestTestCaseToJob
    Successfully converted TestCase to Job
--- PASS: TestTestCaseToJob (0.00s)
=== RUN   TestCombinatorialExplosion
    Combinatorial level 0: 1 cases, level 2: 5 cases
--- PASS: TestCombinatorialExplosion (0.00s)
PASS
ok      claude-squad/jtbd       0.010s
```

## Architecture

The generator uses a template-based approach:

1. **Industry Patterns**: Pre-defined job templates per industry
2. **Job Templates**: Structured definitions with JTBD dimensions
3. **Outcome Templates**: Expected results per job type
4. **Test Case Factory**: Combines templates into complete test cases
5. **Variation Engine**: Generates edge cases and combinations

All generation is deterministic and requires no external dependencies or LLMs.
