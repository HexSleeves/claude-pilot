# Claude Pilot CLI Migration Guide

This guide helps users transition from legacy CLI patterns to the standardized interface introduced in v2.0.0. Follow this guide to update your scripts, automation, and workflows to use the new recommended patterns.

## Table of Contents

- [Migration Overview](#migration-overview)
- [Deprecated Features](#deprecated-features)
- [Migration Timeline](#migration-timeline)
- [Command Changes](#command-changes)
- [Output Format Changes](#output-format-changes)
- [Flag Changes](#flag-changes)
- [Automation Migration](#automation-migration)
- [Breaking Changes](#breaking-changes)
- [Migration Tools](#migration-tools)
- [Rollback Procedures](#rollback-procedures)

## Migration Overview

### What's Changing

The Claude Pilot CLI is transitioning to a more standardized, automation-friendly interface with:

- **Consistent global flags** across all commands
- **Structured output formats** (JSON, NDJSON) with schema versioning
- **Explicit flag-based arguments** instead of positional arguments
- **Standardized error handling** with specific exit codes
- **Enhanced TTY detection** for better automation support

### Why These Changes

1. **Automation Support**: Structured outputs and consistent exit codes enable reliable scripting
2. **Integration Friendly**: JSON schemas provide contracts for external integrations
3. **User Experience**: Better error messages with actionable hints
4. **Future Compatibility**: Extensible architecture for new features

### Compatibility Promise

- **Backward Compatibility**: Legacy patterns continue to work with deprecation warnings
- **Gradual Migration**: Minimum 2 minor versions of deprecation warnings before removal
- **Migration Tools**: Automated tools provided for common migration patterns

## Deprecated Features

### Legacy Flag Patterns

| Legacy Pattern          | Status        | Replacement           | Timeline    |
|-------------------------|---------------|-----------------------|-------------|
| `--json`               | Deprecated    | `--output=json`       | Remove v3.0 |
| Positional session ID  | Deprecated    | `--id session-name`   | Remove v3.0 |
| Mixed flag/positional  | Not supported | Use flags exclusively | Now         |

### Legacy Command Patterns

```bash
# DEPRECATED (but still works with warnings)
claude-pilot create session-name
claude-pilot attach session-name
claude-pilot kill session-name
claude-pilot details session-name
claude-pilot list --json

# RECOMMENDED (new patterns)
claude-pilot create --id session-name
claude-pilot attach --id session-name
claude-pilot kill --id session-name
claude-pilot details --id session-name
claude-pilot list --output=json
```

## Migration Timeline

### Current Release (v2.0.0)

**Status**: Deprecation warnings introduced
**Action Required**: Update scripts to use new patterns
**Timeline**: No breaking changes, all legacy patterns work

**What to do now:**

1. Update new scripts to use new patterns
2. Plan migration timeline for existing scripts
3. Test new patterns in development environments

### Next Minor Release (v2.1.0)

**Status**: Enhanced deprecation warnings
**Action Required**: Begin migrating critical automation
**Timeline**: 3 months after v2.0.0

**Enhanced warnings include:**

- Clear before/after examples in warning messages
- Migration suggestions with exact command replacements
- Documentation links for complex migrations

### Future Minor Release (v2.2.0)

**Status**: Final deprecation warnings
**Action Required**: Complete migration of all scripts
**Timeline**: 6 months after v2.0.0

**Final warnings include:**

- Countdown to removal in major version
- Migration verification tools
- Automated migration script availability

### Major Release (v3.0.0)

**Status**: Legacy patterns removed
**Action Required**: Must have completed migration
**Timeline**: 12 months after v2.0.0

**Breaking changes:**

- Legacy flags removed (`--json`)
- Positional arguments no longer accepted
- Old error message formats removed

## Command Changes

### create Command

**Legacy Patterns:**

```bash
# Positional argument (deprecated)
claude-pilot create my-session
claude-pilot create my-session --description "work"

# Mixed patterns (not recommended)
claude-pilot create my-session --project ~/code
```

**New Patterns:**

```bash
# Explicit flag-based (recommended)
claude-pilot create --id my-session
claude-pilot create --id my-session --description "work"
claude-pilot create --id my-session --project ~/code

# Auto-generated name
claude-pilot create --description "work" --project ~/code
```

**Migration Script:**

```bash
#!/bin/bash
# Migrate create commands
sed -i 's/claude-pilot create \([a-zA-Z0-9_-]\+\)/claude-pilot create --id \1/g' your-script.sh
```

### list Command

**Legacy Patterns:**

```bash
# JSON output (deprecated)
claude-pilot list --json

# Mixed usage
claude-pilot list --active --json
```

**New Patterns:**

```bash
# Structured output (recommended)
claude-pilot list --output=json
claude-pilot list --active --output=json

# Other formats
claude-pilot list --output=table
claude-pilot list --output=ndjson
claude-pilot list --output=quiet
```

**Migration Script:**

```bash
#!/bin/bash
# Migrate list JSON output
sed -i 's/--json/--output=json/g' your-script.sh
```

### details Command

**Legacy Patterns:**

```bash
# Positional argument (deprecated)
claude-pilot details my-session
claude-pilot details my-session --json
```

**New Patterns:**

```bash
# Explicit flag-based (recommended)
claude-pilot details --id my-session
claude-pilot details --id my-session --output=json
```

**Migration Script:**

```bash
#!/bin/bash
# Migrate details commands
sed -i 's/claude-pilot details \([a-zA-Z0-9_-]\+\)/claude-pilot details --id \1/g' your-script.sh
sed -i 's/details --id \([a-zA-Z0-9_-]\+\) --json/details --id \1 --output=json/g' your-script.sh
```

### attach Command

**Legacy Patterns:**

```bash
# Positional argument (deprecated)
claude-pilot attach my-session
```

**New Patterns:**

```bash
# Explicit flag-based (recommended)
claude-pilot attach --id my-session
```

**Migration Script:**

```bash
#!/bin/bash
# Migrate attach commands
sed -i 's/claude-pilot attach \([a-zA-Z0-9_-]\+\)/claude-pilot attach --id \1/g' your-script.sh
```

### kill Command

**Legacy Patterns:**

```bash
# Positional argument (deprecated)
claude-pilot kill my-session
claude-pilot kill my-session --json

# All sessions
claude-pilot kill --all --json
```

**New Patterns:**

```bash
# Explicit flag-based (recommended)
claude-pilot kill --id my-session
claude-pilot kill --id my-session --output=json

# All sessions
claude-pilot kill --all --output=json
```

**Migration Script:**

```bash
#!/bin/bash
# Migrate kill commands
sed -i 's/claude-pilot kill \([a-zA-Z0-9_-]\+\)/claude-pilot kill --id \1/g' your-script.sh
sed -i 's/kill --id \([a-zA-Z0-9_-]\+\) --json/kill --id \1 --output=json/g' your-script.sh
sed -i 's/kill --all --json/kill --all --output=json/g' your-script.sh
```

## Output Format Changes

### JSON Schema Evolution

**Legacy JSON Output (unstructured):**

```json
[
  {
    "id": "abc123",
    "name": "my-session",
    "status": "active"
  }
]
```

**New JSON Output (structured with schema):**

```json
{
  "schemaVersion": "v1",
  "kind": "SessionList",
  "metadata": {
    "backend": "tmux",
    "operation": "list"
  },
  "items": [
    {
      "id": "abc123",
      "name": "my-session",
      "status": "active",
      "createdAt": "2024-01-15T10:30:45Z",
      "updatedAt": "2024-01-15T10:30:45Z"
    }
  ],
  "count": 1
}
```

**Migration for JSON Processing:**

**Before (fragile):**

```bash
# Extract session IDs (legacy)
claude-pilot list --json | jq -r '.[].id'
```

**After (robust):**

```bash
# Extract session IDs (new schema)
claude-pilot list --output=json | jq -r '.items[].id'

# Or use quiet mode for simple ID extraction
claude-pilot list --output=quiet
```

**Automated JSON Migration:**

```bash
#!/bin/bash
# Update JSON processing scripts
sed -i 's/claude-pilot list --json | jq -r '\''\[\]/claude-pilot list --output=json | jq -r '\''.items[]/g' process-sessions.sh
```

### Error Output Changes

**Legacy Error Format:**

```bash
Failed to create session: session already exists
echo $? # Always 1
```

**New Error Format (Human):**

```bash
Error: Session 'my-session' already exists
Category: conflict
Code: session_already_exists
Hint: Use a different session name or attach to the existing session.
echo $? # 4 (conflict exit code)
```

**New Error Format (JSON):**

```json
{
  "schemaVersion": "v1",
  "kind": "Error",
  "error": {
    "code": "session_already_exists",
    "category": "conflict",
    "message": "Session 'my-session' already exists",
    "hint": "Use a different session name or attach to the existing session.",
    "timestamp": "2024-01-15T10:30:45Z"
  }
}
```

## Flag Changes

### Global Flags Introduction

**New Global Flags Available:**

```bash
# Output control
--output=human|table|json|ndjson|quiet  # Replaces command-specific --json
--no-color                              # Disable ANSI colors
--quiet                                 # Shorthand for --output=quiet

# Logging
--debug                                 # Debug logging (overrides --verbose)
--trace                                 # Trace logging (most verbose)

# Interaction
--yes                                   # Accept all prompts (automation)
```

**Migration Examples:**

```bash
# Before: Command-specific JSON
claude-pilot list --json
claude-pilot kill --id session --json

# After: Consistent global flag
claude-pilot list --output=json
claude-pilot kill --id session --output=json

# Before: No non-interactive mode
claude-pilot kill --id session  # Would prompt for confirmation

# After: Non-interactive automation
claude-pilot kill --id session --yes
```

### Command-Specific Flag Changes

**list command:**

- `--json` → `--output=json` (deprecated but works)
- New: `--inactive` flag (was mentioned in help but not implemented)
- New: `--id` flag for filtering single session

**create command:**

- Positional session name → `--id` flag (deprecated but works)
- New: Enhanced validation for flag combinations

**kill command:**

- `--json` → `--output=json` (deprecated but works)
- New: `--force` flag for bypassing confirmations

## Automation Migration

### Script Migration Patterns

**Legacy Automation Script:**

```bash
#!/bin/bash
set -e

# Create session
claude-pilot create work-session --description "Daily work"

# Check if creation was successful (unreliable)
if claude-pilot list --json | jq -e '.[] | select(.name == "work-session")' > /dev/null; then
    echo "Session created successfully"

    # Attach to session
    claude-pilot attach work-session

    # Clean up on exit
    trap 'claude-pilot kill work-session' EXIT
fi
```

**Migrated Automation Script:**

```bash
#!/bin/bash
set -e

# Enable error handling with specific exit codes
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

# Create session with explicit flags and structured output
SESSION_ID=$(claude-pilot create --id work-session \
    --description "Daily work" \
    --output=quiet) || handle_error $?

echo "Created session: $SESSION_ID"

# Verify session exists using structured JSON
SESSION_INFO=$(claude-pilot details --id "$SESSION_ID" --output=json)
SESSION_STATUS=$(echo "$SESSION_INFO" | jq -r '.item.status')

if [ "$SESSION_STATUS" = "active" ]; then
    echo "Session is active and ready"

    # Attach to session (requires TTY)
    if [ -t 0 ]; then
        claude-pilot attach --id "$SESSION_ID"
    else
        echo "Non-interactive environment, skipping attach"
    fi

    # Clean up on exit with non-interactive flag
    trap "claude-pilot kill --id '$SESSION_ID' --yes --output=quiet" EXIT
else
    echo "Session creation failed: status=$SESSION_STATUS" >&2
    exit 1
fi
```

### CI/CD Pipeline Migration

**Legacy CI Script:**

```bash
# Cleanup sessions (unreliable)
claude-pilot list --json | jq -r '.[].id' | xargs -I {} claude-pilot kill {}
```

**Migrated CI Script:**

```bash
# Robust session cleanup with proper error handling
cleanup_sessions() {
    echo "Cleaning up test sessions..."

    # Get all session IDs
    local session_ids
    session_ids=$(claude-pilot list --output=quiet 2>/dev/null || true)

    if [ -n "$session_ids" ]; then
        echo "Found sessions to clean up:"
        echo "$session_ids"

        # Kill all sessions non-interactively
        claude-pilot kill --all --yes --output=quiet || {
            echo "Warning: Some sessions could not be cleaned up" >&2
            # Don't fail the build for cleanup issues
            return 0
        }

        echo "Session cleanup completed"
    else
        echo "No sessions to clean up"
    fi
}

# Use in CI pipeline
cleanup_sessions
```

### JSON Processing Migration

**Legacy JSON Processing:**

```bash
# Fragile: Depends on unstructured format
get_active_sessions() {
    claude-pilot list --json | jq -r '.[] | select(.status == "active") | .id'
}
```

**Migrated JSON Processing:**

```bash
# Robust: Uses structured schema
get_active_sessions() {
    claude-pilot list --active --output=json | jq -r '.items[].id'
}

# Alternative: Use built-in filtering
get_active_sessions_simple() {
    claude-pilot list --active --output=quiet
}

# Schema-aware processing with validation
process_sessions() {
    local session_data
    session_data=$(claude-pilot list --output=json)

    # Validate schema version
    local schema_version
    schema_version=$(echo "$session_data" | jq -r '.schemaVersion')

    if [ "$schema_version" != "v1" ]; then
        echo "Warning: Unexpected schema version: $schema_version" >&2
    fi

    # Process sessions with schema awareness
    echo "$session_data" | jq -r '.items[] | "\(.id): \(.name) (\(.status))"'
}
```

## Breaking Changes

### Changes in v3.0.0

**Removed Features:**

1. `--json` flag on all commands
2. Positional arguments for session IDs
3. Legacy error message formats
4. Unstructured JSON output formats

**Changed Behavior:**

1. Exit codes now follow standard taxonomy (previously all errors were exit 1)
2. Error messages include structured information (category, code, hints)
3. JSON output follows schema with metadata

**Migration Required Before v3.0.0:**

- Update all `--json` flags to `--output=json`
- Convert positional session IDs to `--id` flag
- Update JSON processing to handle new schema structure
- Update error handling to use new exit codes

### Validation Script

Use this script to validate your migration before v3.0.0:

```bash
#!/bin/bash
# validate-migration.sh - Check for legacy patterns

echo "Checking for legacy patterns..."

# Check for deprecated flags
if grep -r "\-\-json" . --include="*.sh" --include="*.bash"; then
    echo "❌ Found deprecated --json flag usage"
    echo "   Replace with --output=json"
fi

# Check for positional session arguments
if grep -r "claude-pilot \(create\|attach\|kill\|details\) [a-zA-Z]" . --include="*.sh" --include="*.bash"; then
    echo "❌ Found positional session arguments"
    echo "   Replace with --id flag"
fi

# Check for legacy JSON processing
if grep -r "jq -r '\[\]'" . --include="*.sh" --include="*.bash"; then
    echo "❌ Found legacy JSON array processing"
    echo "   Update to use new schema structure"
fi

echo "Migration validation complete"
```

## Migration Tools

### Automated Migration Script

```bash
#!/bin/bash
# migrate-claude-pilot.sh - Automated migration for common patterns

migrate_file() {
    local file="$1"
    echo "Migrating $file..."

    # Create backup
    cp "$file" "$file.backup"

    # Apply migrations
    sed -i.tmp \
        -e 's/claude-pilot create \([a-zA-Z0-9_-]\+\)/claude-pilot create --id \1/g' \
        -e 's/claude-pilot attach \([a-zA-Z0-9_-]\+\)/claude-pilot attach --id \1/g' \
        -e 's/claude-pilot kill \([a-zA-Z0-9_-]\+\)/claude-pilot kill --id \1/g' \
        -e 's/claude-pilot details \([a-zA-Z0-9_-]\+\)/claude-pilot details --id \1/g' \
        -e 's/--json/--output=json/g' \
        -e 's/jq -r '\''\[\]/jq -r '\''.items[]/g' \
        "$file"

    # Remove temporary file
    rm "$file.tmp"

    echo "Migrated $file (backup: $file.backup)"
}

# Find and migrate shell scripts
find . -name "*.sh" -o -name "*.bash" | while read -r file; do
    if grep -q "claude-pilot" "$file"; then
        migrate_file "$file"
    fi
done

echo "Migration complete. Review changes and test thoroughly!"
```

### Migration Verification

```bash
#!/bin/bash
# verify-migration.sh - Test migrated scripts

echo "Verifying migration..."

# Test basic commands with new patterns
test_commands() {
    echo "Testing new command patterns..."

    # Test flag-based creation
    if claude-pilot create --id test-migration --output=quiet >/dev/null 2>&1; then
        echo "✅ Create with --id flag works"
        claude-pilot kill --id test-migration --yes --output=quiet >/dev/null 2>&1
    else
        echo "❌ Create with --id flag failed"
    fi

    # Test structured output
    if claude-pilot list --output=json | jq -e '.schemaVersion' >/dev/null 2>&1; then
        echo "✅ Structured JSON output works"
    else
        echo "❌ Structured JSON output failed"
    fi

    # Test exit codes
    claude-pilot details --id non-existent-session --output=quiet >/dev/null 2>&1
    if [ $? -eq 3 ]; then
        echo "✅ Exit codes work correctly"
    else
        echo "❌ Exit codes not working as expected"
    fi
}

# Run tests
test_commands

echo "Verification complete"
```

## Rollback Procedures

### Emergency Rollback

If you encounter issues with the new patterns:

**Immediate Rollback:**

```bash
# Restore from backups created by migration script
find . -name "*.backup" | while read -r backup; do
    original="${backup%.backup}"
    cp "$backup" "$original"
    echo "Restored $original"
done
```

**Version Rollback:**

```bash
# If using a specific CLI version, pin to previous version
# Update your installation/build process to use v1.x.x
git checkout v1.9.0  # or your last known good version
make build
```

### Gradual Rollback

Roll back specific patterns while keeping others:

```bash
# Rollback only JSON flag changes
sed -i 's/--output=json/--json/g' affected-script.sh

# Rollback only positional argument changes
sed -i 's/claude-pilot create --id \([a-zA-Z0-9_-]\+\)/claude-pilot create \1/g' affected-script.sh
```

### Compatibility Shims

Create wrapper functions for legacy compatibility:

```bash
# Add to ~/.bashrc or script headers
claude_pilot_legacy() {
    case "$1" in
        create|attach|kill|details)
            if [[ "$2" =~ ^[a-zA-Z0-9_-]+$ ]] && [[ "$2" != "--"* ]]; then
                # Convert positional to flag
                claude-pilot "$1" --id "$2" "${@:3}"
            else
                claude-pilot "$@"
            fi
            ;;
        list)
            # Convert --json to --output=json
            args=("$@")
            for i in "${!args[@]}"; do
                if [[ "${args[i]}" == "--json" ]]; then
                    args[i]="--output=json"
                fi
            done
            claude-pilot "${args[@]}"
            ;;
        *)
            claude-pilot "$@"
            ;;
    esac
}

# Use legacy wrapper
alias claude-pilot-legacy=claude_pilot_legacy
```

## Summary

This migration guide provides comprehensive guidance for transitioning to the new Claude Pilot CLI interface. Key points:

1. **Start Early**: Begin migration during deprecation period
2. **Test Thoroughly**: Validate changes in development environments
3. **Use Tools**: Leverage automated migration scripts
4. **Plan Rollback**: Maintain backups and rollback procedures
5. **Monitor Timeline**: Track deprecation schedule and plan accordingly

The new interface provides better automation support, structured outputs, and improved error handling. While migration requires effort, the benefits include more reliable scripting, better integration capabilities, and future-proof automation.

For questions or issues during migration:

- Check the [CLI Guide](CLI_GUIDE.md) for detailed usage examples
- Review [Interface Contracts](CONTRACTS.md) for stability guarantees
- Open GitHub issues for migration-specific problems
- Use GitHub Discussions for general migration questions
