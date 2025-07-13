# CLAUDE.md

Claude Pilot is a CLI tool for managing multiple Claude Code sessions with tmux/zellij support.

## Architecture

- **`cmd/`** - CLI commands (Cobra framework)
- **`internal/service/`** - Business logic
- **`internal/multiplexer/`** - tmux/zellij implementations
- **`internal/tui/`** - Terminal UI (Bubble Tea)
- **`internal/storage/`** - JSON session persistence

## Key Patterns

- Interface-based design for multiplexer abstraction
- Repository pattern for session storage
- CommandContext for unified CLI initialization
- Bubble Tea Model-View-Update for TUI

## Development

```bash
make build    # Build application
make test     # Run tests
make tui      # Launch TUI mode
```

## Configuration

- Config: `~/.config/claude-pilot/.claude-pilot.yaml`
- Sessions: `~/.config/claude-pilot/sessions/`
- Backends: auto-detects tmux/zellij (prefers tmux)

## Session States

- `active` - Running and available
- `inactive` - Stopped but metadata exists
- `connected` - User currently attached
- `error` - Error state

Focus on clean interfaces, performance, and user experience consistency.
