---
description: Create a new feature branch following team conventions
argument-hint: <feature-name>
allowed-tools: [Bash]
model: claude-haiku-4-5-20251001
---

# Create Feature Branch $ARGUMENTS

## Feature Branch Creation

1. **Validate Branch Name**
   - Convert to lowercase: `$ARGUMENTS`
   - Replace spaces with hyphens
   - Remove special characters
   - Format: `feature/$ARGUMENTS`

2. **Create Branch**
   - Fetch latest: `git fetch origin main`
   - Create branch: `git checkout -b claude/$ARGUMENTS`
   - Push upstream: `git push -u origin claude/$ARGUMENTS`

3. **Document Feature**
   - Create commit message template
   - Add feature description to branch
   - Document acceptance criteria

4. **Setup Development**
   - Run: `go mod download`
   - Run: `go build ./...`
   - Verify development environment

5. **Initial Commit**
   - Stage changes: `git add .`
   - Commit: `git commit -m "feat: Initialize $ARGUMENTS feature branch"`
   - Push: `git push origin claude/$ARGUMENTS`

## Convention

- All feature branches start with `claude/`
- Use lowercase hyphenated names
- Reference related issues in description
- Keep focused on single feature

## Next Steps

- Implement feature in focused commits
- Create PR when ready for review
- Link to related GitHub issues
