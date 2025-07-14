package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"claude-pilot/core/api"
	"claude-pilot/internal/ui"

	"github.com/spf13/viper"
)

// CommandContext holds common dependencies for all commands
type CommandContext struct {
	Client *api.Client
}

// InitializeCommand handles common initialization for all commands
// This creates an API client that provides access to all functionality
func InitializeCommand() (*CommandContext, error) {
	// Get verbose flag from viper
	verbose := viper.GetBool("verbose")

	// Create API client with configuration
	client, err := api.NewClient(api.ClientConfig{
		ConfigFile: cfgFile,
		Verbose:    verbose,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	return &CommandContext{
		Client: client,
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
