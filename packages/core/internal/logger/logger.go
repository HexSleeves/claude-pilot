package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// HumanHandler provides human-readable console output
type HumanHandler struct {
	writer io.Writer
	opts   *slog.HandlerOptions
}

// NewHumanHandler creates a new human-readable handler
func NewHumanHandler(w io.Writer, opts *slog.HandlerOptions) *HumanHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &HumanHandler{
		writer: w,
		opts:   opts,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *HumanHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes the log record
func (h *HumanHandler) Handle(_ context.Context, r slog.Record) error {
	// Get level icon and color
	levelIcon := getLevelIcon(r.Level)

	// Build the human-readable message
	var buf strings.Builder

	// Level with icon
	buf.WriteString(fmt.Sprintf("%s%s: %s", r.Level.String(), levelIcon, r.Message))

	// Add important attributes in a readable format
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "operation", "session_name", "name", "command", "error":
			buf.WriteString(fmt.Sprintf(" [%s=%v]", a.Key, a.Value))
		case "duration_ms":
			if ms := a.Value.Int64(); ms > 0 {
				buf.WriteString(fmt.Sprintf(" (%dms)", ms))
			}
		}
		return true
	})

	buf.WriteString("\n")

	_, err := h.writer.Write([]byte(buf.String()))
	return err
}

// WithAttrs returns a new handler with the given attributes
func (h *HumanHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return the same handler
	// In a full implementation, you'd store the attributes
	return h
}

// WithGroup returns a new handler with the given group name
func (h *HumanHandler) WithGroup(name string) slog.Handler {
	// For simplicity, return the same handler
	return h
}

// getLevelIcon returns an appropriate icon for each log level
func getLevelIcon(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "⚙" // Gear symbol for debug
	case slog.LevelInfo:
		return "ℹ" // Info symbol
	case slog.LevelWarn:
		return "⚠" // Warning symbol
	case slog.LevelError:
		return "✗" // Cross symbol for error
	default:
		return "•" // Bullet for default
	}
}

// MultiHandler sends logs to multiple handlers
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a handler that writes to multiple handlers
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// Enabled reports whether any handler handles records at the given level
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle sends the record to all handlers
func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithAttrs returns a new handler with the given attributes added to all handlers
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

// WithGroup returns a new handler with the given group added to all handlers
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

// Logger wraps slog.Logger with claude-pilot specific functionality
type Logger struct {
	*slog.Logger
	config Config
	file   *os.File
	mu     sync.RWMutex
}

// Config holds logger configuration
type Config struct {
	// Enabled controls whether logging is active
	Enabled bool

	// Level sets the minimum log level (debug, info, warn, error)
	Level slog.Level

	// FilePath is the path to the log file
	FilePath string

	// MaxSize is the maximum size in MB before rotation (0 = no rotation)
	MaxSize int64

	// TUIMode when true, only logs to file (never stdout)
	TUIMode bool

	// Verbose enables detailed logging output
	Verbose bool
}

// DefaultConfig returns sensible logging defaults
func DefaultConfig() Config {
	homeDir, _ := os.UserHomeDir()
	logPath := filepath.Join(homeDir, ".config", "claude-pilot", "claude-pilot.log")

	return Config{
		Enabled:  false, // Disabled by default per requirements
		Level:    slog.LevelInfo,
		FilePath: logPath,
		MaxSize:  10, // 10MB max log file size
		TUIMode:  false,
		Verbose:  false,
	}
}

// New creates a new logger with the given configuration
func New(config Config) (*Logger, error) {
	if !config.Enabled {
		// Return a disabled logger that discards all output
		return &Logger{
			Logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
				Level: slog.LevelError + 1, // Higher than any valid level
			})),
			config: config,
		}, nil
	}

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Configure handler options for file logging (JSON)
	fileHandlerOpts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.Verbose, // Add source file info in verbose mode
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format for JSON
			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	// Always use JSON handler for file logging
	fileHandler := slog.NewJSONHandler(file, fileHandlerOpts)

	var handler slog.Handler = fileHandler

	// In CLI mode and verbose, also write human-readable logs to stdout
	if !config.TUIMode && config.Verbose {
		// Configure handler options for console logging (human-readable)
		consoleHandlerOpts := &slog.HandlerOptions{
			Level: config.Level,
		}

		// Create human-readable handler for console
		consoleHandler := NewHumanHandler(os.Stdout, consoleHandlerOpts)

		// Use a multi-handler that writes JSON to file and human-readable to console
		handler = NewMultiHandler(fileHandler, consoleHandler)
	}

	logger := &Logger{
		Logger: slog.New(handler),
		config: config,
		file:   file,
	}

	// Add context about the logger setup
	logger.Debug("Logger initialized",
		"enabled", config.Enabled,
		"level", config.Level.String(),
		"file", config.FilePath,
		"tui_mode", config.TUIMode,
		"verbose", config.Verbose,
	)

	return logger, nil
}

// Close closes the log file if it's open
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// WithContext adds context values to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Logger: l.Logger.With(),
		config: l.config,
		file:   l.file,
	}
}

// WithOperation adds operation context to log entries
func (l *Logger) WithOperation(operation string) *Logger {
	return &Logger{
		Logger: l.Logger.With("operation", operation),
		config: l.config,
		file:   l.file,
	}
}

// WithSession adds session context to log entries
func (l *Logger) WithSession(sessionID, sessionName string) *Logger {
	return &Logger{
		Logger: l.Logger.With("session_id", sessionID, "session_name", sessionName),
		config: l.config,
		file:   l.file,
	}
}

// WithDuration logs operation duration
func (l *Logger) WithDuration(start time.Time) *Logger {
	duration := time.Since(start)
	return &Logger{
		Logger: l.Logger.With("duration_ms", duration.Milliseconds()),
		config: l.config,
		file:   l.file,
	}
}

// Performance logs performance metrics
func (l *Logger) Performance(operation string, start time.Time, attrs ...slog.Attr) {
	duration := time.Since(start)
	allAttrs := append([]slog.Attr{
		slog.String("operation", operation),
		slog.Duration("duration", duration),
		slog.Int64("duration_ms", duration.Milliseconds()),
	}, attrs...)

	l.Logger.LogAttrs(context.Background(), slog.LevelDebug, "Performance metric", allAttrs...)
}

// DebugCommand logs command execution details (only in verbose mode)
func (l *Logger) DebugCommand(command string, args []string, workingDir string) {
	if l.config.Verbose {
		l.Debug("Executing command",
			"command", command,
			"args", args,
			"working_dir", workingDir,
		)
	}
}

// ErrorWithContext logs an error with additional context
func (l *Logger) ErrorWithContext(err error, msg string, attrs ...slog.Attr) {
	allAttrs := append([]slog.Attr{
		slog.String("error", err.Error()),
	}, attrs...)

	l.Logger.LogAttrs(context.Background(), slog.LevelError, msg, allAttrs...)
}

// rotateIfNeeded checks if log rotation is needed and performs it
func (l *Logger) rotateIfNeeded() error {
	if l.config.MaxSize <= 0 || l.file == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	stat, err := l.file.Stat()
	if err != nil {
		return err
	}

	// Check if rotation is needed (convert MB to bytes)
	if stat.Size() < l.config.MaxSize*1024*1024 {
		return nil
	}

	// Close current file
	if err := l.file.Close(); err != nil {
		return err
	}

	// Rename current file with timestamp
	timestamp := time.Now().Format("20060102-150405")
	rotatedPath := fmt.Sprintf("%s.%s", l.config.FilePath, timestamp)
	if err := os.Rename(l.config.FilePath, rotatedPath); err != nil {
		return err
	}

	// Create new file
	newFile, err := os.OpenFile(l.config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	l.file = newFile
	return nil
}

// IsEnabled returns true if logging is enabled
func (l *Logger) IsEnabled() bool {
	return l.config.Enabled
}

// IsVerbose returns true if verbose logging is enabled
func (l *Logger) IsVerbose() bool {
	return l.config.Verbose
}
