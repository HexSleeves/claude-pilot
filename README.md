# Claude Pilot 🚀

**Claude Pilot is a command-line interface (CLI) tool designed to manage multiple Claude Code CLI instances simultaneously, enabling developers to create, manage, and interact with multiple AI-powered coding sessions through an intuitive terminal interface and a full-featured TUI.**

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-blue.svg)](https://golang.org/dl/)
[![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/HexSleeves/claude-pilot)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/HexSleeves/claude-pilot/pulls)

-----

## Table of Contents

- [Claude Pilot 🚀](https://www.google.com/search?q=%23claude-pilot-)
  - [Table of Contents](https://www.google.com/search?q=%23table-of-contents)
  - [About The Project](https://www.google.com/search?q=%23about-the-project)
  - [Key Features](https://www.google.com/search?q=%23key-features)
  - [Demo](https://www.google.com/search?q=%23demo)
  - [Getting Started](https://www.google.com/search?q=%23getting-started)
    - [Prerequisites](https://www.google.com/search?q=%23prerequisites)
    - [Installation](https://www.google.com/search?q=%23installation)
  - [Usage](https://www.google.com/search?q=%23usage)
    - [Core Commands](https://www.google.com/search?q=%23core-commands)
  - [Architecture](https://www.google.com/search?q=%23architecture)
  - [Roadmap](https://www.google.com/search?q=%23roadmap)
  - [Contributing](https://www.google.com/search?q=%23contributing)
  - [License](https://www.google.com/search?q=%23license)
  - [Contributors](https://www.google.com/search?q=%23contributors)
  - [Cursor Rules](https://www.google.com/search?q=%23cursor-rules)

-----

## About The Project

Developers often work on multiple projects or tasks simultaneously, but Claude Code CLI instances are typically single-session and don't persist across terminal sessions. This makes it difficult to manage multiple conversations for different contexts and requires manual session management when switching between tasks.

**Claude Pilot** solves this by providing a session management system that allows developers to:

- Create and manage multiple named Claude Code sessions.
- Persist sessions across terminal restarts.
- Switch between different coding contexts seamlessly.
- Maintain session history and metadata.
- Organize work by project or task.

It's built with Go and leverages terminal multiplexers like `tmux` to provide a robust and familiar experience for power users. Support for `zellij` is planned for future releases.

-----

## Key Features

- **Multi-Session Management**: Create, list, and terminate multiple `claude code` CLI instances.
- **Interactive TUI**: Launch a full-featured Terminal User Interface with `claude-pilot tui` for interactive, mouse-supported session management.
- **Advanced Session Creation**: Attach new sessions as panes or windows/tabs to existing sessions, with control over split direction.
- **Session Persistence**: Session metadata is stored on your local machine and persists across application restarts.
- **Terminal Multiplexer Support**: Uses `tmux` for robust session management.
- **Detailed Session Information**: The `list` command shows session ID, name, status, creation time, panes, and more.
- **Unified Theming**: A beautiful, consistent Claude Orange theme is shared across both the CLI and TUI, built with `lipgloss`.
- **Named Sessions**: Give your sessions meaningful names to easily organize your work (e.g., `react-app`, `api-bug-fix`).
- **Intuitive CLI**: A clean, `tmux`-inspired command structure that is easy to learn and use.

-----

## Demo

Watch this short demo to see Claude Pilot in action:

[![Watch the video](assets/claude-pilot-demo.gif)](assets/claude-pilot-demo.gif)

-----

## Getting Started

### Prerequisites

You must have the following tools installed and available in your system's PATH:

1. **Claude CLI**: The `claude` command-line tool.
2. **Terminal Multiplexer**:
      - `tmux` (required)

### Installation

You can install Claude Pilot using `go install`:

```bash
go install github.com/HexSleeves/claude-pilot/packages/claudepilot@latest
```

Alternatively, you can build it from the source using the provided `Makefile`:

```bash
git clone https://github.com/HexSleeves/claude-pilot.git
cd claude-pilot
make build
sudo cp claude-pilot /usr/local/bin/
```

-----

## Usage

Claude Pilot provides a simple and powerful set of commands to manage your sessions.

### Core Commands

**`tui`**
Launches the interactive Terminal User Interface. This is the recommended way to use Claude Pilot for a visual experience.

```bash
claude-pilot tui
```

**`create [session-name]`**
Creates a new Claude session. If no name is provided, a timestamp-based name will be generated.

You can also attach to existing sessions as new panes or windows/tabs:

```bash
# Create a standalone session
claude-pilot create my-go-project

# Create session with description and project path
claude-pilot create --desc "React app" --project ./src

# Attach to existing session as a new pane (default: vertical split)
claude-pilot create debug --attach-to my-go-project --as-pane

# Attach as new pane with horizontal split
claude-pilot create debug --attach-to my-go-project --as-pane --split h

# Attach as new window/tab
claude-pilot create testing --attach-to my-go-project --as-window
```

**Attachment Options:**

- `--attach-to <session-name>`: Target session to attach to
- `--as-pane`: Create as new pane in existing session
- `--as-window`: Create as new window/tab in existing session
- `--split <direction>`: Split direction for panes (`h` for horizontal, `v` for vertical)

**`list`**
Lists all active and inactive sessions in a clean, tabular format.

```bash
claude-pilot list
```

**`attach <session-id|session-name>`**
Attaches to an existing session, allowing you to interact with Claude.

```bash
claude-pilot attach my-go-project
```

Inside the session, you can use standard multiplexer commands to detach (e.g., `Ctrl+B, D` for tmux).

**`kill <session-id|session-name>`**
Terminates a specific session. Use the `--all` flag to kill all sessions.

```bash
# Kill a specific session
claude-pilot kill my-go-project

# Kill all sessions with confirmation
claude-pilot kill --all
```

**`details <session-id|session-name>`**
Shows detailed information for a specific session.

```bash
claude-pilot details my-go-project
```

-----

## Architecture

Claude Pilot is designed with a modular, monorepo architecture to separate concerns and promote code reuse. The project is divided into four main packages: `claudepilot`, `core`, `shared`, and `tui`.

- **`packages/core`**: The heart of the application. This package contains all the core business logic, including:

  - **Service (`service/`)**: Manages session lifecycle (create, read, update, delete).
  - **Multiplexer (`multiplexer/`)**: An interface to communicate with `tmux`.
  - **Storage (`storage/`)**: Handles saving and retrieving session metadata from the filesystem as JSON.
  - **Configuration (`config/`)**: Manages application configuration via Viper.
  - **API (`api/`)**: A clean client-facing API that abstracts the core logic for consumers.

- **`packages/shared`**: Contains code shared between the CLI and TUI frontends.

  - **Interfaces (`interfaces/`)**: Defines the core data structures and service contracts.
  - **Styles (`styles/`)**: A unified, `lipgloss`-based styling system and theme (Claude Orange).
  - **Components (`components/`)**: Reusable UI components, like the session table, for both CLI and TUI.

- **`packages/claudepilot`**: The primary command-line interface.

  - Built with **Cobra**, it provides the main entry point for users to interact with the application (e.g., `claude-pilot create`, `list`, `attach`, `kill`).
  - It consumes the `core` package via its API layer to perform its operations.

- **`packages/tui`**: A rich, interactive Terminal User Interface (TUI).

  - Built with **Bubble Tea**, it offers a more visual and interactive way to manage sessions.
  - Like the `claudepilot` CLI, it is also a consumer of the `core` package's API.

This decoupled architecture allows for flexible development and makes it easier to maintain and extend the application's capabilities.

-----

## Roadmap

We have an exciting roadmap ahead\! Here are some of the features we're planning to implement:

- **Phase 1: Core MVP & TUI (Complete)**

  - [x] Basic CLI structure with Cobra
  - [x] Session data model and storage
  - [x] Tmux integration layer
  - [x] `create`, `list`, `attach`, `kill`, `details` commands
  - [x] Advanced `create` (add panes/windows to existing sessions)
  - [x] `kill --all` command
  - [x] Interactive TUI with `bubbletea`
  - [x] Shared styling and components architecture

- **Phase 2: Advanced Features & Polish**

  - [ ] Zellij backend support
  - [ ] Session templates and presets
  - [ ] Enhanced session filtering and searching in the TUI
  - [ ] Export/import session configurations

- **Phase 3: Integrations & Release**

  - [ ] Cross-platform compatibility testing (Windows, Linux, macOS)
  - [ ] Official package distribution (Homebrew, etc.)
  - [ ] Detailed documentation and tutorials

See the [open issues](https://github.com/HexSleeves/claude-pilot/issues) for a full list of proposed features (and known issues).

-----

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

-----

## License

Distributed under the MIT License. See `LICENSE` file for more information.

-----

## Contributors

A huge thanks to all the people who have contributed to this project.

<!-- Add contributors here -->
<a href="https://github.com/HexSleeves/claude-pilot/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=HexSleeves/claude-pilot" />
</a>

## Cursor Rules

Claude Pilot includes a set of project rules under `.cursor/rules` that guide AI agents:

| Rule File | Purpose |
| :--- | :--- |
| `ReviewGateV3.mdc` | Enforces Review Gate interactive feedback flow. |
| `cursor_rules.mdc` | General guidelines for writing effective rules. |
| `self_improve.mdc` | Continuous improvement heuristics. |
| `debug.mdc` | Systematic debugging workflow. |
| `decompose.mdc` | Breaks PRDs into granular tasks. |
| `prd.mdc` | Generates Product Requirement Documents from ideas. |
| `task.mdc` | Two-phase task planning & execution protocol. |
| `golang_cli_tui.mdc` | Best practices for Go CLI/TUI apps with Cobra & Bubbletea. |

Refer to these rules when interacting with Cursor’s AI agents to get consistent, high-quality assistance.
