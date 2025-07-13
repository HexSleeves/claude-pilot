package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"claude-pilot/internal/interfaces"
)

// FileSessionRepository implements SessionRepository using JSON files
type FileSessionRepository struct {
	sessionsDir string
}

// NewFileSessionRepository creates a new file-based session repository
func NewFileSessionRepository(sessionsDir string) (*FileSessionRepository, error) {
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sessions directory: %w", err)
	}

	return &FileSessionRepository{
		sessionsDir: sessionsDir,
	}, nil
}

// Save stores a session to persistent storage
func (r *FileSessionRepository) Save(session *interfaces.Session) error {
	sessionFile := filepath.Join(r.sessionsDir, session.ID+".json")

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(sessionFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// FindByID retrieves a session by its unique ID
func (r *FileSessionRepository) FindByID(id string) (*interfaces.Session, error) {
	sessionFile := filepath.Join(r.sessionsDir, id+".json")

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session with ID '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session interfaces.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// FindByName retrieves a session by its name
func (r *FileSessionRepository) FindByName(name string) (*interfaces.Session, error) {
	sessions, err := r.List()
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.Name == name {
			return session, nil
		}
	}

	return nil, fmt.Errorf("session with name '%s' not found", name)
}

// List returns all sessions
func (r *FileSessionRepository) List() ([]*interfaces.Session, error) {
	files, err := filepath.Glob(filepath.Join(r.sessionsDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob session files: %w", err)
	}

	var sessions []*interfaces.Session
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			// Skip corrupted files
			continue
		}

		var session interfaces.Session
		if err := json.Unmarshal(data, &session); err != nil {
			// Skip corrupted files
			continue
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// Delete removes a session from storage
func (r *FileSessionRepository) Delete(id string) error {
	sessionFile := filepath.Join(r.sessionsDir, id+".json")

	if err := os.Remove(sessionFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("session with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to remove session file: %w", err)
	}

	return nil
}

// Exists checks if a session exists by ID or name
func (r *FileSessionRepository) Exists(identifier string) bool {
	// Try by ID first
	if _, err := r.FindByID(identifier); err == nil {
		return true
	}

	// Try by name
	if _, err := r.FindByName(identifier); err == nil {
		return true
	}

	return false
}
