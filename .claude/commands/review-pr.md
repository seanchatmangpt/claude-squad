---
description: Review a GitHub pull request for quality and correctness
argument-hint: <pr-number>
allowed-tools: [Bash, Read, Grep, Glob]
model: claude-sonnet-4-5-20250929
---

# Review Pull Request #$ARGUMENTS

## Code Review Checklist

1. **Fetch PR Information**
   - Run: `gh pr view $ARGUMENTS --json title,body,commits,files`
   - Understand the PR scope and intent

2. **Check Build Status**
   - Run: `gh pr checks $ARGUMENTS`
   - Verify all CI checks are passing

3. **Review File Changes**
   - List all modified files
   - Check line counts and complexity
   - Flag overly large diffs

4. **Concurrency Safety** (for Go code)
   - Search for mutex usage patterns
   - Check for goroutine leaks
   - Verify channel safety
   - Look for race condition risks

5. **Error Handling**
   - Verify all errors are handled
   - Check error wrapping with context
   - Ensure no silent failures

6. **Testing Coverage**
   - Review test additions
   - Check edge case coverage
   - Verify test names are descriptive

7. **Documentation**
   - Check for Godoc comments
   - Verify API changes are documented
   - Look for helpful code examples

8. **Performance Impact**
   - Identify allocations
   - Check for memory leaks
   - Review critical path performance

9. **Provide Feedback**
   - Run: `gh pr comment $ARGUMENTS -b "Review summary"`
   - Be specific with file:line references
   - Suggest improvements constructively
