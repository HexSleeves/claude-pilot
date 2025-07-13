# Claude Pilot 🚀

A powerful command-line interface (CLI) tool for managing multiple `claude code` CLI instances, allowing developers to create, view, interact with, and terminate coding sessions through intuitive terminal commands.

## Features

- **Process Management**: Create, list, and terminate multiple `claude code` CLI instances
- **Interactive Mode**: Real-time communication with specific Claude sessions via process piping
- **Terminal UI**: Beautiful, colored terminal output with tabular session lists showing process status
- **Session Persistence**: Session metadata persists across CLI restarts
- **Tmux-inspired Commands**: Familiar command structure for terminal power users
- **Real Claude Integration**: Each session runs an actual `claude code` process

## Installation

```bash
go install github.com/your-username/claude-pilot@latest
```

Or build from source:

```bash
git clone https://github.com/your-username/claude-pilot.git
cd claude-pilot
go build -o claude-pilot
```

## Usage

### Basic Commands

```bash
# Create a new Claude session
claude-pilot create [session-name]

# List all active sessions
claude-pilot list

# Interact with a specific session
claude-pilot attach <session-id>

# Kill a session
claude-pilot kill <session-id>

# Kill all sessions
claude-pilot kill-all
```

### Interactive Mode

```bash
# Attach to a tmux session running Claude
claude-pilot attach my-session

# This opens the tmux session directly - you get the full Claude CLI experience
# Use standard tmux commands:
# Ctrl+B, D    - Detach from session (session keeps running)
# Ctrl+B, ?    - Show tmux help
# exit         - Exit Claude and terminate the session

# Direct communication with Claude (full interactive experience):
Hello Claude, can you help me with this Go project?
```

## Color Scheme

- **Primary**: #FF6B35 (Claude Orange)
- **Success**: #2ECC71 (Green)
- **Error**: #E74C3C (Red)
- **Warning**: #F39C12 (Amber)
- **Text**: #FFFFFF (White)
- **Background**: #1E1E1E (Dark Terminal)

## How It Works

Claude Pilot uses tmux to manage multiple `claude code` CLI sessions:

1. **Session Creation**: Each session creates a new tmux session running `claude code`
2. **Tmux Integration**: Leverages tmux's robust session management capabilities
3. **Interactive Mode**: Attaches directly to the tmux session for full Claude CLI experience
4. **Session Persistence**: Tmux sessions persist independently and can be reattached
5. **Graceful Termination**: Properly terminates tmux sessions when sessions are killed

### Prerequisites

- `tmux` must be installed and available in your PATH
- `claude` CLI must be installed and configured

## Examples

```bash
# Create a new session for a React project
claude-pilot create react-app

# List all sessions with details (shows process status)
claude-pilot list

# Attach to a session (starts claude process if needed)
claude-pilot attach react-app

# Kill a specific session (terminates claude process)
claude-pilot kill react-app
```

## Development

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Build for production
go build -o claude-pilot
```

## License

MIT License - see LICENSE file for details.
