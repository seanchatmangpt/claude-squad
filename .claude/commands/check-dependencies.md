---
description: Check and update project dependencies for security and compatibility
allowed-tools: [Bash]
model: claude-haiku-4-5-20251001
---

# Check Dependencies

## Dependency Analysis

1. **List Outdated Dependencies**
   - Run: `go list -u -m all`
   - Identify available updates
   - Categorize by stability (major/minor/patch)

2. **Security Scanning**
   - Run: `go list -m all | sort -u | while read m; do go get -d $m; done`
   - Check for known CVEs
   - Review security advisories

3. **Verify Module Integrity**
   - Run: `go mod verify`
   - Ensure no tampering of dependencies
   - Validate checksums

4. **Find Unused Dependencies**
   - Run: `go mod tidy`
   - Remove indirect unused modules
   - Optimize go.mod

5. **Update Strategy**
   - Patch updates: Safe, apply immediately
   - Minor updates: Generally safe, test before
   - Major updates: Potential breaking changes, careful review

6. **Compatibility Check**
   - Test with minimum Go version
   - Verify platform compatibility
   - Check for deprecated APIs

## Output

- List of available updates
- Security vulnerability report
- Recommendations for updating
- Compatibility analysis
- Breaking change warnings
