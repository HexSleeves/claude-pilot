package multiplexer

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"claude-pilot/core/internal/interfaces"
)

// ZellijMultiplexer implements the TerminalMultiplexer interface for zellij
type ZellijMultiplexer struct {
	sessionPrefix string
	zellijPath    string
}

// ZellijSession implements the MultiplexerSession interface
type ZellijSession struct {
	id          string
	name        string
	status      interfaces.SessionStatus
	createdAt   time.Time
	workingDir  string
	description string
	isAttached  bool
	isRunning   bool
}

// GetID returns the session ID
func (s *ZellijSession) GetID() string {
	return s.id
}

// GetName returns the session name
func (s *ZellijSession) GetName() string {
	return s.name
}

// GetStatus returns the session status
func (s *ZellijSession) GetStatus() interfaces.SessionStatus {
	return s.status
}

// GetCreatedAt returns the creation time
func (s *ZellijSession) GetCreatedAt() time.Time {
	return s.createdAt
}

// GetWorkingDir returns the working directory
func (s *ZellijSession) GetWorkingDir() string {
	return s.workingDir
}

// GetDescription returns the session description
func (s *ZellijSession) GetDescription() string {
	return s.description
}

// IsAttached returns whether someone is attached to the session
func (s *ZellijSession) IsAttached() bool {
	return s.isAttached
}

// IsRunning returns whether the session is running
func (s *ZellijSession) IsRunning() bool {
	return s.isRunning
}

// NewZellijMultiplexer creates a new zellij multiplexer instance
func NewZellijMultiplexer(sessionPrefix string) (*ZellijMultiplexer, error) {
	if sessionPrefix == "" {
		sessionPrefix = "claude-pilot"
	}

	// Try to find zellij binary
	zellijPath := "zellij" // default
	if _, err := exec.LookPath("zellij"); err != nil {
		// Try common locations
		paths := []string{
			"/opt/homebrew/bin/zellij",
			"/usr/local/bin/zellij",
			"/usr/bin/zellij",
			"~/.cargo/bin/zellij",
		}
		for _, path := range paths {
			if strings.HasPrefix(path, "~/") {
				homeDir, _ := os.UserHomeDir()
				path = strings.Replace(path, "~", homeDir, 1)
			}
			if _, err := os.Stat(path); err == nil {
				zellijPath = path
				break
			}
		}
	}

	return &ZellijMultiplexer{
		sessionPrefix: sessionPrefix,
		zellijPath:    zellijPath,
	}, nil
}

// GetName returns the name of the multiplexer backend
func (zm *ZellijMultiplexer) GetName() string {
	return "zellij"
}

// IsAvailable checks if zellij is available on the system
func (zm *ZellijMultiplexer) IsAvailable() bool {
	_, err := exec.LookPath("zellij")
	if err != nil {
		// Check common paths
		paths := []string{
			"/opt/homebrew/bin/zellij",
			"/usr/local/bin/zellij",
			"/usr/bin/zellij",
		}
		homeDir, _ := os.UserHomeDir()
		if homeDir != "" {
			paths = append(paths, homeDir+"/.cargo/bin/zellij")
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
		return false
	}
	return true
}

// CreateSession creates a new zellij session
func (zm *ZellijMultiplexer) CreateSession(req interfaces.CreateSessionRequest) (interfaces.MultiplexerSession, error) {
	sessionName := fmt.Sprintf("%s-%s", zm.sessionPrefix, req.Name)

	// Check if zellij session already exists
	if zm.HasSession(req.Name) {
		return nil, fmt.Errorf("zellij session '%s' already exists", req.Name)
	}

	// Determine command to run (default to "claude")
	command := req.Command
	if command == "" {
		command = "claude"
	}

	// Create zellij session
	var cmd *exec.Cmd
	if req.WorkingDir != "" {
		// Change to working directory and create session
		cmd = exec.Command(zm.zellijPath, "-s", sessionName, "-c", req.WorkingDir, "--", command)
	} else {
		// Create session in current directory
		cmd = exec.Command(zm.zellijPath, "-s", sessionName, "--", command)
	}

	// Start the session in detached mode
	cmd.Env = append(os.Environ(), "ZELLIJ_AUTO_ATTACH=false")

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to create zellij session: %w", err)
	}

	// Wait a moment for the session to initialize
	time.Sleep(100 * time.Millisecond)

	session := &ZellijSession{
		id:          sessionName,
		name:        req.Name,
		status:      interfaces.StatusActive,
		createdAt:   time.Now(),
		workingDir:  req.WorkingDir,
		description: req.Description,
		isAttached:  false,
		isRunning:   true,
	}

	return session, nil
}

// GetSession retrieves session information by name
func (zm *ZellijMultiplexer) GetSession(name string) (interfaces.MultiplexerSession, error) {
	sessions, err := zm.ListSessions()
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.GetName() == name {
			return session, nil
		}
	}

	return nil, fmt.Errorf("zellij session '%s' not found", name)
}

// ListSessions returns all available zellij sessions
func (zm *ZellijMultiplexer) ListSessions() ([]interfaces.MultiplexerSession, error) {
	// Get list of zellij sessions
	cmd := exec.Command(zm.zellijPath, "list-sessions")
	output, err := cmd.Output()
	if err != nil {
		// If no sessions exist, zellij might return error
		if strings.Contains(err.Error(), "No active zellij sessions") {
			return []interfaces.MultiplexerSession{}, nil
		}
		return nil, fmt.Errorf("failed to list zellij sessions: %w", err)
	}

	var sessions []interfaces.MultiplexerSession
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Parse zellij session list output
	// Expected format varies, but generally includes session names
	sessionNameRegex := regexp.MustCompile(`(\S+)\s*\(.*\)`)

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Active sessions:") {
			continue
		}

		// Try to extract session name
		var sessionName string
		if matches := sessionNameRegex.FindStringSubmatch(line); len(matches) > 1 {
			sessionName = matches[1]
		} else {
			// Fallback: use the first word as session name
			parts := strings.Fields(line)
			if len(parts) > 0 {
				sessionName = parts[0]
			}
		}

		// Only include our claude-pilot sessions
		if !strings.HasPrefix(sessionName, zm.sessionPrefix+"-") {
			continue
		}

		// Extract the user-friendly name
		name := strings.TrimPrefix(sessionName, zm.sessionPrefix+"-")

		// Parse attachment status (if available in output)
		isAttached := strings.Contains(line, "ATTACHED") || strings.Contains(line, "attached")

		session := &ZellijSession{
			id:         sessionName,
			name:       name,
			status:     interfaces.StatusActive,
			createdAt:  time.Now(), // Zellij doesn't provide creation time in list
			isAttached: isAttached,
			isRunning:  true,
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// AttachToSession attaches to an existing zellij session
func (zm *ZellijMultiplexer) AttachToSession(name string) error {
	session, err := zm.GetSession(name)
	if err != nil {
		return err
	}

	zellijSession := session.(*ZellijSession)

	// Attach to the zellij session
	cmd := exec.Command(zm.zellijPath, "attach", zellijSession.id)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DetachFromSession detaches from an existing zellij session
func (zm *ZellijMultiplexer) DetachFromSession(name string) error {
	return fmt.Errorf("detaching from zellij sessions is not supported at the moment")
}

// KillSession terminates a zellij session
func (zm *ZellijMultiplexer) KillSession(name string) error {
	session, err := zm.GetSession(name)
	if err != nil {
		return err
	}

	zellijSession := session.(*ZellijSession)

	// Kill the zellij session
	cmd := exec.Command(zm.zellijPath, "kill-session", zellijSession.id)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill zellij session: %w", err)
	}

	return nil
}

// IsSessionRunning checks if a session is currently running
func (zm *ZellijMultiplexer) IsSessionRunning(name string) bool {
	session, err := zm.GetSession(name)
	if err != nil {
		return false
	}
	return session.IsRunning()
}

// HasSession checks if a session exists
func (zm *ZellijMultiplexer) HasSession(name string) bool {
	sessions, err := zm.ListSessions()
	if err != nil {
		return false
	}

	for _, session := range sessions {
		if session.GetName() == name {
			return true
		}
	}
	return false
}

// KillAllSessions kills all claude-pilot zellij sessions
func (zm *ZellijMultiplexer) KillAllSessions() error {
	sessions, err := zm.ListSessions()
	if err != nil {
		return err
	}

	var errors []string
	for _, session := range sessions {
		if err := zm.KillSession(session.GetName()); err != nil {
			errors = append(errors, fmt.Sprintf("failed to kill session %s: %v", session.GetName(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors killing sessions: %s", strings.Join(errors, "; "))
	}

	return nil
}
