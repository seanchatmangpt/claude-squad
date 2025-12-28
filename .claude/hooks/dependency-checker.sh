#!/bin/bash
# Dependency Checker Hook - Check go.mod changes
# Hook Event: PostToolUse (Edit on go.mod)
# Purpose: Validate dependencies and alert on issues

INPUT=$(cat)

# Extract the file path
FILE_PATH=$(echo "$INPUT" | jq -r '.filePath // empty')

# Only check go.mod changes
if [[ ! "$FILE_PATH" =~ go\.mod$ ]]; then
    exit 0
fi

echo "Checking Go dependencies..."

# Verify module integrity
if ! go mod verify >/dev/null 2>&1; then
    echo "⚠️  Warning: go mod verify failed" >&2
    echo "Run: go mod tidy" >&2
    exit 0
fi

# Run tidy to clean up
echo "Running go mod tidy..."
if ! go mod tidy >/dev/null 2>&1; then
    echo "⚠️  Warning: go mod tidy failed" >&2
    exit 0
fi

# Check for outdated dependencies
echo "Checking for outdated dependencies..."
OUTDATED=$(go list -u -m all 2>/dev/null | grep -c "\[")

if [[ $OUTDATED -gt 0 ]]; then
    echo "⚠️  Found $OUTDATED outdated dependencies. Run: go list -u -m all" >&2
fi

# Verify build still works
echo "Verifying build..."
if ! go build ./... >/dev/null 2>&1; then
    echo "⚠️  Build failed after go.mod changes" >&2
    exit 0
fi

echo "✓ Dependencies verified"
exit 0
