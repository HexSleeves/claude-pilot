package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"claude-pilot/internal/config"
	"claude-pilot/internal/interfaces"
	"claude-pilot/internal/multiplexer"
	"claude-pilot/internal/service"
	"claude-pilot/internal/storage"
	"claude-pilot/internal/ui"
)

// CommandContext holds common dependencies for all commands
type CommandContext struct {
	Config      *config.Config
	Service     interfaces.SessionService
	Multiplexer interfaces.TerminalMultiplexer
}

// InitializeCommand handles common initialization for all commands
// This eliminates the duplicated config loading and session manager creation
// that appears in every command file
func InitializeCommand() (*CommandContext, error) {
	// Load configuration using the global cfgFile variable set by --config flag
	configManager := config.NewConfigManager(cfgFile)
	cfg, err := configManager.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	// Create multiplexer instance based on configuration
	mux, err := multiplexer.CreateMultiplexer(cfg.Backend, cfg.Tmux.SessionPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create multiplexer: %w", err)
	}

	// Create repository
	repository, err := storage.NewFileSessionRepository(cfg.SessionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Create service
	sessionService := service.NewSessionService(repository, mux)

	return &CommandContext{
		Config:      cfg,
		Service:     sessionService,
		Multiplexer: mux,
	}, nil
}

// HandleError provides consistent error handling and exit across all commands
// This eliminates the duplicated error handling pattern that appears in every command
func HandleError(err error, action string) {
	fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to %s: %v", action, err)))
	os.Exit(1)
}

// ConfirmAction handles user confirmation prompts consistently
// This eliminates the duplicated confirmation logic in kill commands
func ConfirmAction(message string) bool {
	fmt.Print(ui.Prompt(message))
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error as empty input is valid
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// GetProjectPath handles project path resolution with fallback to current directory
// This eliminates the duplicated project path logic in create command
func GetProjectPath(projectPath string) string {
	if projectPath == "" {
		cwd, err := os.Getwd()
		if err == nil {
			return cwd
		}
		return ""
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err == nil {
		return absPath
	}
	return projectPath
}
