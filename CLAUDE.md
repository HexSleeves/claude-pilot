# CLAUDE.md

Claude Pilot is a CLI tool for managing multiple Claude Code sessions with tmux/zellij support.

## Architecture

- **`cmd/`** - CLI commands (Cobra framework)
  - `root.go` - Main command setup and TUI mode detection
  - `create.go` - Session creation
  - `list.go` - Session listing
  - `attach.go` - Session attachment
  - `kill.go` - Session termination
  - `common.go` - Shared command initialization logic
- **`internal/config/`** - Configuration management
- **`internal/service/`** - Business logic for session management
- **`internal/multiplexer/`** - tmux/zellij backend implementations
- **`internal/storage/`** - JSON session persistence with file repository pattern
- **`internal/ui/`** - UI components (colors, tables, session formatting)
- **`internal/logger/`** - Structured logging with configurable levels
- **`internal/interfaces/`** - Core interfaces and contracts
- **`internal/utils/`** - Utility functions (filtering, etc.)

## Key Patterns

- Interface-based design for multiplexer abstraction
- Repository pattern for session storage
- CommandContext for unified CLI initialization and dependency injection
- Structured logging with configurable levels
- Clean separation between UI, business logic, and storage layers

## Development

```bash
make build      # Build application
make test       # Run tests
make test-race  # Run tests with race detection
make fmt        # Format code
make lint       # Lint code (requires golangci-lint)
make dev        # Build with debug info
make demo       # Run complete demo
make help       # Show all available targets
```

## Dependencies

- **Cobra** - CLI framework
- **Viper** - Configuration management
- **go-pretty** - Table formatting
- **fatih/color** - Terminal colors
- **google/uuid** - UUID generation

## Configuration

- Config: `~/.config/claude-pilot/.claude-pilot.yaml`
- Sessions: `~/.config/claude-pilot/sessions/`
- Backends: auto-detects tmux/zellij (prefers tmux)
- Environment: `CLAUDE_PILOT_*` variables supported
- Logging: configurable via LOG_LEVEL environment variable

## Session States

- `active` - Running and available
- `inactive` - Stopped but metadata exists
- `connected` - User currently attached
- `error` - Error state

## Commands

- `create` - Create new session with optional description
- `list` - List sessions (supports --all flag)
- `attach` - Attach to existing session
- `kill` - Terminate specific session
- `kill-all` - Terminate all sessions

Focus on clean interfaces, performance, and user experience consistency.
