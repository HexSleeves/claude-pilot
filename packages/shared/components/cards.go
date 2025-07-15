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

// Card defines the interface that all card types must implement
type Card interface {
	Render() string
	GetConfig() CardConfig
}

// BaseCard provides common functionality for all card types
type BaseCard struct {
	config CardConfig
}

// GetConfig returns the card configuration
func (b *BaseCard) GetConfig() CardConfig {
	return b.config
}

// buildHeader constructs the header portion (icon + title) for any card
func (b *BaseCard) buildHeader(content *strings.Builder) {
	b.addIcon(content)
	b.addTitle(content)
}

// addIcon adds the icon to the content if provided
func (b *BaseCard) addIcon(content *strings.Builder) {
	if b.config.Icon != "" {
		content.WriteString(b.config.Icon)
		content.WriteString(" ")
	}
}

// addTitle adds the title to the content if provided
func (b *BaseCard) addTitle(content *strings.Builder) {
	if b.config.Title == "" {
		return
	}

	var titleStyle string
	if b.config.Compact {
		titleStyle = styles.CardHeaderStyle.Render(b.config.Title)
	} else {
		titleStyle = styles.PanelHeaderStyle.Render(b.config.Title)
	}

	content.WriteString(titleStyle)
	content.WriteString("\n")
}

// buildCardStyle constructs the base lipgloss style for any card
func (b *BaseCard) buildCardStyle(borderColor lipgloss.Color) lipgloss.Style {
	cardStyle := styles.CardStyle

	if b.config.Border {
		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
	}

	if b.config.Width > 0 {
		cardStyle = styles.AdaptiveWidth(cardStyle, b.config.Width)
	}

	if b.config.MinHeight > 0 {
		cardStyle = cardStyle.Height(b.config.MinHeight)
	}

	return cardStyle
}

// SummaryCard represents a card showing summary metrics
type SummaryCard struct {
	BaseCard
	value string
	label string
	color lipgloss.Color
}

// StatusCard represents a card showing status information
type StatusCard struct {
	BaseCard
	status    string
	indicator string
	details   []string
}

// InfoCard represents a card showing general information
type InfoCard struct {
	BaseCard
	content []string
}

// Status color mapping for efficient lookups
var statusColorMap = map[string]lipgloss.Color{
	strings.ToLower(string(interfaces.StatusError)):     styles.ErrorColor,
	strings.ToLower(string(interfaces.StatusWarning)):   styles.WarningColor,
	strings.ToLower(string(interfaces.StatusActive)):    styles.SuccessColor,
	strings.ToLower(string(interfaces.StatusInactive)):  styles.WarningColor,
	strings.ToLower(string(interfaces.StatusConnected)): styles.InfoColor,
}

// NewSummaryCard creates a new summary card
func NewSummaryCard(config CardConfig, value, label string, color lipgloss.Color) *SummaryCard {
	return &SummaryCard{
		BaseCard: BaseCard{config: config},
		value:    value,
		label:    label,
		color:    color,
	}
}

// NewStatusCard creates a new status card
func NewStatusCard(config CardConfig, status, indicator string, details []string) *StatusCard {
	return &StatusCard{
		BaseCard:  BaseCard{config: config},
		status:    status,
		indicator: indicator,
		details:   details,
	}
}

// NewInfoCard creates a new info card
func NewInfoCard(config CardConfig, content []string) *InfoCard {
	return &InfoCard{
		BaseCard: BaseCard{config: config},
		content:  content,
	}
}

// Render renders the summary card with enhanced theming
func (c *SummaryCard) Render() string {
	content := c.buildContent()
	cardStyle := c.buildCardStyle(c.color)
	return cardStyle.Render(content)
}

// buildContent constructs the content for the summary card
func (c *SummaryCard) buildContent() string {
	var content strings.Builder
	content.Grow(128) // Pre-allocate reasonable capacity

	c.buildHeader(&content)
	c.addValueAndLabel(&content)

	return content.String()
}

// addValueAndLabel adds the value and label content
func (c *SummaryCard) addValueAndLabel(content *strings.Builder) {
	valueStyle := lipgloss.NewStyle().
		Foreground(c.color).
		Bold(true)

	if c.config.Compact {
		c.addCompactValueLabel(content, valueStyle)
	} else {
		c.addFullValueLabel(content, valueStyle)
	}
}

// addCompactValueLabel adds value and label on the same line
func (c *SummaryCard) addCompactValueLabel(content *strings.Builder, valueStyle lipgloss.Style) {
	content.WriteString(fmt.Sprintf("%s %s",
		valueStyle.Render(c.value),
		styles.CardContentStyle.Render(c.label),
	))
}

// addFullValueLabel adds value and label on separate lines with centering
func (c *SummaryCard) addFullValueLabel(content *strings.Builder, valueStyle lipgloss.Style) {
	width := c.calculateContentWidth()

	content.WriteString(valueStyle.
		Align(lipgloss.Center).
		Width(width).
		Render(c.value))
	content.WriteString("\n")
	content.WriteString(styles.CardContentStyle.
		Align(lipgloss.Center).
		Width(width).
		Render(c.label))
}

// calculateContentWidth calculates the available width for content
func (c *SummaryCard) calculateContentWidth() int {
	width := c.config.Width - 4 // Account for padding and borders
	if width < 1 {
		width = 10 // Minimum width
	}
	return width
}

// Render renders the status card with enhanced theming
func (c *StatusCard) Render() string {
	content := c.buildContent()
	borderColor := c.getBorderColor()
	cardStyle := c.buildCardStyle(borderColor)
	return cardStyle.Render(content)
}

// buildContent constructs the content for the status card
func (c *StatusCard) buildContent() string {
	var content strings.Builder
	content.Grow(256) // Pre-allocate reasonable capacity

	c.buildHeader(&content)
	c.addStatusLine(&content)
	c.addDetails(&content)

	return content.String()
}

// addStatusLine adds the status line with indicator and contextual coloring
func (c *StatusCard) addStatusLine(content *strings.Builder) {
	statusColor := c.getStatusColor()
	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Bold(true)
	statusLine := fmt.Sprintf("%s %s", c.indicator, c.status)
	content.WriteString(statusStyle.Render(statusLine))
}

// addDetails adds detail lines if any are provided
func (c *StatusCard) addDetails(content *strings.Builder) {
	for _, detail := range c.details {
		content.WriteString("\n")
		content.WriteString(styles.CardContentStyle.Render("• " + detail))
	}
}

// getBorderColor determines the appropriate border color based on status
func (c *StatusCard) getBorderColor() lipgloss.Color {
	statusLower := strings.ToLower(c.status)

	// Try contextual color first
	if color := styles.GetContextualColor("status", statusLower); color != "" {
		return color
	}

	// Fallback to status color map
	if color, exists := statusColorMap[statusLower]; exists {
		return color
	}

	// Default fallback
	return styles.InfoColor
}

// getStatusColor determines the appropriate status text color
func (c *StatusCard) getStatusColor() lipgloss.Color {
	return c.getBorderColor() // Reuse the same logic for consistency
}

// Render renders the info card with enhanced theming
func (c *InfoCard) Render() string {
	content := c.buildContent()
	cardStyle := c.buildCardStyle(styles.BackgroundSurface)
	return cardStyle.Render(content)
}

// buildContent constructs the content for the info card
func (c *InfoCard) buildContent() string {
	var content strings.Builder
	content.Grow(128) // Pre-allocate reasonable capacity

	c.buildHeader(&content)
	c.addContentLines(&content)

	return content.String()
}

// addContentLines adds all content lines with proper formatting
func (c *InfoCard) addContentLines(content *strings.Builder) {
	for i, line := range c.content {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(styles.CardContentStyle.Render(line))
	}
}

// Mutator methods for updating card state

// SetValue updates the value for summary cards
func (c *SummaryCard) SetValue(value string) {
	c.value = value
}

// SetColor updates the color for summary cards
func (c *SummaryCard) SetColor(color lipgloss.Color) {
	c.color = color
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

// AddDetail adds a single detail to status cards
func (c *StatusCard) AddDetail(detail string) {
	c.details = append(c.details, detail)
}

// SetContent updates the content for info cards
func (c *InfoCard) SetContent(content []string) {
	c.content = content
}

// AddContentLine adds a single content line to info cards
func (c *InfoCard) AddContentLine(line string) {
	c.content = append(c.content, line)
}

// Factory functions for creating specific card collections

// CreateSessionSummaryCards creates a set of cards for session summary display
func CreateSessionSummaryCards(totalSessions, activeSessions, connectedSessions int, backend string, width int) []string {
	cardWidth := calculateCardWidth(width, 4)

	cards := []*SummaryCard{
		createTotalSessionsCard(totalSessions, cardWidth),
		createActiveSessionsCard(activeSessions, cardWidth),
		createConnectedSessionsCard(connectedSessions, cardWidth),
	}

	backendCard := createBackendCard(backend, cardWidth)

	// Convert cards to rendered strings
	result := make([]string, 0, 4)
	for _, card := range cards {
		result = append(result, card.Render())
	}
	result = append(result, backendCard.Render())

	return result
}

// calculateCardWidth calculates the width for cards in a grid layout
func calculateCardWidth(totalWidth, cardCount int) int {
	// Account for spacing between cards (cardCount-1)*2
	availableWidth := totalWidth - (cardCount-1)*2
	cardWidth := availableWidth / cardCount

	// Ensure minimum width
	if cardWidth < 10 {
		return 10
	}

	return cardWidth
}

// createTotalSessionsCard creates the total sessions summary card
func createTotalSessionsCard(total, width int) *SummaryCard {
	return NewSummaryCard(
		CardConfig{
			Width:     width,
			MinHeight: 1,
			Title:     "Total",
			Icon:      "📊",
			Border:    true,
			Compact:   true,
		},
		fmt.Sprintf("%d", total),
		"sessions",
		styles.ClaudePrimary,
	)
}

// createActiveSessionsCard creates the active sessions summary card
func createActiveSessionsCard(active, width int) *SummaryCard {
	return NewSummaryCard(
		CardConfig{
			Width:     width,
			MinHeight: 1,
			Title:     "Active",
			Icon:      "✅",
			Border:    true,
			Compact:   true,
		},
		fmt.Sprintf("%d", active),
		"running",
		styles.SuccessColor,
	)
}

// createConnectedSessionsCard creates the connected sessions summary card
func createConnectedSessionsCard(connected, width int) *SummaryCard {
	return NewSummaryCard(
		CardConfig{
			Width:     width,
			MinHeight: 1,
			Title:     "Connected",
			Icon:      "🔗",
			Border:    true,
			Compact:   true,
		},
		fmt.Sprintf("%d", connected),
		"attached",
		styles.InfoColor,
	)
}

// createBackendCard creates the backend status card
func createBackendCard(backend string, width int) *StatusCard {
	return NewStatusCard(
		CardConfig{
			Width:     width,
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
}

// CreateSystemStatusCard creates a card showing system status
func CreateSystemStatusCard(backend, version, uptime string, width int) string {
	details := []string{
		fmt.Sprintf("v%s", version),
		fmt.Sprintf("up %s", uptime),
	}

	statusCard := NewStatusCard(
		CardConfig{
			Width:     width,
			MinHeight: 3,
			Title:     "System",
			Icon:      "🖥️",
			Border:    true,
			Compact:   true,
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

// Utility functions

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CardRenderer provides batch rendering capabilities for multiple cards
type CardRenderer struct {
	cards []Card
}

// NewCardRenderer creates a new card renderer
func NewCardRenderer(cards ...Card) *CardRenderer {
	return &CardRenderer{cards: cards}
}

// AddCard adds a card to the renderer
func (r *CardRenderer) AddCard(card Card) {
	r.cards = append(r.cards, card)
}

// RenderAll renders all cards and returns them as a slice of strings
func (r *CardRenderer) RenderAll() []string {
	results := make([]string, 0, len(r.cards))
	for _, card := range r.cards {
		results = append(results, card.Render())
	}
	return results
}

// RenderWithSeparator renders all cards joined by a separator
func (r *CardRenderer) RenderWithSeparator(separator string) string {
	rendered := r.RenderAll()
	return strings.Join(rendered, separator)
}
