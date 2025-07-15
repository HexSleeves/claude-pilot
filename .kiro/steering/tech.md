# Claude Pilot Technical Stack

## Language & Runtime
- **Go 1.24.5+** - Primary language for all components
- Cross-platform support (macOS, Linux, Windows)

## Architecture
Monorepo structure with three main packages:
- **`packages/claudepilot`** - Primary CLI interface using Cobra
- **`packages/tui`** - Terminal User Interface using Bubble Tea  
- **`packages/core`** - Shared business logic and services
- **`packages/shared`** - Common components and utilities

## Key Dependencies

### CLI Framework
- **Cobra** (`github.com/spf13/cobra`) - Command-line interface
- **Viper** (`github.com/spf13/viper`) - Configuration management

### TUI Framework  
- **Bubble Tea** (`github.com/charmbracelet/bubbletea`) - Terminal UI framework
- **Lipgloss** (`github.com/charmbracelet/lipgloss`) - Styling and layout

### Terminal Integration
- **tmux** - Primary terminal multiplexer (preferred)
- **zellij** - Alternative terminal multiplexer
- **fatih/color** - Terminal colors for CLI output

### Storage & Utilities
- **JSON files** - Session metadata storage in `~/.config/claude-pilot/sessions`
- **UUID** (`github.com/google/uuid`) - Unique session identifiers

## Build System

### Makefile Commands
```bash
# Build both CLI and TUI
make build

# Build individual components
make build-claudepilot    # CLI only
make build-tui           # TUI only

# Development
make deps                # Install dependencies
make test                # Run tests
make fmt                 # Format code
make lint                # Lint code (requires golangci-lint)

# Installation
make install             # Install to $GOPATH/bin
make clean              # Clean build artifacts

# Cross-compilation
make build-all          # Build for multiple platforms
```

### Go Module Structure
Each package has its own `go.mod` with local replacements:
```go
replace claude-pilot/core => ../core
replace claude-pilot/shared => ../shared
```

## Configuration
- **YAML-based** configuration in `~/.config/claude-pilot/claude-pilot.yaml`
- Environment variables with `CLAUDE_PILOT_` prefix
- Auto-detection of available terminal multiplexers
- Logging disabled by default (can be enabled with `--verbose`)

## Development Workflow
1. Use `make deps` to install dependencies
2. Use `make build` for development builds
3. Use `make test` before committing
4. Use `make fmt` and `make lint` for code quality
5. Cross-platform testing with `make build-all`

## Testing
- Unit tests in `*_test.go` files
- Race condition testing with `make test-race`
- Cross-platform compatibility testing required