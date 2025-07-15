package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/layout"
	"claude-pilot/shared/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DashboardModel represents the main dashboard view
type DashboardModel struct {
	client *api.Client
	width  int
	height int

	// Child components
	summaryPanel *SummaryPanelModel
	sessionTable *SessionTableModel
	detailPanel  *DetailPanelModel
	createModal  *CreateModalModel

	// Bubbles help component
	help help.Model
	keys KeyMap

	// State
	focused         Component
	showCreateModal bool
	showHelp        bool
	sessions        []*api.Session
	selectedSession *api.Session
	err             error
}

// Component represents focusable components in the dashboard
type Component int

const (
	ComponentSummary Component = iota
	ComponentTable
	ComponentDetail
	ComponentModal
)

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(client *api.Client) *DashboardModel {
	// Create and configure help component
	h := help.New()
	h = styles.ConfigureBubblesHelp(h)

	return &DashboardModel{
		client:       client,
		summaryPanel: NewSummaryPanelModel(client),
		sessionTable: NewSessionTableModel(client),
		detailPanel:  NewDetailPanelModel(client),
		createModal:  NewCreateModalModel(client),
		help:         h,
		keys:         DefaultKeyMap(),
		focused:      ComponentTable,
	}
}

// Init implements tea.Model
func (m *DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.summaryPanel.Init(),
		m.sessionTable.Init(),
		m.detailPanel.Init(),
		m.loadSessions(), // Load initial data
	)
}

// Update implements tea.Model
func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateChildSizes()

	case tea.KeyMsg:
		if m.showCreateModal {
			// Handle modal input
			newModal, modalCmd := m.createModal.Update(msg)
			m.createModal = newModal.(*CreateModalModel)
			cmds = append(cmds, modalCmd)

			// Check if modal was closed
			if msg.String() == "esc" || m.createModal.IsCompleted() {
				m.showCreateModal = false
				if m.createModal.IsCompleted() {
					// Refresh sessions after creation
					cmds = append(cmds, m.loadSessions())
				}
				m.createModal.Reset()
			}
		} else {
			// Handle main dashboard input
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Create):
				m.showCreateModal = true
				m.focused = ComponentModal

			case key.Matches(msg, m.keys.Tab):
				m.cycleFocus()

			case key.Matches(msg, m.keys.Help):
				m.showHelp = !m.showHelp

			case key.Matches(msg, m.keys.Enter):
				if m.focused == ComponentTable && m.selectedSession != nil {
					// Attach to selected session
					return m, m.attachToSession(m.selectedSession.ID)
				}

			case msg.String() == "d":
				if m.focused == ComponentTable && m.selectedSession != nil {
					m.focused = ComponentDetail
				}

			case key.Matches(msg, m.keys.Kill):
				if m.focused == ComponentTable && m.selectedSession != nil {
					return m, m.killSession(m.selectedSession.ID)
				}

			case key.Matches(msg, m.keys.Refresh):
				// Refresh data
				cmds = append(cmds, m.loadSessions())
			}
		}

	case SessionsLoadedMsg:
		m.sessions = msg.Sessions
		m.err = msg.Error
		// Update child components
		m.summaryPanel.SetSessions(msg.Sessions)
		m.sessionTable.SetSessions(msg.Sessions)

	case SessionSelectedMsg:
		m.selectedSession = msg.Session
		m.detailPanel.SetSession(msg.Session)

	case SessionCreatedMsg:
		if msg.Error == nil {
			m.showCreateModal = false
			cmds = append(cmds, m.loadSessions())
		}

	case SessionRefreshedMsg:
		if msg.Session != nil {
			m.selectedSession = msg.Session
			m.detailPanel.SetSession(msg.Session)
		}

	case SessionAttachedMsg:
		if msg.Error != nil {
			m.err = msg.Error
		} else {
			// Successfully attached - exit TUI and let CLI handle the attachment
			return m, tea.Quit
		}

	case SessionKilledMsg:
		if msg.Error != nil {
			m.err = msg.Error
		} else {
			// Successfully killed session - refresh the session list
			cmds = append(cmds, m.loadSessions())
		}
	}

	// Update child components based on focus
	if !m.showCreateModal {
		switch m.focused {
		case ComponentTable:
			newTable, tableCmd := m.sessionTable.Update(msg)
			m.sessionTable = newTable.(*SessionTableModel)
			cmds = append(cmds, tableCmd)

			// Check for selection changes
			if selected := m.sessionTable.GetSelectedSession(); selected != m.selectedSession {
				m.selectedSession = selected
				m.detailPanel.SetSession(selected)
			}

		case ComponentDetail:
			newDetail, detailCmd := m.detailPanel.Update(msg)
			m.detailPanel = newDetail.(*DetailPanelModel)
			cmds = append(cmds, detailCmd)

		case ComponentSummary:
			newSummary, summaryCmd := m.summaryPanel.Update(msg)
			m.summaryPanel = newSummary.(*SummaryPanelModel)
			cmds = append(cmds, summaryCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *DashboardModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading dashboard..."
	}

	// Render main dashboard content
	dashboardContent := m.renderDashboard()

	// Overlay create modal if shown
	if m.showCreateModal {
		modal := m.createModal.View()
		return m.overlayModal(dashboardContent, modal)
	}

	// Overlay help if shown
	if m.showHelp {
		helpView := m.renderHelpView()
		return m.overlayHelp(dashboardContent, helpView)
	}

	return dashboardContent
}

// renderDashboard renders the main dashboard layout
func (m *DashboardModel) renderDashboard() string {
	// Create header with summary cards
	header := m.renderHeader()

	// Create main content area
	// mainContent := m.renderMainContent()

	// Create footer with help
	// footer := m.renderFooter()

	// Use dashboard layout
	return layout.DashboardLayout(m.width, m.height, header, "", "")
	// return layout.DashboardLayout(m.width, m.height, header, mainContent, footer)
}

// renderHeader renders the header with title and summary cards
func (m *DashboardModel) renderHeader() string {
	// Title
	title := styles.TitleStyle.Render("Claude Pilot Dashboard")
	backend := styles.SecondaryTextStyle.Render(fmt.Sprintf("Backend: %s", m.client.GetBackend()))
	spacer := lipgloss.NewStyle().Width(m.width - lipgloss.Width(title) - lipgloss.Width(backend)).Render("")

	titleRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		spacer,
		backend,
	)

	summaryContent := m.summaryPanel.View()
	return lipgloss.JoinVertical(lipgloss.Left, titleRow, summaryContent)
}

// renderMainContent renders the main content area with table and detail panel
func (m *DashboardModel) renderMainContent() string {
	// Get responsive width
	layoutWidth, size := styles.GetResponsiveWidth(m.width)

	// Calculate available height for main content (more compact layout)
	availableHeight := m.height - 10 // Increased from 6 to 10 to make main content shorter
	if availableHeight < 8 {
		availableHeight = 8 // Reduced minimum from 10 to 8
	}

	// Determine layout based on screen size
	switch size {
	case "small":
		// Stack vertically on small screens
		if m.focused == ComponentDetail && m.selectedSession != nil {
			return m.detailPanel.View()
		}
		return m.sessionTable.View()

	case "medium":
		// Two-column layout
		tablePanel := layout.NewPanel(
			layout.LayoutConfig{Width: layoutWidth * 2 / 3, Height: availableHeight, Padding: 1},
			"Sessions",
			m.sessionTable.View(),
			true,
		)
		tablePanel.SetFocused(m.focused == ComponentTable)

		if m.selectedSession != nil {
			detailPanelComp := layout.NewPanel(
				layout.LayoutConfig{Width: layoutWidth / 3, Height: availableHeight, Padding: 1},
				"Session Details",
				m.detailPanel.View(),
				true,
			)
			detailPanelComp.SetFocused(m.focused == ComponentDetail)

			return layout.SidebarLayout(layoutWidth, availableHeight, tablePanel.Render(), detailPanelComp.Render(), layoutWidth/3)
		}

		return tablePanel.Render()

	default: // large
		// Simple side-by-side layout for large screens
		tablePanel := layout.NewPanel(
			layout.LayoutConfig{Width: layoutWidth / 2, Height: availableHeight, Padding: 1},
			"Sessions",
			m.sessionTable.View(),
			true,
		)
		tablePanel.SetFocused(m.focused == ComponentTable)

		if m.selectedSession != nil {
			detailPanelComp := layout.NewPanel(
				layout.LayoutConfig{Width: layoutWidth / 2, Height: availableHeight, Padding: 1},
				"Session Details",
				m.detailPanel.View(),
				true,
			)
			detailPanelComp.SetFocused(m.focused == ComponentDetail)

			// Use horizontal layout for large screens
			return lipgloss.JoinHorizontal(
				lipgloss.Top,
				tablePanel.Render(),
				detailPanelComp.Render(),
			)
		}

		return tablePanel.Render()
	}
}

// renderFooter renders the footer with keyboard shortcuts
func (m *DashboardModel) renderFooter() string {
	if m.showCreateModal {
		// Show modal-specific shortcuts
		shortcuts := []string{
			"Enter: Create",
			"Esc: Cancel",
		}

		helpText := ""
		for i, shortcut := range shortcuts {
			if i > 0 {
				helpText += " • "
			}
			helpText += shortcut
		}

		return styles.FooterStyle.Width(m.width).Render(helpText)
	} else {
		// Use Bubbles help component for main dashboard
		return m.help.View(m.keys)
	}
}

// overlayModal overlays the create modal over the dashboard content
func (m *DashboardModel) overlayModal(background, modal string) string {
	// Calculate modal position (centered)
	modalWidth := 60
	modalHeight := 15

	x := (m.width - modalWidth) / 2
	y := (m.height - modalHeight) / 2

	if x < 0 {
		x = 1
	}
	if y < 0 {
		y = 1
	}

	// Note: x and y are calculated for future use but lipgloss.Place handles centering
	_ = x
	_ = y

	// Create overlay style (simplified since Position is not available)
	overlay := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ClaudePrimary).
		Background(styles.BackgroundPrimary)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		overlay.Render(modal),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.BackgroundSecondary),
	)
}

// renderHelpView renders the help overlay content
func (m *DashboardModel) renderHelpView() string {
	var content strings.Builder

	// Title
	content.WriteString(styles.TitleStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n\n")

	// Main shortcuts
	content.WriteString(styles.HeaderStyle.Render("Main Controls"))
	content.WriteString("\n")
	content.WriteString(m.help.View(m.keys))
	content.WriteString("\n\n")

	// Context-specific shortcuts
	switch m.focused {
	case ComponentTable:
		content.WriteString(styles.HeaderStyle.Render("Table Navigation"))
		content.WriteString("\n")
		content.WriteString(m.help.View(m.keys))
	case ComponentDetail:
		content.WriteString(styles.HeaderStyle.Render("Detail Panel"))
		content.WriteString("\n")
		content.WriteString(m.help.View(m.keys))
	}

	content.WriteString("\n\n")
	content.WriteString(styles.DimTextStyle.Render("Press ? again to close help"))

	return content.String()
}

// overlayHelp overlays the help content over the dashboard
func (m *DashboardModel) overlayHelp(background, helpContent string) string {
	// Calculate help overlay position (centered)
	helpWidth := 60
	helpHeight := 20

	x := (m.width - helpWidth) / 2
	y := (m.height - helpHeight) / 2

	if x < 0 {
		x = 1
	}
	if y < 0 {
		y = 1
	}

	// Note: x and y are calculated for future use but lipgloss.Place handles centering
	_ = x
	_ = y

	// Create help overlay style
	overlay := lipgloss.NewStyle().
		Width(helpWidth).
		Height(helpHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.InfoColor).
		Background(styles.BackgroundPrimary).
		Padding(1)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		overlay.Render(helpContent),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.BackgroundSecondary),
	)
}

// updateChildSizes updates the sizes of child components
func (m *DashboardModel) updateChildSizes() {
	// Calculate available height for main content (more compact)
	availableHeight := m.height - 10 // Increased from 6 to 10 to make main content shorter
	if availableHeight < 8 {
		availableHeight = 8 // Reduced minimum from 10 to 8
	}

	m.summaryPanel.SetSize(m.width, 3) // Keep summary compact
	m.sessionTable.SetSize(m.width*2/3, availableHeight)
	m.detailPanel.SetSize(m.width/3, availableHeight)
	m.createModal.SetSize(60, 15) // Fixed size for modal
}

// cycleFocus cycles through focusable components
func (m *DashboardModel) cycleFocus() {
	switch m.focused {
	case ComponentSummary:
		m.focused = ComponentTable
	case ComponentTable:
		if m.selectedSession != nil {
			m.focused = ComponentDetail
		} else {
			m.focused = ComponentSummary
		}
	case ComponentDetail:
		m.focused = ComponentSummary
	default:
		m.focused = ComponentTable
	}
}

// loadSessions loads sessions from the API
func (m *DashboardModel) loadSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, err := m.client.ListSessions()
		return SessionsLoadedMsg{Sessions: sessions, Error: err}
	}
}

// attachToSession attaches to a session
func (m *DashboardModel) attachToSession(sessionID string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.AttachToSession(sessionID)
		return SessionAttachedMsg{
			SessionID: sessionID,
			Error:     err,
		}
	}
}

// killSession kills a session
func (m *DashboardModel) killSession(sessionID string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.KillSession(sessionID)
		return SessionKilledMsg{
			SessionID: sessionID,
			Error:     err,
		}
	}
}

// SetSize updates the dashboard size
func (m *DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updateChildSizes()
}

// Common message types
type SessionsLoadedMsg struct {
	Sessions []*api.Session
	Error    error
}

type ErrorMsg struct {
	Error error
}

type SessionAttachedMsg struct {
	SessionID string
	Error     error
}

type SessionKilledMsg struct {
	SessionID string
	Error     error
}
