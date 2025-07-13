# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Claude Pilot is a modern CLI tool for managing multiple Claude Code sessions with support for multiple terminal multiplexers (tmux and zellij). It provides session persistence, organization, and a clean interface for AI-assisted development workflows, with built-in support for future TUI (Terminal User Interface) mode.

## Architecture

The codebase follows a clean, layered architecture with proper separation of concerns:

### Core Structure
- **`cmd/`** - CLI commands using Cobra framework (create, list, attach, kill operations)
- **`internal/interfaces/`** - Core interface definitions for abstraction
- **`internal/manager/`** - High-level session orchestration
- **`internal/service/`** - Business logic layer
- **`internal/storage/`** - Data persistence layer
- **`internal/multiplexer/`** - Terminal multiplexer implementations (tmux, zellij)
- **`internal/config/`** - Configuration management
- **`internal/ui/`** - Terminal UI components (colors, table formatting, renderer abstraction)
- **`internal/tui/`** - Bubble Tea TUI foundation (ready for future use)
- **`main.go`** - Entry point that delegates to cmd package

### Key Components
- **SessionManager** (`internal/manager/session_manager.go`): High-level orchestration with dependency injection
- **SessionService** (`internal/service/session_service.go`): Business logic for session operations
- **FileSessionRepository** (`internal/storage/file_repository.go`): JSON-based session persistence
- **TerminalMultiplexer Interface** (`internal/interfaces/multiplexer.go`): Abstraction for tmux/zellij
- **TmuxMultiplexer** (`internal/multiplexer/tmux.go`): Modern tmux implementation
- **ZellijMultiplexer** (`internal/multiplexer/zellij.go`): Complete zellij support
- **MultiplexerFactory** (`internal/multiplexer/factory.go`): Auto-detection and backend creation

## Development Commands

```bash
# Build the application
make build

# Run tests
make test

# Run tests with race detection
make test-race

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Install dependencies
make deps

# Clean build artifacts and session data
make clean

# Development build with debug info
make dev

# Cross-compile for multiple platforms
make build-all

# Run the application
make run

# Run demo workflow
make demo
```

## Key Dependencies

- **Cobra + Viper**: CLI framework and configuration management
- **tmux/zellij**: Terminal multiplexer backends (auto-detected, tmux preferred)
- **go-pretty**: Terminal table formatting
- **fatih/color**: Terminal colors using Claude orange theme (#FF6B35)
- **google/uuid**: Session ID generation
- **Bubble Tea + Bubbles + Lip Gloss**: TUI framework for future terminal interface

## Configuration

Configuration is managed via YAML files and environment variables:

- **Default config location**: `~/.config/claude-pilot/.claude-pilot.yaml`
- **Sessions storage**: Configurable, defaults to `~/.config/claude-pilot/sessions/`
- **Backend selection**: `auto` (default), `tmux`, or `zellij`
- **Environment variables**: Prefixed with `CLAUDE_PILOT_`

Example configuration:
```yaml
backend: auto
sessions_dir: ~/.config/claude-pilot/sessions
default_shell: claude
ui:
  mode: cli
  theme: default
tmux:
  session_prefix: claude-
zellij:
  layout_file: ""
```

## Session Management

Sessions are stored as JSON files with the following states:

- `active`: Session running and available
- `inactive`: Metadata exists but multiplexer session stopped
- `connected`: User currently attached
- `error`: Session in error state

Each session creates a corresponding multiplexer session running the `claude` CLI in the specified project directory.

## Multi-Backend Support

The application supports multiple terminal multiplexers:

- **tmux**: Full-featured, battle-tested, extensive scripting capabilities
- **zellij**: Modern, user-friendly, plugin system
- **auto**: Automatically detects and uses the best available backend (prefers tmux)

Backend selection is handled by the factory pattern with automatic detection and fallback.

## Important Implementation Details

- Session lookup supports both ID and name-based identification
- Multiplexer session names use the session name, not the UUID
- Session persistence handles graceful recovery from multiplexer session failures
- The CLI requires either `tmux` or `zellij` and `claude` to be available in PATH
- Session metadata includes message history for potential future features
- Interface-based design enables easy extension and testing
- TUI foundation is ready for activation via configuration

## Future Features

- **TUI Mode**: Interactive terminal interface using Bubble Tea (foundation complete)
- **Additional Backends**: Screen, Docker, SSH session support
- **Plugin System**: Extensible architecture ready for plugins
- **Web Interface**: Service layer ready for REST API development
- **Database Storage**: Repository pattern supports alternative storage backends