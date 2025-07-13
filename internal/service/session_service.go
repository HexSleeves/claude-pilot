package service

import (
	"fmt"
	"time"

	"claude-pilot/internal/interfaces"

	"github.com/google/uuid"
)

// SessionService implements the SessionService interface
type SessionService struct {
	repository  interfaces.SessionRepository
	multiplexer interfaces.TerminalMultiplexer
}

// NewSessionService creates a new session service
func NewSessionService(repository interfaces.SessionRepository, multiplexer interfaces.TerminalMultiplexer) *SessionService {
	return &SessionService{
		repository:  repository,
		multiplexer: multiplexer,
	}
}

// CreateSession creates a new session with both metadata and multiplexer session
func (s *SessionService) CreateSession(name, description, projectPath string) (*interfaces.Session, error) {
	if name == "" {
		name = fmt.Sprintf("session-%s", time.Now().Format("20060102-150405"))
	}

	// Check if session with same name already exists
	if s.repository.Exists(name) {
		return nil, fmt.Errorf("session with name '%s' already exists", name)
	}

	// Create session metadata
	session := &interfaces.Session{
		ID:          uuid.New().String(),
		Name:        name,
		Status:      interfaces.StatusActive,
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
		ProjectPath: projectPath,
		Description: description,
		Messages:    []interfaces.Message{},
	}

	// Save session metadata first
	if err := s.repository.Save(session); err != nil {
		return nil, fmt.Errorf("failed to save session metadata: %w", err)
	}

	// Create the multiplexer session
	req := interfaces.CreateSessionRequest{
		Name:        name,
		Description: description,
		WorkingDir:  projectPath,
		Command:     "claude",
	}

	_, err := s.multiplexer.CreateSession(req)
	if err != nil {
		// If multiplexer session creation fails, mark session as inactive but keep metadata
		session.Status = interfaces.StatusInactive
		s.repository.Save(session)
		return session, fmt.Errorf("session created but failed to create multiplexer session: %w", err)
	}

	// Update session status
	session.Status = interfaces.StatusActive
	if err := s.repository.Save(session); err != nil {
		return session, fmt.Errorf("session created but failed to update status: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by ID or name
func (s *SessionService) GetSession(identifier string) (*interfaces.Session, error) {
	// Try by ID first
	if session, err := s.repository.FindByID(identifier); err == nil {
		return s.updateSessionStatus(session), nil
	}

	// Try by name
	if session, err := s.repository.FindByName(identifier); err == nil {
		return s.updateSessionStatus(session), nil
	}

	return nil, fmt.Errorf("session '%s' not found", identifier)
}

// ListSessions returns all sessions with their current status
func (s *SessionService) ListSessions() ([]*interfaces.Session, error) {
	sessions, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	// Update status for all sessions
	for i, session := range sessions {
		sessions[i] = s.updateSessionStatus(session)
	}

	return sessions, nil
}

// UpdateSession updates session metadata
func (s *SessionService) UpdateSession(session *interfaces.Session) error {
	if !s.repository.Exists(session.ID) {
		return fmt.Errorf("session '%s' not found", session.ID)
	}

	session.LastActive = time.Now()
	return s.repository.Save(session)
}

// DeleteSession removes a session and its multiplexer session
func (s *SessionService) DeleteSession(identifier string) error {
	session, err := s.GetSession(identifier)
	if err != nil {
		return err
	}

	// Kill the multiplexer session if it's running
	if s.multiplexer.IsSessionRunning(session.Name) {
		if err := s.multiplexer.KillSession(session.Name); err != nil {
			return fmt.Errorf("failed to kill multiplexer session: %w", err)
		}
	}

	// Remove session metadata
	if err := s.repository.Delete(session.ID); err != nil {
		return fmt.Errorf("failed to delete session metadata: %w", err)
	}

	return nil
}

// AttachToSession connects to an existing session
func (s *SessionService) AttachToSession(identifier string) error {
	session, err := s.GetSession(identifier)
	if err != nil {
		return err
	}

	// Update session status to connected
	session.Status = interfaces.StatusConnected
	session.LastActive = time.Now()
	if err := s.repository.Save(session); err != nil {
		// Don't fail attachment due to metadata update failure
		fmt.Printf("Warning: failed to update session metadata: %v\n", err)
	}

	// Attach to the multiplexer session
	return s.multiplexer.AttachToSession(session.Name)
}

// AddMessage adds a message to a session's conversation history
func (s *SessionService) AddMessage(sessionID, role, content string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	message := interfaces.Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	session.Messages = append(session.Messages, message)
	session.LastActive = time.Now()

	return s.repository.Save(session)
}

// IsSessionRunning checks if the session's multiplexer is active
func (s *SessionService) IsSessionRunning(identifier string) bool {
	session, err := s.GetSession(identifier)
	if err != nil {
		return false
	}
	return s.multiplexer.IsSessionRunning(session.Name)
}

// updateSessionStatus updates a session's status based on multiplexer state
func (s *SessionService) updateSessionStatus(session *interfaces.Session) *interfaces.Session {
	if s.multiplexer.IsSessionRunning(session.Name) {
		// Check if someone is attached (this is backend-specific and may not be available)
		if muxSession, err := s.multiplexer.GetSession(session.Name); err == nil {
			if muxSession.IsAttached() {
				session.Status = interfaces.StatusConnected
			} else {
				session.Status = interfaces.StatusActive
			}
		} else {
			session.Status = interfaces.StatusActive
		}
	} else {
		session.Status = interfaces.StatusInactive
	}

	return session
}

// GetMultiplexer returns the underlying multiplexer (for legacy compatibility)
func (s *SessionService) GetMultiplexer() interfaces.TerminalMultiplexer {
	return s.multiplexer
}
