# Claude Pilot Product Overview

Claude Pilot is a command-line interface (CLI) tool designed to manage multiple Claude Code CLI instances simultaneously. It enables developers to create, manage, and interact with multiple AI-powered coding sessions through an intuitive terminal interface.

## Core Problem
- Developers work on multiple projects simultaneously but Claude Code CLI instances are single-session and don't persist across terminal sessions
- No easy way to manage multiple Claude conversations for different contexts
- Manual session management required when switching between tasks

## Solution
Claude Pilot provides a tmux/zellij-based session management system that allows developers to:
- Create and manage multiple named Claude Code sessions
- Persist sessions across terminal restarts  
- Switch between different coding contexts seamlessly
- Maintain session history and metadata
- Organize work by project or task

## Target Users
- **Primary**: Software developers using Claude Code for AI-assisted development
- **Secondary**: DevOps engineers, technical leads, development teams
- **Tertiary**: Students and educators using AI coding assistants

## Key Features
- Multi-session management with named sessions
- Session persistence across application restarts
- Terminal multiplexer support (tmux/zellij)
- Interactive CLI and TUI interfaces
- Detailed session information and metadata
- Colored terminal output with tabular listings