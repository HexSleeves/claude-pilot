package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const DefaultNotificationDuration = 3 * time.Second

// NotificationType represents the type of notification
type NotificationType int

const (
	NotificationSuccess NotificationType = iota
	NotificationError
	NotificationInfo
	NotificationWarning
)

// NotificationModel represents a toast-style notification
type NotificationModel struct {
	Type      NotificationType
	Message   string
	Visible   bool
	Duration  time.Duration
	width     int
	height    int
	startTime time.Time

	// Styling
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
	infoStyle    lipgloss.Style
	warningStyle lipgloss.Style
}

// HideNotificationMsg represents a message to hide the notification
type HideNotificationMsg struct{}

// NewNotificationModel creates a new notification model
func NewNotificationModel() *NotificationModel {
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2ECC71")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#27AE60"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#E74C3C")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#C0392B"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#3498DB")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2980B9"))

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2C3E50")).
		Background(lipgloss.Color("#F39C12")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E67E22"))

	return &NotificationModel{
		successStyle: successStyle,
		errorStyle:   errorStyle,
		infoStyle:    infoStyle,
		warningStyle: warningStyle,
		Duration:     DefaultNotificationDuration,
	}
}

// Init initializes the notification model
func (m NotificationModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the notification
func (m NotificationModel) Update(msg tea.Msg) (NotificationModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case HideNotificationMsg:
		m.Visible = false
	}

	return m, nil
}

// View renders the notification
func (m NotificationModel) View() string {
	if !m.Visible || m.Message == "" {
		return ""
	}

	var style lipgloss.Style
	var icon string

	switch m.Type {
	case NotificationSuccess:
		style = m.successStyle
		icon = "✅"
	case NotificationError:
		style = m.errorStyle
		icon = "❌"
	case NotificationInfo:
		style = m.infoStyle
		icon = "ℹ️"
	case NotificationWarning:
		style = m.warningStyle
		icon = "⚠️"
	default:
		style = m.infoStyle
		icon = "ℹ️"
	}

	content := icon + " " + m.Message

	// Add padding to prevent content from being smashed against the edges
	style = style.PaddingLeft(2).PaddingRight(2).PaddingTop(1).PaddingBottom(1)

	// Position the notification in the top-right corner
	if m.width > 0 {
		maxWidth := min(m.width-8, 50) // Leave margin for padding (increased from 4 to 8)
		style = style.Width(maxWidth)
	}

	return style.Render(content)
}

// Show displays a notification with auto-hide
func (m *NotificationModel) Show(notificationType NotificationType, message string) tea.Cmd {
	m.Type = notificationType
	m.Message = message
	m.Visible = true
	m.startTime = time.Now()

	return tea.Tick(m.Duration, func(t time.Time) tea.Msg {
		return HideNotificationMsg{}
	})
}

// ShowSuccess shows a success notification
func (m *NotificationModel) ShowSuccess(message string) tea.Cmd {
	return m.Show(NotificationSuccess, message)
}

// ShowError shows an error notification
func (m *NotificationModel) ShowError(message string) tea.Cmd {
	return m.Show(NotificationError, message)
}

// ShowInfo shows an info notification
func (m *NotificationModel) ShowInfo(message string) tea.Cmd {
	return m.Show(NotificationInfo, message)
}

// ShowWarning shows a warning notification
func (m *NotificationModel) ShowWarning(message string) tea.Cmd {
	return m.Show(NotificationWarning, message)
}

// Hide hides the notification
func (m *NotificationModel) Hide() {
	m.Visible = false
}

// IsVisible returns true if the notification is visible
func (m NotificationModel) IsVisible() bool {
	return m.Visible
}

// SetDimensions sets the notification dimensions
func (m *NotificationModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
