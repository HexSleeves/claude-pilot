package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// SessionTableModel represents an interactive session table
type SessionTableModel struct {
	client *api.Client
	width  int
	height int

	// Table component
	table *components.Table

	// Data
	sessions    []*api.Session
	sessionData []components.SessionData

	// Selection state
	selectedIndex int

	// Scroll state
	viewportStart int
	viewportSize  int
}

// NewSessionTableModel creates a new session table model
func NewSessionTableModel(client *api.Client) *SessionTableModel {
	table := components.NewTable(components.TableConfig{
		ShowHeaders: true,
		Interactive: true,
		MaxRows:     0, // No limit
	})

	return &SessionTableModel{
		client:        client,
		table:         table,
		selectedIndex: 0,
		viewportStart: 0,
	}
}

// Init implements tea.Model
func (m *SessionTableModel) Init() tea.Cmd {
	return m.loadSessions()
}

// Update implements tea.Model
func (m *SessionTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewport()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveUp()
		case "down", "j":
			m.moveDown()
		case "pgup":
			m.pageUp()
		case "pgdown":
			m.pageDown()
		case "home":
			m.goToTop()
		case "end":
			m.goToBottom()
		case "enter", " ":
			if selected := m.GetSelectedSession(); selected != nil {
				return m, func() tea.Msg {
					return SessionSelectedMsg{Session: selected}
				}
			}
		}

	case SessionsLoadedMsg:
		if msg.Error == nil {
			m.SetSessions(msg.Sessions)
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *SessionTableModel) View() string {
	if len(m.sessions) == 0 {
		return styles.DimTextStyle.Render("No sessions found. Press 'c' to create a new session.")
	}

	// Update table selection
	m.table.SetSelectedRow(m.selectedIndex)

	// Configure table for current viewport
	m.table.SetWidth(m.width)
	m.table.SetMaxRows(m.viewportSize)

	// Create viewport data
	viewportData := m.getViewportData()
	m.table.SetSessionData(viewportData)

	// Render table
	tableView := m.table.RenderTUI()

	// Add selection indicators
	return m.addSelectionIndicators(tableView)
}

// SetSessions updates the sessions data
func (m *SessionTableModel) SetSessions(sessions []*api.Session) {
	m.sessions = sessions
	m.sessionData = m.convertToSessionData(sessions)

	// Adjust selection if needed
	if m.selectedIndex >= len(m.sessions) {
		m.selectedIndex = len(m.sessions) - 1
	}
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}

	m.updateViewport()
}

// SetSize updates the table size
func (m *SessionTableModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updateViewport()
}

// GetSelectedSession returns the currently selected session
func (m *SessionTableModel) GetSelectedSession() *api.Session {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.sessions) {
		return m.sessions[m.selectedIndex]
	}
	return nil
}

// GetSelectedIndex returns the currently selected index
func (m *SessionTableModel) GetSelectedIndex() int {
	return m.selectedIndex
}

// Navigation methods

func (m *SessionTableModel) moveUp() {
	if m.selectedIndex > 0 {
		m.selectedIndex--
		m.adjustViewport()
	}
}

func (m *SessionTableModel) moveDown() {
	if m.selectedIndex < len(m.sessions)-1 {
		m.selectedIndex++
		m.adjustViewport()
	}
}

func (m *SessionTableModel) pageUp() {
	m.selectedIndex -= m.viewportSize
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	m.adjustViewport()
}

func (m *SessionTableModel) pageDown() {
	m.selectedIndex += m.viewportSize
	if m.selectedIndex >= len(m.sessions) {
		m.selectedIndex = len(m.sessions) - 1
	}
	m.adjustViewport()
}

func (m *SessionTableModel) goToTop() {
	m.selectedIndex = 0
	m.adjustViewport()
}

func (m *SessionTableModel) goToBottom() {
	if len(m.sessions) > 0 {
		m.selectedIndex = len(m.sessions) - 1
	}
	m.adjustViewport()
}

// Viewport management

func (m *SessionTableModel) updateViewport() {
	// Calculate viewport size (account for headers and padding)
	availableHeight := m.height - 6 // Account for headers, borders, padding
	if availableHeight < 1 {
		availableHeight = 1
	}
	m.viewportSize = availableHeight
	m.adjustViewport()
}

func (m *SessionTableModel) adjustViewport() {
	if len(m.sessions) == 0 {
		m.viewportStart = 0
		return
	}

	// Ensure selected item is visible
	if m.selectedIndex < m.viewportStart {
		m.viewportStart = m.selectedIndex
	} else if m.selectedIndex >= m.viewportStart+m.viewportSize {
		m.viewportStart = m.selectedIndex - m.viewportSize + 1
	}

	// Ensure viewport doesn't go beyond data
	maxStart := len(m.sessions) - m.viewportSize
	if maxStart < 0 {
		maxStart = 0
	}
	if m.viewportStart > maxStart {
		m.viewportStart = maxStart
	}
	if m.viewportStart < 0 {
		m.viewportStart = 0
	}
}

func (m *SessionTableModel) getViewportData() []components.SessionData {
	end := m.viewportStart + m.viewportSize
	if end > len(m.sessionData) {
		end = len(m.sessionData)
	}

	if m.viewportStart >= len(m.sessionData) {
		return []components.SessionData{}
	}

	return m.sessionData[m.viewportStart:end]
}

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

// Styling

func (m *SessionTableModel) addSelectionIndicators(tableView string) string {
	if len(m.sessions) == 0 {
		return tableView
	}

	// Add scroll indicators if needed
	indicators := ""

	if m.viewportStart > 0 {
		indicators += styles.DimTextStyle.Render("↑ More sessions above") + "\n"
	}

	result := indicators + tableView

	if m.viewportStart+m.viewportSize < len(m.sessions) {
		result += "\n" + styles.DimTextStyle.Render("↓ More sessions below")
	}

	// Add status line
	statusLine := fmt.Sprintf("Session %d of %d", m.selectedIndex+1, len(m.sessions))
	if len(m.sessions) > m.viewportSize {
		statusLine += fmt.Sprintf(" (showing %d-%d)",
			m.viewportStart+1,
			min(m.viewportStart+m.viewportSize, len(m.sessions)))
	}

	result += "\n" + styles.DimTextStyle.Render(statusLine)

	return result
}

// Utility functions

func (m *SessionTableModel) loadSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.client.ListSessions()
		return SessionsLoadedMsg{Sessions: sessions, Error: err}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Message types specific to session table
type SessionSelectedMsg struct {
	Session *api.Session
}

type SessionActionMsg struct {
	Action    string // "attach", "kill", "detail"
	SessionID string
}
