package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingModel handles loading state display
type LoadingModel struct {
	message string
	width   int

	// Styling
	loadingStyle lipgloss.Style
	errorStyle   lipgloss.Style
}

// NewLoadingModel creates a new loading model
func NewLoadingModel() *LoadingModel {
	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498DB"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E74C3C"))

	return &LoadingModel{
		loadingStyle: loadingStyle,
		errorStyle:   errorStyle,
	}
}

// Init initializes the loading model
func (m LoadingModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the loading model
func (m LoadingModel) Update(msg tea.Msg) (LoadingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

// View renders the loading state
func (m LoadingModel) View() string {
	if m.message == "" {
		return ""
	}

	content := "🔄 " + m.message

	// Apply width constraint if available
	if m.width > 0 {
		style := m.loadingStyle.Width(m.width)
		return style.Render(content)
	}

	return m.loadingStyle.Render(content)
}

// ViewError renders an error state
func (m LoadingModel) ViewError(errorMsg string) string {
	if errorMsg == "" {
		return ""
	}

	content := "❌ " + errorMsg

	// Apply width constraint if available
	if m.width > 0 {
		style := m.errorStyle.Width(m.width)
		return style.Render(content)
	}

	return m.errorStyle.Render(content)
}

// SetMessage sets the loading message
func (m *LoadingModel) SetMessage(message string) {
	m.message = message
}

// SetWidth sets the width of the loading component
func (m *LoadingModel) SetWidth(width int) {
	m.width = width
}
