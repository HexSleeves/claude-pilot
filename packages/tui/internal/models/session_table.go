package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/components"
	"claude-pilot/shared/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
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

	// Search functionality
	searchInput         textinput.Model
	searchActive        bool
	searchQuery         string
	filteredSessions    []*api.Session
	filteredSessionData []components.SessionData

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

	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search sessions by name, status, or project..."
	searchInput.Width = 50
	searchInput = styles.ConfigureBubblesTextInput(searchInput)

	return &SessionTableModel{
		client:      client,
		table:       t,
		keys:        DefaultKeyMap(),
		searchInput: searchInput,
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
		// Handle search mode
		if m.searchActive {
			switch {
			case key.Matches(msg, m.keys.Escape):
				// Exit search mode
				m.searchActive = false
				m.searchInput.SetValue("")
				m.searchQuery = ""
				m.clearFilter()
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				// Exit search mode but keep filter
				m.searchActive = false
				return m, nil
			default:
				// Update search input
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				
				// Update search query and filter sessions
				newQuery := m.searchInput.Value()
				if newQuery != m.searchQuery {
					m.searchQuery = newQuery
					m.filterSessions(newQuery)
				}
				return m, cmd
			}
		}

		// Handle normal mode
		switch {
		case key.Matches(msg, m.keys.Search), msg.String() == "/":
			// Enter search mode
			m.searchActive = true
			m.searchInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if selected := m.GetSelectedSession(); selected != nil {
				return m, func() tea.Msg {
					return SessionSelectedMsg{Session: selected}
				}
			}
		case key.Matches(msg, m.keys.Home):
			// Jump to first session
			m.table.GotoTop()
			return m, nil
		case key.Matches(msg, m.keys.End):
			// Jump to last session
			m.table.GotoBottom()
			return m, nil
		case key.Matches(msg, m.keys.PageUp):
			// Page up navigation
			m.table.MoveUp(m.getPageSize())
			return m, nil
		case key.Matches(msg, m.keys.PageDown):
			// Page down navigation
			m.table.MoveDown(m.getPageSize())
			return m, nil
		case key.Matches(msg, m.keys.Escape):
			// Clear search filter if active
			if m.searchQuery != "" {
				m.clearFilter()
				return m, nil
			}
		default:
			// Let the Bubbles table handle standard navigation
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

	var view strings.Builder

	// Show search input if search is active
	if m.searchActive {
		searchHeader := styles.HeaderStyle.Render("Search Sessions")
		view.WriteString(searchHeader + "\n")
		view.WriteString(m.searchInput.View() + "\n")
		
		// Show results count
		currentSessions := m.getCurrentSessionList()
		resultCount := len(currentSessions)
		totalCount := len(m.sessions)
		
		var countText string
		if m.searchQuery == "" {
			countText = fmt.Sprintf("Total: %d sessions", totalCount)
		} else {
			countText = fmt.Sprintf("Found: %d of %d sessions", resultCount, totalCount)
			if resultCount == 0 {
				countText = styles.WarningStyle.Render(countText)
			} else {
				countText = styles.SuccessStyle.Render(countText)
			}
		}
		view.WriteString(countText + "\n\n")
	}

	// Render the Bubbles table
	tableView := m.table.View()
	view.WriteString(tableView)

	// Add session count info
	currentSessions := m.getCurrentSessionList()
	selectedIdx := m.table.Cursor()
	var statusLine string
	if len(currentSessions) > 0 {
		statusLine = fmt.Sprintf("Session %d of %d", selectedIdx+1, len(currentSessions))
	} else {
		statusLine = "No sessions match your search"
	}
	
	// Add search instruction if not in search mode
	if !m.searchActive {
		statusLine += " • Press '/' to search"
	} else {
		statusLine += " • Esc: Clear search • Enter: Apply filter"
	}

	view.WriteString("\n" + styles.DimTextStyle.Render(statusLine))

	return view.String()
}

// SetSessions updates the sessions data
func (m *SessionTableModel) SetSessions(sessions []*api.Session) {
	m.sessions = sessions
	
	// Apply current filter if active
	if m.searchQuery != "" {
		m.filterSessions(m.searchQuery)
	} else {
		m.sessionData = m.convertToSessionData(sessions)
		// Convert session data to Bubbles table rows
		rows := components.ToBubblesSessionRows(m.sessionData)
		m.table.SetRows(rows)
	}
}

// SetSize updates the table size
func (m *SessionTableModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.searchInput.Width = width - 10 // Responsive search input width
	m.updateTableSize()
}

// updateTableSize updates the Bubbles table dimensions
func (m *SessionTableModel) updateTableSize() {
	// Calculate available height for the table (account for headers and status line)
	availableHeight := m.height - 4 // Account for borders and status line
	
	// Account for search header if active
	if m.searchActive {
		availableHeight -= 4 // Search header, input, count, and spacing
	}
	
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
	currentSessions := m.getCurrentSessionList()
	if selectedIdx >= 0 && selectedIdx < len(currentSessions) {
		return currentSessions[selectedIdx]
	}
	return nil
}

// GetSelectedIndex returns the currently selected index
func (m *SessionTableModel) GetSelectedIndex() int {
	return m.table.Cursor()
}

// getPageSize calculates the page size for page up/down navigation
func (m *SessionTableModel) getPageSize() int {
	// Use table height minus 1 to keep one row visible for context
	pageSize := m.table.Height() - 1
	if pageSize < 1 {
		pageSize = 1
	}
	return pageSize
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

// Search and filter methods

// filterSessions filters sessions based on the search query
func (m *SessionTableModel) filterSessions(query string) {
	if query == "" {
		m.filteredSessions = m.sessions
		m.filteredSessionData = m.sessionData
	} else {
		query = strings.ToLower(strings.TrimSpace(query))
		var filtered []*api.Session
		
		for _, session := range m.sessions {
			// Search in session name
			if strings.Contains(strings.ToLower(session.Name), query) {
				filtered = append(filtered, session)
				continue
			}
			
			// Search in session status
			if strings.Contains(strings.ToLower(string(session.Status)), query) {
				filtered = append(filtered, session)
				continue
			}
			
			// Search in project path (if available)
			if session.ProjectPath != "" && strings.Contains(strings.ToLower(session.ProjectPath), query) {
				filtered = append(filtered, session)
				continue
			}
			
			// Search in description (if available)
			if session.Description != "" && strings.Contains(strings.ToLower(session.Description), query) {
				filtered = append(filtered, session)
				continue
			}
		}
		
		m.filteredSessions = filtered
		m.filteredSessionData = m.convertToSessionData(filtered)
	}
	
	// Update table with filtered results
	rows := components.ToBubblesSessionRows(m.filteredSessionData)
	m.table.SetRows(rows)
}

// clearFilter clears the current search filter
func (m *SessionTableModel) clearFilter() {
	m.searchQuery = ""
	m.searchInput.SetValue("")
	m.filteredSessions = nil
	m.filteredSessionData = nil
	
	// Reset table to show all sessions
	m.sessionData = m.convertToSessionData(m.sessions)
	rows := components.ToBubblesSessionRows(m.sessionData)
	m.table.SetRows(rows)
}

// getCurrentSessionList returns the current session list (filtered or all)
func (m *SessionTableModel) getCurrentSessionList() []*api.Session {
	if m.searchQuery != "" && m.filteredSessions != nil {
		return m.filteredSessions
	}
	return m.sessions
}

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
