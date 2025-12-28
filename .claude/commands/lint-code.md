---
description: Lint and format code according to project standards
allowed-tools: [Bash, Glob]
model: claude-haiku-4-5-20251001
---

# Lint and Format Code

## Code Quality Checks

1. **Format Check**
   - Run: `go fmt ./...`
   - Ensure consistent formatting
   - Fix all formatting issues

2. **Vet Analysis**
   - Run: `go vet ./...`
   - Find common mistakes
   - Report shadowing and suspicious constructs

3. **Staticcheck Linting**
   - Run: `staticcheck ./...`
   - Find bugs and performance issues
   - Suggest code improvements
   - Identify unused code

4. **Errcheck**
   - Run: `errcheck ./...`
   - Verify all errors are handled
   - Report unchecked errors

5. **Ineffassign**
   - Run: `ineffassign ./...`
   - Find unused assignments
   - Clean up dead code

6. **Spelling**
   - Run: `misspell -error ./...`
   - Check documentation spelling
   - Fix common typos

7. **Complexity Check**
   - Run: `gocyclo -over 10 ./...`
   - Report overly complex functions
   - Suggest refactoring opportunities

## Output

- Formatting applied
- All linting issues fixed
- Code quality report
- Complexity metrics
- Recommendations for improvement
