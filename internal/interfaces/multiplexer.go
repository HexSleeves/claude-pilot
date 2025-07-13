package interfaces

import (
	"time"
)

// SessionStatus represents the current state of a session
type SessionStatus string

const (
	StatusActive    SessionStatus = "active"
	StatusInactive  SessionStatus = "inactive"
	StatusConnected SessionStatus = "connected"
	StatusError     SessionStatus = "error"
)

// MultiplexerSession represents a session managed by a terminal multiplexer
type MultiplexerSession interface {
	GetID() string
	GetName() string
	GetStatus() SessionStatus
	GetCreatedAt() time.Time
	GetWorkingDir() string
	GetDescription() string
	IsAttached() bool
	IsRunning() bool
}

// CreateSessionRequest contains parameters for creating a new session
type CreateSessionRequest struct {
	Name        string
	Description string
	WorkingDir  string
	Command     string // Command to run in the session (default: "claude")
}

// TerminalMultiplexer defines the interface for terminal multiplexer backends
type TerminalMultiplexer interface {
	// GetName returns the name of the multiplexer backend (e.g., "tmux", "zellij")
	GetName() string

	// IsAvailable checks if the multiplexer binary is available on the system
	IsAvailable() bool

	// CreateSession creates a new terminal session
	CreateSession(req CreateSessionRequest) (MultiplexerSession, error)

	// GetSession retrieves session information by name
	GetSession(name string) (MultiplexerSession, error)

	// ListSessions returns all available sessions
	ListSessions() ([]MultiplexerSession, error)

	// AttachToSession attaches to an existing session (blocking operation)
	AttachToSession(name string) error

	// KillSession terminates a session
	KillSession(name string) error

	// IsSessionRunning checks if a session is currently running
	IsSessionRunning(name string) bool

	// HasSession checks if a session exists
	HasSession(name string) bool
}

// MultiplexerFactory creates multiplexer instances
type MultiplexerFactory interface {
	// CreateMultiplexer creates a multiplexer instance for the given backend
	CreateMultiplexer(backend string) (TerminalMultiplexer, error)

	// GetAvailableBackends returns list of available multiplexer backends
	GetAvailableBackends() []string

	// GetDefaultBackend returns the preferred backend (first available)
	GetDefaultBackend() string

	// GetBackendInfo returns information about a specific backend
	GetBackendInfo(backend string) (map[string]any, error)

	// ValidateBackend checks if a backend name is valid
	ValidateBackend(backend string) error
}
