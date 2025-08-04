# Claude Pilot CLI JSON Schemas

This directory contains JSON Schema definitions for all structured CLI outputs from Claude Pilot. These schemas serve as contracts for automation and scripting use cases, ensuring predictable data formats across different CLI output modes.

## Schema Files

### Core Response Schemas

- **`session.json`** - Schema for individual session object responses
- **`session-list.json`** - Schema for session list responses  
- **`error.json`** - Schema for error responses
- **`operation-result.json`** - Schema for command operation results

## Schema Structure

All schemas follow a consistent structure:

```json
{
  "schemaVersion": "v1",
  "kind": "ResponseType",
  "metadata": { /* optional */ },
  "item|items|error|result": { /* response data */ }
}
```

### Common Fields

- **`schemaVersion`**: Version identifier for backward compatibility (currently "v1")
- **`kind`**: Response type identifier (`Session`, `SessionList`, `Error`, `OperationResult`)
- **`metadata`**: Optional key-value pairs with request context (requestId, timestamp, etc.)

## Versioning Strategy

### Current Version: v1

The initial release uses version "v1" across all schemas. This version includes:

- Complete field definitions for all CLI data structures
- Comprehensive validation rules and constraints
- Example data for each schema type
- Clear documentation for all fields

### Backward Compatibility Policy

Claude Pilot follows semantic versioning principles for schema changes:

#### Patch Changes (v1.0.x)
- Documentation updates
- Additional examples
- Non-breaking clarifications
- New optional validation rules

#### Minor Changes (v1.x.0)
- New optional fields
- Additional enum values
- Extended metadata fields
- Non-breaking enhancements

#### Major Changes (v2.0.0)
- Breaking field changes
- Required field modifications
- Enum value removals
- Structural changes

### Migration Guidance

When schema versions change, migration will be handled as follows:

1. **Automated Migration**: CLI will automatically migrate older schema responses when possible
2. **Deprecation Warnings**: Deprecated schema versions will include warnings in output
3. **Support Timeline**: Previous major versions supported for minimum 6 months
4. **Migration Tools**: Scripts provided for bulk data migration when needed

## Usage Examples

### Validating CLI Output

Use these schemas to validate JSON output from Claude Pilot commands:

```bash
# Generate JSON output
claude-pilot list --output json > sessions.json

# Validate against schema (using ajv-cli as example)
ajv validate -s docs/schemas/session-list.json -d sessions.json
```

### Automation Scripts

Example of parsing structured output in automation:

```bash
#!/bin/bash
# Get active sessions and process with jq
claude-pilot list --output json | jq -r '
  select(.kind == "SessionList") |
  .items[] |
  select(.status == "active") |
  .id
'
```

### Programming Language Integration

#### Python Example

```python
import json
import subprocess
from jsonschema import validate

# Run command and get JSON output
result = subprocess.run(['claude-pilot', 'list', '--output', 'json'], 
                       capture_output=True, text=True)
data = json.loads(result.stdout)

# Load and validate schema
with open('docs/schemas/session-list.json') as f:
    schema = json.load(f)

validate(instance=data, schema=schema)  # Raises exception if invalid

# Process validated data
for session in data['items']:
    print(f"Session: {session['name']} ({session['status']})")
```

#### JavaScript Example

```javascript
const { exec } = require('child_process');
const Ajv = require('ajv');
const addFormats = require('ajv-formats');

const ajv = new Ajv();
addFormats(ajv);

// Load schema
const schema = require('./docs/schemas/session-list.json');
const validate = ajv.compile(schema);

// Execute command
exec('claude-pilot list --output json', (error, stdout, stderr) => {
  if (error) throw error;
  
  const data = JSON.parse(stdout);
  
  // Validate response
  if (!validate(data)) {
    console.error('Invalid response:', validate.errors);
    return;
  }
  
  // Process validated data
  data.items.forEach(session => {
    console.log(`Session: ${session.name} (${session.status})`);
  });
});
```

## Schema Development

### Adding New Schemas

When adding new response types:

1. Create new schema file following naming convention: `{type}.json`
2. Include all required JSON Schema metadata (`$schema`, `$id`, `title`, `description`)
3. Define complete `properties` with validation rules
4. Add comprehensive `examples` section
5. Update this README with new schema documentation

### Schema Validation

All schemas are validated against JSON Schema Draft 07 specification. Use standard JSON Schema validators to ensure compliance:

```bash
# Validate schema structure
ajv compile -s docs/schemas/session.json
```

### Testing

Schemas should be tested against real CLI output:

```bash
# Test all output formats
for format in json ndjson; do
  claude-pilot list --output $format | ajv validate -s docs/schemas/session-list.json
done
```

## Support

For questions about schema usage or to report schema issues:

1. Check existing CLI output matches schema definitions
2. Verify schema version matches CLI version compatibility
3. Review migration documentation for version changes
4. Open issue with schema validation details and CLI version info