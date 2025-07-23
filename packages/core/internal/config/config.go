package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Backend specifies the terminal multiplexer to use (tmux only)
	Backend string `mapstructure:"backend" yaml:"backend"`

	// BackendPath specifies the custom path to the multiplexer binary
	BackendPath string `mapstructure:"backend_path" yaml:"backend_path"`

	// SessionsDir is the directory where session metadata is stored
	SessionsDir string `mapstructure:"sessions_dir" yaml:"sessions_dir"`

	// DefaultShell is the command to run in new sessions
	DefaultShell string `mapstructure:"default_shell" yaml:"default_shell"`

	// Logging configuration
	Logging LoggingConfig `mapstructure:"logging" yaml:"logging"`

	// UI configuration
	UI UIConfig `mapstructure:"ui" yaml:"ui"`

	// Tmux-specific configuration
	Tmux TmuxConfig `mapstructure:"tmux" yaml:"tmux"`
}

// UIConfig contains user interface configuration
type UIConfig struct {
	// Mode specifies the UI mode (cli, tui)
	Mode string `mapstructure:"mode" yaml:"mode"`

	// Theme specifies the color theme
	Theme string `mapstructure:"theme" yaml:"theme"`

	// ShowIcons enables/disables icon display
	ShowIcons bool `mapstructure:"show_icons" yaml:"show_icons"`
}

// TmuxConfig contains tmux-specific configuration
type TmuxConfig struct {
	// SessionPrefix is prepended to all tmux session names
	SessionPrefix string `mapstructure:"session_prefix" yaml:"session_prefix"`

	// DefaultLayout specifies the default tmux layout
	DefaultLayout string `mapstructure:"default_layout" yaml:"default_layout"`

	// StatusBar configuration
	StatusBar bool `mapstructure:"status_bar" yaml:"status_bar"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	// Enabled controls whether logging is active (disabled by default)
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`

	// Level sets the minimum log level (debug, info, warn, error)
	Level string `mapstructure:"level" yaml:"level"`

	// File is the path to the log file
	File string `mapstructure:"file" yaml:"file"`

	// MaxSize is the maximum size in MB before rotation (0 = no rotation)
	MaxSize int64 `mapstructure:"max_size" yaml:"max_size"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
        return &Config{
                Backend:      "auto", // Auto-detect available backend
		BackendPath:  "",     // Use system PATH
		SessionsDir:  filepath.Join(homeDir, ".config", "claude-pilot", "sessions"),
		DefaultShell: "claude",
		Logging: LoggingConfig{
			Enabled: false, // Disabled by default per requirements
			Level:   "info",
			File:    filepath.Join(homeDir, ".config", "claude-pilot", "claude-pilot.log"),
			MaxSize: 10, // 10MB max log file size
		},
		UI: UIConfig{
			Mode:      "cli",
			Theme:     "default",
			ShowIcons: true,
		},
		Tmux: TmuxConfig{
			SessionPrefix: "claude-",
			DefaultLayout: "main-horizontal",
			StatusBar:     true,
		},
	}
}

// ConfigManager handles configuration loading and saving
type ConfigManager struct {
	configFile string
	config     *Config
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configFile string) *ConfigManager {
	return &ConfigManager{
		configFile: configFile,
		config:     DefaultConfig(),
	}
}

// Load reads configuration from file or creates default config
func (cm *ConfigManager) Load() (*Config, error) {
	viper.SetConfigType("yaml")

	if cm.configFile != "" {
		// Ensure the directory for the custom config file exists
		configDir := filepath.Dir(cm.configFile)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory for custom config: %w", err)
		}
		viper.SetConfigFile(cm.configFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		claudePilotConfigDir := filepath.Join(homeDir, ".config", "claude-pilot")

		// Ensure config directory exists
		if err := os.MkdirAll(claudePilotConfigDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		viper.AddConfigPath(claudePilotConfigDir)
		viper.SetConfigName("claude-pilot")
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("CLAUDE_PILOT")
	viper.AutomaticEnv()

	// Set defaults
	cm.setDefaults()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create a default one
			var configPath string
			if cm.configFile != "" {
				configPath = cm.configFile
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return cm.config, nil // Use defaults if we can't determine config dir
				}
				configPath = filepath.Join(homeDir, ".config", "claude-pilot", "claude-pilot.yaml")
			}

			if err := cm.createDefaultConfigFileAt(configPath); err != nil {
				// If we can't create the config file, just use defaults without error
				// This ensures the application works even if filesystem is read-only
				return cm.config, nil
			}
			// Try to read the newly created config file
			if err := viper.ReadInConfig(); err != nil {
				// If we still can't read it, just use defaults
				return cm.config, nil
			}
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	if err := viper.Unmarshal(cm.config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand home directory paths
	if err := cm.expandHomePaths(); err != nil {
		return nil, fmt.Errorf("failed to expand home directory paths: %w", err)
	}

	// Validate and set computed values
	if err := cm.validateAndSetDefaults(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cm.config, nil
}

// Save writes the current configuration to file
func (cm *ConfigManager) Save() error {
	viper.Set("backend", cm.config.Backend)
	viper.Set("backend_path", cm.config.BackendPath)
	viper.Set("sessions_dir", cm.config.SessionsDir)
	viper.Set("default_shell", cm.config.DefaultShell)
	viper.Set("logging", cm.config.Logging)
	viper.Set("ui", cm.config.UI)
	viper.Set("tmux", cm.config.Tmux)

	return viper.WriteConfig()
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// UpdateConfig updates the configuration
func (cm *ConfigManager) UpdateConfig(config *Config) {
	cm.config = config
}

// setDefaults sets default values in viper
func (cm *ConfigManager) setDefaults() {
	defaults := DefaultConfig()

	viper.SetDefault("backend", defaults.Backend)
	viper.SetDefault("backend_path", defaults.BackendPath)
	viper.SetDefault("sessions_dir", defaults.SessionsDir)
	viper.SetDefault("default_shell", defaults.DefaultShell)
	viper.SetDefault("logging.enabled", defaults.Logging.Enabled)
	viper.SetDefault("logging.level", defaults.Logging.Level)
	viper.SetDefault("logging.file", defaults.Logging.File)
	viper.SetDefault("logging.max_size", defaults.Logging.MaxSize)
	viper.SetDefault("ui.mode", defaults.UI.Mode)
	viper.SetDefault("ui.theme", defaults.UI.Theme)
	viper.SetDefault("ui.show_icons", defaults.UI.ShowIcons)
	viper.SetDefault("tmux.session_prefix", defaults.Tmux.SessionPrefix)
	viper.SetDefault("tmux.default_layout", defaults.Tmux.DefaultLayout)
	viper.SetDefault("tmux.status_bar", defaults.Tmux.StatusBar)
}

// validateAndSetDefaults validates configuration and sets computed defaults
func (cm *ConfigManager) validateAndSetDefaults() error {
	// Ensure sessions directory exists
	if err := os.MkdirAll(cm.config.SessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Validate backend selection - auto resolves to tmux
	if cm.config.Backend == "auto" {
		cm.config.Backend = "tmux"
	}
	
	validBackends := []string{"tmux"}
	isValid := false
	for _, backend := range validBackends {
		if cm.config.Backend == backend {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid backend '%s', must be tmux (zellij support planned for future release)", cm.config.Backend)
	}

	// Validate UI mode
	validModes := []string{"cli", "tui"}
	isValid = false
	for _, mode := range validModes {
		if cm.config.UI.Mode == mode {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid UI mode '%s', must be one of: %v", cm.config.UI.Mode, validModes)
	}

	// Validate log level
	validLevels := []string{"debug", "info", "warn", "error"}
	isValid = false
	for _, level := range validLevels {
		if cm.config.Logging.Level == level {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid log level '%s', must be one of: %v", cm.config.Logging.Level, validLevels)
	}

	return nil
}

// createDefaultConfigFileAt creates a default config file at the specified path
func (cm *ConfigManager) createDefaultConfigFileAt(configFilePath string) error {
	// Ensure the directory for the config file exists
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create the default config content with comments
	defaultConfigContent := `# Claude Pilot Configuration
# Configuration file for Claude Pilot - AI session manager
#
# For more information, visit: https://github.com/HexSleeves/claude-pilot

# Backend selection: auto
backend: auto

# Directory where session metadata is stored
# Will be created automatically if it doesn't exist
sessions_dir: ` + filepath.Join(homeDir, ".config", "claude-pilot", "sessions") + `

# Default shell command to run (claude CLI)
default_shell: claude

# Logging configuration
logging:
  # Enable/disable logging (disabled by default)
  # Logging can also be enabled with the --verbose/-v flag
  enabled: false

  # Log level: debug, info, warn, error
  level: info

  # Path to log file (will be created automatically)
  file: ` + filepath.Join(homeDir, ".config", "claude-pilot", "claude-pilot.log") + `

  # Maximum log file size in MB before rotation (0 = no rotation)
  max_size: 10

# UI configuration
ui:
  # Interface mode: cli or tui
  # cli: Traditional command-line interface
  # tui: Interactive terminal user interface
  mode: cli

  # Theme settings (reserved for future use)
  theme: default

# Backend-specific configurations
tmux:
  # Prefix for tmux session names (optional)
  session_prefix: claude-
`

	// Write the default config file
	if err := os.WriteFile(configFilePath, []byte(defaultConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	return nil
}

// expandHomePaths expands ~ to the user's home directory in all path fields
func (cm *ConfigManager) expandHomePaths() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Expand SessionsDir
	cm.config.SessionsDir = ExpandHomePath(cm.config.SessionsDir, homeDir)

	// Expand BackendPath if it starts with ~
	cm.config.BackendPath = ExpandHomePath(cm.config.BackendPath, homeDir)

	// Expand Logging.File if it starts with ~
	cm.config.Logging.File = ExpandHomePath(cm.config.Logging.File, homeDir)

	return nil
}

// ExpandHomePath expands ~ to homeDir if path starts with ~/
func ExpandHomePath(path, homeDir string) string {
	if path == "~" {
		return homeDir
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
