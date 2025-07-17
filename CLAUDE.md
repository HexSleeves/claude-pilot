# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Claude Pilot is a CLI tool for managing multiple Claude Code sessions with terminal multiplexer (tmux) support, built with a modular monorepo architecture.

## Architecture

### Monorepo Structure

- **`packages/claudepilot/`** - Main CLI application (Go 1.24.5)
- **`packages/core/`** - Core business logic and shared components
- **`packages/tui/`** - Terminal User Interface using Charm's Bubbletea

### Package Organization

#### packages/claudepilot/ (Main CLI)

- **`cmd/`** - Cobra command definitions (create, list, attach, kill)
- **`internal/styles/`** - Lipgloss styling and compatibility layer
- **`internal/ui/`** - CLI UI components (colors, formatting, tables)
- **`main.go`** - Entry point with Charm Fang signal handling

#### packages/core/ (Business Logic)

- **`internal/config/`** - Viper-based configuration management
- **`internal/service/`** - Session management business logic
- **`internal/multiplexer/`** - tmux backend implementation
- **`internal/storage/`** - JSON session persistence with repository pattern
- **`internal/logger/`** - Structured logging with configurable levels
- **`internal/interfaces/`** - Core interfaces and contracts

#### packages/tui/ (Terminal UI)

- **`internal/models/`** - Bubbletea model implementations
- **`internal/ui/`** - Interactive TUI interface components

## Development Commands

### Core Build Commands

```bash
make build          # Build main binary from claudepilot package
make dev            # Build with debug info (-gcflags="all=-N -l")
make run            # Build and run application
make run-dev        # Run directly with 'go run .'
make clean          # Clean artifacts and session data
```

## Key Dependencies

### CLI Framework

- **Cobra** - Command structure and CLI framework
- **Viper** - Configuration management with YAML support

### UI & Styling

- **Charm Lipgloss** - Terminal styling and colors
- **Charm Bubbletea** - Interactive TUI framework
- **Charm Fang** - Signal handling
- **go-pretty** - Table formatting for CLI output
- **fatih/color** - Additional color support

### Core Libraries

- **google/uuid** - Session ID generation

## Architecture Patterns

### Design Principles

- **Interface-based design** for multiplexer abstraction
- **Repository pattern** for session storage
- **CommandContext pattern** for unified CLI initialization and dependency injection
- **Clean architecture** with separation between UI, business logic, and storage

### Command Structure (Cobra)

- Place command definitions in `cmd/` following factory pattern (`newXCmd()`)
- Keep business logic in `internal/` packages, avoid logic inside command files
- Use `RunE`/`PreRunE` for proper error propagation
- Provide both `Short` and `Long` descriptions for commands

### TUI Development (Bubbletea)

- Follow `Init`, `Update`, `View` triad strictly
- Keep side-effects (network, FS) in separate goroutines; use messages for results
- Style output using Lipgloss with centralized theme definitions
- Maintain immutable model updates; avoid mutating shared state

## Configuration

### File Locations

- Config: `~/.config/claude-pilot/claude-pilot.yaml`
- Sessions: `~/.config/claude-pilot/sessions/`
- Template: `templates/claude-pilot.yaml`

### Environment Variables

- `CLAUDE_PILOT_*` - Configuration overrides
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

### Backend Support

- **tmux** - Current backend with session prefix support
- **zellij** - Planned future backend with layout file support
- Auto-detection with tmux as current default

## Session Management

### Session States

- `active` - Running and available
- `inactive` - Stopped but metadata exists
- `connected` - User currently attached
- `error` - Error state

### CLI Commands

- `create` - Create session with optional description
- `list` - List sessions (use `--all` for detailed view)
- `attach` - Attach to existing session
- `kill` - Terminate specific session
- `kill-all` - Terminate all sessions

## Testing & Quality

### Testing Commands

```bash
make test           # Standard test suite
make test-race     # Race condition detection
```

### Code Quality

- Static analysis with Datadog rulesets (go-best-practices, go-security)
- golangci-lint for comprehensive linting
- Cross-platform build verification

## Error Handling

- Wrap errors using `%w` with context (`fmt.Errorf`)
- Use structured logging with configurable levels
- Proper error propagation through command chain
- Graceful handling of multiplexer backend failures
