#!/bin/bash
# Git Safety Guard Hook - Prevent dangerous git operations
# Hook Event: PreToolUse (Bash tool)
# Purpose: Prevent accidental force pushes and destructive operations

INPUT=$(cat)

# Extract command being executed
COMMAND=$(echo "$INPUT" | jq -r '.command // empty')

# Dangerous patterns to block
DANGEROUS_PATTERNS=(
    "git push.*--force"
    "git push.*-f"
    "git reset.*--hard"
    "git clean.*-fd"
    "git rebase.*-i"
    "git branch.*-D"
    "git checkout.*-"
    "rm -rf .git"
    "git update-ref"
    "git reflog delete"
)

# Check if command matches any dangerous pattern
for pattern in "${DANGEROUS_PATTERNS[@]}"; do
    if [[ "$COMMAND" =~ $pattern ]]; then
        echo "⚠️  Git Safety Guard: Dangerous operation detected" >&2
        echo "Command: $COMMAND" >&2
        echo "This operation is protected. Use with explicit approval only." >&2

        # For force push to main/master, always block
        if [[ "$COMMAND" =~ "git push.*--force" && "$COMMAND" =~ "(main|master)" ]]; then
            echo "BLOCKED: Force push to main/master is not allowed" >&2
            exit 2  # Block the operation
        fi

        # For other dangerous operations, warn but allow with caution
        echo "⚠️  Proceeding with caution..." >&2
    fi
done

exit 0
