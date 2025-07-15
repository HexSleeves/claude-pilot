package components

import (
	"claude-pilot/shared/interfaces"
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

// Render renders the summary card - Enhanced with standardized theming
func (c *SummaryCard) Render() string {
	var content strings.Builder

	// Add icon if provided with consistent spacing
	if c.config.Icon != "" {
		content.WriteString(c.config.Icon + " ")
	}

	// Add title with enhanced styling
	if c.config.Title != "" {
		if c.config.Compact {
			content.WriteString(styles.CardHeaderStyle.Render(c.config.Title) + "\n")
		} else {
			content.WriteString(styles.PanelHeaderStyle.Render(c.config.Title) + "\n")
		}
	}

	// Enhanced value styling with theme-aware colors
	valueStyle := lipgloss.NewStyle().
		Foreground(c.color).
		Bold(true)

	if c.config.Compact {
		// Compact layout: value and label on same line with better spacing
		content.WriteString(fmt.Sprintf("%s %s",
			valueStyle.Render(c.value),
			styles.CardContentStyle.Render(c.label),
		))
	} else {
		// Full layout: value and label on separate lines with responsive width
		width := c.config.Width - 4 // Account for padding and borders
		if width < 1 {
			width = 10 // Minimum width
		}

		content.WriteString(valueStyle.
			Align(lipgloss.Center).
			Width(width).
			Render(c.value) + "\n")
		content.WriteString(styles.CardContentStyle.
			Align(lipgloss.Center).
			Width(width).
			Render(c.label))
	}

	// Apply enhanced card styling using theme styles
	cardStyle := styles.CardStyle

	if c.config.Border {
		// Use theme-aware border with contextual coloring
		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.color)
	}

	// Apply responsive sizing
	if c.config.Width > 0 {
		cardStyle = styles.AdaptiveWidth(cardStyle, c.config.Width)
	}

	if c.config.MinHeight > 0 {
		cardStyle = cardStyle.Height(c.config.MinHeight)
	}

	return cardStyle.Render(content.String())
}

// Render renders the status card - Enhanced with standardized theming
func (c *StatusCard) Render() string {
	var content strings.Builder

	// Add icon with consistent spacing
	if c.config.Icon != "" {
		content.WriteString(c.config.Icon + " ")
	}

	// Add title with enhanced styling
	if c.config.Title != "" {
		if c.config.Compact {
			content.WriteString(styles.CardHeaderStyle.Render(c.config.Title) + "\n")
		} else {
			content.WriteString(styles.PanelHeaderStyle.Render(c.config.Title) + "\n")
		}
	}

	// Enhanced status with indicator using contextual colors
	statusColor := styles.GetContextualColor("status", strings.ToLower(c.status))
	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Bold(true)
	statusLine := fmt.Sprintf("%s %s", c.indicator, c.status)
	content.WriteString(statusStyle.Render(statusLine))

	// Add details with improved formatting
	if len(c.details) > 0 {
		for _, detail := range c.details {
			content.WriteString("\n" + styles.CardContentStyle.Render("• "+detail))
		}
	}

	// Apply enhanced card styling using theme styles
	cardStyle := styles.CardStyle

	if c.config.Border {
		// Use contextual border color based on status
		borderColor := styles.GetContextualColor("status", strings.ToLower(c.status))

		// Fallback to interface-based matching if contextual color doesn't work
		statusLower := strings.ToLower(c.status)
		switch statusLower {
		case strings.ToLower(string(interfaces.StatusError)):
			borderColor = styles.ErrorColor
		case strings.ToLower(string(interfaces.StatusWarning)):
			borderColor = styles.WarningColor
		case strings.ToLower(string(interfaces.StatusActive)):
			borderColor = styles.SuccessColor
		case strings.ToLower(string(interfaces.StatusInactive)):
			borderColor = styles.WarningColor
		case strings.ToLower(string(interfaces.StatusConnected)):
			borderColor = styles.InfoColor
		}

		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
	}

	// Apply responsive sizing
	if c.config.Width > 0 {
		cardStyle = styles.AdaptiveWidth(cardStyle, c.config.Width)
	}

	if c.config.MinHeight > 0 {
		cardStyle = cardStyle.Height(c.config.MinHeight)
	}

	return cardStyle.Render(content.String())
}

// Render renders the info card - Enhanced with standardized theming
func (c *InfoCard) Render() string {
	var content strings.Builder

	// Add icon with consistent spacing
	if c.config.Icon != "" {
		content.WriteString(c.config.Icon + " ")
	}

	// Add title with enhanced styling
	if c.config.Title != "" {
		if c.config.Compact {
			content.WriteString(styles.CardHeaderStyle.Render(c.config.Title) + "\n")
		} else {
			content.WriteString(styles.PanelHeaderStyle.Render(c.config.Title) + "\n")
		}
	}

	// Add content lines with improved formatting
	for i, line := range c.content {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(styles.CardContentStyle.Render(line))
	}

	// Apply enhanced card styling using theme styles
	cardStyle := styles.CardStyle

	if c.config.Border {
		// Use neutral border color for info cards
		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.BackgroundSurface)
	}

	// Apply responsive sizing
	if c.config.Width > 0 {
		cardStyle = styles.AdaptiveWidth(cardStyle, c.config.Width)
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
			MinHeight: 1,
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
			MinHeight: 1,
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
			MinHeight: 1,
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
			MinHeight: 1,
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
			MinHeight: 3,
			Title:     "System",
			Icon:      "🖥️",
			Border:    true,
			Compact:   true, // Make it compact
		},
		backend,
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
