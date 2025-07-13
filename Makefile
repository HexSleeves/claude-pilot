# Claude Pilot CLI Tool Makefile

# Variables
BINARY_NAME=claude-pilot
BINARY_PATH=./$(BINARY_NAME)
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf ~/.claude-pilot

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

# Development build with debug info
.PHONY: dev
dev:
	@echo "Building development version..."
	@go build -gcflags="all=-N -l" -o $(BINARY_NAME) .

# Cross-compile for different platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

# Run the application
.PHONY: run
run: build
	@$(BINARY_PATH)

# Demo commands
.PHONY: demo
demo: build
	@echo "Running Claude Pilot demo..."
	@echo "1. Creating test sessions..."
	@$(BINARY_PATH) create demo-session -d "Demo session for testing"
	@$(BINARY_PATH) create react-demo -d "React development demo"
	@echo "\n2. Listing sessions..."
	@$(BINARY_PATH) list
	@echo "\n3. Showing session details..."
	@$(BINARY_PATH) list --all
	@echo "\n4. Cleaning up demo sessions..."
	@echo "y" | $(BINARY_PATH) kill-all

# Help target
.PHONY: help
help:
	@echo "Claude Pilot CLI Tool - Available Make targets:"
	@echo ""
	@echo "  build      - Build the binary"
	@echo "  deps       - Install dependencies"
	@echo "  clean      - Clean build artifacts and session data"
	@echo "  test       - Run tests"
	@echo "  test-race  - Run tests with race detection"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code (requires golangci-lint)"
	@echo "  install    - Install binary to GOPATH/bin"
	@echo "  dev        - Build development version with debug info"
	@echo "  build-all  - Cross-compile for multiple platforms"
	@echo "  run        - Build and run the application"
	@echo "  demo       - Run a complete demo of the CLI features"
	@echo "  help       - Show this help message"
