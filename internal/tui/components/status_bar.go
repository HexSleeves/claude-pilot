package components

import (
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
)

// StatusBarModel handles status display and feedback
type StatusBarModel struct {
	state   StatusState
	message string
	width   int

	// Styling
	loadingStyle lipgloss.Style
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
}

// NewStatusBarModel creates a new status bar model
func NewStatusBarModel() *StatusBarModel {
	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498DB"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ECC71"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E74C3C"))

	return &StatusBarModel{
		state:        StatusIdle,
		loadingStyle: loadingStyle,
		successStyle: successStyle,
		errorStyle:   errorStyle,
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
	}

	return m, nil
}

// View renders the status bar
func (m StatusBarModel) View() string {
	if m.state == StatusIdle || m.message == "" {
		return ""
	}

	var style lipgloss.Style
	switch m.state {
	case StatusLoading:
		style = m.loadingStyle
	case StatusSuccess:
		style = m.successStyle
	case StatusError:
		style = m.errorStyle
	default:
		style = m.loadingStyle
	}

	// Apply width constraint if available
	if m.width > 0 {
		style = style.Width(m.width)
	}

	return style.Render(m.message)
}

// SetLoading sets the status bar to loading state
func (m *StatusBarModel) SetLoading(message string) {
	m.state = StatusLoading
	m.message = message
}

// SetSuccess sets the status bar to success state
func (m *StatusBarModel) SetSuccess(message string) {
	m.state = StatusSuccess
	m.message = message
}

// SetError sets the status bar to error state
func (m *StatusBarModel) SetError(message string) {
	m.state = StatusError
	m.message = message
}

// Clear clears the status bar
func (m *StatusBarModel) Clear() {
	m.state = StatusIdle
	m.message = ""
}

// IsVisible returns true if the status bar should be displayed
func (m StatusBarModel) IsVisible() bool {
	return m.state != StatusIdle && m.message != ""
}

// SetWidth sets the width of the status bar
func (m *StatusBarModel) SetWidth(width int) {
	m.width = width
}
