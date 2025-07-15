package components

import (
	"claude-pilot/shared/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CardConfig holds configuration for card rendering
type CardConfig struct {
	Width     int
	MinHeight int
	Title     string
	Icon      string
	Border    bool
	Compact   bool
}

// SummaryCard represents a card showing summary metrics
type SummaryCard struct {
	config CardConfig
	value  string
	label  string
	color  lipgloss.Color
}

// StatusCard represents a card showing status information
type StatusCard struct {
	config    CardConfig
	status    string
	indicator string
	details   []string
}

// InfoCard represents a card showing general information
type InfoCard struct {
	config  CardConfig
	content []string
}

// NewSummaryCard creates a new summary card
func NewSummaryCard(config CardConfig, value, label string, color lipgloss.Color) *SummaryCard {
	return &SummaryCard{
		config: config,
		value:  value,
		label:  label,
		color:  color,
	}
}

// NewStatusCard creates a new status card
func NewStatusCard(config CardConfig, status, indicator string, details []string) *StatusCard {
	return &StatusCard{
		config:    config,
		status:    status,
		indicator: indicator,
		details:   details,
	}
}

// NewInfoCard creates a new info card
func NewInfoCard(config CardConfig, content []string) *InfoCard {
	return &InfoCard{
		config:  config,
		content: content,
	}
}

// Render renders the summary card
func (c *SummaryCard) Render() string {
	var content strings.Builder

	// Add icon if provided
	if c.config.Icon != "" {
		content.WriteString(c.config.Icon + " ")
	}

	// Add title if provided
	if c.config.Title != "" {
		if c.config.Compact {
			content.WriteString(styles.Bold(c.config.Title) + "\n")
		} else {
			content.WriteString(styles.PanelHeaderStyle.Render(c.config.Title) + "\n")
		}
	}

	// Value styling based on color
	valueStyle := lipgloss.NewStyle().
		Foreground(c.color).
		Bold(true)

	if c.config.Compact {
		// Compact layout: value and label on same line
		content.WriteString(fmt.Sprintf("%s %s",
			valueStyle.Render(c.value),
			styles.SecondaryTextStyle.Render(c.label),
		))
	} else {
		// Full layout: value and label on separate lines
		content.WriteString(valueStyle.
			Align(lipgloss.Center).
			Width(c.config.Width-2).
			Render(c.value) + "\n")
		content.WriteString(styles.SecondaryTextStyle.
			Align(lipgloss.Center).
			Width(c.config.Width - 2).
			Render(c.label))
	}

	// Apply card styling with minimal padding
	cardStyle := lipgloss.NewStyle().
		Padding(0, 1) // Reduced padding from (1,2) to (0,1)

	if c.config.Border {
		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.color)
	}

	if c.config.Width > 0 {
		cardStyle = cardStyle.Width(c.config.Width)
	}

	if c.config.MinHeight > 0 {
		cardStyle = cardStyle.Height(c.config.MinHeight)
	}

	return cardStyle.Render(content.String())
}

// Render renders the status card
func (c *StatusCard) Render() string {
	var content strings.Builder

	// Add icon if provided
	if c.config.Icon != "" {
		content.WriteString(c.config.Icon + " ")
	}

	// Add title if provided
	if c.config.Title != "" {
		if c.config.Compact {
			content.WriteString(styles.Bold(c.config.Title) + "\n")
		} else {
			content.WriteString(styles.PanelHeaderStyle.Render(c.config.Title) + "\n")
		}
	}

	// Status with indicator
	statusLine := fmt.Sprintf("%s %s", c.indicator, c.status)
	content.WriteString(styles.HighlightStyle.Render(statusLine))

	// Add details if provided (more compact)
	if len(c.details) > 0 {
		for _, detail := range c.details {
			content.WriteString("\n" + styles.SecondaryTextStyle.Render("• "+detail))
		}
	}

	// Apply card styling with minimal padding
	cardStyle := lipgloss.NewStyle().
		Padding(0, 1) // Reduced padding from (1,2) to (0,1)

	if c.config.Border {
		// Determine border color based on status
		borderColor := styles.InfoColor
		if strings.Contains(strings.ToLower(c.status), "error") {
			borderColor = styles.ErrorColor
		} else if strings.Contains(strings.ToLower(c.status), "warning") {
			borderColor = styles.WarningColor
		} else if strings.Contains(strings.ToLower(c.status), "active") || strings.Contains(strings.ToLower(c.status), "connected") {
			borderColor = styles.SuccessColor
		}

		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
	}

	if c.config.Width > 0 {
		cardStyle = cardStyle.Width(c.config.Width)
	}

	if c.config.MinHeight > 0 {
		cardStyle = cardStyle.Height(c.config.MinHeight)
	}

	return cardStyle.Render(content.String())
}

// Render renders the info card
func (c *InfoCard) Render() string {
	var content strings.Builder

	// Add icon if provided
	if c.config.Icon != "" {
		content.WriteString(c.config.Icon + " ")
	}

	// Add title if provided
	if c.config.Title != "" {
		if c.config.Compact {
			content.WriteString(styles.Bold(c.config.Title) + "\n")
		} else {
			content.WriteString(styles.PanelHeaderStyle.Render(c.config.Title) + "\n")
		}
	}

	// Add content lines
	for i, line := range c.content {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(styles.SecondaryTextStyle.Render(line))
	}

	// Apply card styling with minimal padding
	cardStyle := lipgloss.NewStyle().
		Padding(0, 1) // Reduced padding from (1,2) to (0,1)

	if c.config.Border {
		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.TextMuted)
	}

	if c.config.Width > 0 {
		cardStyle = cardStyle.Width(c.config.Width)
	}

	if c.config.MinHeight > 0 {
		cardStyle = cardStyle.Height(c.config.MinHeight)
	}

	return cardStyle.Render(content.String())
}

// SetValue updates the value for summary cards
func (c *SummaryCard) SetValue(value string) {
	c.value = value
}

// SetStatus updates the status for status cards
func (c *StatusCard) SetStatus(status, indicator string) {
	c.status = status
	c.indicator = indicator
}

// SetDetails updates the details for status cards
func (c *StatusCard) SetDetails(details []string) {
	c.details = details
}

// SetContent updates the content for info cards
func (c *InfoCard) SetContent(content []string) {
	c.content = content
}

// CreateSessionSummaryCards creates a set of cards for session summary display
func CreateSessionSummaryCards(totalSessions, activeSessions, connectedSessions int, backend string, width int) []string {
	cardWidth := max(
		// 4 cards with spacing, reduced spacing
		(width-8)/4,
		// reduced minimum width
		10)

	// Total sessions card
	totalCard := NewSummaryCard(
		CardConfig{
			Width:     cardWidth,
			MinHeight: 3, // Reduced height
			Title:     "Total",
			Icon:      "📊",
			Border:    true,
			Compact:   true,
		},
		fmt.Sprintf("%d", totalSessions),
		"sessions",
		styles.ClaudePrimary,
	)

	// Active sessions card
	activeCard := NewSummaryCard(
		CardConfig{
			Width:     cardWidth,
			MinHeight: 3, // Reduced height
			Title:     "Active",
			Icon:      "✅",
			Border:    true,
			Compact:   true,
		},
		fmt.Sprintf("%d", activeSessions),
		"running",
		styles.SuccessColor,
	)

	// Connected sessions card
	connectedCard := NewSummaryCard(
		CardConfig{
			Width:     cardWidth,
			MinHeight: 3, // Reduced height
			Title:     "Connected",
			Icon:      "🔗",
			Border:    true,
			Compact:   true,
		},
		fmt.Sprintf("%d", connectedSessions),
		"attached",
		styles.InfoColor,
	)

	// Backend card
	backendCard := NewStatusCard(
		CardConfig{
			Width:     cardWidth,
			MinHeight: 3, // Reduced height
			Title:     "Backend",
			Icon:      "⚙️",
			Border:    true,
			Compact:   true,
		},
		backend,
		"●",
		nil,
	)

	return []string{
		totalCard.Render(),
		activeCard.Render(),
		connectedCard.Render(),
		backendCard.Render(),
	}
}

// CreateSystemStatusCard creates a card showing system status
func CreateSystemStatusCard(backend, version string, uptime string, width int) string {
	details := []string{
		fmt.Sprintf("v%s", version),  // Shortened version display
		fmt.Sprintf("up %s", uptime), // Shortened uptime display
	}

	statusCard := NewStatusCard(
		CardConfig{
			Width:     width,
			MinHeight: 3, // Reduced height
			Title:     "System",
			Icon:      "🖥️",
			Border:    true,
			Compact:   true, // Make it compact
		},
		fmt.Sprintf("%s", backend), // Shortened status
		"●",
		details,
	)

	return statusCard.Render()
}

// CreateQuickActionsCard creates a card with quick action hints
func CreateQuickActionsCard(width int) string {
	actions := []string{
		"c - Create new session",
		"Enter - Attach to session",
		"d - View session details",
		"k - Kill session",
		"q - Quit application",
	}

	infoCard := NewInfoCard(
		CardConfig{
			Width:   width,
			Title:   "Quick Actions",
			Icon:    "⌨️",
			Border:  true,
			Compact: false,
		},
		actions,
	)

	return infoCard.Render()
}
