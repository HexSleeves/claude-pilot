package api

import (
	"os"
	"path/filepath"
)

// DefaultConfigFile returns the default configuration file path
func DefaultConfigFile() string {
	const DEFAULT_CONFIG_FILE = "claude-pilot.yaml"
	const DEFAULT_CONFIG_DIR = ".config/claude-pilot"

	// Check if config file is specified via environment variable
	if configFile := os.Getenv("CLAUDE_PILOT_CONFIG"); configFile != "" {
		return configFile
	}

	// Check current directory first
	if _, err := os.Stat(DEFAULT_CONFIG_FILE); err == nil {
		return DEFAULT_CONFIG_FILE
	}

	// Check user's home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		configPath := filepath.Join(homeDir, DEFAULT_CONFIG_DIR, DEFAULT_CONFIG_FILE)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Return default path (will be created if needed)
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, DEFAULT_CONFIG_DIR, DEFAULT_CONFIG_FILE)
	}

	return DEFAULT_CONFIG_FILE
}

// NewDefaultClient creates a client with default configuration
func NewDefaultClient(verbose bool) (*Client, error) {
	return NewClient(ClientConfig{
		ConfigFile: DefaultConfigFile(),
		Verbose:    verbose,
	})
}

// GetProjectPath resolves a project path, using current directory if empty
func GetProjectPath(projectPath string) string {
	if projectPath == "" {
		if cwd, err := os.Getwd(); err == nil {
			return cwd
		}
		return "."
	}

	// Expand ~ to home directory
	if projectPath == "~" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			return homeDir
		}
	} else if len(projectPath) > 2 && projectPath[:2] == "~/" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			return filepath.Join(homeDir, projectPath[2:])
		}
	}

	// Convert to absolute path
	if absPath, err := filepath.Abs(projectPath); err == nil {
		return absPath
	}

	return projectPath
}
