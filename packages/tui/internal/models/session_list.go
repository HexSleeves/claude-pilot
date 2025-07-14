package models

import (
	"fmt"
	"strings"

	"claude-pilot/core/api"
	"claude-pilot/tui/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
)

// SessionListModel handles the session list view
type SessionListModel struct {
	client   *api.Client
	sessions []*api.Session
	cursor   int
	width    int
	height   int
	loading  bool
	err      error
}

// SessionSelectedMsg is sent when a session is selected
type SessionSelectedMsg struct {
	Session *api.Session
}

// SessionsLoadedMsg is sent when sessions are loaded
type SessionsLoadedMsg struct {
	Sessions []*api.Session
	Error    error
}

// NewSessionListModel creates a new session list model
func NewSessionListModel(client *api.Client) *SessionListModel {
	return &SessionListModel{
		client:  client,
		loading: true,
	}
}

// Init implements tea.Model
func (m *SessionListModel) Init() tea.Cmd {
	return m.loadSessions()
}

// Update implements tea.Model
func (m *SessionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SessionsLoadedMsg:
		m.loading = false
		m.sessions = msg.Sessions
		m.err = msg.Error
		if len(m.sessions) > 0 && m.cursor >= len(m.sessions) {
			m.cursor = len(m.sessions) - 1
		}

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.sessions) > 0 && m.cursor < len(m.sessions) {
				return m, func() tea.Msg {
					return SessionSelectedMsg{Session: m.sessions[m.cursor]}
				}
			}

		case "r":
			// Refresh sessions
			m.loading = true
			return m, m.loadSessions()
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *SessionListModel) View() string {
	if m.loading {
		return styles.InfoStyle.Render("Loading sessions...")
	}

	if m.err != nil {
		return styles.ErrorStyle.Render(fmt.Sprintf("Error loading sessions: %v", m.err))
	}

	if len(m.sessions) == 0 {
		return styles.MutedTextStyle.Render("No sessions found. Press 'tab' to create a new session.")
	}

	var content strings.Builder

	// Header
	content.WriteString(styles.HeaderStyle.Render("Sessions"))
	content.WriteString("\n\n")

	// Session list
	for i, session := range m.sessions {
		var line strings.Builder

		// Cursor indicator
		if i == m.cursor {
			line.WriteString(styles.SelectedStyle.Render("> "))
		} else {
			line.WriteString("  ")
		}

		// Session name
		nameStyle := styles.SessionNameStyle
		if i == m.cursor {
			nameStyle = styles.SelectedStyle
		}
		line.WriteString(nameStyle.Render(session.Name))

		// Session status
		line.WriteString(" ")
		statusStyle := styles.FormatSessionStatus(string(session.Status))
		if i == m.cursor {
			statusStyle = styles.SelectedStyle
		}
		line.WriteString(statusStyle.Render(fmt.Sprintf("[%s]", session.Status)))

		// Session ID (truncated)
		line.WriteString(" ")
		idStyle := styles.SessionIDStyle
		if i == m.cursor {
			idStyle = styles.SelectedStyle
		}
		truncatedID := styles.TruncateText(session.ID, 8)
		line.WriteString(idStyle.Render(fmt.Sprintf("(%s...)", truncatedID)))

		// Project path (truncated)
		if session.ProjectPath != "" {
			line.WriteString(" ")
			pathStyle := styles.MutedTextStyle
			if i == m.cursor {
				pathStyle = styles.SelectedStyle
			}
			truncatedPath := styles.TruncateText(session.ProjectPath, 30)
			line.WriteString(pathStyle.Render(truncatedPath))
		}

		content.WriteString(line.String())
		content.WriteString("\n")
	}

	// Instructions
	content.WriteString("\n")
	content.WriteString(styles.MutedTextStyle.Render("↑/↓: navigate • enter: details • r: refresh • tab: create"))

	return content.String()
}

// SetSize updates the model's size
func (m *SessionListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSelectedSession returns the currently selected session
func (m *SessionListModel) GetSelectedSession() *api.Session {
	if len(m.sessions) == 0 || m.cursor >= len(m.sessions) {
		return nil
	}
	return m.sessions[m.cursor]
}

// RefreshSessions reloads the session list
func (m *SessionListModel) RefreshSessions() tea.Cmd {
	m.loading = true
	return m.loadSessions()
}

// loadSessions loads sessions from the API
func (m *SessionListModel) loadSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.client.ListSessions()
		return SessionsLoadedMsg{
			Sessions: sessions,
			Error:    err,
		}
	}
}
