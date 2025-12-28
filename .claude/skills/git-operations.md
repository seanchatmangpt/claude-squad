---
name: git-operations
description: Manage git workflows including commits, branches, and pull requests
allowed-tools: [Bash, Read, Glob]
---

# Git Operations Skill

Specialized skill for managing git workflows and version control operations.

## Git Workflows

### Branch Management
- Create feature branches following `claude/<name>` convention
- Keep branches focused on single features
- Pull latest main before starting work
- Delete branches after merging to main

### Commit Standards
- Write descriptive commit messages
- Reference issue numbers: "Fixes #123"
- Format: `<type>: <description>` (feat, fix, refactor, docs)
- One logical change per commit

### Pull Request Process
- Create PR after feature is complete
- Link to related issues
- Include test results
- Request review from team
- Address review feedback

### Version Control
- Tag releases with semantic versioning (v1.0.0)
- Keep main branch stable
- Use worktrees for parallel work
- Document breaking changes

## Safe Operations

### Protected Operations
- Main branch requires PR review
- No direct pushes to main
- All tests must pass before merge
- Requires at least one approval

### Cleanup
- Remove merged branches
- Clean unused worktrees
- Archive old feature branches
- Keep repository clean

## Usage

Provides safe git command execution with validation. Used for branch creation, commit management, and PR workflows. Enforces team conventions and prevents accidental main branch modifications.

## Output

Confirms successful git operations, reports branch status, and verifies commits follow standards.
