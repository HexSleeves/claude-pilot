# Golden Test Integration Guide

This document explains how to use the golden test framework integrated into the Claude Pilot project.

## Overview

Golden tests (also known as snapshot tests) ensure that CLI output remains consistent across changes. They compare actual command output against stored "golden" reference files, helping catch unintended changes to user-facing output.

## Project Structure

```
claude-pilot/
├── testdata/
│   └── golden/                    # Golden test files
│       ├── README.md             # Golden files documentation
│       ├── help/                 # Help command outputs
│       ├── output/               # Success scenario outputs
│       └── error/                # Error scenario outputs
├── packages/claudepilot/
│   ├── internal/testing/
│   │   └── golden.go             # Golden test framework
│   └── cmd_contract_test.go      # CLI contract tests
└── Makefile                      # Includes golden test targets
```

## Running Golden Tests

### Using Make Targets

```bash
# Run all golden tests
make test-golden

# Build CLI binary and run specific golden tests
make build-claudepilot
cd packages/claudepilot && go test -v -run "TestHelpOutputs"

# Run all contract tests
cd packages/claudepilot && go test -v -run "Test.*Contract"
```

### Using Go Test Directly

```bash
# Run all golden tests
cd packages/claudepilot
go test -v ./...

# Run specific test categories
go test -v -run "TestHelpOutputs"
go test -v -run "TestErrorScenarios"
go test -v -run "TestListOutputFormats"

# Run golden file integrity check
go test -v -run "TestGoldenFileIntegrity"
```

## Test Categories

### 1. Help Output Tests (`TestHelpOutputs`)

Tests all CLI help variations:
- Root command help (`--help`)
- Subcommand help (`create --help`, `list --help`, etc.)
- Version output (`--version`)

### 2. Output Format Tests (`TestListOutputFormats`)

Tests different output formats:
- Human-readable format (`--output human`)
- Table format (`--output table`)
- JSON format (`--output json`)
- NDJSON format (`--output ndjson`)
- Quiet format (`--output quiet`)

### 3. Command Tests

- **Create Command** (`TestCreateCommand`): Tests session creation
- **Details Command** (`TestDetailsCommand`): Tests session details display
- **Kill Command** (`TestKillCommand`): Tests session termination

### 4. Error Scenario Tests (`TestErrorScenarios`)

Tests error conditions:
- Session not found errors
- Validation errors
- Permission denied errors
- Network errors

### 5. Utility Tests

- **Golden File Integrity** (`TestGoldenFileIntegrity`): Validates all golden files exist
- **File Discovery** (`TestDiscoverAllGoldenFiles`): Tests golden file discovery

## Adding New Golden Tests

### 1. Create Test Case

```go
func TestMyNewCommand(t *testing.T) {
    tests := []goldentest.GoldenTestCase{
        {
            Name:       "my command success",
            Args:       []string{"my-command", "--flag", "value"},
            GoldenFile: "output/my-command-success.txt",
            WantCode:   0,
            Setup:      setupMyCommandTest,    // Optional
            Cleanup:    cleanupMyCommandTest, // Optional
        },
    }

    suite.RunTests(t, tests)
}
```

### 2. Create Golden File

Create the expected output file in `testdata/golden/`:

```
testdata/golden/output/my-command-success.txt
```

### 3. Set Up Test Helpers (Optional)

```go
func setupMyCommandTest() error {
    // Set up test environment
    return nil
}

func cleanupMyCommandTest() error {
    // Clean up after test
    return nil
}
```

## Golden File Conventions

### File Organization

- **`help/`**: Help text outputs
- **`output/`**: Success scenario outputs
- **`error/`**: Error scenario outputs

### File Naming

- Format: `[category]-[scenario]-[format].[ext]`
- Examples:
  - `list-success-human.txt`
  - `create-success-json.json`
  - `error-validation.txt`

### Content Guidelines

1. **Width Constraint**: All outputs constrained to 100 columns
2. **Consistent Data**: Use realistic but sanitized test data
3. **Timestamps**: Use consistent format: `2025-08-04 10:30:15`
4. **Paths**: Use generic examples: `/Users/user/Projects/...`
5. **UUIDs**: Use consistent patterns for reproducibility

## Updating Golden Files

### When Output Changes Intentionally

1. Update the relevant golden files manually
2. Or use the update helper:

```go
// In test setup
suite.UpdateGoldenFile("output/my-command.txt", newExpectedOutput)
```

### When Adding New Commands

1. Create golden files for all supported output formats
2. Add common error scenarios
3. Include help output if applicable

## Framework Features

### Environment Consistency

The framework automatically sets:
- `NO_COLOR=1`: Disables ANSI colors
- `TERM=xterm-256color`: Consistent terminal type
- `COLUMNS=100`: Consistent terminal width
- `LINES=24`: Consistent terminal height

### Error Filtering

Automatically filters common CI/environment errors:
- "Error getting terminal size: operation not supported on socket"
- "inappropriate ioctl for device"
- "not a terminal"

### JSON Validation

Special handling for JSON output:
- Validates JSON syntax
- Pretty-prints for comparison
- Structural comparison (order-independent)

### Output Normalization

- Trims trailing whitespace
- Normalizes line endings
- Handles terminal size differences

## Best Practices

### 1. Test Coverage

- Cover all output formats for each command
- Include both success and error scenarios
- Test edge cases (empty results, validation failures)

### 2. Golden File Maintenance

- Keep golden files up to date with output changes
- Use meaningful test data but keep it sanitized
- Document format changes in commit messages

### 3. CI Integration

- Golden tests run automatically in CI
- Failures indicate potential breaking changes
- Review golden file changes carefully during PR review

### 4. Test Isolation

- Use setup/cleanup functions for test isolation
- Don't rely on external state
- Mock external dependencies when possible

## Troubleshooting

### Common Issues

1. **Terminal Size Errors**: These are automatically filtered
2. **JSON Format Differences**: Framework handles JSON pretty-printing
3. **Whitespace Differences**: Output normalization handles most cases
4. **Path Differences**: Use relative paths or environment variables

### Debugging Failed Tests

1. Check the detailed diff output in test results
2. Verify golden file exists and is readable
3. Run command manually to see actual output
4. Check for environment-specific differences

### Updating After Changes

```bash
# Regenerate golden files (manual process)
./claude-pilot --help > testdata/golden/help/help-root.txt
./claude-pilot list --output json > testdata/golden/output/list-json.json

# Or use the framework's update functions in test code
```

## Integration with CI/CD

Golden tests are included in the standard test suite and run in CI. Any golden file mismatches will cause CI failures, ensuring output consistency is maintained across all changes.

```yaml
# Example GitHub Actions integration
- name: Run Golden Tests
  run: make test-golden
```

This framework helps maintain CLI contract stability while allowing for systematic updates when changes are intentional.