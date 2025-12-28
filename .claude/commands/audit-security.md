---
description: Audit codebase for security vulnerabilities
allowed-tools: [Bash, Read, Grep, Glob]
model: claude-sonnet-4-5-20250929
---

# Security Audit

## Vulnerability Scanning

1. **Path Traversal Vulnerabilities**
   - Search for: `filepath.Join`, `filepath.Clean`, user input concatenation
   - Verify: All user-controlled paths are validated
   - Check: No `..` sequences allowed

2. **Command Injection**
   - Search for: `exec.Command`, `os/exec` usage
   - Verify: Arguments are never constructed from user input
   - Check: Using arrays, not shell strings

3. **SSRF (Server-Side Request Forgery)**
   - Search for: `http.Get`, `http.Post`, URL construction
   - Verify: URLs are from whitelist or validated
   - Check: No arbitrary scheme/host combinations

4. **SQL Injection**
   - Search for: SQL query construction
   - Verify: Using parameterized queries
   - Check: No string concatenation in queries

5. **Cryptographic Issues**
   - Search for: `crypto/rand`, TLS configuration
   - Verify: Using secure random sources
   - Check: TLS version >= 1.2

6. **Authentication/Authorization**
   - Verify: Password hashing uses bcrypt/scrypt
   - Check: JWT uses strong signing algorithms
   - Review: Permission checks before sensitive operations

7. **Dependency Vulnerabilities**
   - Run: `go list -u -m all | grep -v direct`
   - Run: `go mod tidy && go mod verify`
   - Check for known CVEs in dependencies

8. **Sensitive Data**
   - Search for: Hardcoded credentials, API keys
   - Verify: .env files are git-ignored
   - Check: Secrets not logged or serialized

## Output

Provide a risk assessment with:
- **Critical**: Exploitable vulnerabilities requiring immediate fix
- **High**: Security gaps that should be addressed
- **Medium**: Best practice improvements
- **Low**: Recommendations for future hardening
