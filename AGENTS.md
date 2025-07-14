# AGENTS.md - Claude Pilot Development Guide

## Build/Test Commands

- `make build` - Build the binary
- `make test` - Run all tests
- `make test-race` - Run tests with race detection
- `go test -v ./internal/config` - Run single package tests
- `make fmt` - Format code with gofmt
- `make lint` - Lint code (requires golangci-lint)

## Code Style Guidelines

- **Imports**: Use absolute imports (`claude-pilot/internal/...`), group standard/external/internal
- **Naming**: Use camelCase for variables, PascalCase for exported types/functions
- **Types**: Define interfaces in `internal/interfaces/`, use pointer receivers for methods
- **Error Handling**: Return errors explicitly, wrap with context using `fmt.Errorf`
- **Testing**: Use table-driven tests, descriptive test names, test both success/error cases
- **Logging**: Use structured logging via `internal/logger`, include context fields
- **Comments**: Document exported functions/types, avoid obvious comments

## Architecture Patterns

- Repository pattern for data persistence (`internal/storage/`)
- Interface-based design for multiplexer abstraction
- Dependency injection via constructors (e.g., `NewSessionServiceWithLogger`)
- CommandContext pattern for CLI initialization (`cmd/common.go`)
