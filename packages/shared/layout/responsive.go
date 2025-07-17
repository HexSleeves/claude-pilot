package layout

import (
	"claude-pilot/shared/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LayoutConfig holds basic configuration for layout components
type LayoutConfig struct {
	Width   int
	Height  int
	Padding int
}

// Panel represents a simple panel with basic styling
type Panel struct {
	config  LayoutConfig
	title   string
	content string
	border  bool
	focused bool
}

// NewPanel creates a new panel
func NewPanel(config LayoutConfig, title, content string, border bool) *Panel {
	return &Panel{
		config:  config,
		title:   title,
		content: content,
		border:  border,
		focused: false,
	}
}

// SetFocused sets the focused state of the panel
func (p *Panel) SetFocused(focused bool) {
	p.focused = focused
}

// SetContent updates the panel content
func (p *Panel) SetContent(content string) {
	p.content = content
}

// SetTitle updates the panel title
func (p *Panel) SetTitle(title string) {
	p.title = title
}

// Render renders the panel with simple lipgloss styling
func (p *Panel) Render() string {
	// Prepare content
	var content strings.Builder

	// Add title if provided
	if p.title != "" {
		titleStyle := styles.PanelHeaderStyle
		if p.focused {
			titleStyle = titleStyle.Foreground(styles.ClaudePrimary)
		}
		content.WriteString(titleStyle.Render(p.title) + "\n")
	}

	// Add main content
	content.WriteString(p.content)

	// Apply panel styling
	panelStyle := lipgloss.NewStyle().
		Padding(p.config.Padding)

	if p.config.Width > 0 {
		panelStyle = panelStyle.Width(p.config.Width)
	}

	if p.config.Height > 0 {
		panelStyle = panelStyle.Height(p.config.Height)
	}

	if p.border {
		borderColor := styles.TextMuted
		if p.focused {
			borderColor = styles.ClaudePrimary
		}
		panelStyle = panelStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
	}

	return panelStyle.Render(content.String())
}
