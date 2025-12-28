---
description: Analyze test coverage and identify gaps for improvement
allowed-tools: [Bash, Read, Glob, Grep]
model: claude-sonnet-4-5-20250929
---

# Analyze Test Coverage

## Coverage Analysis Workflow

1. **Generate Coverage Profile**
   - Run: `go test ./... -coverprofile=coverage.out`
   - Collect coverage data from all packages
   - Generate detailed coverage profile

2. **Identify Coverage Gaps**
   - Run: `go tool cover -html=coverage.out -o coverage.html`
   - Review visual coverage report
   - Identify uncovered lines and blocks

3. **Package Analysis**
   - Run: `go tool cover -func=coverage.out`
   - Sort by package coverage percentage
   - Identify lowest coverage areas

4. **Critical Path Coverage**
   - Find functions with < 50% coverage
   - Flag error handling paths
   - Highlight public API coverage

5. **Test Gap Identification**
   - Analyze uncovered conditional branches
   - Identify missing edge case tests
   - Note error scenario gaps

6. **Coverage by Module**
   - Main logic: Target 80%+ coverage
   - Utils/helpers: Target 70%+ coverage
   - Integration code: Target 60%+ coverage

7. **Recommendations**
   - Prioritize high-risk untested code
   - Suggest specific test cases
   - Estimate effort to improve coverage

## Output

- Current coverage percentage
- Coverage by package
- List of uncovered critical code
- Prioritized test recommendations
- HTML coverage report
