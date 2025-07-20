package interfaces

import "time"

// SessionStatus represents the current state of a session
type SessionStatus string

const (
	StatusActive    SessionStatus = "active"
	StatusInactive  SessionStatus = "inactive"
	StatusConnected SessionStatus = "connected"
	StatusError     SessionStatus = "error"
	StatusWarning   SessionStatus = "warning"
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
	Backend     string        `json:"backend"`
	CreatedAt   time.Time     `json:"created_at"`
	LastActive  time.Time     `json:"last_active"`
	ProjectPath string        `json:"project_path"`
	Description string        `json:"description"`
	Messages    []Message     `json:"messages"`
}

// AttachmentType represents how to attach to an existing session
type AttachmentType string

const (
	AttachmentNone   AttachmentType = ""       // Create standalone session
	AttachmentPane   AttachmentType = "pane"   // Create as new pane
	AttachmentWindow AttachmentType = "window" // Create as new window/tab
)

// SplitDirection represents the direction for pane splits
type SplitDirection string

const (
	SplitVertical   SplitDirection = "v" // Split vertically (left/right)
	SplitHorizontal SplitDirection = "h" // Split horizontally (top/bottom)
)

// CreateSessionRequest contains parameters for creating a new session
type CreateSessionRequest struct {
	Name           string
	Description    string
	WorkingDir     string
	Command        string         // Command to run in the session (default: "claude")
	AttachTo       string         // Target session name to attach to
	AttachmentType AttachmentType // How to attach (pane, window, or standalone)
	SplitDirection SplitDirection // Direction for pane splits (v/h)
}

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

// TerminalMultiplexer defines the interface for terminal multiplexer backends
type TerminalMultiplexer interface {
	// GetName returns the name of the multiplexer backend (e.g., "tmux")
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

	// SaveIndex saves any in-memory indexes to disk (for performance optimization)
	SaveIndex() error
}

// SessionService defines the business logic interface for session management
type SessionService interface {
	// CreateSession creates a new session with both metadata and multiplexer session
	CreateSession(name, description, projectPath string) (*Session, error)
	
	// CreateSessionAdvanced creates a new session with advanced attachment options
	CreateSessionAdvanced(req CreateSessionRequest) (*Session, error)

	// GetSession retrieves a session by ID or name
	GetSession(identifier string) (*Session, error)

	// ListSessions returns all sessions with their current status
	ListSessions() ([]*Session, error)

	// ListActiveSessions returns all active sessions
	ListFilteredSessions(filter string) ([]*Session, error)

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

	// KillAllSessions terminates all sessions
	KillAllSessions() error
}
