---
description: Analyze and fix a GitHub issue with full validation
argument-hint: <issue-number>
allowed-tools: [Bash, Read, Edit, Grep, Glob]
model: claude-sonnet-4-5-20250929
---

# Fix GitHub Issue #$ARGUMENTS

## Workflow

1. **Fetch Issue Details**
   - Run: `gh issue view $ARGUMENTS --json title,body,labels`
   - Parse the issue requirements and acceptance criteria

2. **Understand the Problem**
   - Analyze the issue description
   - Identify root cause and affected components
   - Note any related issues or PRs

3. **Search Codebase**
   - Use Grep to find related code patterns
   - Use Glob to locate relevant files
   - Understand current implementation

4. **Implement Fix**
   - Make necessary code changes
   - Follow CLAUDE.md coding standards
   - Update related tests as needed

5. **Validate**
   - Run: `go build ./...`
   - Run: `go test ./...`
   - Ensure code follows formatting standards

6. **Create Commit**
   - Write descriptive commit message
   - Reference the issue number: "Fixes #$ARGUMENTS"
   - Include testing summary

7. **Create Pull Request** (if needed)
   - Run: `gh pr create --title "fix: ..." --body "Fixes #$ARGUMENTS"`
   - Link to the original issue
