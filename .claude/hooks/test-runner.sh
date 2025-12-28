#!/bin/bash
# Test Runner Hook - Run tests after code changes
# Hook Event: PostToolUse (Edit/Write tools on .go files)
# Purpose: Validate that changes don't break tests

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

# Skip tests for test files themselves
if [[ "$FILE_PATH" =~ _test\.go$ ]]; then
    echo "Test file changed, run tests manually with /run-tests" >&2
    exit 0
fi

# Get the package directory
PKG_DIR=$(dirname "$FILE_PATH")
PKG=$(go list "$PKG_DIR" 2>/dev/null)

if [[ -z "$PKG" ]]; then
    echo "Warning: Could not determine package for $FILE_PATH" >&2
    exit 0
fi

# Run tests for the changed package only
echo "Running tests for $PKG..."
if go test "$PKG" -v -timeout 30s 2>&1 | tail -20; then
    echo "✓ Tests passed for $PKG"
    exit 0
else
    echo "⚠️  Some tests failed for $PKG - review and fix" >&2
    exit 0  # Don't block on test failures (developer may want to fix manually)
fi
