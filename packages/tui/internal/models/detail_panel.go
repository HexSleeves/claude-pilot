package models

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/styles"
	"claude-pilot/shared/utils"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailPanelModel represents the session detail panel
type DetailPanelModel struct {
	client  *api.Client
	width   int
	height  int
	session *api.Session

	// Bubbles viewport for scrolling
	viewport viewport.Model

	// View mode
	showFullMessages bool

	// Key bindings
	keys KeyMap
}

// NewDetailPanelModel creates a new detail panel model
func NewDetailPanelModel(client *api.Client) *DetailPanelModel {
	// Create and configure viewport
	vp := viewport.New(80, 20)
	vp = styles.ConfigureBubblesViewport(vp)

	return &DetailPanelModel{
		client:           client,
		viewport:         vp,
		showFullMessages: false,
		keys:             DefaultKeyMap(),
	}
}

// Init implements tea.Model
func (m *DetailPanelModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *DetailPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewportSize()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ScrollUp):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.ScrollDown):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.PageUp):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.PageDown):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.Home):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.End):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case msg.String() == "f":
			m.toggleFullMessages()
		case msg.String() == "a":
			// Attach to session
			if m.session != nil {
				return m, m.attachToSession()
			}
		case msg.String() == "x":
			// Kill session
			if m.session != nil {
				return m, m.killSession()
			}
		case msg.String() == "r":
			// Refresh session details
			if m.session != nil {
				return m, m.refreshSession()
			}
		}
	}

	// Update viewport for other messages
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *DetailPanelModel) View() string {
	if m.session == nil {
		return m.renderEmpty()
	}

	return m.renderDetail()
}

// renderEmpty renders the empty state
func (m *DetailPanelModel) renderEmpty() string {
	emptyMsg := styles.DimTextStyle.Render("Select a session to view details")

	if m.height > 0 {
		// Center the message vertically
		padding := (m.height - 3) / 2
		if padding > 0 {
			emptyMsg = strings.Repeat("\n", padding) + emptyMsg
		}
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(emptyMsg)
}

// renderDetail renders the session details
func (m *DetailPanelModel) renderDetail() string {
	var content strings.Builder

	// Session info section
	content.WriteString(m.renderSessionInfo())
	content.WriteString("\n\n")

	// Action buttons section
	content.WriteString(m.renderActionButtons())
	content.WriteString("\n\n")

	// Messages section
	content.WriteString(m.renderMessages())

	// Set viewport content and return the viewport view
	m.viewport.SetContent(content.String())
	return m.viewport.View()
}

// renderSessionInfo renders the basic session information
func (m *DetailPanelModel) renderSessionInfo() string {
	var info strings.Builder

	// Header with session name
	info.WriteString(styles.SessionNameStyle.Render(m.session.Name))
	info.WriteString("\n")
	info.WriteString(styles.SessionIDStyle.Render(fmt.Sprintf("ID: %s", m.session.ID[:8])))
	info.WriteString("\n\n")

	// Status with styling
	statusText := m.formatStatus(m.session.Status)
	info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Status:"), statusText))

	// Backend
	backend := m.client.GetBackend()
	info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Backend:"), styles.DimTextStyle.Render(backend)))

	// Timestamps
	info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Created:"),
		styles.DimTextStyle.Render(m.session.CreatedAt.Format("2006-01-02 15:04"))))
	info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Last Active:"),
		styles.DimTextStyle.Render(m.formatTimeAgo(m.session.LastActive))))

	// Message count
	info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Messages:"),
		styles.HighlightStyle.Render(fmt.Sprintf("%d", len(m.session.Messages)))))

	// Project path if available
	if m.session.ProjectPath != "" {
		info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Project:"),
			styles.DimTextStyle.Render(m.truncatePath(m.session.ProjectPath))))
	}

	// Description if available
	if m.session.Description != "" {
		info.WriteString(fmt.Sprintf("%-12s %s\n", styles.Bold("Description:"),
			styles.SecondaryTextStyle.Render(m.session.Description)))
	}

	return info.String()
}

// renderActionButtons renders the action buttons
func (m *DetailPanelModel) renderActionButtons() string {
	var buttons []string

	switch m.session.Status {
	case api.StatusActive, api.StatusInactive:
		buttons = append(buttons, styles.ButtonPrimaryStyle.Render("a) Attach"))
	case api.StatusConnected:
		buttons = append(buttons, styles.DimTextStyle.Render("Already connected"))
	case api.StatusError:
		buttons = append(buttons, styles.ErrorStyle.Render("Session has errors"))
	}

	buttons = append(buttons, styles.ButtonDangerStyle.Render("k) Kill"))
	buttons = append(buttons, styles.ButtonPrimaryStyle.Render("r) Refresh"))

	if len(m.session.Messages) > 0 {
		if m.showFullMessages {
			buttons = append(buttons, styles.ButtonPrimaryStyle.Render("f) Summary"))
		} else {
			buttons = append(buttons, styles.ButtonPrimaryStyle.Render("f) Full Messages"))
		}
	}

	return strings.Join(buttons, "  ")
}

// renderMessages renders the message history
func (m *DetailPanelModel) renderMessages() string {
	if len(m.session.Messages) == 0 {
		return styles.DimTextStyle.Render("No messages in this session")
	}

	var content strings.Builder
	content.WriteString(styles.HeaderStyle.Render("Recent Messages"))
	content.WriteString("\n")
	content.WriteString(styles.HorizontalLine(m.width - 4))
	content.WriteString("\n\n")

	if m.showFullMessages {
		content.WriteString(m.renderFullMessages())
	} else {
		content.WriteString(m.renderMessageSummary())
	}

	return content.String()
}

// renderMessageSummary renders a summary of recent messages
func (m *DetailPanelModel) renderMessageSummary() string {
	var content strings.Builder

	// Show last 5 messages in summary format
	start := len(m.session.Messages) - 5
	if start < 0 {
		start = 0
	}

	for i := start; i < len(m.session.Messages); i++ {
		msg := m.session.Messages[i]

		// Role indicator
		roleStyle := styles.InfoStyle
		roleIcon := "🤖"
		if msg.Role == "user" {
			roleStyle = styles.HighlightStyle
			roleIcon = "👤"
		}

		// Timestamp
		timestamp := styles.DimTextStyle.Render(msg.Timestamp.Format("15:04:05"))

		// Content preview (truncated)
		contentPreview := styles.TruncateText(msg.Content, 60)

		content.WriteString(fmt.Sprintf("%s %s %s\n",
			roleIcon,
			roleStyle.Render(fmt.Sprintf("[%s]", msg.Role)),
			timestamp,
		))
		content.WriteString(fmt.Sprintf("  %s\n\n",
			styles.SecondaryTextStyle.Render(contentPreview)))
	}

	if len(m.session.Messages) > 5 {
		content.WriteString(styles.DimTextStyle.Render(
			fmt.Sprintf("... and %d more messages (press 'f' for full view)",
				len(m.session.Messages)-5)))
	}

	return content.String()
}

// renderFullMessages renders all messages in full detail
func (m *DetailPanelModel) renderFullMessages() string {
	var content strings.Builder

	for i, msg := range m.session.Messages {
		// Role header
		roleStyle := styles.InfoStyle
		roleIcon := "🤖"
		if msg.Role == "user" {
			roleStyle = styles.HighlightStyle
			roleIcon = "👤"
		}

		content.WriteString(fmt.Sprintf("%s %s %s %s\n",
			roleIcon,
			roleStyle.Render(strings.ToUpper(msg.Role)),
			styles.DimTextStyle.Render(msg.Timestamp.Format("2006-01-02 15:04:05")),
			styles.DimTextStyle.Render(fmt.Sprintf("(%d/%d)", i+1, len(m.session.Messages))),
		))

		// Message content
		content.WriteString(styles.PrimaryTextStyle.Render(msg.Content))
		content.WriteString("\n")

		// Separator
		if i < len(m.session.Messages)-1 {
			content.WriteString(styles.DimTextStyle.Render(strings.Repeat("─", 40)))
			content.WriteString("\n\n")
		}
	}

	return content.String()
}

// SetSession sets the session to display
func (m *DetailPanelModel) SetSession(session *api.Session) {
	m.session = session
	m.showFullMessages = false
	// Reset viewport to top
	m.viewport.GotoTop()
}

// SetSize updates the panel size
func (m *DetailPanelModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updateViewportSize()
}

// updateViewportSize updates the viewport dimensions
func (m *DetailPanelModel) updateViewportSize() {
	// Calculate available space for viewport
	availableHeight := m.height - 2 // Account for borders
	if availableHeight < 1 {
		availableHeight = 1
	}

	availableWidth := m.width - 2 // Account for borders
	if availableWidth < 1 {
		availableWidth = 1
	}

	m.viewport.Width = availableWidth
	m.viewport.Height = availableHeight
}

// Scrolling is now handled by the Bubbles viewport component

func (m *DetailPanelModel) toggleFullMessages() {
	m.showFullMessages = !m.showFullMessages
	// Reset viewport to top when toggling
	m.viewport.GotoTop()
}

// Action methods

func (m *DetailPanelModel) attachToSession() tea.Cmd {
	return func() tea.Msg {
		if m.session == nil {
			return SessionAttachedMsg{
				SessionID: "",
				Error:     fmt.Errorf("no session selected"),
			}
		}

		err := m.client.AttachToSession(m.session.ID)
		return SessionAttachedMsg{
			SessionID: m.session.ID,
			Error:     err,
		}
	}
}

func (m *DetailPanelModel) killSession() tea.Cmd {
	return func() tea.Msg {
		if m.session == nil {
			return SessionKilledMsg{
				SessionID: "",
				Error:     fmt.Errorf("no session selected"),
			}
		}

		err := m.client.KillSession(m.session.ID)
		return SessionKilledMsg{
			SessionID: m.session.ID,
			Error:     err,
		}
	}
}

func (m *DetailPanelModel) refreshSession() tea.Cmd {
	return func() tea.Msg {
		session, err := m.client.GetSession(m.session.ID)
		if err != nil {
			// Handle error, return existing for now
			return SessionRefreshedMsg{Session: m.session}
		}
		return SessionRefreshedMsg{Session: session}
	}
}

// Utility methods

func (m *DetailPanelModel) formatStatus(status api.SessionStatus) string {
	return utils.FormatSessionStatus(string(status)).Render(string(status))
}

func (m *DetailPanelModel) formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}

func (m *DetailPanelModel) truncatePath(path string) string {
	maxLen := 30
	if len(path) <= maxLen {
		return path
	}

	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return ".../" + parts[len(parts)-1]
	}

	return path[:maxLen-3] + "..."
}

// Message types
type SessionRefreshedMsg struct {
	Session *api.Session
}
