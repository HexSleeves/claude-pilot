package service

import (
	"fmt"
	"time"

	"claude-pilot/core/internal/logger"
	"claude-pilot/shared/interfaces"

	"log/slog"

	"github.com/google/uuid"
)

// SessionService implements the SessionService interface
type SessionService struct {
	repository  interfaces.SessionRepository
	multiplexer interfaces.TerminalMultiplexer
	logger      *logger.Logger
}

// NewSessionService creates a new session service
func NewSessionService(repository interfaces.SessionRepository, multiplexer interfaces.TerminalMultiplexer) *SessionService {
	// Create a disabled logger by default for backward compatibility
	disabledLogger, _ := logger.Setup.Disabled().Build()
	return &SessionService{
		repository:  repository,
		multiplexer: multiplexer,
		logger:      disabledLogger,
	}
}

// NewSessionServiceWithLogger creates a new session service with a logger
func NewSessionServiceWithLogger(repository interfaces.SessionRepository, multiplexer interfaces.TerminalMultiplexer, log *logger.Logger) *SessionService {
	return &SessionService{
		repository:  repository,
		multiplexer: multiplexer,
		logger:      log,
	}
}

// CreateSession creates a new session with both metadata and multiplexer session
func (s *SessionService) CreateSession(name, description, projectPath string) (*interfaces.Session, error) {
	// Use the advanced method with default parameters
	req := interfaces.CreateSessionRequest{
		Name:           name,
		Description:    description,
		WorkingDir:     projectPath,
		Command:        "claude",
		AttachTo:       "",
		AttachmentType: interfaces.AttachmentNone,
		SplitDirection: interfaces.SplitVertical,
	}
	return s.CreateSessionAdvanced(req)
}

// CreateSessionAdvanced creates a new session with advanced attachment options
func (s *SessionService) CreateSessionAdvanced(req interfaces.CreateSessionRequest) (*interfaces.Session, error) {
	start := time.Now()

	if req.Name == "" {
		req.Name = fmt.Sprintf("session-%s", time.Now().Format("20060102-150405"))
	}

	s.logger.Debug("Creating session",
		"name", req.Name,
		"description", req.Description,
		"project_path", req.WorkingDir,
		"attach_to", req.AttachTo,
		"attachment_type", req.AttachmentType,
		"split_direction", req.SplitDirection)

	// For attached sessions, we don't create separate metadata since they are part of existing sessions
	if req.AttachTo != "" && req.AttachmentType != interfaces.AttachmentNone {
		return s.createAttachedSession(req, start)
	}

	// Check if session with same name already exists
	if s.repository.Exists(req.Name) {
		s.logger.Warn("Session creation failed: name already exists",
			"name", req.Name)
		return nil, fmt.Errorf("session with name '%s' already exists", req.Name)
	}

	// Create session metadata
	session := &interfaces.Session{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Backend:     s.multiplexer.GetName(),
		Status:      interfaces.StatusActive,
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
		ProjectPath: req.WorkingDir,
		Description: req.Description,
	}

	// Save session metadata first
	if err := s.repository.Save(session); err != nil {
		s.logger.Error("Failed to save session metadata",
			"session_id", session.ID,
			"name", req.Name,
			"error", err)
		return nil, fmt.Errorf("failed to save session metadata: %w", err)
	}

	// Set default command if not provided
	if req.Command == "" {
		req.Command = "claude"
	}

	s.logger.Debug("Creating multiplexer session",
		"session_id", session.ID,
		"name", req.Name,
		"command", req.Command,
		"working_dir", req.WorkingDir)

	_, err := s.multiplexer.CreateSession(req)
	if err != nil {
		s.logger.Error("Failed to create multiplexer session",
			"session_id", session.ID,
			"name", req.Name,
			"error", err)

		// If multiplexer session creation fails, mark session as inactive but keep metadata
		session.Status = interfaces.StatusInactive
		if err := s.repository.Save(session); err != nil {
			s.logger.Error("Failed to update session status",
				"session_id", session.ID,
				"name", req.Name,
				"error", err)
		}

		return session, fmt.Errorf("session created but failed to create multiplexer session: %w", err)
	}

	// Update session status
	session.Status = interfaces.StatusActive
	if err := s.repository.Save(session); err != nil {
		s.logger.Error("Failed to update session status",
			"session_id", session.ID,
			"name", req.Name,
			"error", err)
		return session, fmt.Errorf("session created but failed to update status: %w", err)
	}

	// Save index after session creation (important operations)
	if err := s.repository.SaveIndex(); err != nil {
		// Index save failure is not critical, just log it
		s.logger.Warn("Failed to save name index after session creation",
			"session_id", session.ID,
			"name", req.Name,
			"error", err)
	}

	s.logger.Performance("CreateSession", start,
		slog.String("session_id", session.ID),
		slog.String("name", req.Name),
		slog.String("project_path", req.WorkingDir))

	s.logger.Info("Session created successfully",
		"session_id", session.ID,
		"name", req.Name,
		"status", string(session.Status))

	return session, nil
}

// createAttachedSession creates a session attached to an existing session as pane or window
func (s *SessionService) createAttachedSession(req interfaces.CreateSessionRequest, start time.Time) (*interfaces.Session, error) {
	s.logger.Debug("Creating attached session",
		"name", req.Name,
		"attach_to", req.AttachTo,
		"attachment_type", req.AttachmentType)

	// Verify target session exists
	targetSession, err := s.GetSession(req.AttachTo)
	if err != nil {
		s.logger.Error("Target session not found",
			"attach_to", req.AttachTo,
			"error", err)
		return nil, fmt.Errorf("target session '%s' not found: %w", req.AttachTo, err)
	}

	// Set default command if not provided
	if req.Command == "" {
		req.Command = "claude"
	}

	// Create the attached multiplexer session (pane or window)
	_, err = s.multiplexer.CreateSession(req)
	if err != nil {
		s.logger.Error("Failed to create attached session",
			"name", req.Name,
			"attach_to", req.AttachTo,
			"attachment_type", req.AttachmentType,
			"error", err)
		return nil, fmt.Errorf("failed to create attached session: %w", err)
	}

	// Create a virtual session object representing the attachment
	// This doesn't get saved to repository since it's part of the target session
	attachedSession := &interfaces.Session{
		ID:          uuid.New().String(),
		Name:        fmt.Sprintf("%s-attached-to-%s", req.Name, req.AttachTo),
		Backend:     s.multiplexer.GetName(),
		Status:      interfaces.StatusActive,
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
		ProjectPath: req.WorkingDir,
		Description: fmt.Sprintf("%s (attached as %s to %s)", req.Description, req.AttachmentType, req.AttachTo),
	}

	s.logger.Performance("CreateAttachedSession", start,
		slog.String("name", req.Name),
		slog.String("attach_to", req.AttachTo),
		slog.String("attachment_type", string(req.AttachmentType)))

	s.logger.Info("Attached session created successfully",
		"name", req.Name,
		"attach_to", req.AttachTo,
		"attachment_type", req.AttachmentType,
		"target_session_id", targetSession.ID)

	return attachedSession, nil
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
	start := time.Now()

	s.logger.Debug("Listing sessions")

	sessions, err := s.repository.List()
	if err != nil {
		s.logger.Error("Failed to list sessions from repository", "error", err)
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	s.logger.Debug("Retrieved sessions from repository", "count", len(sessions))

	// Batch update status for all sessions
	s.batchUpdateSession(sessions)

	s.logger.Performance("ListSessions", start, slog.Int("session_count", len(sessions)))

	s.logger.Debug("Sessions listed successfully", "count", len(sessions))

	return sessions, nil
}

// ListFilteredSessions returns all sessions with the given filter
func (s *SessionService) ListFilteredSessions(filter string) ([]*interfaces.Session, error) {
	start := time.Now()

	s.logger.Debug("Listing sessions with filter", "filter", filter)

	sessions, err := s.ListSessions()
	if err != nil {
		s.logger.Error("Failed to list sessions from repository", "error", err)
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	// Apply filters
	if filter != "" {
		switch filter {
		case "active":
			sessions = filterByStatus(sessions, interfaces.StatusActive)
		case "inactive":
			sessions = filterByStatus(sessions, interfaces.StatusInactive)
		}
	}

	s.logger.Debug("Sessions listed successfully", "count", len(sessions))

	s.logger.Performance("ListFilteredSessions", start, slog.String("filter", filter), slog.Int("session_count", len(sessions)))

	return sessions, nil
}

// filterActiveSessions filters the sessions to only include active sessions

func filterByStatus(sessions []*interfaces.Session, status interfaces.SessionStatus) []*interfaces.Session {
	filteredSessions := make([]*interfaces.Session, 0)
	for _, session := range sessions {
		if session.Status == status {
			filteredSessions = append(filteredSessions, session)
		}
	}
	return filteredSessions
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
	start := time.Now()

	s.logger.Debug("Deleting session", "identifier", identifier)

	session, err := s.GetSession(identifier)
	if err != nil {
		s.logger.Error("Failed to find session for deletion",
			"identifier", identifier,
			"error", err)
		return err
	}

	sessionLogger := s.logger.WithSession(session.ID, session.Name)

	// Kill the multiplexer session if it's running
	if s.multiplexer.IsSessionRunning(session.Name) {
		sessionLogger.Debug("Killing running multiplexer session")
		if err := s.multiplexer.KillSession(session.Name); err != nil {
			sessionLogger.Error("Failed to kill multiplexer session", "error", err)
			return fmt.Errorf("failed to kill multiplexer session: %w", err)
		}
	}

	// Remove session metadata
	if err := s.repository.Delete(session.ID); err != nil {
		sessionLogger.Error("Failed to delete session metadata", "error", err)
		return fmt.Errorf("failed to delete session metadata: %w", err)
	}

	// Save index after deletion (important operations)
	if err := s.repository.SaveIndex(); err != nil {
		// Index save failure is not critical, just log it
		sessionLogger.Warn("Failed to save name index after session deletion", "error", err)
	}

	s.logger.Performance("DeleteSession", start,
		slog.String("session_id", session.ID),
		slog.String("name", session.Name))

	sessionLogger.Info("Session deleted successfully")

	return nil
}

// AttachToSession connects to an existing session
func (s *SessionService) AttachToSession(identifier string) error {
	start := time.Now()

	s.logger.Debug("Attaching to session", "identifier", identifier)

	session, err := s.GetSession(identifier)
	if err != nil {
		s.logger.Error("Failed to find session for attachment",
			"identifier", identifier,
			"error", err)
		return err
	}

	sessionLogger := s.logger.WithSession(session.ID, session.Name)

	// Update session status to connected
	session.Status = interfaces.StatusConnected
	session.LastActive = time.Now()
	if err := s.repository.Save(session); err != nil {
		// Don't fail attachment due to metadata update failure
		sessionLogger.Warn("Failed to update session metadata before attachment", "error", err)
	}

	sessionLogger.Info("Attaching to multiplexer session")

	// Attach to the multiplexer session
	err = s.multiplexer.AttachToSession(session.Name)
	if err != nil {
		sessionLogger.Error("Failed to attach to multiplexer session", "error", err)
		return err
	}

	s.logger.Performance("AttachToSession", start,
		slog.String("session_id", session.ID),
		slog.String("name", session.Name))

	return nil
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

// batchUpdateSessionStatus efficiently updates status for multiple sessions
func (s *SessionService) batchUpdateSession(sessions []*interfaces.Session) {
	// Get all multiplexer sessions once
	muxSessions, err := s.multiplexer.ListSessions()
	if err != nil {
		// If we can't get multiplexer sessions, fall back to individual checks
		for i, session := range sessions {
			sessions[i] = s.updateSessionStatus(session)
		}
		return
	}

	// Create a map of session names to multiplexer sessions for O(1) lookup
	// Pre-allocate map with expected capacity
	muxSessionMap := make(map[string]interfaces.MultiplexerSession, len(muxSessions))
	for _, muxSession := range muxSessions {
		muxSessionMap[muxSession.GetName()] = muxSession
	}

	// Update all sessions using the batch data
	for _, session := range sessions {
		if muxSession, exists := muxSessionMap[session.Name]; exists {
			// Session exists in multiplexer
			if muxSession.IsAttached() {
				session.Status = interfaces.StatusConnected
			} else {
				session.Status = interfaces.StatusActive
			}

			// Update the session with the multiplexer session
			session.Panes, err = s.multiplexer.GetSessionPaneCount(session.Name)
			if err != nil {
				s.logger.Error("Failed to get session pane count", "error", err)
			}
		} else {
			// Session not found in multiplexer
			session.Status = interfaces.StatusInactive
			session.Panes = 0
		}
	}
}

// KillAllSessions terminates all sessions
func (s *SessionService) KillAllSessions() error {
	sessions, err := s.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	var errors []string
	for _, session := range sessions {
		if err := s.DeleteSession(session.ID); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete session %s: %v", session.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors deleting sessions: %v", errors)
	}

	return nil
}

// GetMultiplexer returns the underlying multiplexer (for legacy compatibility)
func (s *SessionService) GetMultiplexer() interfaces.TerminalMultiplexer {
	return s.multiplexer
}

// GetRepository returns the underlying repository
func (s *SessionService) GetRepository() interfaces.SessionRepository {
	return s.repository
}

// GetSessionPaneCount returns the number of panes in a session
func (s *SessionService) GetSessionPaneCount(identifier string) (int, error) {
	session, err := s.GetSession(identifier)
	if err != nil {
		return 0, err
	}

	// Get pane count from multiplexer
	return s.multiplexer.GetSessionPaneCount(session.Name)
}
