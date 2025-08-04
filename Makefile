# Claude Pilot CLI & TUI Tool Makefile

# Variables
BINARY_NAME         := claude-pilot
TUI_BINARY_NAME     := claude-pilot-tui
BINARY_PATH         := ./$(BINARY_NAME)
TUI_BINARY_PATH     := ./$(TUI_BINARY_NAME)
CLAUDEPILOT_DIR     := packages/claudepilot
TUI_DIR             := packages/tui
GO_CLAUDEPILOT      := $(shell find $(CLAUDEPILOT_DIR) -name "*.go" -type f)
GO_TUI              := $(shell find $(TUI_DIR) -name "*.go" -type f)

# Default target: build both CLI & TUI
.PHONY: all
all: build

# Build both binaries
.PHONY: build
build: build-claudepilot build-tui

# Build claude-pilot CLI
.PHONY: build-claudepilot
build-claudepilot: $(GO_CLAUDEPILOT)
	@echo "Building $(BINARY_NAME)..."
	@cd $(CLAUDEPILOT_DIR) && go build -o ../../$(BINARY_NAME) .

# Build tui
.PHONY: build-tui
build-tui: $(GO_TUI)
	@echo "Building $(TUI_BINARY_NAME)..."
	@cd $(TUI_DIR)/cmd && go build -o ../../../$(TUI_BINARY_NAME) .

# Install dependencies for both
.PHONY: deps
deps: deps-claudepilot deps-tui

.PHONY: deps-claudepilot
deps-claudepilot:
	@echo "Installing dependencies for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go mod tidy && go mod download

.PHONY: deps-tui
deps-tui:
	@echo "Installing dependencies for tui..."
	@cd $(TUI_DIR) && go mod tidy && go mod download

.PHONY: update-deps
update-deps:
	@echo "Updating dependencies for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go get -u ./... && go mod tidy
	@echo "Updating dependencies for tui..."
	@cd $(TUI_DIR) && go get -u ./... && go mod tidy
# Run tests for both
.PHONY: test
test: test-claudepilot test-tui

.PHONY: test-claudepilot
test-claudepilot:
	@echo "Running tests for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go test -v ./...

.PHONY: test-tui
test-tui:
	@echo "Running tests for tui..."
	@cd $(TUI_DIR) && go test -v ./...

# Run golden contract tests
.PHONY: test-golden
test-golden: build-claudepilot
	@echo "Running golden file contract tests..."
	@cd $(CLAUDEPILOT_DIR) && go test -v -run "Test.*Contract" ./...

# Run CLI tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go test -race -v ./...

# Format code for both
.PHONY: fmt
fmt: fmt-claudepilot fmt-tui

.PHONY: fmt-claudepilot
fmt-claudepilot:
	@echo "Formatting code for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && go fmt ./...

.PHONY: fmt-tui
fmt-tui:
	@echo "Formatting code for tui..."
	@cd $(TUI_DIR) && go fmt ./...

# Lint code for both
.PHONY: lint
lint: lint-claudepilot lint-tui

.PHONY: lint-claudepilot
lint-claudepilot:
	@echo "Linting code for claudepilot..."
	@cd $(CLAUDEPILOT_DIR) && golangci-lint run

.PHONY: lint-tui
lint-tui:
	@echo "Linting code for tui..."
	@cd $(TUI_DIR) && golangci-lint run

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(TUI_BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -rf ~/.claude-pilot

# Install binaries to GOPATH/bin
.PHONY: install
install: build install-cli install-tui

.PHONY: install-cli
install-cli:
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

.PHONY: install-tui
install-tui:
	@echo "Installing $(TUI_BINARY_NAME) to $(GOPATH)/bin..."
	@cp $(TUI_BINARY_NAME) $(GOPATH)/bin/$(TUI_BINARY_NAME)

# Development builds
.PHONY: dev
dev: dev-claudepilot dev-tui

.PHONY: dev-claudepilot
dev-claudepilot:
	@echo "Building development version of $(BINARY_NAME)..."
	@cd $(CLAUDEPILOT_DIR) && go build -gcflags="all=-N -l" -o ../../$(BINARY_NAME) .

.PHONY: dev-tui
dev-tui:
	@echo "Building development version of $(TUI_BINARY_NAME)..."
	@cd $(TUI_DIR) && go build -gcflags="all=-N -l" -o ../../$(TUI_BINARY_NAME) .

# Cross-compile for different platforms (CLI only)
.PHONY: build-all
build-all:
	@echo "Cross-compiling $(BINARY_NAME) for multiple platforms..."
	@cd $(CLAUDEPILOT_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../$(BINARY_NAME)-linux-amd64 .
	@cd $(CLAUDEPILOT_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../../$(BINARY_NAME)-darwin-amd64 .
	@cd $(CLAUDEPILOT_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../../$(BINARY_NAME)-darwin-arm64 .
	@cd $(CLAUDEPILOT_DIR) && GOOS=windows GOARCH=amd64 go build -o ../../$(BINARY_NAME)-windows-amd64.exe .

# Run targets
.PHONY: run
run: run-cli

.PHONY: run-cli
run-cli: build-claudepilot
	@$(BINARY_PATH)

.PHONY: run-tui
run-tui: build-tui
	@cd $(TUI_DIR) && go run ./cmd/main.go

# Help message
.PHONY: help
help:
	@echo "Claude Pilot CLI & TUI Tool - Available Make targets:"
	@echo ""
	@echo "  all               - Build both claude-pilot & tui"
	@echo "  build             - Build both binaries"
	@echo "  deps              - Install dependencies for all packages"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests for both packages"
	@echo "  test-golden       - Run golden file contract tests"
	@echo "  test-race         - Run CLI tests with race detection"
	@echo "  fmt               - Format code for both packages"
	@echo "  lint              - Lint code for both packages"
	@echo "  install           - Install both binaries to \$$GOPATH/bin"
	@echo "  dev               - Development build for both"
	@echo "  build-all         - Cross-compile CLI for multiple platforms"
	@echo "  run               - Build & run CLI"
	@echo "  run-tui           - Build & run TUI"
	@echo ""
	@echo "Package-specific targets:"
	@echo "  build-claudepilot - Build only the CLI"
	@echo "  test-claudepilot  - Test only the CLI"
	@echo "  deps-claudepilot  - Install CLI deps"
	@echo "  fmt-claudepilot   - Format CLI code"
	@echo "  lint-claudepilot  - Lint CLI code"
	@echo ""
	@echo "  build-tui         - Build only the TUI"
	@echo "  test-tui          - Test only the TUI"
	@echo "  deps-tui          - Install TUI deps"
	@echo "  fmt-tui           - Format TUI code"
	@echo "  lint-tui          - Lint TUI code"
	@echo "  run-tui           - Build & run the TUI"
