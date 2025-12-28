---
description: Run comprehensive test suite with coverage analysis
allowed-tools: [Bash]
model: claude-haiku-4-5-20251001
---

# Run Test Suite

## Test Execution

1. **Unit Tests**
   - Run: `go test ./... -v -race`
   - Captures all test output with verbose logging
   - Race condition detector enabled

2. **Coverage Report**
   - Run: `go test ./... -coverprofile=coverage.out`
   - Run: `go tool cover -html=coverage.out -o coverage.html`
   - Generates visual coverage report

3. **Benchmark Tests**
   - Run: `go test ./... -bench=. -benchmem`
   - Compare memory allocations
   - Identify performance regressions

4. **Integration Tests**
   - Run: `go test -tags=integration ./...`
   - Tests requiring external services
   - Validates end-to-end workflows

5. **Failure Analysis**
   - Parse test failures
   - Identify flaky tests
   - Report specific failed test cases

## Output

Summary of:
- Tests passed/failed
- Coverage percentage
- Critical test failures
- Performance baselines
- Recommendations for failing tests
