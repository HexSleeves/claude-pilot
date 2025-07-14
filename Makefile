# Claude Pilot CLI Tool Makefile

# Variables
BINARY_NAME=claude-pilot
BINARY_PATH=./$(BINARY_NAME)
CLAUDEPILOT_DIR=packages/claudepilot
TUI_DIR=packages/tui
GO_FILES=$(shell find $(CLAUDEPILOT_DIR) -name "*.go" -type f)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	@cd $(CLAUDEPILOT_DIR) && go build -o ../../$(BINARY_NAME) .

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go mod tidy && go mod download
	@echo "Installing dependencies for tui..."
	@cd $(TUI_DIR) && go mod tidy && go mod download

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -rf ~/.claude-pilot

# Run tests
.PHONY: test
test:
	@echo "Running tests for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go test -v ./...

# Run with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go test -race -v ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && golangci-lint run

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

# Development build with debug info
.PHONY: dev
dev:
	@echo "Building development version..."
	@cd $(CLAUDEPILOT_DIR) && go build -gcflags="all=-N -l" -o ../../$(BINARY_NAME) .

# Cross-compile for different platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@cd $(CLAUDEPILOT_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../$(BINARY_NAME)-linux-amd64 .
	@cd $(CLAUDEPILOT_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../../$(BINARY_NAME)-darwin-amd64 .
	@cd $(CLAUDEPILOT_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../../$(BINARY_NAME)-darwin-arm64 .
	@cd $(CLAUDEPILOT_DIR) && GOOS=windows GOARCH=amd64 go build -o ../../$(BINARY_NAME)-windows-amd64.exe .

# Run the application
.PHONY: run
run: build
	@$(BINARY_PATH)

# Package-specific targets
.PHONY: build-claudepilot
build-claudepilot:
	@echo "Building claudepilot package..."
	@cd $(CLAUDEPILOT_DIR) && go build -o ../../$(BINARY_NAME) .

.PHONY: test-claudepilot
test-claudepilot:
	@echo "Running tests for claudepilot package..."
	@cd $(CLAUDEPILOT_DIR) && go test -v ./...

.PHONY: deps-claudepilot
deps-claudepilot:
	@echo "Installing dependencies for claudepilot package..."
	@cd $(CLAUDEPILOT_DIR) && go mod tidy && go mod download

.PHONY: fmt-claudepilot
fmt-claudepilot:
	@echo "Formatting claudepilot package..."
	@cd $(CLAUDEPILOT_DIR) && go fmt ./...

.PHONY: lint-claudepilot
lint-claudepilot:
	@echo "Linting claudepilot package..."
	@cd $(CLAUDEPILOT_DIR) && golangci-lint run

.PHONY: deps-tui
deps-tui:
	@echo "Installing dependencies for tui package..."
	@cd $(TUI_DIR) && go mod tidy && go mod download

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
	@echo "Main targets:"
	@echo "  build      - Build the binary from claudepilot package"
	@echo "  deps       - Install dependencies for all packages"
	@echo "  clean      - Clean build artifacts and session data"
	@echo "  test       - Run tests for claudepilot package"
	@echo "  test-race  - Run tests with race detection for claudepilot package"
	@echo "  fmt        - Format code for claudepilot package"
	@echo "  lint       - Lint code for claudepilot package (requires golangci-lint)"
	@echo "  install    - Install binary to GOPATH/bin"
	@echo "  dev        - Build development version with debug info"
	@echo "  build-all  - Cross-compile for multiple platforms"
	@echo "  run        - Build and run the application"
	@echo "  demo       - Run a complete demo of the CLI features"
	@echo ""
	@echo "Package-specific targets:"
	@echo "  build-claudepilot - Build only the claudepilot package"
	@echo "  test-claudepilot  - Run tests only for claudepilot package"
	@echo "  deps-claudepilot  - Install dependencies only for claudepilot package"
	@echo "  fmt-claudepilot   - Format code only for claudepilot package"
	@echo "  lint-claudepilot  - Lint code only for claudepilot package"
	@echo "  deps-tui          - Install dependencies only for tui package"
	@echo ""
	@echo "  help       - Show this help message"
