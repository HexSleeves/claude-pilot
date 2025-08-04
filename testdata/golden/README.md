# Golden Test Data

This directory contains golden files (expected output snapshots) for CLI contract testing. These files ensure that the CLI output format remains stable and catches unintended changes to user-facing messages and formats.

## Directory Structure

```bash
testdata/golden/
├── README.md                      # This documentation
├── help/                         # Help command outputs
│   ├── help-root.txt            # Root command help
│   ├── help-create.txt          # Create command help
│   ├── help-list.txt            # List command help
│   ├── help-details.txt         # Details command help
│   ├── help-kill.txt            # Kill command help
│   ├── help-attach.txt          # Attach command help
│   └── help-version.txt         # Version output
├── output/                       # Success scenario outputs
│   ├── list-human.txt           # Human-readable list format
│   ├── list-table.txt           # Table format list
│   ├── list-json.json           # JSON format list
│   ├── list-ndjson.ndjson       # NDJSON format list
│   ├── list-quiet.txt           # Quiet format list
│   ├── list-empty-human.txt     # Empty list human format
│   ├── list-empty-json.json     # Empty list JSON format
│   ├── create-success-human.txt # Successful create output
│   ├── create-success-json.json # Successful create JSON
│   ├── create-success-quiet.txt # Successful create quiet
│   ├── details-success-human.txt# Session details human format
│   ├── details-success-json.json# Session details JSON format
│   ├── kill-success-human.txt   # Kill success human format
│   ├── kill-success-json.json   # Kill success JSON format
│   ├── kill-all-success-human.txt # Kill all success
│   └── kill-partial-failure-human.txt # Partial kill failure
└── error/                        # Error scenario outputs
    ├── error-session-not-found.txt    # Session not found error
    ├── error-session-not-found-json.json # Session not found JSON
    ├── error-validation.txt       # Validation error format
    ├── error-validation-json.json # Validation error JSON
    ├── error-network.txt          # Network error format
    ├── error-network-json.json    # Network error JSON
    ├── error-permission-denied.txt # Permission error
    └── error-permission-denied-json.json # Permission error JSON
```

## Usage in Tests

### Test Structure

Golden files should be used in contract tests that verify:

1. **Help Output Stability**: Ensure help text format remains consistent
2. **Success Output Formats**: Verify all output formats (human, table, json, ndjson, quiet)
3. **Error Message Consistency**: Ensure error messages are helpful and consistent
4. **Schema Compliance**: JSON outputs match expected schemas

### Testing Pattern

```go
func TestCommandOutput(t *testing.T) {
    tests := []struct {
        name       string
        args       []string
        goldenFile string
        wantCode   int
    }{
        {
            name:       "list sessions human format",
            args:       []string{"list"},
            goldenFile: "output/list-human.txt",
            wantCode:   0,
        },
        {
            name:       "session not found error",
            args:       []string{"details", "--id", "nonexistent"},
            goldenFile: "error/error-session-not-found.txt",
            wantCode:   3,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Run command and compare with golden file
            assertGoldenOutput(t, tt.args, tt.goldenFile, tt.wantCode)
        })
    }
}
```

## Format Guidelines

### Width Constraints

All outputs are constrained to **100 columns** for consistency across different terminal sizes.

### Human Format

- Uses ANSI colors and Unicode symbols (● ○ ✓ ✗ ⚠ ℹ →)
- Includes helpful next steps and suggestions
- Tables use box-drawing characters for clarity

### JSON Format

- Follows the standardized schema structure
- Includes metadata for operation context
- Error responses include structured error objects
- Schema version for compatibility tracking

### Quiet Format

- Minimal output, typically just IDs or essential data
- No formatting, colors, or extra text
- Suitable for scripting and automation

### Error Formats

- Consistent error message structure
- Helpful troubleshooting information
- Suggestions for resolution
- Structured error codes and categories

## Maintenance

### Updating Golden Files

When CLI output format changes intentionally:

1. Update the relevant golden files
2. Ensure backward compatibility considerations
3. Update schema version if JSON structure changes
4. Verify all output formats are updated consistently

### Adding New Commands

For new commands, create golden files for:

- Help output (`help/help-[command].txt`)
- All supported output formats
- Common error scenarios
- Edge cases (empty results, validation errors)

### Schema Evolution

When updating JSON schemas:

1. Increment `schemaVersion` in JSON outputs
2. Maintain backward compatibility when possible
3. Document breaking changes in release notes
4. Update corresponding golden files

## Notes

- Golden files use realistic but sanitized data
- Timestamps use consistent format: `2025-08-04 10:30:15`
- File paths use generic examples: `/Users/user/Projects/...`
- UUIDs use consistent format patterns for reproducibility
- Terminal size errors are expected in CI environments
