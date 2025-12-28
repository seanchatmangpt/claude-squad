---
name: code-review
description: Review Go code for quality, safety, and best practices
allowed-tools: [Read, Grep, Glob]
---

# Code Review Skill

Specialized skill for reviewing Go code against project standards and best practices.

## Review Criteria

### Concurrency Safety
- Check for proper mutex usage with `sync.Mutex` and `sync.RWMutex`
- Verify channel operations use timeouts and proper closing
- Look for goroutine leaks and improper WaitGroup usage
- Identify race conditions from shared state access

### Error Handling
- All errors should be handled or wrapped with context
- Check for unchecked type assertions (panic risk)
- Verify error messages provide actionable information
- Look for silent failures and ignored errors

### Code Quality
- Follow CLAUDE.md naming conventions
- Verify error wrapping uses `fmt.Errorf(...%w...)`
- Check function/variable naming is clear and consistent
- Flag unnecessarily complex code (cyclomatic complexity > 10)

### Performance
- Identify memory allocations in hot paths
- Check for unbounded collections
- Look for unnecessary copying of large structs
- Verify lock contention is minimized

### Testing
- Public APIs should have corresponding tests
- Error cases should be tested
- Edge cases should have explicit tests
- Flaky tests should be identified

## Usage

Activated automatically when reviewing code changes or examining Go files. Can be manually invoked for detailed analysis of specific packages or modules.

## Output

Provides file:line specific findings with severity levels and recommendations.
