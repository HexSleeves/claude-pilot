package components

import (
	"claude-pilot/internal/interfaces"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionListModel handles the session table display and operations
type SessionListModel struct {
	table           table.Model
	sessions        []*interfaces.Session
	selectedSession *interfaces.Session
	width           int
	height          int

	// Styling
	baseStyle lipgloss.Style
}

// NewSessionListModel creates a new session list model
func NewSessionListModel() *SessionListModel {
	// Create table columns for session display
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Project", Width: 30},
		{Title: "Created", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	return &SessionListModel{
		table:     t,
		baseStyle: baseStyle,
	}
}

// Init initializes the session list model
func (m SessionListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the session list
func (m SessionListModel) Update(msg tea.Msg) (SessionListModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateDimensions()
	}

	// Update table
	m.table, cmd = m.table.Update(msg)

	// Update selected session based on table cursor
	if len(m.sessions) > 0 && m.table.Cursor() < len(m.sessions) {
		m.selectedSession = m.sessions[m.table.Cursor()]
	}

	return m, cmd
}

// View renders the session list
func (m SessionListModel) View() string {
	if len(m.sessions) == 0 {
		return "No sessions found. Press 'c' to create a new session."
	}
	return m.baseStyle.Render(m.table.View())
}

// SetSessions updates the sessions data
func (m *SessionListModel) SetSessions(sessions []*interfaces.Session) {
	m.sessions = sessions
	m.updateTable()
}

// GetSelectedSession returns the currently selected session
func (m SessionListModel) GetSelectedSession() *interfaces.Session {
	return m.selectedSession
}

// updateTable updates the table with current session data
func (m *SessionListModel) updateTable() {
	rows := make([]table.Row, len(m.sessions))

	for i, session := range m.sessions {
		var status string
		switch session.Status {
		case interfaces.StatusActive:
			status = "● Active"
		case interfaces.StatusInactive:
			status = "⏸ Inactive"
		case interfaces.StatusConnected:
			status = "🔗 Connected"
		default:
			status = string(session.Status)
		}

		rows[i] = table.Row{
			session.Name,
			status,
			session.ProjectPath,
			session.CreatedAt.Format("2006-01-02 15:04"),
		}
	}

	m.table.SetRows(rows)
}

// updateDimensions updates the table dimensions
func (m *SessionListModel) updateDimensions() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	// Calculate available space for the table
	availableHeight := m.height - 2 // Account for some padding
	if availableHeight < 3 {
		availableHeight = 3
	}

	tableWidth := max(m.width-4, 10) // Account for border padding

	m.table.SetWidth(tableWidth)
	m.table.SetHeight(availableHeight)
}

// SetDimensions sets the component dimensions
func (m *SessionListModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.updateDimensions()
}
