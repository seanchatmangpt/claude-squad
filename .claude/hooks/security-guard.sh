#!/bin/bash
# Security Guard Hook - Prevent editing sensitive files
# Hook Event: PreToolUse
# Purpose: Protect sensitive files from accidental modification

# List of protected file patterns
PROTECTED_PATTERNS=(
    ".env"
    ".env.*"
    "secrets/"
    "credentials.json"
    "config/private"
    ".ssh/"
    ".git/config"
    "package-lock.json"
    "go.sum"
)

# Read the incoming hook data
INPUT=$(cat)

# Extract the tool name and file path (format depends on hook context)
# This is a simplified example - actual parsing depends on your setup
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName // empty')
FILE_PATH=$(echo "$INPUT" | jq -r '.filePath // empty')

# Only check for Read/Edit/Write operations
if [[ "$TOOL_NAME" != "Read" && "$TOOL_NAME" != "Edit" && "$TOOL_NAME" != "Write" ]]; then
    exit 0
fi

# Check against protected patterns
for pattern in "${PROTECTED_PATTERNS[@]}"; do
    if [[ "$FILE_PATH" == *"$pattern"* ]]; then
        echo "🔒 Security Guard: File protection active - '$FILE_PATH' is protected" >&2
        echo "This file is protected from modification. Contact admin if needed." >&2
        exit 2  # Block the operation
    fi
done

# Allow the operation
exit 0
