# Claude Pilot 🚀

**Claude Pilot is a command-line interface (CLI) tool designed to manage multiple Claude Code CLI instances simultaneously, enabling developers to create, manage, and interact with multiple AI-powered coding sessions through an intuitive terminal interface.**

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/your-username/claude-pilot)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/your-username/claude-pilot/pulls)

---

## Table of Contents

- [Claude Pilot 🚀](#claude-pilot-)
  - [Table of Contents](#table-of-contents)
  - [About The Project](#about-the-project)
  - [Key Features](#key-features)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
    - [Core Commands](#core-commands)
  - [Architecture](#architecture)
  - [Roadmap](#roadmap)
  - [Contributing](#contributing)
  - [License](#license)
  - [Contributors](#contributors)
  - [Cursor Rules](#cursor-rules)

---

## About The Project

Developers often work on multiple projects or tasks simultaneously, but Claude Code CLI instances are typically single-session and don't persist across terminal sessions. This makes it difficult to manage multiple conversations for different contexts and requires manual session management when switching between tasks.

**Claude Pilot** solves this by providing a session management system that allows developers to:

- Create and manage multiple named Claude Code sessions.
- Persist sessions across terminal restarts.
- Switch between different coding contexts seamlessly.
- Maintain session history and metadata.
- Organize work by project or task.

It's built with Go and leverages terminal multiplexers like `tmux` and `zellij` to provide a robust and familiar experience for power users.

---

## Key Features

- **Multi-Session Management**: Create, list, and terminate multiple `claude code` CLI instances.
- **Session Persistence**: Session metadata is stored on your local machine and persists across application restarts.
- **Terminal Multiplexer Support**: Choose between `tmux` and `zellij` for session management.
- **Interactive Mode**: Attach to any session to interact directly with the Claude CLI.
- **Detailed Session Information**: The `list` command shows session ID, name, status, creation time, and more.
- **Named Sessions**: Give your sessions meaningful names to easily organize your work (e.g., `react-app`, `api-bug-fix`).
- **Intuitive CLI**: A clean, `tmux`-inspired command structure that is easy to learn and use.
- **Colored UI**: A beautiful, colored terminal output with tabular session lists for clarity.

---

## Getting Started

### Prerequisites

You must have the following tools installed and available in your system's PATH:

1. **Claude CLI**: The `claude` command-line tool.
2. **Terminal Multiplexer**:
    - `tmux` (recommended)
    - or `zellij`

### Installation

You can install Claude Pilot using `go install`:

```bash
# Replace with the actual repository path when available
go install github.com/your-username/claude-pilot@latest
```

Alternatively, you can build it from the source:

```bash
# Replace with the actual repository URL
git clone https://github.com/your-username/claude-pilot.git
cd claude-pilot
go build -o claude-pilot
sudo mv claude-pilot /usr/local/bin/
```

---

## Usage

Claude Pilot provides a simple and powerful set of commands to manage your sessions.

### Core Commands

**`create [session-name]`**
Creates a new Claude session. If no name is provided, a random one will be generated.

```bash
claude-pilot create my-go-project
```

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
Terminates a specific session.

```bash
claude-pilot kill my-go-project
```

**`kill-all`**
Terminates all running sessions.

```bash
claude-pilot kill-all
```

---

## Architecture

Claude Pilot is built with a clean, modular architecture in Go. It uses the Cobra framework for the CLI and interfaces with terminal multiplexers for session management.

```**mermaid**
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Interface │    │  Session Manager │    │  Multiplexer    │
│   (Cobra)       │────│                  │────│  (Tmux/Zellij)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                │
                       ┌──────────────────┐
                       │  JSON Storage    │
                       │  ~/.claude-pilot │
                       └──────────────────┘
```

- **CLI Layer (`cmd/`)**: Defines the command structure and handles user input.
- **Session Management (`internal/service/`, `internal/manager/`)**: Core logic for CRUD operations on sessions.
- **Multiplexer (`internal/multiplexer/`)**: An interface that communicates with `tmux` or `zellij` to create, attach to, and kill sessions.
- **Storage (`internal/storage/`)**: Handles saving and retrieving session metadata from JSON files.

---

## Roadmap

We have an exciting roadmap ahead! Here are some of the features we're planning to implement:

- **Phase 1: Core MVP (Complete)**
  - [x] Basic CLI structure with Cobra
  - [x] Session data model and storage
  - [x] Tmux integration layer
  - [x] `create`, `list`, `attach`, `kill` commands
  - [x] Enhanced UI with colors and tables
  - [x] Session metadata persistence

- **Phase 2: Advanced Features**
  - [ ] `kill-all` command implementation
  - [ ] Session templates and presets
  - [ ] Session search and filtering
  - [ ] Export/import session configurations

- **Phase 3: Integrations & Release**
  - [ ] Comprehensive testing suite
  - [ ] Cross-platform compatibility testing (Windows, Linux, macOS)
  - [ ] Official package distribution (Homebrew, etc.)
  - [ ] Detailed documentation and tutorials

See the [open issues](https://github.com/your-username/claude-pilot/issues) for a full list of proposed features (and known issues).

---

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m '''Add some AmazingFeature'''`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## License

Distributed under the MIT License. See `LICENSE` file for more information.

---

## Contributors

A huge thanks to all the people who have contributed to this project.

<!-- Add contributors here -->
<a href="https://github.com/HexSleeves/claude-pilot/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=HexSleeves/claude-pilot" />
</a>

## Cursor Rules

Claude Pilot includes a set of project rules under `.cursor/rules` that guide AI agents:

| Rule File | Purpose |
|-----------|---------|
| `ReviewGateV2.mdc` | Enforces Review Gate interactive feedback flow |
| `cursor_rules.mdc` | General guidelines for writing effective rules |
| `self_improve.mdc` | Continuous improvement heuristics |
| `debug.mdc` | Systematic debugging workflow |
| `decompose.mdc` | Breaks PRDs into granular tasks |
| `prd.mdc` | Generates Product Requirement Documents from ideas |
| `task.mdc` | Two-phase task planning & execution protocol |
| `golang_cli_tui.mdc` | Best practices for Go CLI/TUI apps with Cobra & Bubbletea |

Refer to these rules when interacting with Cursor’s AI agents to get consistent, high-quality assistance.
