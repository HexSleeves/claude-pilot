package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpModel handles the help screen display
type HelpModel struct {
	width  int
	height int

	// Styling
	titleStyle   lipgloss.Style
	contentStyle lipgloss.Style
}

// NewHelpModel creates a new help model
func NewHelpModel() *HelpModel {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B35")).
		Bold(true).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Padding(0, 1)

	return &HelpModel{
		titleStyle:   titleStyle,
		contentStyle: contentStyle,
	}
}

// Init initializes the help model
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the help model
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the help screen
func (m HelpModel) View() string {
	title := m.titleStyle.Render("Claude Pilot - Help")

	helpText := `
Key Bindings:
  ↑/k, ↓/j    Navigate sessions
  enter       View session details
  c           Create new session
  a           Attach to selected session
  d           Delete selected session
  r           Refresh session list
  ?           Toggle this help
  q           Quit application

Session Status:
  ● Active    Session is running
  ⏸ Inactive  Session exists but not running
  🔗 Connected Someone is attached to session

Press ? again to return to session list.
`

	// Apply width constraint if terminal size is available
	if m.width > 0 {
		m.contentStyle = m.contentStyle.Width(m.width)
	}

	content := m.contentStyle.Render(helpText)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)
}

// SetDimensions sets the help model dimensions
func (m *HelpModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}
