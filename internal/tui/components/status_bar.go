package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusState represents the state of the status bar
type StatusState int

const (
	StatusIdle StatusState = iota
	StatusLoading
	StatusSuccess
	StatusError
	StatusInfo
)

// ClearStatusMsg represents a message to clear the status bar
type ClearStatusMsg struct{}

// StatusBarModel handles status display and feedback
type StatusBarModel struct {
	state        StatusState
	message      string
	width        int
	autoClear    bool
	clearTimeout time.Duration

	// Styling
	loadingStyle lipgloss.Style
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
	infoStyle    lipgloss.Style
}

// NewStatusBarModel creates a new status bar model
func NewStatusBarModel() *StatusBarModel {
	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498DB")).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ECC71")).
		Bold(true)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E74C3C")).
		Bold(true)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F39C12")).
		Bold(true)

	return &StatusBarModel{
		state:        StatusIdle,
		clearTimeout: 3 * time.Second, // Auto-clear after 3 seconds
		loadingStyle: loadingStyle,
		successStyle: successStyle,
		errorStyle:   errorStyle,
		infoStyle:    infoStyle,
	}
}

// Init initializes the status bar model
func (m StatusBarModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the status bar
func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case ClearStatusMsg:
		m.state = StatusIdle
		m.message = ""
		m.autoClear = false
	}

	return m, nil
}

// View renders the status bar
func (m StatusBarModel) View() string {
	if m.state == StatusIdle || m.message == "" {
		return ""
	}

	var style lipgloss.Style
	var icon string

	switch m.state {
	case StatusLoading:
		style = m.loadingStyle
		icon = "⏳"
	case StatusSuccess:
		style = m.successStyle
		icon = "✅"
	case StatusError:
		style = m.errorStyle
		icon = "❌"
	case StatusInfo:
		style = m.infoStyle
		icon = "ℹ️"
	default:
		style = m.loadingStyle
		icon = ""
	}

	content := icon + " " + m.message

	// Add padding to prevent content from being smashed against the left side
	style = style.PaddingLeft(0).PaddingRight(2)

	// Apply width constraint if available
	if m.width > 0 {
		style = style.Width(m.width)
	}

	return style.Render(content)
}

// SetLoading sets the status bar to loading state
func (m *StatusBarModel) SetLoading(message string) {
	m.state = StatusLoading
	m.message = message
	m.autoClear = false
}

// SetSuccess sets the status bar to success state with auto-clear
func (m *StatusBarModel) SetSuccess(message string) tea.Cmd {
	m.state = StatusSuccess
	m.message = message
	m.autoClear = true
	return m.autoClearCmd()
}

// SetError sets the status bar to error state with auto-clear
func (m *StatusBarModel) SetError(message string) tea.Cmd {
	m.state = StatusError
	m.message = message
	m.autoClear = true
	return m.autoClearCmd()
}

// SetInfo sets the status bar to info state with auto-clear
func (m *StatusBarModel) SetInfo(message string) tea.Cmd {
	m.state = StatusInfo
	m.message = message
	m.autoClear = true
	return m.autoClearCmd()
}

// SetPersistent sets a persistent message that won't auto-clear
func (m *StatusBarModel) SetPersistent(state StatusState, message string) {
	m.state = state
	m.message = message
	m.autoClear = false
}

// Clear clears the status bar
func (m *StatusBarModel) Clear() {
	m.state = StatusIdle
	m.message = ""
	m.autoClear = false
}

// IsVisible returns true if the status bar should be displayed
func (m StatusBarModel) IsVisible() bool {
	return m.state != StatusIdle && m.message != ""
}

// SetWidth sets the width of the status bar
func (m *StatusBarModel) SetWidth(width int) {
	m.width = width
}

// GetState returns the current status state
func (m StatusBarModel) GetState() StatusState {
	return m.state
}

// autoClearCmd returns a command that clears the status after a timeout
func (m StatusBarModel) autoClearCmd() tea.Cmd {
	return tea.Tick(m.clearTimeout, func(t time.Time) tea.Msg {
		return ClearStatusMsg{}
	})
}
