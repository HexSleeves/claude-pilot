package tui

import (
	"fmt"
	"os"
	"runtime/debug"

	"claude-pilot/internal/interfaces"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppState represents the current state of the TUI application
type AppState int

const (
	StateSessionList AppState = iota
	StateSessionDetail
	StateSessionCreate
	StateHelp
)

// LoadingState represents the loading state of async operations
type LoadingState int

const (
	LoadingIdle LoadingState = iota
	LoadingInProgress
	LoadingSuccess
	LoadingError
)

// Model represents the main TUI application model
type Model struct {
	service interfaces.SessionService
	state   AppState
	width   int
	height  int

	// Components
	sessionTable table.Model
	sessionList  list.Model

	// Data
	sessions        []*interfaces.Session
	selectedSession *interfaces.Session

	// Loading and error states
	loadingState  LoadingState
	errorMessage  string
	statusMessage string

	// Styling
	baseStyle    lipgloss.Style
	titleStyle   lipgloss.Style
	errorStyle   lipgloss.Style
	successStyle lipgloss.Style
	loadingStyle lipgloss.Style
}

// keyMap defines the key bindings for the application
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Space   key.Binding
	Create  key.Binding
	Attach  key.Binding
	Delete  key.Binding
	Refresh key.Binding
	Help    key.Binding
	Quit    key.Binding
	Escape  key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create session"),
		),
		Attach: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "attach to session"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete session"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "escape"),
		),
	}
}

var keys = DefaultKeyMap()

// NewModel creates a new TUI model
func NewModel(service interfaces.SessionService) *Model {
	// Create table columns for session display
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Created", Width: 20},
		{Title: "Project", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Define styles
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B35")).
		Bold(true).
		Padding(0, 1)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E74C3C"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ECC71"))

	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498DB"))

	return &Model{
		service:      service,
		state:        StateSessionList,
		sessionTable: t,
		loadingState: LoadingIdle,
		baseStyle:    baseStyle,
		titleStyle:   titleStyle,
		errorStyle:   errorStyle,
		successStyle: successStyle,
		loadingStyle: loadingStyle,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	m.loadingState = LoadingInProgress
	m.statusMessage = "Loading sessions..."
	return m.loadSessionsCmd()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defer func() {
		if r := recover(); r != nil {
			// Log panic but don't crash the application
			fmt.Fprintf(os.Stderr, "Panic in Update method: %v\n", r)
			// Return the model in a safe state
		}
	}()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update table dimensions based on available space
		m.updateTableDimensions()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Refresh):
			m.loadingState = LoadingInProgress
			m.statusMessage = "Refreshing sessions..."
			return m, m.loadSessionsCmd()

		case key.Matches(msg, keys.Create):
			// TODO: Implement session creation
			return m, nil

		case key.Matches(msg, keys.Attach):
			if m.selectedSession != nil {
				return m, m.attachToSession(m.selectedSession.Name)
			}

		case key.Matches(msg, keys.Delete):
			if m.selectedSession != nil {
				return m, m.deleteSession(m.selectedSession.ID)
			}

		case key.Matches(msg, keys.Escape):
			if m.state == StateHelp {
				m.state = StateSessionList
			}
			return m, nil

		case key.Matches(msg, keys.Help):
			if m.state == StateHelp {
				m.state = StateSessionList
			} else {
				m.state = StateHelp
			}
			return m, nil
		}

	case sessionsLoadedMsg:
		m.sessions = msg.sessions
		m.loadingState = LoadingSuccess
		m.statusMessage = "Sessions loaded successfully"
		if msg.err != nil {
			m.loadingState = LoadingError
			m.errorMessage = msg.err.Error()
			m.statusMessage = "Failed to load sessions"
		}
		m.updateTable()
		m.updateTableDimensions()

	case sessionDeletedMsg:
		m.loadingState = LoadingInProgress
		m.statusMessage = "Refreshing sessions..."
		return m, m.loadSessionsCmd()

	case sessionAttachedMsg:
		// Attachment successful, exit TUI
		return m, tea.Quit

	case sessionErrorMsg:
		m.loadingState = LoadingError
		m.errorMessage = fmt.Sprintf("Failed to %s session: %v", msg.operation, msg.err)
		m.statusMessage = fmt.Sprintf("Error during %s operation", msg.operation)
		return m, nil
	}

	// Update table
	m.sessionTable, cmd = m.sessionTable.Update(msg)

	// Update selected session based on table cursor
	if len(m.sessions) > 0 && m.sessionTable.Cursor() < len(m.sessions) {
		m.selectedSession = m.sessions[m.sessionTable.Cursor()]
	}

	return m, cmd
}

// View renders the model
func (m Model) View() string {
	defer func() {
		if r := recover(); r != nil {
			// Return a safe error message instead of panicking
			fmt.Fprintf(os.Stderr, "Panic in View rendering: %v\n", r)
		}
	}()

	switch m.state {
	case StateHelp:
		return m.helpView()
	case StateSessionList:
		return m.sessionListView()
	default:
		return m.sessionListView()
	}
}

// sessionListView renders the main session list view
func (m Model) sessionListView() string {
	title := m.titleStyle.Render("Claude Pilot - Session Manager")

	var content string
	var statusBar string
	var errorMessage string

	// Show loading state
	if m.loadingState == LoadingInProgress {
		content = m.loadingStyle.Render("🔄 Loading sessions...")
	} else if len(m.sessions) == 0 {
		if m.loadingState == LoadingError {
			content = m.errorStyle.Render("❌ " + m.errorMessage)
		} else {
			content = "No sessions found. Press 'c' to create a new session."
		}
	} else {
		content = m.baseStyle.Render(m.sessionTable.View())
	}

	// Status bar with current state
	switch m.loadingState {
	case LoadingInProgress:
		statusBar = m.loadingStyle.Render(m.statusMessage)
	case LoadingSuccess:
		statusBar = m.successStyle.Render(m.statusMessage)
	case LoadingError:
		statusBar = m.errorStyle.Render(m.statusMessage)
	default:
		statusBar = ""
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press ? for help, q to quit")

	// Build layout responsively
	var parts []string
	parts = append(parts, title, "")

	// Add content with proper width constraints
	if m.width > 0 {
		contentStyle := lipgloss.NewStyle().Width(m.width)
		content = contentStyle.Render(content)
	}
	parts = append(parts, content)

	// Add error message if present
	if errorMessage != "" {
		parts = append(parts, "", errorMessage)
	}

	// Add status bar if present
	if statusBar != "" {
		if m.width > 0 {
			statusStyle := lipgloss.NewStyle().Width(m.width)
			statusBar = statusStyle.Render(statusBar)
		}
		parts = append(parts, "", statusBar)
	}

	// Add help text
	if m.width > 0 {
		helpStyle := lipgloss.NewStyle().Width(m.width)
		help = helpStyle.Render(help)
	}
	parts = append(parts, "", help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// helpView renders the help screen
func (m Model) helpView() string {
	title := m.titleStyle.Render("Claude Pilot - Help")

	helpText := `
Key Bindings:
  ↑/k, ↓/j    Navigate sessions
  enter       View session details
  c           Create new session
  a           Attach to selected session
  d           Delete selected session
  r           Refresh session list
  ?           Toggle this help
  q           Quit application

Session Status:
  ● Active    Session is running
  ⏸ Inactive  Session exists but not running
  🔗 Connected Someone is attached to session

Press ? again to return to session list.
`

	// Apply width constraint if terminal size is available
	if m.width > 0 {
		contentStyle := lipgloss.NewStyle().
			Width(m.width).
			Padding(0, 1)
		helpText = contentStyle.Render(helpText)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		helpText,
	)
}

// updateTable updates the table with current session data
func (m *Model) updateTable() {
	rows := make([]table.Row, len(m.sessions))

	for i, session := range m.sessions {
		status := string(session.Status)
		switch session.Status {
		case interfaces.StatusActive:
			status = "● Active"
		case interfaces.StatusInactive:
			status = "⏸ Inactive"
		case interfaces.StatusConnected:
			status = "🔗 Connected"
		}

		rows[i] = table.Row{
			session.Name,
			status,
			session.ProjectPath,
			session.CreatedAt.Format("2006-01-02 15:04"),
		}
	}

	m.sessionTable.SetRows(rows)
}

// updateTableDimensions updates the dimensions of the table based on available space
func (m *Model) updateTableDimensions() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	// Calculate dimensions based on actual component sizes
	title := m.titleStyle.Render("Claude Pilot - Session Manager")
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press ? for help, q to quit")

	// Calculate used height for fixed components
	usedHeight := lipgloss.Height(title) + 2 // title + empty lines
	usedHeight += lipgloss.Height(help) + 1  // help + empty line

	// Add status bar height if present
	if m.loadingState != LoadingIdle {
		statusBar := m.loadingStyle.Render("Status")
		usedHeight += lipgloss.Height(statusBar) + 1 // status + empty line
	}

	// Calculate available height for content
	availableHeight := max(m.height-usedHeight,
		// Minimum height
		3)

	// Set table dimensions with proper padding
	tableWidth := max(
		// Account for border padding
		m.width-4,
		// Minimum width
		10)

	m.sessionTable.SetWidth(tableWidth)
	m.sessionTable.SetHeight(availableHeight)
}

// Commands for handling async operations

type sessionsLoadedMsg struct {
	sessions []*interfaces.Session
	err      error
}

type sessionDeletedMsg struct {
	sessionID string
}

type sessionAttachedMsg struct {
	sessionName string
}

type sessionErrorMsg struct {
	operation string
	err       error
}

func (m Model) loadSessionsCmd() tea.Cmd {
	return func() tea.Msg {
		defer func() {
			if r := recover(); r != nil {
				// Panics in commands are handled by returning an error message
				// The function will return the error below
			}
		}()

		sessions, err := m.service.ListSessions()
		if err != nil {
			return sessionsLoadedMsg{sessions: []*interfaces.Session{}, err: err}
		}
		return sessionsLoadedMsg{sessions: sessions, err: nil}
	}
}

func (m Model) deleteSession(sessionID string) tea.Cmd {
	return func() tea.Msg {
		defer func() {
			if r := recover(); r != nil {
				// Panics in commands are handled by returning an error message
				// The function will return the error below
			}
		}()

		err := m.service.DeleteSession(sessionID)
		if err != nil {
			return sessionErrorMsg{
				operation: "delete",
				err:       err,
			}
		}
		return sessionDeletedMsg{sessionID: sessionID}
	}
}

func (m Model) attachToSession(sessionName string) tea.Cmd {
	return func() tea.Msg {
		defer func() {
			if r := recover(); r != nil {
				// Panics in commands are handled by returning an error message
				// The function will return the error below
			}
		}()

		err := m.service.AttachToSession(sessionName)
		if err != nil {
			return sessionErrorMsg{
				operation: "attach",
				err:       err,
			}
		}
		return sessionAttachedMsg{sessionName: sessionName}
	}
}

// RunTUI starts the TUI application
func RunTUI(service interfaces.SessionService) error {
	// Set up panic recovery to prevent terminal corruption
	defer func() {
		if r := recover(); r != nil {
			// Try to restore terminal state
			fmt.Print("\033[?25h") // Show cursor
			fmt.Print("\033[0m")   // Reset colors
			fmt.Print("\033[2J")   // Clear screen
			fmt.Print("\033[H")    // Move cursor to top-left

			// Print panic information
			fmt.Fprintf(os.Stderr, "\nPanic in TUI application: %v\n", r)
			fmt.Fprintf(os.Stderr, "Terminal state has been restored. You may need to run 'reset' if display issues persist.\n")

			// Print stack trace for debugging
			debug.PrintStack()

			os.Exit(1)
		}
	}()

	model := NewModel(service)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		// Ensure terminal is restored on error
		fmt.Print("\033[?25h") // Show cursor
		fmt.Print("\033[0m")   // Reset colors
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
