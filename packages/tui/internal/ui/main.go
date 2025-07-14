package ui

import (
	"fmt"

	"claude-pilot/core/api"
	"claude-pilot/tui/internal/models"
	"claude-pilot/tui/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewState represents the current view in the TUI
type ViewState int

const (
	ViewSessionList ViewState = iota
	ViewSessionDetail
	ViewCreateSession
)

// MainModel is the root model for the TUI application
type MainModel struct {
	client      *api.Client
	currentView ViewState
	width       int
	height      int

	// Sub-models for different views
	sessionList   *models.SessionListModel
	sessionDetail *models.SessionDetailModel
	createSession *models.CreateSessionModel

	// Error handling
	err error
}

// NewMainModel creates a new main TUI model
func NewMainModel(client *api.Client) *MainModel {
	return &MainModel{
		client:        client,
		currentView:   ViewSessionList,
		sessionList:   models.NewSessionListModel(client),
		sessionDetail: models.NewSessionDetailModel(client),
		createSession: models.NewCreateSessionModel(client),
	}
}

// Init implements tea.Model
func (m *MainModel) Init() tea.Cmd {
	// Initialize the session list
	return m.sessionList.Init()
}

// Update implements tea.Model
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update sub-models with new size
		m.sessionList.SetSize(msg.Width, msg.Height)
		m.sessionDetail.SetSize(msg.Width, msg.Height)
		m.createSession.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			// Switch between views
			switch m.currentView {
			case ViewSessionList:
				m.currentView = ViewCreateSession
			case ViewCreateSession:
				m.currentView = ViewSessionList
			}

		case "enter":
			if m.currentView == ViewSessionList {
				// Get selected session and switch to detail view
				if selectedSession := m.sessionList.GetSelectedSession(); selectedSession != nil {
					m.sessionDetail.SetSession(selectedSession)
					m.currentView = ViewSessionDetail
				}
			}

		case "esc":
			// Go back to session list from any other view
			if m.currentView != ViewSessionList {
				m.currentView = ViewSessionList
			}
		}

	case models.SessionSelectedMsg:
		// Handle session selection from list
		m.sessionDetail.SetSession(msg.Session)
		m.currentView = ViewSessionDetail

	case models.SessionCreatedMsg:
		// Handle session creation
		if msg.Error == nil {
			// Refresh session list and go back to list view
			m.currentView = ViewSessionList
			cmd = m.sessionList.RefreshSessions()
		} else {
			m.err = msg.Error
		}

	case models.ErrorMsg:
		m.err = msg.Error
	}

	// Update the current view's model
	switch m.currentView {
	case ViewSessionList:
		var newModel tea.Model
		newModel, cmd = m.sessionList.Update(msg)
		m.sessionList = newModel.(*models.SessionListModel)

	case ViewSessionDetail:
		var newModel tea.Model
		newModel, cmd = m.sessionDetail.Update(msg)
		m.sessionDetail = newModel.(*models.SessionDetailModel)

	case ViewCreateSession:
		var newModel tea.Model
		newModel, cmd = m.createSession.Update(msg)
		m.createSession = newModel.(*models.CreateSessionModel)
	}

	return m, cmd
}

// View implements tea.Model
func (m *MainModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Header
	header := m.renderHeader()

	// Main content based on current view
	var content string
	switch m.currentView {
	case ViewSessionList:
		content = m.sessionList.View()
	case ViewSessionDetail:
		content = m.sessionDetail.View()
	case ViewCreateSession:
		content = m.createSession.View()
	}

	// Footer with navigation hints
	footer := m.renderFooter()

	// Error display
	errorDisplay := ""
	if m.err != nil {
		errorDisplay = styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		errorDisplay,
		footer,
	)
}

// renderHeader renders the application header
func (m *MainModel) renderHeader() string {
	title := "Claude Pilot TUI"
	backend := fmt.Sprintf("Backend: %s", m.client.GetBackend())

	titleStyle := styles.TitleStyle.Width(m.width - len(backend) - 2)
	backendStyle := styles.SubtleStyle

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		titleStyle.Render(title),
		backendStyle.Render(backend),
	)
}

// renderFooter renders the navigation footer
func (m *MainModel) renderFooter() string {
	var hints []string

	switch m.currentView {
	case ViewSessionList:
		hints = []string{
			"↑/↓: navigate",
			"enter: details",
			"tab: create",
			"q: quit",
		}
	case ViewSessionDetail:
		hints = []string{
			"esc: back",
			"q: quit",
		}
	case ViewCreateSession:
		hints = []string{
			"tab: back to list",
			"enter: create",
			"esc: cancel",
			"q: quit",
		}
	}

	return styles.FooterStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Left, hints...),
	)
}
