#!/bin/bash
# Code Formatter Hook - Auto-format code after edits
# Hook Event: PostToolUse (Edit/Write tools)
# Purpose: Automatically format code according to project standards

INPUT=$(cat)

# Extract the tool name and file path
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName // empty')
FILE_PATH=$(echo "$INPUT" | jq -r '.filePath // empty')

# Only run for Go files after Edit/Write
if [[ ! "$FILE_PATH" =~ \.go$ ]]; then
    exit 0
fi

if [[ "$TOOL_NAME" != "Edit" && "$TOOL_NAME" != "Write" ]]; then
    exit 0
fi

# Get the directory containing the file
DIR=$(dirname "$FILE_PATH")

# Run Go formatter
echo "Formatting $FILE_PATH..."
if ! go fmt "$FILE_PATH" 2>/dev/null; then
    echo "Warning: go fmt encountered an error" >&2
    exit 0  # Don't block on formatting errors
fi

# Run Go vet on the package
echo "Running go vet on package..."
if ! go vet ./... 2>/dev/null; then
    echo "Warning: go vet found issues" >&2
fi

echo "✓ Code formatting complete"
exit 0
