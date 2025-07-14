package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"claude-pilot/core/internal/interfaces"
	"claude-pilot/core/internal/utils"
)

var IGNORE_FILES = []string{".name_index.json"}

// NameIndex maps session names to IDs for fast lookup
type NameIndex struct {
	NameToID map[string]string `json:"name_to_id"`
}

// FileSessionRepository implements SessionRepository using JSON files
type FileSessionRepository struct {
	sessionsDir string
	nameIndex   *NameIndex
	indexMutex  sync.RWMutex
}

// NewFileSessionRepository creates a new file-based session repository
func NewFileSessionRepository(sessionsDir string) (*FileSessionRepository, error) {
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sessions directory: %w", err)
	}

	repo := &FileSessionRepository{
		sessionsDir: sessionsDir,
		nameIndex:   &NameIndex{NameToID: make(map[string]string)},
	}

	// Load or rebuild the name index
	if err := repo.loadNameIndex(); err != nil {
		// If index is corrupted or missing, rebuild it
		if err := repo.rebuildNameIndex(); err != nil {
			return nil, fmt.Errorf("failed to initialize name index: %w", err)
		}
	}

	return repo, nil
}

// Save stores a session to persistent storage
func (r *FileSessionRepository) Save(session *interfaces.Session) error {
	// Update name index first (in memory)
	r.indexMutex.Lock()
	r.nameIndex.NameToID[session.Name] = session.ID
	r.indexMutex.Unlock()

	sessionFile := filepath.Join(r.sessionsDir, session.ID+".json")

	// Use compact JSON marshaling for better performance
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Use more efficient file writing with proper error handling
	if err := r.writeFileAtomic(sessionFile, data, 0644); err != nil {
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

// FindByName retrieves a session by its name using the index
func (r *FileSessionRepository) FindByName(name string) (*interfaces.Session, error) {
	r.indexMutex.RLock()
	id, exists := r.nameIndex.NameToID[name]
	r.indexMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("session with name '%s' not found", name)
	}

	return r.FindByID(id)
}

// List returns all sessions
func (r *FileSessionRepository) List() ([]*interfaces.Session, error) {
	files, err := filepath.Glob(filepath.Join(r.sessionsDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob session files: %w", err)
	}

	// Filter out ignored files
	files = utils.Filter(files, func(file string) bool {
		return !slices.Contains(IGNORE_FILES, filepath.Base(file))
	})

	// Pre-allocate slice with known capacity to avoid reallocation
	sessions := make([]*interfaces.Session, 0, len(files))
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
	// Get the session to find its name for index cleanup
	session, err := r.FindByID(id)
	if err != nil {
		return err // Session not found
	}

	sessionFile := filepath.Join(r.sessionsDir, id+".json")

	if err := os.Remove(sessionFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("session with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to remove session file: %w", err)
	}

	// Remove from name index
	r.indexMutex.Lock()
	delete(r.nameIndex.NameToID, session.Name)
	r.indexMutex.Unlock()

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

// getIndexPath returns the path to the name index file
func (r *FileSessionRepository) getIndexPath() string {
	return filepath.Join(r.sessionsDir, ".name_index.json")
}

// loadNameIndex loads the name index from disk
func (r *FileSessionRepository) loadNameIndex() error {
	indexPath := r.getIndexPath()

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	r.indexMutex.Lock()
	defer r.indexMutex.Unlock()

	return json.Unmarshal(data, r.nameIndex)
}

// saveNameIndex saves the name index to disk
func (r *FileSessionRepository) saveNameIndex() error {
	indexPath := r.getIndexPath()

	r.indexMutex.RLock()
	data, err := json.MarshalIndent(r.nameIndex, "", "  ")
	r.indexMutex.RUnlock()

	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0644)
}

// rebuildNameIndex rebuilds the name index by scanning all session files
func (r *FileSessionRepository) rebuildNameIndex() error {
	sessions, err := r.List()
	if err != nil {
		return err
	}

	r.indexMutex.Lock()
	r.nameIndex.NameToID = make(map[string]string)
	for _, session := range sessions {
		r.nameIndex.NameToID[session.Name] = session.ID
	}
	r.indexMutex.Unlock()

	return r.saveNameIndex()
}

// writeFileAtomic writes data to a file atomically by writing to a temp file first
func (r *FileSessionRepository) writeFileAtomic(filename string, data []byte, perm os.FileMode) error {
	// Write to a temporary file first
	tempFile := filename + ".tmp"

	if err := os.WriteFile(tempFile, data, perm); err != nil {
		return err
	}

	// Atomically move the temp file to the final location
	return os.Rename(tempFile, filename)
}

// SaveIndex saves the name index to disk (call this periodically or at shutdown)
func (r *FileSessionRepository) SaveIndex() error {
	return r.saveNameIndex()
}
