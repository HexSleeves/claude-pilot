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

// renderDashboard renders the main dashboard layout using enhanced flexbox
func (m *DashboardModel) renderDashboard() string {
	// Create header with summary cards
	header := m.renderHeader()

	// Create main content area with enhanced flexbox layout
	mainContent := m.renderMainContent()

	// Create footer with help
	footer := m.renderFooter()

	// Use dashboard layout with all sections
	return layout.DashboardLayout(m.width, m.height, header, mainContent, footer)
}

// renderHeader renders the header with title and summary cards using flexbox
func (m *DashboardModel) renderHeader() string {
	// Create title row using flexbox for better alignment
	title := styles.TitleStyle.Render("Claude Pilot Dashboard")
	backend := styles.SecondaryTextStyle.Render(fmt.Sprintf("Backend: %s", m.client.GetBackend()))

	// Use flexbox container for title row with space-between justification
	titleContainer := layout.NewFlexContainer(
		layout.LayoutConfig{Width: m.width, Height: 1, Padding: 0, Gap: 0},
		layout.FlexRow,
	).SetJustifyContent(layout.SpaceBetween).SetAlignItems(layout.AlignCenter)

	// Add title and backend info as flex items
	titleContainer.AddItem(layout.FlexItem{
		Content:    title,
		FlexGrow:   1, // Don't grow
		FlexShrink: 1, // Don't shrink
		Order:      1,
	})

	titleContainer.AddItem(layout.FlexItem{
		Content:    backend,
		FlexGrow:   1, // Don't grow
		FlexShrink: 1, // Don't shrink
		Order:      2,
	})

	titleRow := titleContainer.Render()

	// Create summary section using flexbox for responsive cards
	summaryContent := m.summaryPanel.View()

	// Combine title and summary vertically
	headerContainer := layout.NewFlexContainer(
		layout.LayoutConfig{Width: m.width, Height: 3, Padding: 0, Gap: 0},
		layout.FlexColumn,
	).SetJustifyContent(layout.FlexStart).SetAlignItems(layout.AlignStretch)

	headerContainer.AddItem(layout.FlexItem{
		Content:    titleRow,
		FlexGrow:   0, // Fixed height for title
		FlexShrink: 0,
		Order:      1,
	})

	headerContainer.AddItem(layout.FlexItem{
		Content:    summaryContent,
		FlexGrow:   1, // Take remaining space
		FlexShrink: 0,
		Order:      2,
	})

	return headerContainer.Render()
}

// renderMainContent renders the main content area with enhanced flexbox layout
func (m *DashboardModel) renderMainContent() string {
	// Get responsive width
	layoutWidth, size := styles.GetResponsiveWidth(m.width)

	// Calculate available height for main content
	availableHeight := max(8, m.height-10)

	// Create main content container with flexbox
	mainContainer := layout.NewFlexContainer(
		layout.LayoutConfig{Width: layoutWidth, Height: availableHeight, Padding: 1, Gap: 2},
		layout.FlexRow,
	).SetJustifyContent(layout.FlexStart).SetAlignItems(layout.AlignStretch)

	// Determine layout based on screen size
	switch size {
	case "small":
		// Stack vertically on small screens - use column layout
		mainContainer = layout.NewFlexContainer(
			layout.LayoutConfig{Width: layoutWidth, Height: availableHeight, Padding: 1, Gap: 1},
			layout.FlexColumn,
		).SetJustifyContent(layout.FlexStart).SetAlignItems(layout.AlignStretch)

		if m.focused == ComponentDetail && m.selectedSession != nil {
			// Show only detail panel when focused
			mainContainer.AddItem(layout.FlexItem{
				Content:    m.renderDetailPanel(layoutWidth, availableHeight),
				FlexGrow:   1,
				FlexShrink: 0,
				Order:      1,
			})
		} else {
			// Show only table when not in detail mode
			mainContainer.AddItem(layout.FlexItem{
				Content:    m.renderSessionTable(layoutWidth, availableHeight),
				FlexGrow:   1,
				FlexShrink: 0,
				Order:      1,
			})
		}

	case "medium":
		// Two-column layout with 2:1 ratio
		mainContainer.AddItem(layout.FlexItem{
			Content:    m.renderSessionTable(layoutWidth*2/3, availableHeight),
			FlexGrow:   2, // Takes 2/3 of space
			FlexShrink: 1,
			Order:      1,
		})

		if m.selectedSession != nil {
			mainContainer.AddItem(layout.FlexItem{
				Content:    m.renderDetailPanel(layoutWidth/3, availableHeight),
				FlexGrow:   1, // Takes 1/3 of space
				FlexShrink: 1,
				Order:      2,
			})
		}

	default: // large
		// Equal split layout for large screens
		mainContainer.AddItem(layout.FlexItem{
			Content:    m.renderSessionTable(layoutWidth/2, availableHeight),
			FlexGrow:   1, // Equal flex-grow
			FlexShrink: 1,
			Order:      1,
		})

		if m.selectedSession != nil {
			mainContainer.AddItem(layout.FlexItem{
				Content:    m.renderDetailPanel(layoutWidth/2, availableHeight),
				FlexGrow:   1, // Equal flex-grow
				FlexShrink: 1,
				Order:      2,
			})
		}
	}

	return mainContainer.Render()
}

// renderSessionTable renders the session table panel with proper styling
func (m *DashboardModel) renderSessionTable(width, height int) string {
	tablePanel := layout.NewPanel(
		layout.LayoutConfig{Width: width, Height: height, Padding: 1},
		"Sessions",
		m.sessionTable.View(),
		true,
	)
	tablePanel.SetFocused(m.focused == ComponentTable)
	return tablePanel.Render()
}

// renderDetailPanel renders the detail panel with proper styling
func (m *DashboardModel) renderDetailPanel(width, height int) string {
	detailPanelComp := layout.NewPanel(
		layout.LayoutConfig{Width: width, Height: height, Padding: 1},
		"Session Details",
		m.detailPanel.View(),
		true,
	)
	detailPanelComp.SetFocused(m.focused == ComponentDetail)
	return detailPanelComp.Render()
}

// renderFooter renders the footer with keyboard shortcuts using flexbox
func (m *DashboardModel) renderFooter() string {
	var footerContent string

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
		footerContent = helpText
	} else {
		// Use Bubbles help component for main dashboard
		footerContent = m.help.View(m.keys)
	}

	// Create footer container with flexbox for better alignment
	footerContainer := layout.NewFlexContainer(
		layout.LayoutConfig{Width: m.width, Height: 2, Padding: 0, Gap: 0},
		layout.FlexRow,
	).SetJustifyContent(layout.Center).SetAlignItems(layout.AlignCenter)

	footerContainer.AddItem(layout.FlexItem{
		Content:    styles.FooterStyle.Render(footerContent),
		FlexGrow:   0, // Don't grow
		FlexShrink: 0, // Don't shrink
		Order:      1,
	})

	return footerContainer.Render()
}

// overlayModal overlays the create modal over the dashboard content using flexbox
func (m *DashboardModel) overlayModal(background, modal string) string {
	// Modal dimensions - increased width for better form layout
	modalWidth := 80
	modalHeight := 18

	// Create modal content with styling
	overlay := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ClaudePrimary)
		// Background(styles.BackgroundPrimary)

	styledModal := overlay.Render(modal)

	// Overlay the modal on top of the background using lipgloss.Place
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styledModal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.BackgroundSecondary),
		lipgloss.WithWhitespaceBackground(lipgloss.Color(background)),
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

// overlayHelp overlays the help content over the dashboard using flexbox
func (m *DashboardModel) overlayHelp(background, helpContent string) string {
	// Help overlay dimensions
	helpWidth := 60
	helpHeight := 20

	// Create help overlay style
	overlay := lipgloss.NewStyle().
		Width(helpWidth).
		Height(helpHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.InfoColor).
		Background(styles.BackgroundPrimary).
		Padding(1)

	styledHelp := overlay.Render(helpContent)

	// Use flexbox container for perfect centering
	helpContainer := layout.NewFlexContainer(
		layout.LayoutConfig{Width: m.width, Height: m.height, Padding: 0, Gap: 0},
		layout.FlexRow,
	).SetJustifyContent(layout.Center).SetAlignItems(layout.AlignCenter)

	helpContainer.AddItem(layout.FlexItem{
		Content:    styledHelp,
		FlexGrow:   0, // Don't grow
		FlexShrink: 0, // Don't shrink
		Order:      1,
	})

	// Render with background styling
	result := helpContainer.Render()

	// Apply background styling using lipgloss.Place for whitespace handling
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		result,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.BackgroundSecondary),
	)
}

// updateChildSizes updates the sizes of child components
func (m *DashboardModel) updateChildSizes() {
	// Calculate available height for main content using max function
	availableHeight := max(8, m.height-10)

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
