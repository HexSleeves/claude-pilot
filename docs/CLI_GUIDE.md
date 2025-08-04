# Claude Pilot CLI Guide

Complete reference guide for the Claude Pilot command-line interface.

## Table of Contents

- [Installation](#installation)
- [Global Flags](#global-flags)
- [Commands](#commands)
- [Output Formats](#output-formats)
- [Error Handling](#error-handling)
- [Environment Variables](#environment-variables)
- [Usage Patterns](#usage-patterns)
- [Examples](#examples)

## Installation

```bash
# Build from source
git clone https://github.com/your-org/claude-pilot
cd claude-pilot
make build

# The binary will be available at ./bin/claude-pilot
```

## Global Flags

These flags are available for all commands and control overall CLI behavior:

### Output Control

- `-o, --output string` - Output format (human|table|json|ndjson|quiet) (default: human)
- `--no-color` - Disable ANSI colors in output
- `--quiet` - Equivalent to `--output=quiet`

### Logging and Debug

- `-v, --verbose` - Enable verbose output
- `--debug` - Enable debug logging (overrides --verbose)
- `--trace` - Enable trace logging (overrides --debug and --verbose)

### Interaction

- `--yes` - Accept defaults for prompts (non-interactive mode)

### Configuration

- `--config string` - Config file path (default: ~/.config/claude-pilot/claude-pilot.yaml)

### Example Usage

```bash
# JSON output with debug logging
claude-pilot list --output=json --debug

# Non-interactive mode without colors
claude-pilot create my-session --yes --no-color

# Trace logging with custom config
claude-pilot attach --id my-session --trace --config ~/custom-config.yaml
```

## Commands

### `create` - Create a new session

Create a new Claude coding session with optional configuration.

```bash
# Basic usage
claude-pilot create [session-name]
claude-pilot create --id session-name

# With options
claude-pilot create --id my-session \
  --description "Working on feature X" \
  --project ~/my-project

# Attach to existing session
claude-pilot create --id new-session \
  --attach-to existing-session \
  --as-pane
```

**Flags:**

- `--id string` - Session name (alternative to positional argument)
- `--description string` - Session description
- `--project string` - Project directory path
- `--attach-to string` - Attach to existing session
- `--as-pane` - Attach as new pane
- `--as-window` - Attach as new window
- `--split string` - Split direction (horizontal|vertical)

**Output Formats:**

```bash
# Human-readable (default)
claude-pilot create my-session

# JSON with session details
claude-pilot create my-session --output=json

# Quiet mode (session ID only)
claude-pilot create my-session --output=quiet
```

### `list` - List sessions

Display all available sessions with filtering and sorting options.

```bash
# List all sessions
claude-pilot list

# Filter active sessions only
claude-pilot list --active

# Sort by creation date
claude-pilot list --sort=created
```

**Flags:**

- `--active` - Show only active sessions
- `--inactive` - Show only inactive sessions
- `--sort string` - Sort by field (name|created|updated) (default: name)
- `--id string` - Filter by specific session ID
- `--json` - JSON output (deprecated, use --output=json)

**Output Formats:**

```bash
# Human-readable table (default)
claude-pilot list

# Table format
claude-pilot list --output=table

# JSON format with schema
claude-pilot list --output=json

# NDJSON format for streaming
claude-pilot list --output=ndjson

# Quiet mode (IDs only)
claude-pilot list --output=quiet
```

### `details` - Show session details

Display detailed information about a specific session.

```bash
# Show session details
claude-pilot details --id my-session

# JSON output
claude-pilot details --id my-session --output=json
```

**Flags:**

- `--id string` - Session ID (required)

**Note:** Positional arguments are deprecated. Use `--id` flag instead.

### `attach` - Attach to session

Attach to an existing session for interactive use.

```bash
# Attach to session
claude-pilot attach --id my-session

# With deprecation warning (legacy)
claude-pilot attach my-session
```

**Flags:**

- `--id string` - Session ID (preferred over positional argument)

**Requirements:**

- Requires interactive terminal (TTY)
- Session must be running
- Use Ctrl+B, D to detach (tmux default)

### `kill` - Terminate sessions

Terminate one or more sessions safely.

```bash
# Kill specific session
claude-pilot kill --id my-session

# Kill all sessions
claude-pilot kill --all

# Force kill without confirmation
claude-pilot kill --id my-session --force

# Non-interactive kill
claude-pilot kill --all --yes
```

**Flags:**

- `--id string` - Session ID (alternative to positional argument)
- `--all` - Kill all sessions
- `--force` - Skip confirmation prompts
- `--json` - JSON output (deprecated, use --output=json)

**Safety Features:**

- Shows session details before confirmation
- Respects `--yes` flag for automation
- Provides operation summary
- Handles partial failures gracefully

### `tui` - Terminal UI mode

Launch the interactive terminal user interface.

```bash
# Launch TUI
claude-pilot tui

# With global flags forwarded
claude-pilot tui --debug --no-color
```

**Global Flag Forwarding:**
All global flags are automatically forwarded to the TUI mode:

- `--config` - Configuration file
- `--no-color` - Disable colors
- `--debug/--trace` - Logging levels
- `--yes` - Non-interactive mode

## Output Formats

### Human (`--output=human`)

Default format with colors, styling, and helpful messages.

```
Session Details
ID:          abc123
Name:        my-session
Status:      active
Created:     2024-01-15 10:30:45
Project:     /home/user/project
```

### Table (`--output=table`)

Structured table format suitable for terminals.

```
┌─────────┬────────────┬────────┬─────────┬─────────────────────┐
│ ID      │ Name       │ Status │ Project │ Created             │
├─────────┼────────────┼────────┼─────────┼─────────────────────┤
│ abc123  │ my-session │ active │ project │ 2024-01-15 10:30    │
└─────────┴────────────┴────────┴─────────┴─────────────────────┘
```

### JSON (`--output=json`)

Structured JSON format with schema version for automation.

```json
{
  "schemaVersion": "v1",
  "kind": "Session",
  "metadata": {
    "backend": "tmux",
    "operation": "details"
  },
  "item": {
    "id": "abc123",
    "name": "my-session",
    "status": "active",
    "createdAt": "2024-01-15T10:30:45Z",
    "project": "/home/user/project"
  }
}
```

### NDJSON (`--output=ndjson`)

Newline-delimited JSON for streaming and processing.

```json
{"schemaVersion":"v1","kind":"SessionListHeader","count":2}
{"schemaVersion":"v1","kind":"Session","item":{"id":"abc123","name":"session1"}}
{"schemaVersion":"v1","kind":"Session","item":{"id":"def456","name":"session2"}}
```

### Quiet (`--output=quiet`)

Minimal output with essential information only.

```
abc123
```

## Error Handling

### Exit Codes

| Code | Category     | Description                    |
|------|--------------|--------------------------------|
| 0    | Success      | Command completed successfully |
| 1    | Internal     | Internal errors                |
| 2    | Validation   | Invalid arguments or flags     |
| 3    | Not Found    | Resource not found             |
| 4    | Conflict     | Resource already exists        |
| 5    | Auth         | Permission denied              |
| 6    | Network      | Connection failed              |
| 7    | Timeout      | Operation timed out            |
| 8    | Unsupported  | Operation not supported        |

### Error Output

**Human Format:**

```
Error: Session 'non-existent' not found
Category: not_found
Code: session_not_found
Hint: Use 'claude-pilot list' to see available sessions.
```

**JSON Format:**

```json
{
  "schemaVersion": "v1",
  "kind": "Error",
  "error": {
    "code": "session_not_found",
    "category": "not_found",
    "message": "Session 'non-existent' not found",
    "hint": "Use 'claude-pilot list' to see available sessions.",
    "timestamp": "2024-01-15T10:30:45Z"
  }
}
```

## Environment Variables

### Behavior Control

- `FORCE_TTY=1` - Force TTY detection on
- `NO_TTY=1` - Force TTY detection off
- `NO_COLOR=1` - Disable colors (standard)
- `FORCE_COLOR=1` - Force colors on

### Configuration

- `CLAUDE_PILOT_CONFIG` - Configuration file path
- `CLAUDE_PILOT_VERBOSE=1` - Enable verbose logging
- `CLAUDE_PILOT_DEBUG=1` - Enable debug logging

### Example Usage

```bash
# Force non-interactive mode
NO_TTY=1 claude-pilot list --output=json

# Disable colors for automation
NO_COLOR=1 claude-pilot create my-session

# Custom configuration
CLAUDE_PILOT_CONFIG=~/my-config.yaml claude-pilot list
```

## Usage Patterns

### Interactive Usage

```bash
# Create and attach to session
claude-pilot create my-work-session --attach-to main --as-pane
claude-pilot attach --id my-work-session

# Browse and manage sessions
claude-pilot list --active
claude-pilot details --id my-session
claude-pilot kill --id old-session
```

### Automation and Scripting

```bash
#!/bin/bash

# Create session and capture ID
SESSION_ID=$(claude-pilot create work-session --output=quiet)

# Check if creation was successful
if [ $? -eq 0 ]; then
    echo "Created session: $SESSION_ID"

    # Get session details as JSON
    claude-pilot details --id "$SESSION_ID" --output=json > session-info.json

    # Clean up on exit
    trap "claude-pilot kill --id '$SESSION_ID' --force --quiet" EXIT
else
    echo "Failed to create session" >&2
    exit 1
fi
```

### CI/CD Integration

```bash
# Non-interactive session management
claude-pilot list --output=json | jq '.items[].id' | \
  xargs -I {} claude-pilot kill --id {} --yes --quiet

# Validate session creation
claude-pilot create test-session --yes --output=json | \
  jq -e '.item.status == "active"'
```

### Monitoring and Logging

```bash
# Stream session events
claude-pilot list --output=ndjson | while read line; do
    echo "$(date): $line" >> session-log.txt
done

# Debug session issues
claude-pilot attach --id problematic-session --debug 2> debug.log
```

## Examples

### Basic Workflow

```bash
# 1. Create a new session
claude-pilot create my-project --description "Working on feature X"

# 2. List active sessions
claude-pilot list --active

# 3. Attach to the session
claude-pilot attach --id my-project

# 4. Detach (Ctrl+B, D in tmux)

# 5. Check session details
claude-pilot details --id my-project

# 6. Clean up when done
claude-pilot kill --id my-project
```

### Advanced Usage

```bash
# Create session with custom project path
claude-pilot create work \
  --project ~/projects/my-app \
  --description "Development work" \
  --output=json

# Attach new pane to existing session
claude-pilot create logs \
  --attach-to work \
  --as-pane \
  --split horizontal

# Bulk operations
claude-pilot list --inactive --output=quiet | \
  xargs -I {} claude-pilot kill --id {} --force

# Monitor sessions in real-time
watch -n 5 'claude-pilot list --output=table'
```

### Integration Examples

**Shell Function:**

```bash
# Add to ~/.bashrc or ~/.zshrc
cpsession() {
    local name=${1:-$(basename $(pwd))}
    local id=$(claude-pilot create "$name" --output=quiet)
    if [ $? -eq 0 ]; then
        claude-pilot attach --id "$id"
    fi
}
```

**JSON Processing:**

```bash
# Get session count by status
claude-pilot list --output=json | \
  jq '.items | group_by(.status) | map({status: .[0].status, count: length})'

# Find sessions by project
claude-pilot list --output=json | \
  jq '.items[] | select(.project | contains("my-project"))'
```

**Error Handling:**

```bash
#!/bin/bash
set -e

# Function to handle errors
handle_error() {
    local exit_code=$1
    case $exit_code in
        2) echo "Invalid arguments provided" >&2 ;;
        3) echo "Session not found" >&2 ;;
        4) echo "Session already exists" >&2 ;;
        *) echo "Unexpected error (code: $exit_code)" >&2 ;;
    esac
    exit $exit_code
}

# Try to create session
claude-pilot create "$1" --output=quiet || handle_error $?
```

## Migration from Legacy Usage

### Deprecated Patterns

```bash
# OLD (deprecated but still works)
claude-pilot create session-name
claude-pilot attach session-name
claude-pilot kill session-name
claude-pilot list --json

# NEW (recommended)
claude-pilot create --id session-name
claude-pilot attach --id session-name
claude-pilot kill --id session-name
claude-pilot list --output=json
```

### Migration Timeline

- **Current release**: Deprecation warnings shown for legacy usage
- **Next minor release**: Enhanced warnings and migration guidance
- **Future major release**: Legacy patterns removed (with migration tools)

For detailed migration guidance, see [MIGRATION.md](MIGRATION.md).

## Troubleshooting

### Common Issues

**Command not found:**

```bash
# Ensure binary is in PATH
export PATH="$PATH:/path/to/claude-pilot/bin"
```

**Permission denied:**

```bash
# Check file permissions
ls -la ~/.config/claude-pilot/
chmod 600 ~/.config/claude-pilot/claude-pilot.yaml
```

**TTY errors:**

```bash
# For non-interactive environments
claude-pilot command --output=json --yes
```

**Session not found:**

```bash
# List available sessions
claude-pilot list --output=table

# Check exact session ID
claude-pilot list --output=quiet | grep partial-name
```

### Debug Information

```bash
# Enable debug logging
claude-pilot command --debug 2> debug.log

# Trace all operations
claude-pilot command --trace

# Check configuration
claude-pilot config show --output=json
```

### Getting Help

```bash
# Command-specific help
claude-pilot create --help
claude-pilot list --help

# Global help
claude-pilot --help

# Version information
claude-pilot version
```

For more detailed information, see:

- [Interface Contracts](CONTRACTS.md)
- [Migration Guide](MIGRATION.md)
- [JSON Schemas](schemas/)
