---
name: security-audit
description: Audit code for security vulnerabilities and best practices
allowed-tools: [Bash, Read, Grep, Glob]
---

# Security Audit Skill

Specialized skill for identifying and remediating security vulnerabilities.

## Vulnerability Categories

### Input Validation
- Path traversal prevention
- Command injection prevention
- SQL injection prevention
- Cross-site scripting (XSS) prevention

### Cryptography
- Enforce minimum TLS 1.2
- Use secure random generation
- Proper key management
- Secure hashing algorithms

### Authentication & Authorization
- Password hashing (bcrypt/scrypt)
- JWT signature validation
- Session management
- Permission enforcement

### Network Security
- SSRF prevention
- Rate limiting
- DDoS protection
- Secure communication

### Data Protection
- No hardcoded secrets
- Secure file permissions
- Encrypted sensitive data
- Safe deserialization

### Dependency Security
- Track vulnerable dependencies
- Use dependency scanning
- Keep dependencies updated
- Monitor security advisories

## Audit Process

### Code Review
- Search for dangerous functions
- Identify user input sources
- Trace to usage points
- Verify sanitization

### Configuration Review
- Check TLS settings
- Verify authentication config
- Review permission models
- Check logging settings

### Dependency Audit
- Run vulnerability scanner
- Check for known CVEs
- Evaluate new versions
- Plan security updates

## Severity Levels

- **Critical**: Exploitable vulnerability, immediate fix required
- **High**: Security gap, should be fixed soon
- **Medium**: Best practice improvement
- **Low**: Hardening recommendation

## Usage

Analyzes code and configuration for security vulnerabilities. Provides specific remediation guidance. Tracks vulnerability status and remediation timeline.

## Output

Prioritized list of vulnerabilities with detailed remediation steps and risk assessment.
