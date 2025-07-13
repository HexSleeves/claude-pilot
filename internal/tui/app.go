package tui

import (
	"fmt"

	"claude-pilot/internal/interfaces"
	"claude-pilot/internal/manager"

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

// Model represents the main TUI application model
type Model struct {
	sessionManager *manager.SessionManager
	state          AppState
	width          int
	height         int

	// Components
	sessionTable table.Model
	sessionList  list.Model

	// Data
	sessions        []*interfaces.Session
	selectedSession *interfaces.Session

	// Styling
	baseStyle    lipgloss.Style
	titleStyle   lipgloss.Style
	errorStyle   lipgloss.Style
	successStyle lipgloss.Style
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
func NewModel(sessionManager *manager.SessionManager) *Model {
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

	return &Model{
		sessionManager: sessionManager,
		state:          StateSessionList,
		sessionTable:   t,
		baseStyle:      baseStyle,
		titleStyle:     titleStyle,
		errorStyle:     errorStyle,
		successStyle:   successStyle,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.loadSessions
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sessionTable.SetWidth(msg.Width - 4)
		m.sessionTable.SetHeight(msg.Height - 8)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Refresh):
			return m, m.loadSessions

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
		m.updateTable()

	case sessionDeletedMsg:
		return m, m.loadSessions

	case sessionAttachedMsg:
		// Attachment successful, exit TUI
		return m, tea.Quit
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
	if len(m.sessions) == 0 {
		content = "No sessions found. Press 'c' to create a new session."
	} else {
		content = m.baseStyle.Render(m.sessionTable.View())
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press ? for help, q to quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		help,
	)
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
		if session.Status == interfaces.StatusActive {
			status = "● Active"
		} else if session.Status == interfaces.StatusInactive {
			status = "⏸ Inactive"
		} else if session.Status == interfaces.StatusConnected {
			status = "🔗 Connected"
		}

		rows[i] = table.Row{
			session.Name,
			status,
			session.CreatedAt.Format("2006-01-02 15:04"),
			session.ProjectPath,
		}
	}

	m.sessionTable.SetRows(rows)
}

// Commands for handling async operations

type sessionsLoadedMsg struct {
	sessions []*interfaces.Session
}

type sessionDeletedMsg struct {
	sessionID string
}

type sessionAttachedMsg struct {
	sessionName string
}

func (m Model) loadSessions() tea.Msg {
	sessions, err := m.sessionManager.ListSessions()
	if err != nil {
		// Handle error - for now, return empty list
		return sessionsLoadedMsg{sessions: []*interfaces.Session{}}
	}
	return sessionsLoadedMsg{sessions: sessions}
}

func (m Model) deleteSession(sessionID string) tea.Cmd {
	return func() tea.Msg {
		err := m.sessionManager.DeleteSession(sessionID)
		if err != nil {
			// Handle error
			return nil
		}
		return sessionDeletedMsg{sessionID: sessionID}
	}
}

func (m Model) attachToSession(sessionName string) tea.Cmd {
	return func() tea.Msg {
		err := m.sessionManager.AttachToSession(sessionName)
		if err != nil {
			// Handle error
			return nil
		}
		return sessionAttachedMsg{sessionName: sessionName}
	}
}

// RunTUI starts the TUI application
func RunTUI(sessionManager *manager.SessionManager) error {
	model := NewModel(sessionManager)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
