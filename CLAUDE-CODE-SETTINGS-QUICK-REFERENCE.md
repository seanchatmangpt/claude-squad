# Claude Code Settings Quick Reference

## File Locations

```bash
# User settings (global)
~/.claude/settings.json

# Project settings (team-shared, commit to git)
.claude/settings.json

# Local project settings (personal, gitignored)
.claude/settings.local.json
```

## Common Configurations

### Basic Development Setup

```json
{
  "model": "claude-sonnet-4-5-20250929",
  "permissions": {
    "allow": [
      "Grep", "Glob",
      "Bash(npm run:*)",
      "Bash(yarn:*)",
      "Read(./src/**)",
      "Read(./tests/**)"
    ],
    "ask": [
      "Bash(git push:*)",
      "Bash(git commit:*)",
      "Edit(./src/**)",
      "Write(./src/**)"
    ],
    "deny": [
      "Read(./.env)",
      "Read(./secrets/**)",
      "Bash(rm -rf:*)"
    ]
  }
}
```

### Secure Production Environment

```json
{
  "model": "claude-opus-4-5-20251101",
  "permissions": {
    "ask": ["*"],
    "deny": [
      "Read(./.env*)",
      "Read(./secrets/**)",
      "Read(./config/credentials.*)",
      "Bash(rm:*)",
      "Bash(sudo:*)",
      "Write(/etc/**)",
      "Write(/usr/**)"
    ]
  },
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": false,
    "allowUnsandboxedCommands": false
  }
}
```

### Auto-Format on Save

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "prettier --write .",
            "timeout": 60
          }
        ]
      }
    ]
  }
}
```

### Auto-Test After Changes

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "npm test -- --bail --findRelatedTests",
            "timeout": 300
          }
        ]
      }
    ]
  }
}
```

### Sandboxed Auto-Allow Mode

```json
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "network": {
      "allowLocalBinding": true
    }
  }
}
```

### Corporate Proxy Setup

```json
{
  "env": {
    "HTTP_PROXY": "http://proxy.company.com:8080",
    "HTTPS_PROXY": "https://proxy.company.com:8443",
    "NO_PROXY": "localhost,127.0.0.1,.company.internal"
  }
}
```

### Git Attribution

```json
{
  "attribution": {
    "commit": "Co-Authored-By: Claude Code <claude-code@anthropic.com>",
    "pr": "\n\n---\n*Generated with Claude Code*"
  }
}
```

### MCP Server Allowlist

```json
{
  "enabledMcpjsonServers": [
    "filesystem",
    "github",
    "postgres"
  ],
  "disabledMcpjsonServers": [
    "experimental-api"
  ]
}
```

## Permission Pattern Examples

```json
{
  "allow": [
    "Grep",                    // Exact tool match
    "Bash(npm run:*)",        // Command prefix match
    "Read(./src/**)",         // Path glob pattern
    "Edit|Write"              // Multiple tools (OR)
  ],
  "ask": [
    "Bash(git push:*)",
    "WebFetch",
    "Write(./config/**)"
  ],
  "deny": [
    "Read(./.env*)",          // Env files
    "Read(./secrets/**)",     // Secrets directory
    "Bash(rm -rf:*)",         // Destructive commands
    "Write(/etc/**)"          // System directories
  ]
}
```

## Hook Event Types

| Event | When | Use For |
|-------|------|---------|
| `SessionStart` | Session begins | Setup, load env |
| `SessionEnd` | Session ends | Cleanup |
| `PreToolUse` | Before tool runs | Validation |
| `PostToolUse` | After tool runs | Format, test |
| `PermissionRequest` | Permission asked | Auto-approve logic |
| `Stop` | Task complete | Generate summary |

## Hook Matcher Examples

```json
{
  "PreToolUse": [
    {
      "matcher": "Edit|Write|MultiEdit",
      "hooks": [{"type": "command", "command": "lint-check.sh"}]
    }
  ],
  "PostToolUse": [
    {
      "matcher": "Bash(git commit:*)",
      "hooks": [{"type": "command", "command": "git-hooks.sh"}]
    }
  ]
}
```

## Environment Variables Quick Ref

```bash
# Authentication
ANTHROPIC_API_KEY=sk-ant-api03-...
ANTHROPIC_AUTH_TOKEN=Bearer ...

# Model Configuration
CLAUDE_CODE_MAX_OUTPUT_TOKENS=8192
MAX_THINKING_TOKENS=10000

# Privacy
DISABLE_TELEMETRY=true
DISABLE_PROMPT_CACHING=true

# Networking
HTTP_PROXY=http://proxy:8080
HTTPS_PROXY=https://proxy:8443
NO_PROXY=localhost,.internal

# Auto-set by Claude Code
CLAUDE_PROJECT_DIR=/path/to/project
CLAUDE_CODE_REMOTE=true  # (web only)
CLAUDE_ENV_FILE=/tmp/env.sh  # (SessionStart only)
```

## CLI Commands

```bash
# View active permissions
/permissions

# Configure hooks interactively
/hooks

# Switch model
/model

# Clear conversation
/clear

# View help
/help
```

## Troubleshooting Commands

```bash
# Validate JSON syntax
jq . ~/.claude/settings.json

# Check file exists
ls -la .claude/settings.json

# View effective settings (run inside Claude Code)
# Settings are applied in order: Enterprise > CLI > Local > Project > User
```

## Security Checklist

- [ ] Deny access to `.env` files
- [ ] Deny access to `secrets/` directory
- [ ] Block `rm -rf` commands
- [ ] Block `sudo` commands
- [ ] Block writes to system directories (`/etc`, `/usr`)
- [ ] Enable sandbox
- [ ] Review MCP server allowlist
- [ ] Set permission precedence correctly
- [ ] Test destructive command blocking

## Model Selection Guide

| Model | ID | Best For |
|-------|-------|----------|
| **Opus 4.5** | `claude-opus-4-5-20251101` | Complex tasks, production code |
| **Sonnet 4.5** | `claude-sonnet-4-5-20250929` | Balanced (default) |
| **Haiku 3.5** | `claude-haiku-3-5-20241022` | Fast iteration, simple tasks |

## Settings Hierarchy (Precedence)

```
1. Enterprise managed settings  ← HIGHEST
2. File-based managed settings
3. Command-line arguments
4. Local project (.claude/settings.local.json)
5. Shared project (.claude/settings.json)
6. User (~/.claude/settings.json)  ← LOWEST
```

**Rule**: Higher precedence ALWAYS wins, even for `deny` overriding `allow`.

## Common Mistakes

❌ **Don't**: Store API keys in committed files
✅ **Do**: Use `.claude/settings.local.json` (gitignored)

❌ **Don't**: Enable `allowUnsandboxedCommands` in production
✅ **Do**: Use `excludedCommands` for specific tools

❌ **Don't**: Use `"allow": ["*"]` globally
✅ **Do**: Explicitly allow needed operations

❌ **Don't**: Forget to test hooks before committing
✅ **Do**: Test in `.local.json` first

## Quick Copy Templates

### Minimal Secure Setup
```json
{
  "permissions": {
    "deny": ["Read(./.env*)", "Read(./secrets/**)", "Bash(rm -rf:*)"]
  }
}
```

### Developer Friendly
```json
{
  "sandbox": {"enabled": true, "autoAllowBashIfSandboxed": true},
  "permissions": {"allow": ["Grep", "Glob", "Bash(npm run:*)"]}
}
```

### Locked Down
```json
{
  "permissions": {"ask": ["*"], "deny": ["Read(./.env*)", "Bash(rm:*)", "Bash(sudo:*)"]},
  "sandbox": {"enabled": true, "allowUnsandboxedCommands": false}
}
```

---

**Full Documentation**: See `CLAUDE-CODE-SETTINGS-DOCUMENTATION.md`

**Template**: See `claude-code-settings-template.json`

**JSON Reference**: See `claude-code-settings-complete-reference.json`
