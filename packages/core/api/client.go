package api

import (
	"fmt"
	"os"

	"claude-pilot/core/internal/config"
	"claude-pilot/core/internal/logger"
	"claude-pilot/core/internal/multiplexer"
	"claude-pilot/core/internal/service"
	"claude-pilot/core/internal/storage"
	"claude-pilot/shared/interfaces"
)

// Client provides a high-level API for both CLI and TUI to consume
type Client struct {
	config      *config.Config
	logger      *logger.Logger
	service     interfaces.SessionService
	multiplexer interfaces.TerminalMultiplexer
}

// ClientConfig holds configuration options for creating a client
type ClientConfig struct {
	ConfigFile string
	Verbose    bool
}

// NewClient creates a new API client with the specified configuration
func NewClient(cfg ClientConfig) (*Client, error) {
	// Load configuration using the provided config file
	configManager := config.NewConfigManager(cfg.ConfigFile)
	config, err := configManager.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	// Check for LOG_LEVEL environment variable to override config
	logLevel := config.Logging.Level
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel = envLogLevel
	}

	// Create logger based on configuration and flags
	loggerBuilder := logger.NewBuilder().
		WithEnabled(config.Logging.Enabled || cfg.Verbose).
		WithLevel(logLevel).
		WithFile(config.Logging.File).
		WithMaxSize(config.Logging.MaxSize).
		WithTUIMode(config.UI.Mode == "tui").
		WithVerbose(cfg.Verbose)

	log, err := loggerBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Create multiplexer instance based on configuration
	mux, err := multiplexer.CreateMultiplexer(config.Backend, config.Tmux.SessionPrefix)
	if err != nil {
		log.Error("Failed to create multiplexer",
			"backend", config.Backend,
			"prefix", config.Tmux.SessionPrefix,
			"error", err)
		return nil, fmt.Errorf("failed to create multiplexer: %w", err)
	}

	// Create repository
	repository, err := storage.NewFileSessionRepository(config.SessionsDir)
	if err != nil {
		log.Error("Failed to create repository",
			"sessions_dir", config.SessionsDir,
			"error", err)
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Create service with logger
	sessionService := service.NewSessionServiceWithLogger(repository, mux, log)

	log.Info("Client initialized successfully",
		"backend", config.Backend,
		"sessions_dir", config.SessionsDir,
		"ui_mode", config.UI.Mode,
		"logging_enabled", config.Logging.Enabled,
		"verbose", cfg.Verbose)

	return &Client{
		config:      config,
		logger:      log,
		service:     sessionService,
		multiplexer: mux,
	}, nil
}

// GetConfig returns the current configuration
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// GetLogger returns the logger instance
func (c *Client) GetLogger() *logger.Logger {
	return c.logger
}

// GetBackend returns the name of the current multiplexer backend
func (c *Client) GetBackend() string {
	return c.multiplexer.GetName()
}

// CreateSessionRequest contains parameters for creating a new session
type CreateSessionRequest struct {
	Name           string
	Description    string
	ProjectPath    string
	AttachTo       string                    // Target session to attach to
	AttachmentType interfaces.AttachmentType // How to attach (pane, window, or standalone)
	SplitDirection interfaces.SplitDirection // Direction for pane splits
}

// CreateSession creates a new session with the specified parameters
func (c *Client) CreateSession(req CreateSessionRequest) (*interfaces.Session, error) {
	// Convert API request to service request
	serviceReq := interfaces.CreateSessionRequest{
		Name:           req.Name,
		Description:    req.Description,
		WorkingDir:     req.ProjectPath,
		Command:        "claude",
		AttachTo:       req.AttachTo,
		AttachmentType: req.AttachmentType,
		SplitDirection: req.SplitDirection,
	}

	return c.service.CreateSessionAdvanced(serviceReq)
}

// ListSessions returns all sessions
func (c *Client) ListSessions() ([]*interfaces.Session, error) {
	return c.service.ListSessions()
}

// ListFilteredSessions returns all sessions with the given filter
func (c *Client) ListFilteredSessions(filter string) ([]*interfaces.Session, error) {
	return c.service.ListFilteredSessions(filter)
}

// GetSession retrieves a session by ID or name
func (c *Client) GetSession(identifier string) (*interfaces.Session, error) {
	return c.service.GetSession(identifier)
}

// AttachToSession connects to an existing session
func (c *Client) AttachToSession(identifier string) error {
	return c.service.AttachToSession(identifier)
}

// KillSession terminates a specific session
func (c *Client) KillSession(identifier string) error {
	return c.service.DeleteSession(identifier)
}

// KillAllSessions terminates all sessions
func (c *Client) KillAllSessions() error {
	return c.service.KillAllSessions()
}

// IsSessionRunning checks if a session's multiplexer is active
func (c *Client) IsSessionRunning(identifier string) bool {
	return c.service.IsSessionRunning(identifier)
}

// GetSessionPaneCount returns the number of panes in a session
func (c *Client) GetSessionPaneCount(identifier string) (int, error) {
	return c.service.GetSessionPaneCount(identifier)
}

// Session represents a session with all its data (re-exported for convenience)
type Session = interfaces.Session

// SessionStatus represents the status of a session (re-exported for convenience)
type SessionStatus = interfaces.SessionStatus

// Message represents a message in a session (re-exported for convenience)
type Message = interfaces.Message

// Status constants (re-exported for convenience)
const (
	StatusActive    = interfaces.StatusActive
	StatusInactive  = interfaces.StatusInactive
	StatusConnected = interfaces.StatusConnected
	StatusError     = interfaces.StatusError
)
