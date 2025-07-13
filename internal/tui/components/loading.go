package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingModel handles loading state display with animations
type LoadingModel struct {
	message    string
	width      int
	isLoading  bool
	spinnerIdx int

	// Styling
	loadingStyle lipgloss.Style
	errorStyle   lipgloss.Style
}

// Spinner frames for loading animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// TickMsg represents a tick for the loading animation
type TickMsg time.Time

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
	case TickMsg:
		if m.isLoading {
			m.spinnerIdx = (m.spinnerIdx + 1) % len(spinnerFrames)
			return m, m.tick()
		}
	}

	return m, nil
}

// View renders the loading state with animation
func (m LoadingModel) View() string {
	if m.message == "" || !m.isLoading {
		return ""
	}

	spinner := spinnerFrames[m.spinnerIdx]
	content := spinner + " " + m.message

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

// ViewSuccess renders a success state
func (m LoadingModel) ViewSuccess(successMsg string) string {
	if successMsg == "" {
		return ""
	}

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ECC71"))

	content := "✅ " + successMsg

	// Apply width constraint if available
	if m.width > 0 {
		style := successStyle.Width(m.width)
		return style.Render(content)
	}

	return successStyle.Render(content)
}

// StartLoading starts the loading animation
func (m *LoadingModel) StartLoading(message string) tea.Cmd {
	m.message = message
	m.isLoading = true
	m.spinnerIdx = 0
	return m.tick()
}

// StopLoading stops the loading animation
func (m *LoadingModel) StopLoading() {
	m.isLoading = false
	m.message = ""
}

// SetMessage sets the loading message
func (m *LoadingModel) SetMessage(message string) {
	m.message = message
}

// SetWidth sets the width of the loading component
func (m *LoadingModel) SetWidth(width int) {
	m.width = width
}

// IsLoading returns true if currently loading
func (m LoadingModel) IsLoading() bool {
	return m.isLoading
}

// tick returns a command that sends a TickMsg after a delay
func (m LoadingModel) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
