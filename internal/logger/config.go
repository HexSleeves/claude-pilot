package logger

import (
	"log/slog"
	"path/filepath"
)

// Builder provides a fluent interface for logger configuration
type Builder struct {
	config Config
}

// NewBuilder creates a new logger configuration builder with defaults
func NewBuilder() *Builder {
	return &Builder{
		config: DefaultConfig(),
	}
}

// WithEnabled enables or disables logging
func (b *Builder) WithEnabled(enabled bool) *Builder {
	b.config.Enabled = enabled
	return b
}

// WithLevel sets the log level from string
func (b *Builder) WithLevel(level string) *Builder {
	switch level {
	case "debug":
		b.config.Level = slog.LevelDebug
	case "info":
		b.config.Level = slog.LevelInfo
	case "warn":
		b.config.Level = slog.LevelWarn
	case "error":
		b.config.Level = slog.LevelError
	default:
		b.config.Level = slog.LevelInfo
	}
	return b
}

// WithFile sets the log file path
func (b *Builder) WithFile(filePath string) *Builder {
	b.config.FilePath = filePath
	return b
}

// WithFileInDir sets the log file in a specific directory
func (b *Builder) WithFileInDir(dir string) *Builder {
	b.config.FilePath = filepath.Join(dir, "claude-pilot.log")
	return b
}

// WithMaxSize sets the maximum log file size in MB
func (b *Builder) WithMaxSize(sizeMB int64) *Builder {
	b.config.MaxSize = sizeMB
	return b
}

// WithTUIMode enables TUI mode (file-only logging)
func (b *Builder) WithTUIMode(tuiMode bool) *Builder {
	b.config.TUIMode = tuiMode
	return b
}

// WithVerbose enables verbose logging
func (b *Builder) WithVerbose(verbose bool) *Builder {
	b.config.Verbose = verbose
	if verbose {
		// In verbose mode, enable debug level logging
		b.config.Level = slog.LevelDebug
	}
	return b
}

// Build creates a new logger with the configured settings
func (b *Builder) Build() (*Logger, error) {
	return New(b.config)
}

// GetConfig returns the current configuration
func (b *Builder) GetConfig() Config {
	return b.config
}

// FromConfig creates a builder from an existing config
func FromConfig(config Config) *Builder {
	return &Builder{config: config}
}

// QuickSetup provides common logger configurations
type QuickSetup struct{}

// Disabled returns a disabled logger configuration
func (QuickSetup) Disabled() *Builder {
	return NewBuilder().WithEnabled(false)
}

// FileOnly returns a file-only logger configuration
func (QuickSetup) FileOnly(configDir string) *Builder {
	return NewBuilder().
		WithEnabled(true).
		WithFileInDir(configDir).
		WithTUIMode(true)
}

// CLIWithFile returns a CLI logger that logs to both file and stdout when verbose
func (QuickSetup) CLIWithFile(configDir string, verbose bool) *Builder {
	return NewBuilder().
		WithEnabled(true).
		WithFileInDir(configDir).
		WithTUIMode(false).
		WithVerbose(verbose)
}

// TUIWithFile returns a TUI logger that only logs to file
func (QuickSetup) TUIWithFile(configDir string, verbose bool) *Builder {
	return NewBuilder().
		WithEnabled(true).
		WithFileInDir(configDir).
		WithTUIMode(true).
		WithVerbose(verbose)
}

// Setup provides quick setup methods
var Setup QuickSetup
