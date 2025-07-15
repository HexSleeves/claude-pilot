# Claude Pilot Project Structure

## Repository Organization

### Root Level
```
├── packages/           # Monorepo packages
├── docs/              # Documentation (PRD, PLAN)
├── assets/            # Demo videos, tapes
├── templates/         # Configuration templates
├── scripts/           # Build and utility scripts
├── Makefile           # Build system
├── README.md          # Main documentation
└── LICENSE            # MIT license
```

### Package Structure
The project follows a monorepo pattern with clear separation of concerns:

#### `packages/claudepilot/` - Primary CLI
```
├── cmd/               # Cobra command definitions
│   ├── root.go       # Root command and initialization
│   ├── create.go     # Session creation
│   ├── list.go       # Session listing
│   ├── attach.go     # Session attachment
│   ├── kill.go       # Session termination
│   └── common.go     # Shared command utilities
├── internal/         # Private implementation
│   ├── ui/           # CLI-specific UI components
│   └── styles/       # CLI styling (legacy)
├── main.go           # Entry point
├── go.mod            # Module definition
└── go.sum            # Dependency checksums
```

#### `packages/tui/` - Terminal User Interface
```
├── internal/
│   ├── models/       # Bubble Tea models
│   │   ├── dashboard.go      # Main dashboard
│   │   ├── session_table.go  # Session table view
│   │   ├── detail_panel.go   # Session details
│   │   ├── summary_panel.go  # Summary information
│   │   └── create_modal.go   # Session creation modal
│   └── ui/           # TUI-specific UI logic
├── main.go           # TUI entry point
├── go.mod
└── go.sum
```

#### `packages/core/` - Business Logic
```
├── api/              # External API clients
├── internal/         # Core business logic
│   ├── config/       # Configuration management
│   ├── logger/       # Logging utilities
│   ├── multiplexer/  # Terminal multiplexer abstraction
│   │   ├── tmux.go   # Tmux implementation
│   │   └── zellij.go # Zellij implementation
│   ├── service/      # Business services
│   ├── storage/      # Data persistence
│   └── utils/        # Core utilities
├── go.mod
└── go.sum
```

#### `packages/shared/` - Common Components
```
├── components/       # Reusable UI components
│   ├── cards.go     # Card components
│   └── table.go     # Table components
├── interfaces/       # Shared interfaces
├── styles/          # Comprehensive theming system
│   ├── theme.go     # Main theme definitions
│   └── THEMING_STANDARDS.md # Styling guidelines
├── utils/           # Shared utilities
├── go.mod
└── go.sum
```

## Architectural Patterns

### Dependency Flow
```
claudepilot (CLI) ──┐
                    ├──> core (business logic)
tui (TUI) ──────────┘
                    └──> shared (common components)
```

### Configuration Hierarchy
1. **Default values** - Hardcoded sensible defaults
2. **Config file** - `~/.config/claude-pilot/claude-pilot.yaml`
3. **Environment variables** - `CLAUDE_PILOT_*` prefix
4. **Command flags** - Runtime overrides

### Data Storage
- **Session metadata**: JSON files in `~/.config/claude-pilot/sessions/`
- **Configuration**: YAML file in `~/.config/claude-pilot/`
- **Logs**: Optional logging to `~/.config/claude-pilot/claude-pilot.log`

## Naming Conventions

### Files & Directories
- **Snake_case** for directories: `internal/multiplexer/`
- **Lowercase** for Go files: `session_service.go`
- **CamelCase** for Go types: `SessionManager`
- **kebab-case** for executables: `claude-pilot`

### Go Code Style
- **Exported** functions start with capital letter
- **Private** functions start with lowercase letter
- **Interfaces** end with `-er` suffix when possible
- **Constants** in ALL_CAPS with underscores
- **Package names** are lowercase, single word when possible

### Session Management
- **Session IDs**: UUID format for uniqueness
- **Session names**: User-provided, alphanumeric with hyphens
- **Tmux sessions**: Prefixed with `claude-` by default
- **Metadata files**: Named by session ID with `.json` extension

## Development Guidelines

### Package Dependencies
- `core` package should have minimal external dependencies
- `shared` package provides common utilities to all other packages
- CLI and TUI packages consume `core` and `shared` but don't depend on each other
- Use local module replacements for development

### Error Handling
- Return errors explicitly, don't panic
- Wrap errors with context using `fmt.Errorf`
- Use structured logging when logging is enabled
- Graceful degradation when external tools (tmux/zellij) are unavailable

### Testing Structure
- Unit tests alongside source files (`*_test.go`)
- Integration tests in separate `integration/` directories
- Mock external dependencies (tmux, zellij, filesystem)
- Test both success and error paths