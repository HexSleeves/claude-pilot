# Claude Pilot CLI Contracts

This document defines the formal external interface contracts and stability guarantees for the Claude Pilot CLI. These contracts ensure consistent behavior across releases and provide the foundation for automation and integration.

## Table of Contents

- [Interface Contracts](#interface-contracts)
- [Stability Guarantees](#stability-guarantees)
- [Output Schemas](#output-schemas)
- [Error Contracts](#error-contracts)
- [Performance Contracts](#performance-contracts)
- [Testing Contracts](#testing-contracts)
- [Versioning Policy](#versioning-policy)

## Interface Contracts

### Command Signatures

All commands follow consistent patterns and maintain backward compatibility within major versions.

#### Global Flags Contract

These flags are guaranteed to be available for all commands:

```bash
Global Flags:
  -o, --output string    Output format (human|table|json|ndjson|quiet) (default "human")
      --debug           Enable debug logging (overrides --verbose)
      --trace           Enable trace logging (overrides --debug and --verbose)
      --no-color        Disable ANSI colors in output
      --yes             Accept defaults for prompts (non-interactive mode)
  -v, --verbose         Verbose output
      --config string   Config file path
```

**Contract Guarantees:**

- All global flags work with all commands
- Flag precedence: `--trace` > `--debug` > `--verbose`
- `--no-color` disables all ANSI escape sequences
- `--yes` makes all commands non-interactive
- Default output format is always `human`

#### Command-Specific Flags

**create command:**

```bash
claude-pilot create [session-name] [flags]

Required: None (session name can be auto-generated)
Optional:
  --id string           Session name (preferred over positional)
  --description string  Session description
  --project string      Project directory path
  --attach-to string    Attach to existing session
  --as-pane            Attach as new pane
  --as-window          Attach as new window
  --split string       Split direction (horizontal|vertical)
```

**list command:**

```bash
claude-pilot list [flags]

Optional:
  --active      Show only active sessions
  --inactive    Show only inactive sessions
  --sort string Sort by field (name|created|updated) (default "name")
  --id string   Filter by specific session ID
  --json        JSON output (deprecated, use --output=json)
```

**details command:**

```bash
claude-pilot details --id <session-id> [flags]

Required:
  --id string   Session ID

Note: Positional arguments deprecated but supported with warnings
```

**attach command:**

```bash
claude-pilot attach --id <session-id> [flags]

Required:
  --id string   Session ID (positional supported with deprecation warning)

Requirements:
  - Interactive terminal (TTY)
  - Session must exist and be running
```

**kill command:**

```bash
claude-pilot kill [--id <session-id> | --all] [flags]

Required (one of):
  --id string   Session ID
  --all         Kill all sessions

Optional:
  --force       Skip confirmation prompts
  --json        JSON output (deprecated, use --output=json)
```

**tui command:**

```bash
claude-pilot tui [flags]

All global flags are forwarded to the TUI process
```

### Exit Code Contract

Standard exit codes are guaranteed across all commands:

| Code | Category     | Description                    | Usage                           |
|------|--------------|--------------------------------|---------------------------------|
| 0    | Success      | Command completed successfully | Normal operation completion     |
| 1    | Internal     | Internal errors                | Unexpected failures             |
| 2    | Validation   | Invalid arguments or flags     | User input errors               |
| 3    | Not Found    | Resource not found             | Session/resource doesn't exist  |
| 4    | Conflict     | Resource already exists        | Duplicate creation attempts     |
| 5    | Auth         | Permission denied              | Access control failures         |
| 6    | Network      | Connection failed              | Multiplexer connection issues   |
| 7    | Timeout      | Operation timed out            | Long-running operation limits   |
| 8    | Unsupported  | Operation not supported        | TTY required, feature disabled  |

**Contract Guarantees:**

- Exit codes are consistent across all commands
- Code 0 always indicates success
- Codes 2-8 indicate specific error categories
- Automation can rely on exit codes for error handling

### Environment Variables Contract

These environment variables are guaranteed to be respected:

**Standard Environment Variables:**

- `NO_COLOR` - Disable colors (follows [no-color.org](https://no-color.org/) standard)
- `FORCE_COLOR` - Force colors even in non-TTY environments
- `TERM` - Terminal type detection for color support

**CLI-Specific Environment Variables:**

- `FORCE_TTY=1` - Force TTY detection on (testing/debugging)
- `NO_TTY=1` - Force TTY detection off (automation)
- `CLAUDE_PILOT_CONFIG` - Configuration file path
- `CLAUDE_PILOT_VERBOSE=1` - Enable verbose logging
- `CLAUDE_PILOT_DEBUG=1` - Enable debug logging

**Contract Guarantees:**

- Environment variables take precedence over default values
- CLI-specific variables use `CLAUDE_PILOT_` prefix
- Boolean variables use `1`/`true` for enabled, any other value for disabled

## Stability Guarantees

### Semantic Versioning Policy

Claude Pilot CLI follows [Semantic Versioning 2.0.0](https://semver.org/):

**MAJOR.MINOR.PATCH (e.g., 2.1.3)**

- **MAJOR**: Breaking changes to command interfaces, flag behavior, or output schemas
- **MINOR**: New features, commands, or flags (backward compatible)
- **PATCH**: Bug fixes, performance improvements (backward compatible)

### Backward Compatibility Windows

**Within Major Versions:**

- All command signatures remain stable
- All flag names and behaviors remain consistent
- All output formats maintain schema compatibility
- Exit codes remain unchanged
- Environment variables continue to work

**Deprecation Process:**

1. **Deprecation Announcement**: Feature marked as deprecated with warning messages
2. **Deprecation Period**: Minimum 2 minor versions of warning messages and migration guidance
3. **Removal**: Breaking change in next major version only

**Example Deprecation Timeline:**

- v2.1.0: `--json` flag deprecated, warnings shown, `--output=json` recommended
- v2.2.0: Enhanced deprecation warnings with migration examples
- v2.3.0: Final deprecation warnings
- v3.0.0: `--json` flag removed

### Schema Evolution Policy

**JSON Schema Versioning:**

- All JSON outputs include `schemaVersion` field
- Schema versions follow independent versioning (v1, v2, etc.)
- New fields can be added without version bump
- Field removal or type changes require version bump
- Multiple schema versions supported simultaneously during transitions

**Schema Compatibility Matrix:**

| CLI Version | Supported Schema Versions | Default Version |
|-------------|---------------------------|-----------------|
| 2.0.x       | v1                        | v1              |
| 2.1.x       | v1, v2                    | v1              |
| 2.2.x       | v1, v2                    | v2              |
| 3.0.x       | v2, v3                    | v2              |

## Output Schemas

### Schema Structure Guarantees

All structured outputs follow consistent patterns:

**Top-level Structure:**

```json
{
  "schemaVersion": "v1",
  "kind": "ResourceType",
  "metadata": {},
  "item": {} | "items": []
}
```

**Required Fields:**

- `schemaVersion`: Always present, indicates schema version
- `kind`: Resource type identifier (Session, SessionList, Error, OperationResult)

**Metadata Fields:**

- `backend`: Multiplexer backend (e.g., "tmux")
- `operation`: Command that generated the output
- `timestamp`: ISO 8601 timestamp when applicable
- `requestId`: Unique request identifier for debugging

### Session Object Schema (v1)

```json
{
  "id": "string (UUID format)",
  "name": "string",
  "description": "string (optional)",
  "project": "string (optional)",
  "status": "active|inactive|error",
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)",
  "attachedTo": "string (optional)",
  "windowCount": "integer (optional)",
  "paneCount": "integer (optional)"
}
```

**Field Contracts:**

- `id`: Always present, unique identifier
- `name`: Human-readable session name
- `status`: One of predefined enum values
- `createdAt`/`updatedAt`: Always valid ISO 8601 timestamps
- Optional fields may be null or absent

### Session List Schema (v1)

```json
{
  "schemaVersion": "v1",
  "kind": "SessionList",
  "metadata": {
    "backend": "tmux",
    "operation": "list",
    "filters": {
      "active": true,
      "sort": "name"
    }
  },
  "items": [],
  "count": 0
}
```

**Field Contracts:**

- `items`: Array of Session objects (may be empty)
- `count`: Always matches `items.length`
- `metadata.filters`: Applied filters and sorting

### Error Schema (v1)

```json
{
  "schemaVersion": "v1",
  "kind": "Error",
  "error": {
    "code": "session_not_found",
    "category": "not_found",
    "message": "Session 'xyz' not found",
    "hint": "Use 'claude-pilot list' to see available sessions.",
    "details": {},
    "timestamp": "2024-01-15T10:30:45Z",
    "requestId": "abc123"
  }
}
```

**Field Contracts:**

- `code`: Machine-readable error identifier
- `category`: Error category from defined enum
- `message`: Human-readable error description
- `hint`: Actionable remediation guidance (optional)
- `timestamp`: Always present in ISO 8601 format

### Operation Result Schema (v1)

```json
{
  "schemaVersion": "v1",
  "kind": "OperationResult",
  "result": {
    "success": true,
    "message": "Session created successfully",
    "data": {},
    "errors": [],
    "metadata": {
      "operation": "create",
      "duration": "1.23s"
    }
  }
}
```

**Field Contracts:**

- `success`: Boolean indicating operation outcome
- `message`: Human-readable result description
- `data`: Operation-specific result data (optional)
- `errors`: Array of error objects for partial failures
- `metadata`: Operation context and metrics

## Error Contracts

### Error Categories

Fixed set of error categories with guaranteed meanings:

| Category     | Exit Code | Description                  | Examples                    |
|--------------|-----------|------------------------------|-----------------------------|
| validation   | 2         | Invalid user input           | Missing required flags      |
| not_found    | 3         | Resources don't exist        | Session not found           |
| conflict     | 4         | Resource conflicts           | Session already exists      |
| auth         | 5         | Permission issues            | Access denied               |
| network      | 6         | Connection problems          | Multiplexer unreachable     |
| timeout      | 7         | Operation timeouts           | Command took too long       |
| unsupported  | 8         | Feature not available        | TTY required but not found  |
| internal     | 1         | Internal system errors       | Unexpected failures         |

### Error Code Patterns

Error codes follow consistent naming patterns:

**Pattern:** `{resource}_{action}_{condition}`

**Examples:**

- `session_not_found` - Session resource doesn't exist
- `session_already_exists` - Session creation conflict
- `multiplexer_connection_failed` - Network connectivity issue
- `invalid_flag_value` - Validation error
- `permission_denied` - Authorization failure

### Hint Guidelines

Error hints provide actionable remediation guidance:

**Hint Patterns:**

- "Use 'command --flag' to resolve this issue"
- "Check [resource] and try again"
- "Ensure [condition] is met before retrying"

**Examples:**

```
Error: Session 'xyz' not found
Hint: Use 'claude-pilot list' to see available sessions.

Error: Permission denied accessing session directory
Hint: Check file permissions and user access rights.

Error: Cannot attach to session in non-interactive environment
Hint: This command requires an interactive terminal. Use --yes flag for non-interactive mode.
```

## Performance Contracts

### Response Time Budgets

Guaranteed maximum response times for common operations:

| Operation            | Budget | Measurement                    |
|----------------------|--------|--------------------------------|
| `--help` commands    | 100ms  | Time to display help text      |
| `--version` command  | 50ms   | Time to display version info   |
| `list` (empty)       | 200ms  | Time to display "no sessions"  |
| `list` (<10 sessions)| 500ms  | Time to display session table  |
| Command validation   | 50ms   | Time to validate flags/args    |

**Contract Guarantees:**

- Response times measured on standard hardware (GitHub Actions runners)
- Budgets apply to 95th percentile of measurements
- Regression tests fail if budgets are exceeded
- Network-dependent operations excluded from budgets

### Resource Usage Guidelines

**Memory Usage:**

- CLI commands should not exceed 50MB resident memory
- TUI mode may use up to 100MB for interface rendering
- Memory usage should not grow over time (no leaks)

**CPU Usage:**

- CLI commands should not use >50% CPU for more than 1 second
- Background operations should be minimal
- No CPU-intensive operations during help display

**Disk Usage:**

- Configuration files should not exceed 1MB
- Session metadata should not exceed 10KB per session
- Log files should rotate and not exceed 10MB total

## Testing Contracts

### Golden File Testing

**Contract:** All user-facing output is validated against golden files to prevent unintended changes.

**Coverage:**

- All command help outputs
- All output formats (human, table, json, ndjson, quiet)
- All error message formats
- Success and failure scenarios

**Process:**

1. Golden files stored in `testdata/golden/`
2. CI validates current output against golden files
3. Intentional changes require golden file updates
4. Breaking changes trigger major version bumps

### Schema Validation

**Contract:** All JSON outputs validate against published schemas.

**Process:**

1. JSON Schema files in `docs/schemas/`
2. All structured outputs validated in CI
3. Schema changes follow versioning policy
4. Backward compatibility maintained

### Cross-Platform Testing

**Contract:** CLI works consistently across supported platforms.

**Coverage:**

- Linux (Ubuntu latest)
- macOS (latest)
- Go versions: current and previous stable

**Validation:**

- All commands execute successfully
- Output formats are consistent
- Exit codes are identical
- File paths work correctly

### Shell Compatibility

**Contract:** CLI works in common shell environments.

**Coverage:**

- bash
- zsh
- fish (basic compatibility)

**Validation:**

- Commands execute without shell-specific errors
- Output displays correctly
- Environment variables work as expected

## Versioning Policy

### Release Schedule

**Regular Releases:**

- **Patch releases**: Monthly or as needed for critical fixes
- **Minor releases**: Quarterly with new features
- **Major releases**: Annually or when breaking changes accumulate

**Emergency Releases:**

- Security vulnerabilities
- Data loss bugs
- Critical functionality failures

### Version Compatibility

**Supported Versions:**

- Current major version: Full support
- Previous major version: Security fixes only
- Older versions: No support

**Example Timeline:**

- v3.0.0 released: v3.x receives full support, v2.x security only
- v2.x EOL: 12 months after v3.0.0 release
- v4.0.0 released: v2.x support ends, v3.x security only

### Migration Support

**Breaking Changes:**

- Advance notice in release notes
- Migration guides provided
- Automated migration tools when possible
- Deprecation warnings in preceding minor versions

**Schema Migrations:**

- New schema versions introduced gradually
- Old versions supported during transition period
- Clear migration timeline communicated
- Tools provided for bulk data migration

### API Stability

**Guaranteed Stable:**

- Command names and basic syntax
- Global flag names and behavior
- Exit codes and meanings
- Core output schema fields

**May Change Between Majors:**

- New commands and flags (additive)
- Enhanced error messages
- Additional output fields (non-breaking)
- Performance optimizations

**Requires Major Version:**

- Command removal or renaming
- Flag behavior changes
- Exit code changes
- Output schema breaking changes
- Minimum system requirement changes

## Contract Validation

### Automated Testing

All contracts are validated through automated testing:

1. **Unit Tests**: Validate individual component contracts
2. **Integration Tests**: Validate end-to-end command behavior
3. **Contract Tests**: Validate output schemas and formats
4. **Performance Tests**: Validate response time budgets
5. **Compatibility Tests**: Validate cross-platform behavior

### Manual Review Process

Contract changes require:

1. **Design Review**: Contract changes reviewed by maintainers
2. **Deprecation Planning**: Breaking changes require migration plan
3. **Documentation Updates**: All changes documented
4. **User Communication**: Breaking changes announced in advance

### Monitoring and Feedback

**User Feedback Channels:**

- GitHub Issues for bug reports
- GitHub Discussions for feature requests
- Release notes for change announcements

**Metrics Tracking:**

- Command usage patterns
- Error rates by category
- Performance regressions
- Schema validation failures

This contract document is living and evolves with the CLI. All changes are versioned and tracked to maintain transparency and predictability for users and integrations.
