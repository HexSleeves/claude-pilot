package manager

import (
	"fmt"

	"claude-pilot/internal/config"
	"claude-pilot/internal/interfaces"
	"claude-pilot/internal/multiplexer"
	"claude-pilot/internal/service"
	"claude-pilot/internal/storage"
)

// SessionManager orchestrates session management using the new architecture
type SessionManager struct {
	config      *config.Config
	factory     interfaces.MultiplexerFactory
	repository  interfaces.SessionRepository
	service     interfaces.SessionService
	multiplexer interfaces.TerminalMultiplexer
}

// NewSessionManager creates a new session manager with dependency injection
func NewSessionManager(cfg *config.Config) (*SessionManager, error) {
	// Create multiplexer factory
	factory := multiplexer.NewFactory(cfg.Tmux.SessionPrefix)

	// Create multiplexer instance based on configuration
	var mux interfaces.TerminalMultiplexer
	var err error

	if cfg.Backend == "auto" {
		mux, err = factory.CreateMultiplexer("auto")
	} else {
		mux, err = factory.CreateMultiplexer(cfg.Backend)
	}

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

	return &SessionManager{
		config:      cfg,
		factory:     factory,
		repository:  repository,
		service:     sessionService,
		multiplexer: mux,
	}, nil
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession(name, description, projectPath string) (*interfaces.Session, error) {
	return sm.service.CreateSession(name, description, projectPath)
}

// GetSession retrieves a session by ID or name
func (sm *SessionManager) GetSession(identifier string) (*interfaces.Session, error) {
	return sm.service.GetSession(identifier)
}

// ListSessions returns all sessions
func (sm *SessionManager) ListSessions() ([]*interfaces.Session, error) {
	return sm.service.ListSessions()
}

// UpdateSession updates session metadata
func (sm *SessionManager) UpdateSession(session *interfaces.Session) error {
	return sm.service.UpdateSession(session)
}

// DeleteSession removes a session
func (sm *SessionManager) DeleteSession(identifier string) error {
	return sm.service.DeleteSession(identifier)
}

// AttachToSession connects to an existing session
func (sm *SessionManager) AttachToSession(identifier string) error {
	return sm.service.AttachToSession(identifier)
}

// AddMessage adds a message to a session's conversation history
func (sm *SessionManager) AddMessage(sessionID, role, content string) error {
	return sm.service.AddMessage(sessionID, role, content)
}

// IsSessionRunning checks if the session's multiplexer is active
func (sm *SessionManager) IsSessionRunning(identifier string) bool {
	return sm.service.IsSessionRunning(identifier)
}

// GetMultiplexer returns the underlying multiplexer (for compatibility)
func (sm *SessionManager) GetMultiplexer() interfaces.TerminalMultiplexer {
	return sm.multiplexer
}

// GetConfig returns the configuration
func (sm *SessionManager) GetConfig() *config.Config {
	return sm.config
}

// GetAvailableBackends returns available multiplexer backends
func (sm *SessionManager) GetAvailableBackends() []string {
	return sm.factory.GetAvailableBackends()
}

// SwitchBackend changes the multiplexer backend
func (sm *SessionManager) SwitchBackend(backend string) error {
	// Validate backend
	if err := sm.factory.ValidateBackend(backend); err != nil {
		return err
	}

	// Create new multiplexer
	newMux, err := sm.factory.CreateMultiplexer(backend)
	if err != nil {
		return fmt.Errorf("failed to create new multiplexer: %w", err)
	}

	// Update configuration
	sm.config.Backend = backend
	sm.multiplexer = newMux

	// Create new service with new multiplexer
	sm.service = service.NewSessionService(sm.repository, newMux)

	return nil
}

// GetBackendInfo returns information about available backends
func (sm *SessionManager) GetBackendInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	info["current"] = sm.config.Backend
	info["available"] = sm.factory.GetAvailableBackends()
	info["default"] = sm.factory.GetDefaultBackend()

	backends := make(map[string]interface{})
	for _, backend := range []string{"tmux", "zellij", "auto"} {
		if backendInfo, err := sm.factory.GetBackendInfo(backend); err == nil {
			backends[backend] = backendInfo
		}
	}
	info["backends"] = backends

	return info, nil
}

// KillAllSessions terminates all sessions
func (sm *SessionManager) KillAllSessions() error {
	sessions, err := sm.service.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	var errors []string
	for _, session := range sessions {
		if err := sm.service.DeleteSession(session.ID); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete session %s: %v", session.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors deleting sessions: %v", errors)
	}

	return nil
}

// ValidateConfiguration validates the current configuration
func (sm *SessionManager) ValidateConfiguration() error {
	// Check if current backend is available
	if !sm.multiplexer.IsAvailable() {
		return fmt.Errorf("configured backend '%s' is not available", sm.config.Backend)
	}

	// Check if sessions directory is accessible
	if _, err := sm.repository.List(); err != nil {
		return fmt.Errorf("sessions directory is not accessible: %w", err)
	}

	return nil
}

// Cleanup performs cleanup operations
func (sm *SessionManager) Cleanup() error {
	// This could include cleanup of orphaned sessions, temporary files, etc.
	// For now, it's a placeholder for future cleanup logic
	return nil
}
