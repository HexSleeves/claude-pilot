package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// SessionTableModel represents an interactive session table
type SessionTableModel struct {
	client *api.Client
	width  int
	height int

	// Bubbles table component
	table table.Model

	// Data
	sessions    []*api.Session
	sessionData []components.SessionData

	// Key bindings
	keys KeyMap
}

// NewSessionTableModel creates a new session table model
func NewSessionTableModel(client *api.Client) *SessionTableModel {
	// Create Bubbles table with predefined columns
	columns := components.GetBubblesTableColumns()

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Apply Claude theme styling
	t = styles.ConfigureBubblesTable(t)

	return &SessionTableModel{
		client: client,
		table:  t,
		keys:   DefaultKeyMap(),
	}
}

// Init implements tea.Model
func (m *SessionTableModel) Init() tea.Cmd {
	return m.loadSessions()
}

// Update implements tea.Model
func (m *SessionTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateTableSize()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if selected := m.GetSelectedSession(); selected != nil {
				return m, func() tea.Msg {
					return SessionSelectedMsg{Session: selected}
				}
			}
		default:
			// Let the Bubbles table handle navigation
			m.table, cmd = m.table.Update(msg)
		}

	case SessionsLoadedMsg:
		if msg.Error == nil {
			m.SetSessions(msg.Sessions)
		}
	}

	return m, cmd
}

// View implements tea.Model
func (m *SessionTableModel) View() string {
	if len(m.sessions) == 0 {
		return styles.DimTextStyle.Render("No sessions found. Press 'c' to create a new session.")
	}

	// Render the Bubbles table
	tableView := m.table.View()

	// Add session count info
	selectedIdx := m.table.Cursor()
	statusLine := fmt.Sprintf("Session %d of %d", selectedIdx+1, len(m.sessions))

	return tableView + "\n" + styles.DimTextStyle.Render(statusLine)
}

// SetSessions updates the sessions data
func (m *SessionTableModel) SetSessions(sessions []*api.Session) {
	m.sessions = sessions
	m.sessionData = m.convertToSessionData(sessions)

	// Convert session data to Bubbles table rows
	rows := components.ToBubblesSessionRows(m.sessionData)
	m.table.SetRows(rows)
}

// SetSize updates the table size
func (m *SessionTableModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updateTableSize()
}

// updateTableSize updates the Bubbles table dimensions
func (m *SessionTableModel) updateTableSize() {
	// Calculate available height for the table (account for headers and status line)
	availableHeight := m.height - 4 // Account for borders and status line
	if availableHeight < 3 {
		availableHeight = 3
	}

	// Update table dimensions
	m.table.SetWidth(m.width)
	m.table.SetHeight(availableHeight)
}

// GetSelectedSession returns the currently selected session
func (m *SessionTableModel) GetSelectedSession() *api.Session {
	selectedIdx := m.table.Cursor()
	if selectedIdx >= 0 && selectedIdx < len(m.sessions) {
		return m.sessions[selectedIdx]
	}
	return nil
}

// GetSelectedIndex returns the currently selected index
func (m *SessionTableModel) GetSelectedIndex() int {
	return m.table.Cursor()
}

// Navigation methods are now handled by the Bubbles table component

// Data conversion

func (m *SessionTableModel) convertToSessionData(sessions []*api.Session) []components.SessionData {
	data := make([]components.SessionData, len(sessions))

	for i, session := range sessions {
		// Convert backend status
		backend := m.getBackendDisplay(session)

		data[i] = components.SessionData{
			ID:          session.ID,
			Name:        session.Name,
			Status:      string(session.Status),
			Backend:     backend,
			Created:     session.CreatedAt,
			LastActive:  session.LastActive,
			Messages:    len(session.Messages),
			ProjectPath: session.ProjectPath,
		}
	}

	return data
}

func (m *SessionTableModel) getBackendDisplay(session *api.Session) string {
	backend := m.client.GetBackend()

	switch session.Status {
	case api.StatusConnected:
		return fmt.Sprintf("%s (attached)", backend)
	case api.StatusActive:
		return fmt.Sprintf("%s (running)", backend)
	case api.StatusInactive:
		return fmt.Sprintf("%s (stopped)", backend)
	case api.StatusError:
		return fmt.Sprintf("%s (error)", backend)
	default:
		return backend
	}
}

// Styling is now handled by the Bubbles table component

// Utility functions

func (m *SessionTableModel) loadSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.client.ListSessions()
		return SessionsLoadedMsg{Sessions: sessions, Error: err}
	}
}

// Message types specific to session table
type SessionSelectedMsg struct {
	Session *api.Session
}

type SessionActionMsg struct {
	Action    string // "attach", "kill", "detail"
	SessionID string
}
