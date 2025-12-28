# Claude Code Settings Complete Reference (2025)

## Overview

This comprehensive reference documents all available settings for Claude Code configuration as of December 2025. Settings are configured through JSON files at different hierarchy levels, with enterprise settings taking precedence over user settings.

---

## Table of Contents

- [Configuration Files](#configuration-files)
- [Settings Hierarchy](#settings-hierarchy)
- [Core Settings](#core-settings)
- [Permissions](#permissions)
- [Sandbox Configuration](#sandbox-configuration)
- [Hooks](#hooks)
- [MCP Servers](#mcp-servers)
- [Plugins](#plugins)
- [Environment Variables](#environment-variables)
- [Best Practices](#best-practices)

---

## Configuration Files

### File Locations

| Scope | File Path | Description | Committed to Git? |
|-------|-----------|-------------|-------------------|
| **User** | `~/.claude/settings.json` | Global user settings for all projects | N/A |
| **User (Legacy)** | `~/.claude.json` | Legacy location, still supported | N/A |
| **Project** | `.claude/settings.json` | Team-shared project configuration | ✅ Yes |
| **Local Project** | `.claude/settings.local.json` | Personal project overrides | ❌ No (gitignored) |
| **Enterprise** | System-specific paths | Enterprise managed settings | N/A |

### Enterprise Locations

- **macOS**: `/Library/Application Support/ClaudeCode/`
- **Linux/WSL**: `/etc/claude-code/`
- **Windows**: `C:\Program Files\ClaudeCode\`

---

## Settings Hierarchy

Settings are applied in the following order (highest to lowest precedence):

1. **Enterprise managed settings** (highest)
2. File-based managed settings
3. Command-line arguments
4. Local project settings (`.claude/settings.local.json`)
5. Shared project settings (`.claude/settings.json`)
6. User settings (`~/.claude/settings.json`) (lowest)

---

## Core Settings

### Model Configuration

```json
{
  "model": "claude-sonnet-4-5-20250929",
  "outputStyle": "concise"
}
```

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `model` | string | `"claude-sonnet-4-5-20250929"` | Override default model (sonnet, opus, haiku) |
| `outputStyle` | string | `null` | Adjust system prompt for response style |

**Available Models (2025):**
- `claude-opus-4-5-20251101` - Most capable
- `claude-sonnet-4-5-20250929` - Balanced (default)
- `claude-haiku-3-5-20241022` - Fast and efficient

### Authentication

```json
{
  "apiKeyHelper": "/path/to/auth-script.sh",
  "forceLoginMethod": "claudeai",
  "forceLoginOrgUUID": "org_1234567890abcdef"
}
```

| Setting | Type | Description |
|---------|------|-------------|
| `apiKeyHelper` | string | Custom shell script for generating auth values |
| `forceLoginMethod` | string | Restrict to `"claudeai"` or `"console"` |
| `forceLoginOrgUUID` | string | Auto-select organization UUID |

### Session Management

```json
{
  "cleanupPeriodDays": 60,
  "companyAnnouncements": [
    "Security policy: All code requires review",
    "Reminder: Use approved MCP servers only"
  ]
}
```

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `cleanupPeriodDays` | number | `30` | Delete sessions inactive for N days |
| `companyAnnouncements` | string[] | `[]` | Startup messages (enterprise) |

### Environment Variables

```json
{
  "env": {
    "NODE_ENV": "production",
    "API_BASE_URL": "https://api.company.com",
    "LOG_LEVEL": "info",
    "MAX_RETRIES": "3"
  }
}
```

All variables in `env` are available to tools and hooks during execution.

### Attribution

```json
{
  "attribution": {
    "commit": "Co-Authored-By: Claude Code <claude-code@anthropic.com>",
    "pr": "\n\n---\n*Generated with assistance from Claude Code*"
  }
}
```

| Setting | Type | Description |
|---------|------|-------------|
| `attribution.commit` | string | Git commit trailer (replaces deprecated `includeCoAuthoredBy`) |
| `attribution.pr` | string | Pull request description footer |

---

## Permissions

### Permission Structure

```json
{
  "permissions": {
    "allow": ["Bash(npm run:*)", "Read(~/.config/**)"],
    "ask": ["Bash(git push:*)", "WebFetch"],
    "deny": ["Read(./.env)", "Write(/etc/**)"],
    "additionalDirectories": ["../docs/"],
    "defaultMode": "acceptEdits",
    "disableBypassPermissionsMode": false
  }
}
```

### Permission Rules

| Rule Type | Description | Behavior |
|-----------|-------------|----------|
| `allow` | Auto-approved operations | Runs without prompts |
| `ask` | Requires approval | Shows permission dialog |
| `deny` | Blocked operations | Rejected, user notified |

### Pattern Matching

Permissions use **prefix matching** with wildcards:

```json
{
  "allow": [
    "Bash(npm run:*)",        // Matches: npm run test, npm run build, etc.
    "Read(./src/**)",         // Matches all files in src/ recursively
    "Grep",                   // Matches Grep tool exactly
    "Edit|Write"              // Matches Edit OR Write tools
  ]
}
```

### Common Permission Patterns

**Development Workflow:**
```json
{
  "allow": [
    "Bash(npm run:*)",
    "Bash(yarn:*)",
    "Bash(pnpm:*)",
    "Read(./src/**)",
    "Read(./tests/**)",
    "Grep",
    "Glob"
  ],
  "ask": [
    "Bash(git commit:*)",
    "Bash(git push:*)",
    "Edit(./src/**)",
    "Write(./src/**)"
  ],
  "deny": [
    "Read(./.env)",
    "Read(./.env.*)",
    "Read(./secrets/**)",
    "Bash(rm -rf:*)"
  ]
}
```

**Security-Focused:**
```json
{
  "deny": [
    "Read(./.env)",
    "Read(./.env.*)",
    "Read(./secrets/**)",
    "Read(./config/credentials.*)",
    "Write(/etc/**)",
    "Write(/usr/**)",
    "Write(/bin/**)",
    "Bash(rm -rf:*)",
    "Bash(chmod:*)",
    "Bash(sudo:*)"
  ]
}
```

### Additional Directories

```json
{
  "additionalDirectories": [
    "../shared-utils/",
    "~/company-config/",
    "/opt/project-templates/"
  ]
}
```

Expands access beyond the project root. Use with caution.

### Permission Precedence

**Important**: Project-level `deny` overrides user-level `allow`:

```json
// User settings (~/.claude/settings.json)
{"permissions": {"allow": ["Write(./src/**)"]}}

// Project settings (.claude/settings.json)
{"permissions": {"deny": ["Write(./src/critical/**)"]}}}

// Result: Write to ./src/critical/** is DENIED (project wins)
```

---

## Sandbox Configuration

### Overview

Sandboxing provides isolated execution for Bash commands with restricted file system and network access.

```json
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["docker", "systemctl"],
    "allowUnsandboxedCommands": false,
    "network": {
      "httpProxyPort": 8080,
      "socksProxyPort": 8081,
      "allowUnixSockets": ["~/.ssh/agent-socket"],
      "allowLocalBinding": true
    },
    "enableWeakerNestedSandbox": false
  }
}
```

### Sandbox Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `enabled` | boolean | `true` | Enable sandbox isolation |
| `autoAllowBashIfSandboxed` | boolean | `false` | Auto-approve sandboxed bash commands |
| `excludedCommands` | string[] | `[]` | Commands that run outside sandbox |
| `allowUnsandboxedCommands` | boolean | `false` | Allow escape hatch (security risk) |
| `enableWeakerNestedSandbox` | boolean | `false` | Docker compatibility mode |

### Network Configuration

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `httpProxyPort` | number | `8080` | HTTP proxy port for outbound traffic |
| `socksProxyPort` | number | `8081` | SOCKS proxy port |
| `allowUnixSockets` | string[] | `[]` | Unix socket paths (e.g., SSH agent) |
| `allowLocalBinding` | boolean | `false` | Allow binding to local interfaces |

### File System Isolation

By default, sandbox provides:
- **Read/Write**: Current working directory and subdirectories
- **Controlled via**: `permissions.deny` with `Read()` patterns

```json
{
  "permissions": {
    "deny": [
      "Read(/etc/**)",
      "Read(/usr/**)",
      "Read(~/.ssh/**)"
    ]
  }
}
```

### Docker Compatibility

```json
{
  "sandbox": {
    "excludedCommands": ["docker", "docker-compose"],
    "enableWeakerNestedSandbox": true
  }
}
```

⚠️ **Warning**: `enableWeakerNestedSandbox` reduces security guarantees.

---

## Hooks

### Overview

Hooks execute custom commands or LLM prompts at specific lifecycle events.

### Hook Events

| Event | Trigger Point | Use Cases |
|-------|---------------|-----------|
| `SessionStart` | Session initialization | Load env, setup workspace |
| `SessionEnd` | Session termination | Cleanup, save state |
| `PreToolUse` | Before tool execution | Validation, pre-formatting |
| `PostToolUse` | After tool completion | Linting, testing, formatting |
| `PermissionRequest` | Permission dialog | Auto-approve/deny logic |
| `UserPromptSubmit` | Before processing input | Context injection |
| `Stop` | Main agent finishes | Summary generation |
| `SubagentStop` | Subagent finishes | Result aggregation |
| `PreCompact` | Before context compaction | Save important context |
| `Notification` | Alert sent | External logging |

### Hook Configuration Structure

```json
{
  "hooks": {
    "EventName": [
      {
        "matcher": "ToolPattern",
        "hooks": [
          {
            "type": "command",
            "command": "bash-command-here",
            "timeout": 60
          },
          {
            "type": "prompt",
            "prompt": "LLM evaluation prompt with $ARGUMENTS",
            "timeout": 30
          }
        ]
      }
    ]
  }
}
```

### Matcher Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| Exact string | Matches tool exactly | `"Write"` |
| Pipe-separated | Matches any listed | `"Edit\|Write\|MultiEdit"` |
| Regex | Regular expression | `"Notebook.*"` |
| Wildcard | Matches all tools | `"*"` |
| Empty/omitted | For non-tool events | `""` |
| Bash pattern | Matches commands | `"Bash(git push:*)"` |

### Hook Types

**Command Hooks:**
```json
{
  "type": "command",
  "command": "prettier --write $FILE",
  "timeout": 60
}
```

**Prompt Hooks:**
```json
{
  "type": "prompt",
  "prompt": "Evaluate if this operation is safe: $ARGUMENTS. Return JSON with {\"decision\": \"approve|block\", \"reason\": \"...\"}",
  "timeout": 30
}
```

### Exit Codes

| Code | Name | Behavior |
|------|------|----------|
| `0` | Success | Output processed, execution continues |
| `2` | Blocking Error | stderr used as error, operation blocked |
| Other | Non-blocking Error | stderr logged, execution continues |

### Environment Variables

Available in all hook commands:

| Variable | Description | Available In |
|----------|-------------|--------------|
| `CLAUDE_PROJECT_DIR` | Absolute project root path | All hooks |
| `CLAUDE_CODE_REMOTE` | `"true"` for web, empty for CLI | All hooks |
| `CLAUDE_ENV_FILE` | Path to session env file | `SessionStart` only |

### Practical Examples

**Auto-format on edits:**
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

**Run tests after changes:**
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

**Pre-commit validation:**
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash(git commit:*)",
        "hooks": [
          {
            "type": "command",
            "command": "git diff --cached --check",
            "timeout": 30
          }
        ]
      }
    ]
  }
}
```

**Session initialization:**
```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "source .env.local 2>/dev/null && echo 'export API_KEY=$API_KEY' > $CLAUDE_ENV_FILE",
            "timeout": 10
          }
        ]
      }
    ]
  }
}
```

### Hook Management Settings

```json
{
  "disableAllHooks": false,
  "allowManagedHooksOnly": false
}
```

| Setting | Type | Description |
|---------|------|-------------|
| `disableAllHooks` | boolean | Globally disable all hooks |
| `allowManagedHooksOnly` | boolean | (Enterprise) Only allow managed/SDK hooks |

---

## MCP Servers

### Model Context Protocol (MCP) Configuration

MCP servers extend Claude Code with external integrations (databases, APIs, file systems, etc.).

```json
{
  "enableAllProjectMcpServers": false,
  "enabledMcpjsonServers": ["filesystem", "github", "postgres"],
  "disabledMcpjsonServers": ["experimental-server"],
  "allowedMcpServers": [
    {
      "name": "approved-filesystem",
      "command": "/usr/local/bin/mcp-fs-server"
    }
  ],
  "deniedMcpServers": [
    {
      "name": "insecure-server",
      "command": "*"
    }
  ]
}
```

### Settings

| Setting | Type | Scope | Description |
|---------|------|-------|-------------|
| `enableAllProjectMcpServers` | boolean | User/Project | Auto-approve all `.mcp.json` servers |
| `enabledMcpjsonServers` | string[] | User/Project | Pre-approved server names |
| `disabledMcpjsonServers` | string[] | User/Project | Explicitly rejected servers |
| `allowedMcpServers` | object[] | Enterprise | Allowlist (only these permitted) |
| `deniedMcpServers` | object[] | Enterprise | Blocklist (always rejected) |

### MCP Server Objects

```json
{
  "name": "server-name",
  "command": "/path/to/server-executable"
}
```

**Wildcards supported:**
```json
{
  "name": "*",
  "command": "*/untrusted/*"
}
```

---

## Plugins

### Plugin Configuration

```json
{
  "enabledPlugins": {
    "formatter@acme-tools": true,
    "deployer@acme-tools": false,
    "linter@community": true
  },
  "extraKnownMarketplaces": {
    "acme-tools": {
      "source": {
        "source": "github",
        "repo": "acme-corp/claude-plugins"
      }
    }
  },
  "strictKnownMarketplaces": [
    {
      "source": "github",
      "repo": "company/approved-plugins"
    }
  ]
}
```

### Settings

| Setting | Type | Description |
|---------|------|-------------|
| `enabledPlugins` | object | Map of plugin names to enabled status |
| `extraKnownMarketplaces` | object | Additional plugin sources beyond official |
| `strictKnownMarketplaces` | array | (Enterprise) ONLY these marketplaces allowed |

---

## Environment Variables

### Core Variables

| Variable | Type | Description |
|----------|------|-------------|
| `ANTHROPIC_API_KEY` | string | API authentication key |
| `ANTHROPIC_AUTH_TOKEN` | string | Custom Authorization header |
| `CLAUDE_CODE_MAX_OUTPUT_TOKENS` | number | Limit response length |
| `MAX_THINKING_TOKENS` | number | Enable extended thinking |
| `DISABLE_TELEMETRY` | boolean | Opt out of analytics |
| `DISABLE_PROMPT_CACHING` | boolean | Disable caching globally |
| `CLAUDE_ENV_FILE` | string | Session env persistence file |
| `HTTP_PROXY` | string | HTTP proxy URL |
| `HTTPS_PROXY` | string | HTTPS proxy URL |
| `NO_PROXY` | string | Comma-separated bypass list |
| `CLAUDE_PROJECT_DIR` | string | Project root (auto-set) |
| `CLAUDE_CODE_REMOTE` | string | `"true"` for web environment |

### Usage in settings.json

```json
{
  "env": {
    "ANTHROPIC_API_KEY": "sk-ant-api03-...",
    "CLAUDE_CODE_MAX_OUTPUT_TOKENS": "8192",
    "MAX_THINKING_TOKENS": "10000",
    "DISABLE_TELEMETRY": "true",
    "HTTP_PROXY": "http://proxy.company.com:8080",
    "NO_PROXY": "localhost,127.0.0.1,.internal"
  }
}
```

---

## Best Practices

### 1. Layer Your Settings

**User settings** (`~/.claude/settings.json`):
- Personal preferences
- Global defaults
- Development conveniences

**Project settings** (`.claude/settings.json`):
- Team standards
- Project-specific permissions
- Shared hooks and MCP servers

**Local settings** (`.claude/settings.local.json`):
- Personal overrides
- Machine-specific paths
- API keys (DO NOT commit)

### 2. Security First

**Always deny sensitive files:**
```json
{
  "permissions": {
    "deny": [
      "Read(./.env)",
      "Read(./.env.*)",
      "Read(./secrets/**)",
      "Read(./config/credentials.*)",
      "Read(~/.ssh/id_*)"
    ]
  }
}
```

**Restrict destructive operations:**
```json
{
  "permissions": {
    "deny": [
      "Bash(rm -rf:*)",
      "Bash(sudo:*)",
      "Write(/etc/**)",
      "Write(/usr/**)"
    ]
  }
}
```

### 3. Use Sandbox Auto-Allow Wisely

```json
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true
  }
}
```

✅ **Good**: When you trust the sandbox isolation
❌ **Bad**: When running on production systems or handling sensitive data

### 4. Hook Best Practices

**Keep hooks fast:**
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "quick-lint.sh",
            "timeout": 10
          }
        ]
      }
    ]
  }
}
```

**Run expensive operations post-execution:**
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "npm test",
            "timeout": 300
          }
        ]
      }
    ]
  }
}
```

### 5. MCP Server Security

**Prefer allowlist approach (enterprise):**
```json
{
  "allowedMcpServers": [
    {"name": "approved-fs", "command": "/opt/mcp/fs-server"},
    {"name": "company-github", "command": "/opt/mcp/github"}
  ],
  "enableAllProjectMcpServers": false
}
```

**Development environments:**
```json
{
  "enabledMcpjsonServers": ["filesystem", "github"],
  "disabledMcpjsonServers": ["experimental", "beta"]
}
```

### 6. Attribution Configuration

**Minimal (just commits):**
```json
{
  "attribution": {
    "commit": "Co-Authored-By: Claude Code <claude-code@anthropic.com>"
  }
}
```

**Full (commits + PRs):**
```json
{
  "attribution": {
    "commit": "Co-Authored-By: Claude Code <claude-code@anthropic.com>",
    "pr": "\n\n---\n*This PR was created with assistance from Claude Code*\n*Model: claude-sonnet-4-5*"
  }
}
```

### 7. Model Selection

**Development:**
```json
{"model": "claude-sonnet-4-5-20250929"}
```

**Production/critical:**
```json
{"model": "claude-opus-4-5-20251101"}
```

**Fast iteration:**
```json
{"model": "claude-haiku-3-5-20241022"}
```

### 8. Environment Organization

```json
{
  "env": {
    "// Core": "",
    "NODE_ENV": "development",
    "LOG_LEVEL": "info",

    "// API Configuration": "",
    "API_BASE_URL": "https://api.company.com",
    "API_TIMEOUT": "30000",

    "// Feature Flags": "",
    "ENABLE_EXPERIMENTAL": "false",
    "USE_CACHE": "true"
  }
}
```

---

## Complete Example Configuration

See `claude-code-settings-template.json` for a complete, copy-ready template.

**Recommended minimal setup:**

```json
{
  "model": "claude-sonnet-4-5-20250929",
  "permissions": {
    "allow": ["Grep", "Glob", "Bash(npm run:*)"],
    "ask": ["Bash(git push:*)", "WebFetch"],
    "deny": ["Read(./.env)", "Read(./secrets/**)", "Bash(rm -rf:*)"]
  },
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true
  },
  "attribution": {
    "commit": "Co-Authored-By: Claude Code <claude-code@anthropic.com>"
  }
}
```

---

## Troubleshooting

### Settings Not Taking Effect

1. **Check precedence**: Enterprise > CLI args > Local > Project > User
2. **Verify file location**: Correct `.claude/` directory?
3. **JSON syntax**: Use `jq` or JSON validator
4. **Restart Claude Code**: Some settings require restart

### Permission Issues

```bash
# View active permissions
/permissions

# Test specific permission
# (create test scenario in Claude Code)
```

### Hook Not Executing

1. **Check matcher pattern**: Case-sensitive, exact match?
2. **Verify command exists**: `which command-name`
3. **Check timeout**: Increase if command is slow
4. **Review exit code**: Only 0 processes output
5. **Enable verbose mode**: Ctrl+O to see hook output

### Sandbox Problems

```json
{
  "sandbox": {
    "excludedCommands": ["problematic-command"],
    "allowUnsandboxedCommands": true  // temporary, security risk
  }
}
```

---

## Resources

### Official Documentation

- **Settings Reference**: https://code.claude.com/docs/en/settings
- **Hooks Reference**: https://code.claude.com/docs/en/hooks
- **Sandboxing Guide**: https://code.claude.com/docs/en/sandboxing
- **Admin Controls**: https://www.anthropic.com/news/claude-code-on-team-and-enterprise

### Community Resources

- **Claude Code Best Practices**: https://www.anthropic.com/engineering/claude-code-best-practices
- **Settings Guide (eesel AI)**: https://www.eesel.ai/blog/settings-json-claude-code
- **Admin Controls Guide**: https://www.eesel.ai/blog/admin-controls-claude-code

### Examples & Templates

- GitHub: Search for "claude-code-settings" for community examples
- Official examples in this repository

---

## Version History

- **2025-12-28**: Initial comprehensive reference
- Based on: Claude Code 2025 release
- Model IDs: Opus 4.5, Sonnet 4.5, Haiku 3.5

---

*This documentation reflects the state of Claude Code as of December 2025. Always refer to the official documentation for the most up-to-date information.*
