package interfaces

import (
	"time"
)

// Message represents a message in a Claude session
type Message struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Session represents a Claude coding session with persistence
type Session struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Status      SessionStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	LastActive  time.Time     `json:"last_active"`
	ProjectPath string        `json:"project_path"`
	Description string        `json:"description"`
	Messages    []Message     `json:"messages"`
}

// SessionRepository handles persistence of session metadata
type SessionRepository interface {
	// Save stores a session to persistent storage
	Save(session *Session) error

	// FindByID retrieves a session by its unique ID
	FindByID(id string) (*Session, error)

	// FindByName retrieves a session by its name
	FindByName(name string) (*Session, error)

	// List returns all sessions
	List() ([]*Session, error)

	// Delete removes a session from storage
	Delete(id string) error

	// Exists checks if a session exists by ID or name
	Exists(identifier string) bool
}

// SessionService defines the business logic interface for session management
type SessionService interface {
	// CreateSession creates a new session with both metadata and multiplexer session
	CreateSession(name, description, projectPath string) (*Session, error)

	// GetSession retrieves a session by ID or name
	GetSession(identifier string) (*Session, error)

	// ListSessions returns all sessions with their current status
	ListSessions() ([]*Session, error)

	// UpdateSession updates session metadata
	UpdateSession(session *Session) error

	// DeleteSession removes a session and its multiplexer session
	DeleteSession(identifier string) error

	// AttachToSession connects to an existing session
	AttachToSession(identifier string) error

	// AddMessage adds a message to a session's conversation history
	AddMessage(sessionID, role, content string) error

	// IsSessionRunning checks if the session's multiplexer is active
	IsSessionRunning(identifier string) bool
}
