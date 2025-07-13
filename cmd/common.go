package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"claude-pilot/internal/config"
	"claude-pilot/internal/interfaces"
	"claude-pilot/internal/logger"
	"claude-pilot/internal/multiplexer"
	"claude-pilot/internal/service"
	"claude-pilot/internal/storage"
	"claude-pilot/internal/ui"

	"github.com/spf13/viper"
)

// CommandContext holds common dependencies for all commands
type CommandContext struct {
	Config      *config.Config
	Logger      *logger.Logger
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

	// Get verbose flag from viper
	verbose := viper.GetBool("verbose")

	// Create logger based on configuration and flags
	loggerBuilder := logger.NewBuilder().
		WithEnabled(cfg.Logging.Enabled || verbose). // Enable if configured OR verbose flag is set
		WithLevel(cfg.Logging.Level).
		WithFile(cfg.Logging.File).
		WithMaxSize(cfg.Logging.MaxSize).
		WithTUIMode(cfg.UI.Mode == "tui").
		WithVerbose(verbose)

	log, err := loggerBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Create multiplexer instance based on configuration
	mux, err := multiplexer.CreateMultiplexer(cfg.Backend, cfg.Tmux.SessionPrefix)
	if err != nil {
		log.Error("Failed to create multiplexer",
			"backend", cfg.Backend,
			"prefix", cfg.Tmux.SessionPrefix,
			"error", err)
		return nil, fmt.Errorf("failed to create multiplexer: %w", err)
	}

	// Create repository
	repository, err := storage.NewFileSessionRepository(cfg.SessionsDir)
	if err != nil {
		log.Error("Failed to create repository",
			"sessions_dir", cfg.SessionsDir,
			"error", err)
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Create service with logger
	sessionService := service.NewSessionServiceWithLogger(repository, mux, log)

	log.Info("Command context initialized successfully",
		"backend", cfg.Backend,
		"sessions_dir", cfg.SessionsDir,
		"ui_mode", cfg.UI.Mode,
		"logging_enabled", cfg.Logging.Enabled,
		"verbose", verbose)

	return &CommandContext{
		Config:      cfg,
		Logger:      log,
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
