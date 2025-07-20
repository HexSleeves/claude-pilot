# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Project Name:** [Your Project Name]
**Technology Stack:** [List main technologies - e.g., React, Node.js, PostgreSQL]
**Architecture:** [Brief description - e.g., Microservices, Monolithic, JAMstack]

## Core Development Principles

### 1. Code Quality Standards

- **Immutable Rule:** All code must pass linting and type checking before commit
- **Testing Required:** Minimum 80% test coverage for new features
- **Documentation:** All public APIs must be documented
- **Security:** Never commit secrets, always use environment variables

### 2. Modular Command Structure

This project uses the Claude Code modular command system for consistent workflows:

- Commands are organized in `.claude/commands/` by category
- Each command follows XML-structured format for clarity
- Use `/[category]:[command]` syntax for execution
- Commands are environment-aware and security-focused

### 3. Emergency Procedures

- **Build Failures:** Run `/dev:debug-session` for systematic troubleshooting
- **Test Failures:** Use `/test:coverage-analysis` to identify issues
- **Deployment Issues:** Execute `/deploy:rollback-procedure` for emergency rollback
- **Security Concerns:** Immediately run security scans and notify team

## Command Categories

### Project Management

- `/project:create-feature` - Full feature development with tests and docs
- `/project:scaffold-component` - Component creation with boilerplate
- `/project:setup-environment` - Development environment initialization

### Development Workflow

- `/dev:code-review` - Structured code review with quality checks
- `/dev:refactor-analysis` - Code improvement recommendations
- `/dev:debug-session` - Systematic debugging and problem solving

### Testing

- `/test:generate-tests` - Comprehensive test suite generation
- `/test:coverage-analysis` - Test coverage assessment and improvement
- `/test:integration-tests` - Integration test creation and execution

### Deployment

- `/deploy:prepare-release` - Release preparation with quality gates
- `/deploy:deploy-staging` - Staging deployment with validation
- `/deploy:rollback-procedure` - Emergency rollback execution

### Documentation

- `/docs:api-docs` - API documentation generation
- `/docs:update-readme` - README maintenance and updates
- `/docs:architecture-review` - Architecture documentation and review

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

## Quality Gates

### Pre-commit Requirements

- [ ] All tests pass
- [ ] Linting passes
- [ ] Type checking passes
- [ ] Security scan clean
- [ ] Code coverage maintained

## Token Management

### Context Optimization

- Use progressive disclosure for large codebases
- Load commands just-in-time based on current task
- Clear context when switching between major tasks
- Monitor token usage and optimize accordingly

### Memory Management

- Leverage MCP memory server for session continuity
- Store architectural decisions in external documentation
- Use modular instructions to reduce token overhead
- Implement context compression for repeated patterns

### Debug Commands

- Use `/dev:debug-session` for systematic debugging
- Check logs: `npm run logs` or `yarn logs`
- Monitor performance: `npm run monitor` or `yarn monitor`
