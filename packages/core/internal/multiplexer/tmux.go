package multiplexer

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"claude-pilot/core/internal/logger"
	"claude-pilot/shared/interfaces"
)

// TmuxMultiplexer implements the TerminalMultiplexer interface for tmux
type TmuxMultiplexer struct {
	sessionPrefix string
	tmuxPath      string
	logger        *logger.Logger
}

// TmuxSession implements the MultiplexerSession interface
type TmuxSession struct {
	id          string
	name        string
	tmuxName    string
	status      interfaces.SessionStatus
	createdAt   time.Time
	workingDir  string
	description string
	isAttached  bool
	isRunning   bool
}

// GetID returns the session ID
func (s *TmuxSession) GetID() string {
	return s.id
}

// GetName returns the session name
func (s *TmuxSession) GetName() string {
	return s.name
}

// GetStatus returns the session status
func (s *TmuxSession) GetStatus() interfaces.SessionStatus {
	return s.status
}

// GetCreatedAt returns the creation time
func (s *TmuxSession) GetCreatedAt() time.Time {
	return s.createdAt
}

// GetWorkingDir returns the working directory
func (s *TmuxSession) GetWorkingDir() string {
	return s.workingDir
}

// GetDescription returns the session description
func (s *TmuxSession) GetDescription() string {
	return s.description
}

// IsAttached returns whether someone is attached to the session
func (s *TmuxSession) IsAttached() bool {
	return s.isAttached
}

// IsRunning returns whether the session is running
func (s *TmuxSession) IsRunning() bool {
	return s.isRunning
}

// NewTmuxMultiplexer creates a new tmux multiplexer instance
func NewTmuxMultiplexer(sessionPrefix string) (*TmuxMultiplexer, error) {
	// Create a disabled logger by default for backward compatibility
	disabledLogger, _ := logger.Setup.Disabled().Build()
	return NewTmuxMultiplexerWithLogger(sessionPrefix, disabledLogger)
}

// NewTmuxMultiplexerWithLogger creates a new tmux multiplexer instance with logger
func NewTmuxMultiplexerWithLogger(sessionPrefix string, log *logger.Logger) (*TmuxMultiplexer, error) {
	if sessionPrefix == "" {
		sessionPrefix = "claude-pilot"
	}

	log.Debug("Initializing tmux multiplexer", "session_prefix", sessionPrefix)

	// Try to find tmux binary
	tmuxPath := "tmux" // default
	if _, err := exec.LookPath("tmux"); err != nil {
		log.Debug("tmux not found in PATH, trying common locations")
		// Try common locations
		paths := []string{
			"/opt/homebrew/bin/tmux",
			"/usr/local/bin/tmux",
			"/usr/bin/tmux",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				tmuxPath = path
				log.Debug("Found tmux binary", "path", path)
				break
			}
		}
	} else {
		log.Debug("Found tmux in PATH")
	}

	tm := &TmuxMultiplexer{
		sessionPrefix: sessionPrefix,
		tmuxPath:      tmuxPath,
		logger:        log,
	}

	log.Info("Tmux multiplexer initialized",
		"session_prefix", sessionPrefix,
		"tmux_path", tmuxPath)

	return tm, nil
}

// GetName returns the name of the multiplexer backend
func (tm *TmuxMultiplexer) GetName() string {
	return "tmux"
}

// IsAvailable checks if tmux is available on the system
func (tm *TmuxMultiplexer) IsAvailable() bool {
	_, err := exec.LookPath("tmux")
	if err != nil {
		// Check common paths
		paths := []string{
			"/opt/homebrew/bin/tmux",
			"/usr/local/bin/tmux",
			"/usr/bin/tmux",
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

// CreateSession creates a new tmux session, or attaches to existing session as pane/window
func (tm *TmuxMultiplexer) CreateSession(req interfaces.CreateSessionRequest) (interfaces.MultiplexerSession, error) {
	start := time.Now()

	// Handle attachment to existing session
	if req.AttachTo != "" && req.AttachmentType != interfaces.AttachmentNone {
		return tm.createAttachedSession(req)
	}

	// Create standalone session (original behavior)
	return tm.createStandaloneSession(req, start)
}

// createStandaloneSession creates a new standalone tmux session
func (tm *TmuxMultiplexer) createStandaloneSession(req interfaces.CreateSessionRequest, start time.Time) (interfaces.MultiplexerSession, error) {
	tmuxName := fmt.Sprintf("%s-%s", tm.sessionPrefix, req.Name)

	tm.logger.Debug("Creating tmux session",
		"name", req.Name,
		"tmux_name", tmuxName,
		"command", req.Command,
		"working_dir", req.WorkingDir)

	// Check if tmux session already exists
	if tm.HasSession(req.Name) {
		tm.logger.Warn("Tmux session creation failed: already exists",
			"name", req.Name,
			"tmux_name", tmuxName)
		return nil, fmt.Errorf("tmux session '%s' already exists", req.Name)
	}

	// Determine command to run (default to "claude")
	command := req.Command
	if command == "" {
		command = "claude"
	}

	// Create tmux session with specified command
	var cmd *exec.Cmd
	if req.WorkingDir != "" {
		// Create session in specific directory
		cmd = exec.Command(tm.tmuxPath, "new-session", "-d", "-s", tmuxName, "-c", req.WorkingDir, command)
	} else {
		// Create session in current directory
		cmd = exec.Command(tm.tmuxPath, "new-session", "-d", "-s", tmuxName, command)
	}

	tm.logger.DebugCommand(tm.tmuxPath, cmd.Args[1:], req.WorkingDir)

	if err := cmd.Run(); err != nil {
		tm.logger.Error("Failed to create tmux session",
			"name", req.Name,
			"tmux_name", tmuxName,
			"command", strings.Join(cmd.Args, " "),
			"error", err)
		return nil, fmt.Errorf("failed to create tmux session: %w", err)
	}

	session := &TmuxSession{
		id:          tmuxName,
		name:        req.Name,
		tmuxName:    tmuxName,
		status:      interfaces.StatusActive,
		createdAt:   time.Now(),
		workingDir:  req.WorkingDir,
		description: req.Description,
		isAttached:  false,
		isRunning:   true,
	}

	tm.logger.Performance("CreateSession", start,
		slog.String("name", req.Name),
		slog.String("tmux_name", tmuxName))

	tm.logger.Info("Tmux session created successfully",
		"name", req.Name,
		"tmux_name", tmuxName,
		"command", command)

	return session, nil
}

// createAttachedSession creates a new pane or window in an existing session
func (tm *TmuxMultiplexer) createAttachedSession(req interfaces.CreateSessionRequest) (interfaces.MultiplexerSession, error) {
	targetTmuxName := fmt.Sprintf("%s-%s", tm.sessionPrefix, req.AttachTo)

	tm.logger.Debug("Creating attached session",
		"name", req.Name,
		"attach_to", req.AttachTo,
		"target_tmux_name", targetTmuxName,
		"attachment_type", req.AttachmentType,
		"split_direction", req.SplitDirection,
		"command", req.Command,
		"working_dir", req.WorkingDir)

	// Verify target session exists
	if !tm.HasSession(req.AttachTo) {
		return nil, fmt.Errorf("target session '%s' does not exist", req.AttachTo)
	}

	// Determine command to run (default to "claude")
	command := req.Command
	if command == "" {
		command = "claude"
	}

	var cmd *exec.Cmd
	var err error

	switch req.AttachmentType {
	case interfaces.AttachmentPane:
		cmd, err = tm.buildSplitPaneCommand(targetTmuxName, req.WorkingDir, req.SplitDirection, command)
	case interfaces.AttachmentWindow:
		cmd, err = tm.buildNewWindowCommand(targetTmuxName, req.WorkingDir, req.Name, command)
	default:
		return nil, fmt.Errorf("unsupported attachment type: %s", req.AttachmentType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build tmux command: %w", err)
	}

	tm.logger.DebugCommand(tm.tmuxPath, cmd.Args[1:], req.WorkingDir)

	if err := cmd.Run(); err != nil {
		tm.logger.Error("Failed to create attached session",
			"name", req.Name,
			"attach_to", req.AttachTo,
			"attachment_type", req.AttachmentType,
			"command", strings.Join(cmd.Args, " "),
			"error", err)
		return nil, fmt.Errorf("failed to create attached session: %w", err)
	}

	// Create a virtual session object representing the attachment
	// Note: This doesn't create a separate tmux session, but represents the pane/window
	attachedSession := &TmuxSession{
		id:          fmt.Sprintf("%s-%s-attached", targetTmuxName, req.Name),
		name:        fmt.Sprintf("%s-attached-to-%s", req.Name, req.AttachTo),
		tmuxName:    targetTmuxName, // Points to the parent session
		status:      interfaces.StatusActive,
		createdAt:   time.Now(),
		workingDir:  req.WorkingDir,
		description: fmt.Sprintf("%s (attached as %s to %s)", req.Description, req.AttachmentType, req.AttachTo),
		isAttached:  false,
		isRunning:   true,
	}

	tm.logger.Info("Attached session created successfully",
		"name", req.Name,
		"attach_to", req.AttachTo,
		"attachment_type", req.AttachmentType,
		"command", command)

	return attachedSession, nil
}

// buildSplitPaneCommand builds the tmux command for creating a new pane
func (tm *TmuxMultiplexer) buildSplitPaneCommand(targetSession, workingDir string, splitDir interfaces.SplitDirection, command string) (*exec.Cmd, error) {
	args := []string{"split-window", "-t", targetSession}

	// Add split direction
	if splitDir == interfaces.SplitHorizontal {
		args = append(args, "-v") // tmux -v means horizontal split (top/bottom)
	} else {
		args = append(args, "-h") // tmux -h means vertical split (left/right)
	}

	// Add working directory if specified
	if workingDir != "" {
		args = append(args, "-c", workingDir)
	}

	// Add command
	args = append(args, command)

	return exec.Command(tm.tmuxPath, args...), nil
}

// buildNewWindowCommand builds the tmux command for creating a new window
func (tm *TmuxMultiplexer) buildNewWindowCommand(targetSession, workingDir, windowName, command string) (*exec.Cmd, error) {
	args := []string{"new-window", "-t", targetSession}

	// Add window name if provided
	if windowName != "" {
		args = append(args, "-n", windowName)
	}

	// Add working directory if specified
	if workingDir != "" {
		args = append(args, "-c", workingDir)
	}

	// Add command
	args = append(args, command)

	return exec.Command(tm.tmuxPath, args...), nil
}

// GetSession retrieves session information by name
func (tm *TmuxMultiplexer) GetSession(name string) (interfaces.MultiplexerSession, error) {
	sessions, err := tm.ListSessions()
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.GetName() == name {
			return session, nil
		}
	}

	return nil, fmt.Errorf("tmux session '%s' not found", name)
}

// ListSessions returns all available tmux sessions
func (tm *TmuxMultiplexer) ListSessions() ([]interfaces.MultiplexerSession, error) {
	// Get list of tmux sessions
	cmd := exec.Command(tm.tmuxPath, "list-sessions", "-F", "#{session_name},#{session_created},#{session_attached}")
	output, err := cmd.Output()
	if err != nil {
		// If no sessions exist, tmux returns exit code 1
		if strings.Contains(err.Error(), "no server running") || strings.Contains(err.Error(), "no sessions") {
			return []interfaces.MultiplexerSession{}, nil
		}
		return nil, fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	var sessions []interfaces.MultiplexerSession
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}

		sessionName := parts[0]

		// Only include our claude-pilot sessions
		if !strings.HasPrefix(sessionName, tm.sessionPrefix+"-") {
			continue
		}

		// Extract the user-friendly name
		name := strings.TrimPrefix(sessionName, tm.sessionPrefix+"-")

		// Parse created time (unix timestamp)
		var createdAt time.Time
		if createdTime := parts[1]; createdTime != "" {
			if timestamp, err := strconv.ParseInt(createdTime, 10, 64); err == nil {
				createdAt = time.Unix(timestamp, 0)
			} else {
				createdAt = time.Now() // fallback
			}
		}

		// Parse attached status
		isAttached := parts[2] == "1"

		session := &TmuxSession{
			id:         sessionName,
			name:       name,
			tmuxName:   sessionName,
			status:     interfaces.StatusActive,
			createdAt:  createdAt,
			isAttached: isAttached,
			isRunning:  true,
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// AttachToSession attaches to an existing tmux session
func (tm *TmuxMultiplexer) AttachToSession(name string) error {
	session, err := tm.GetSession(name)
	if err != nil {
		return err
	}

	tmuxSession := session.(*TmuxSession)

	// Attach to the tmux session
	cmd := exec.Command(tm.tmuxPath, "attach-session", "-t", tmuxSession.tmuxName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// KillSession terminates a tmux session
func (tm *TmuxMultiplexer) KillSession(name string) error {
	session, err := tm.GetSession(name)
	if err != nil {
		return err
	}

	tmuxSession := session.(*TmuxSession)

	// Kill the tmux session
	cmd := exec.Command(tm.tmuxPath, "kill-session", "-t", tmuxSession.tmuxName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill tmux session: %w", err)
	}

	return nil
}

// IsSessionRunning checks if a session is currently running
func (tm *TmuxMultiplexer) IsSessionRunning(name string) bool {
	session, err := tm.GetSession(name)
	if err != nil {
		return false
	}
	return session.IsRunning()
}

// HasSession checks if a session exists
func (tm *TmuxMultiplexer) HasSession(name string) bool {
	tmuxName := fmt.Sprintf("%s-%s", tm.sessionPrefix, name)
	cmd := exec.Command(tm.tmuxPath, "has-session", "-t", tmuxName)
	return cmd.Run() == nil
}

// GetTmuxSessionInfo gets detailed info about a tmux session (legacy compatibility)
func (tm *TmuxMultiplexer) GetTmuxSessionInfo(name string) (map[string]string, error) {
	session, err := tm.GetSession(name)
	if err != nil {
		return nil, err
	}

	tmuxSession := session.(*TmuxSession)

	// Get detailed session info
	cmd := exec.Command(tm.tmuxPath, "display-message", "-t", tmuxSession.tmuxName, "-p",
		"#{session_name},#{session_created},#{session_attached},#{session_windows},#{session_activity}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get session info: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(parts) < 5 {
		return nil, fmt.Errorf("unexpected tmux output format")
	}

	info := map[string]string{
		"name":     parts[0],
		"created":  parts[1],
		"attached": parts[2],
		"windows":  parts[3],
		"activity": parts[4],
	}

	return info, nil
}

// KillAllSessions kills all claude-pilot tmux sessions
func (tm *TmuxMultiplexer) KillAllSessions() error {
	sessions, err := tm.ListSessions()
	if err != nil {
		return err
	}

	var errors []string
	for _, session := range sessions {
		if err := tm.KillSession(session.GetName()); err != nil {
			errors = append(errors, fmt.Sprintf("failed to kill session %s: %v", session.GetName(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors killing sessions: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetSessionPaneCount returns the number of panes in a tmux session
func (tm *TmuxMultiplexer) GetSessionPaneCount(name string) (int, error) {
	tmuxName := fmt.Sprintf("%s-%s", tm.sessionPrefix, name)

	// Check if session exists first
	if !tm.HasSession(name) {
		return 0, fmt.Errorf("session '%s' not found", name)
	}

	// Get pane count using tmux list-panes command
	cmd := exec.Command(tm.tmuxPath, "list-panes", "-t", tmuxName, "-F", "#{pane_id}")
	output, err := cmd.Output()
	if err != nil {
		tm.logger.Error("Failed to get pane count for session",
			"name", name,
			"tmux_name", tmuxName,
			"error", err)
		return 0, fmt.Errorf("failed to get pane count: %w", err)
	}

	// Count the number of lines (each line represents a pane)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}

	paneCount := len(lines)
	tm.logger.Debug("Retrieved pane count for session",
		"name", name,
		"tmux_name", tmuxName,
		"pane_count", paneCount)

	return paneCount, nil
}
