package styles

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	evertrastable "github.com/evertras/bubble-table/table"
)

// Text utility functions
func TruncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}

// Responsive utilities - Enhanced with comprehensive adaptive behavior
func GetResponsiveWidth(terminalWidth int) (int, string) {
	switch {
	case terminalWidth < BreakpointSmall:
		return terminalWidth - 4, "small"
	case terminalWidth < BreakpointMedium:
		return terminalWidth - 8, "medium"
	default:
		return terminalWidth - 12, "large"
	}
}

// Adaptive width helper for different component types
func AdaptiveWidth(base lipgloss.Style, width int) lipgloss.Style {
	switch {
	case width < BreakpointSmall:
		return base.Width(width - 4)
	case width < BreakpointMedium:
		return base.Width(width - 8)
	default:
		return base.Width(width - 12)
	}
}

// Adaptive height helper for different component types
func AdaptiveHeight(base lipgloss.Style, height int) lipgloss.Style {
	switch {
	case height < 24:
		return base.Height(height - 2)
	case height < 40:
		return base.Height(height - 4)
	default:
		return base.Height(height - 6)
	}
}

// Responsive padding based on terminal size
func ResponsivePadding(width int) (horizontal, vertical int) {
	switch {
	case width < BreakpointSmall:
		return 1, 0
	case width < BreakpointMedium:
		return 2, 1
	default:
		return 3, 1
	}
}

// Responsive margin based on terminal size
func ResponsiveMargin(width int) (horizontal, vertical int) {
	switch {
	case width < BreakpointSmall:
		return 0, 0
	case width < BreakpointMedium:
		return 1, 0
	default:
		return 2, 1
	}
}

// Get appropriate table column widths based on terminal size
func GetTableColumnWidths(terminalWidth int, numColumns int) []int {
	availableWidth := terminalWidth - (numColumns * 3) // Account for borders and padding
	baseWidth := availableWidth / numColumns

	switch {
	case terminalWidth < BreakpointSmall:
		// Compact layout for small terminals
		return []int{baseWidth, baseWidth, baseWidth}[:numColumns]
	case terminalWidth < BreakpointMedium:
		// Balanced layout for medium terminals
		if numColumns >= 3 {
			return []int{baseWidth * 2, baseWidth, baseWidth}
		}
		return []int{baseWidth, baseWidth}[:numColumns]
	default:
		// Full layout for large terminals
		if numColumns >= 4 {
			return []int{baseWidth * 2, baseWidth, baseWidth, baseWidth}
		} else if numColumns >= 3 {
			return []int{baseWidth * 2, baseWidth, baseWidth}
		}
		return []int{baseWidth, baseWidth}[:numColumns]
	}
}

// Theme-aware color selection based on context
func GetContextualColor(context string, state string) lipgloss.Color {
	switch context {
	case "status":
		switch state {
		case "success", "active", "connected":
			return SuccessColor
		case "warning", "inactive":
			return WarningColor
		case "error", "failed":
			return ErrorColor
		case "info", "pending":
			return InfoColor
		default:
			return TextSecondary
		}
	case "action":
		switch state {
		case "primary":
			return ActionPrimary
		case "secondary":
			return ActionSecondary
		case "success":
			return ActionSuccess
		case "warning":
			return ActionWarning
		case "danger":
			return ActionDanger
		default:
			return ActionNeutral
		}
	case "interactive":
		switch state {
		case "hover":
			return HoverColor
		case "focus":
			return FocusColor
		case "active":
			return ActiveColor
		case "selected":
			return SelectedColor
		case "disabled":
			return DisabledColor
		default:
			return TextPrimary
		}
	default:
		return TextPrimary
	}
}

// Style rendering utility functions
func Title(text string) string {
	return TitleStyle.Render(text)
}

func Subtitle(text string) string {
	return SubtitleStyle.Render(text)
}

func Header(text string) string {
	return HeaderStyle.Render(text)
}

func Success(text string) string {
	return SuccessStyle.Render("✓ " + text)
}

func Error(text string) string {
	return ErrorStyle.Render("✗ " + text)
}

func Warning(text string) string {
	return WarningStyle.Render("⚠ " + text)
}

func Info(text string) string {
	return InfoStyle.Render("ℹ " + text)
}

func Bold(text string) string {
	return BoldStyle.Render(text)
}

func Dim(text string) string {
	return MutedTextStyle.Render(text)
}

func Highlight(text string) string {
	return HighlightStyle.Render(text)
}

func Arrow() string {
	return ArrowStyle.Render("→")
}

func Prompt(text string) string {
	return PromptStyle.Render(text)
}

// Status formatting functions with icons
func StatusActive(text string) string {
	return SessionStatusActiveStyle.Render("● " + text)
}

func StatusInactive(text string) string {
	return SessionStatusInactiveStyle.Render("⏸ " + text)
}

func StatusConnected(text string) string {
	return SessionStatusConnectedStyle.Render("🔗 " + text)
}

func StatusError(text string) string {
	return SessionStatusErrorStyle.Render("✗ " + text)
}

func Spinner(text string) string {
	return SpinnerStyle.Render("⠋ " + text)
}

// Box rendering functions
func MainBox(content string) string {
	return MainBoxStyle.Render(content)
}

func InfoBox(content string) string {
	return InfoBoxStyle.Render(content)
}

func WarningBox(content string) string {
	return WarningBoxStyle.Render(content)
}

func ErrorBox(content string) string {
	return ErrorBoxStyle.Render(content)
}

// Panel rendering for dashboard layouts
func Panel(content string) string {
	return PanelStyle.Render(content)
}

func PanelWithHeader(header, content string) string {
	headerRendered := PanelHeaderStyle.Render(header)
	return PanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, headerRendered, content))
}

// Horizontal line with styling
func HorizontalLine(width int) string {
	return lipgloss.NewStyle().
		Foreground(TextMuted).
		Render(lipgloss.PlaceHorizontal(width, lipgloss.Center, "─"))
}

// Create a banner with title and subtitle
func Banner(title, subtitle string) string {
	titleRendered := TitleStyle.Render(title)
	subtitleRendered := SubtitleStyle.Render(subtitle)

	content := lipgloss.JoinVertical(lipgloss.Left, titleRendered, subtitleRendered)

	return MainBoxStyle.
		Width(60).
		Align(lipgloss.Center).
		Render(content)
}

// Helper for background color application
func WithBackground(text string, background lipgloss.Color) string {
	return lipgloss.NewStyle().Background(background).Render(text)
}

// Bubbles Component Styling Functions
// These functions provide styled configurations for Bubbles components
// using the existing Claude orange theme

// GetBubblesTableStyles returns a table.Styles configuration with Claude theme
func GetBubblesTableStyles() table.Styles {
	return table.Styles{
		Header: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(ClaudePrimary).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center),
		Cell: lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1).
			Align(lipgloss.Left),
		Selected: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(SelectedColor).
			Bold(true),
	}
}

// ConfigureBubblesTextInputStyles applies styling to a textinput model
func ConfigureBubblesTextInputStyles(ti textinput.Model) textinput.Model {
	// Configure cursor style
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(ClaudePrimary)

	// Configure text styles
	ti.TextStyle = lipgloss.NewStyle().Foreground(TextPrimary)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(TextMuted)

	// Configure prompt style
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ClaudePrimary)

	return ti
}

// Evertras Table Styling Functions
// These functions provide styled configurations for evertras/bubble-table components
// using the existing Claude orange theme

// GetEvertrasTableStyles returns styling configuration for evertras table
func GetEvertrasTableStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ClaudePrimary).
		Foreground(TextPrimary)
}

// ConfigureEvertrasTable applies comprehensive Claude theme styling to an evertras table model
// with support for interactive features. Individual styling is handled through column definitions
// and row data formatting since evertras table uses a different styling approach
func ConfigureEvertrasTable(t evertrastable.Model) evertrastable.Model {
	// Apply base styling with borders and colors
	// Note: Row-level styling (hover, selection) is handled through column styles and data formatting
	return t.WithBaseStyle(GetEvertrasTableStyles())
}

// EvertrasColumnStyles holds styling configurations for different column types
type EvertrasColumnStyles struct {
	Header    lipgloss.Style
	ID        lipgloss.Style
	Name      lipgloss.Style
	Status    lipgloss.Style
	Backend   lipgloss.Style
	Timestamp lipgloss.Style
	Project   lipgloss.Style
	Messages  lipgloss.Style
}

// GetEvertrasColumnStyles returns column-specific styles for different data types
func GetEvertrasColumnStyles() EvertrasColumnStyles {
	return EvertrasColumnStyles{
		Header: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(ClaudePrimary).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center),
		ID: lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(0, 1).
			Align(lipgloss.Left),
		Name: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Left),
		Status: lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1).
			Align(lipgloss.Center),
		Backend: lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1).
			Align(lipgloss.Center),
		Timestamp: lipgloss.NewStyle().
			Foreground(TextDim).
			Padding(0, 1).
			Align(lipgloss.Center),
		Project: lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(0, 1).
			Align(lipgloss.Left),
		Messages: lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Right),
	}
}

// EvertrasRowStyles holds styling configurations for different row states
type EvertrasRowStyles struct {
	Normal    lipgloss.Style
	Selected  lipgloss.Style
	Hover     lipgloss.Style
	Alternate lipgloss.Style
}

// GetEvertrasRowStyles returns row state styles for normal, selected, hover, and alternate states
func GetEvertrasRowStyles() EvertrasRowStyles {
	return EvertrasRowStyles{
		Normal: lipgloss.NewStyle().
			Foreground(TextSecondary),
		Selected: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(SelectedColor).
			Bold(true),
		Hover: lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(HoverColor).
			Bold(true),
		Alternate: lipgloss.NewStyle().
			Foreground(TextSecondary).
			Background(BackgroundSurface),
	}
}

// EvertrasPaginationStyles holds styling configurations for pagination elements
type EvertrasPaginationStyles struct {
	Footer     lipgloss.Style
	PageInfo   lipgloss.Style
	Navigation lipgloss.Style
	Disabled   lipgloss.Style
}

// GetEvertrasPaginationStyles returns pagination footer styling configurations
func GetEvertrasPaginationStyles() EvertrasPaginationStyles {
	return EvertrasPaginationStyles{
		Footer: lipgloss.NewStyle().
			Foreground(TextMuted).
			Background(BackgroundSecondary).
			Padding(0, 1).
			Align(lipgloss.Center),
		PageInfo: lipgloss.NewStyle().
			Foreground(TextSecondary).
			Bold(true),
		Navigation: lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true),
		Disabled: lipgloss.NewStyle().
			Foreground(DisabledColor),
	}
}

// EvertrasSortStyles holds styling configurations for sort indicators
type EvertrasSortStyles struct {
	Indicator  lipgloss.Style
	Ascending  lipgloss.Style
	Descending lipgloss.Style
	Neutral    lipgloss.Style
}

// GetEvertrasSortStyles returns sort indicator styling configurations
func GetEvertrasSortStyles() EvertrasSortStyles {
	return EvertrasSortStyles{
		Indicator: lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true),
		Ascending: lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true),
		Descending: lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true),
		Neutral: lipgloss.NewStyle().
			Foreground(TextMuted),
	}
}

// ConfigureResponsiveEvertrasTable adjusts table dimensions and features based on terminal size
// Disables certain features on very small terminals and adapts column widths for different screen sizes
func ConfigureResponsiveEvertrasTable(model evertrastable.Model, width, height int) evertrastable.Model {
	// Apply base configuration
	configuredModel := ConfigureEvertrasTable(model)

	// Adjust based on terminal size
	switch {
	case width < BreakpointSmall:
		// Small terminal: minimal features, compact layout
		configuredModel = configuredModel.
			WithTargetWidth(width - 4).
			WithMinimumHeight(height - 6)
		// Disable pagination on very small screens
		if height < 20 {
			configuredModel = configuredModel.WithPageSize(0)
		}

	case width < BreakpointMedium:
		// Medium terminal: balanced features
		configuredModel = configuredModel.
			WithTargetWidth(width - 8).
			WithMinimumHeight(height - 8)

	default:
		// Large terminal: full features
		configuredModel = configuredModel.
			WithTargetWidth(width - 12).
			WithMinimumHeight(height - 10)
	}

	return configuredModel
}

// StyleSortIndicator renders sort arrows for column headers based on sort direction
func StyleSortIndicator(column, direction string) string {
	styles := GetEvertrasSortStyles()

	switch direction {
	case "asc":
		return styles.Ascending.Render("↑")
	case "desc":
		return styles.Descending.Render("↓")
	default:
		return styles.Neutral.Render("•")
	}
}

// StylePaginationInfo renders page information with consistent styling
func StylePaginationInfo(current, total int) string {
	styles := GetEvertrasPaginationStyles()

	if total <= 1 {
		return ""
	}

	info := fmt.Sprintf("Page %d of %d", current, total)
	return styles.PageInfo.Render(info)
}

// StyleRowSelection handles row state styling for selection and hover states
func StyleRowSelection(content string, isSelected, isHovered bool) string {
	styles := GetEvertrasRowStyles()

	switch {
	case isSelected:
		return styles.Selected.Render(content)
	case isHovered:
		return styles.Hover.Render(content)
	default:
		return styles.Normal.Render(content)
	}
}

// StyleFilterIndicator shows active filter status with visual feedback
func StyleFilterIndicator(filterText string) string {
	if filterText == "" {
		return ""
	}

	filterStyle := lipgloss.NewStyle().
		Foreground(InfoColor).
		Background(BackgroundSecondary).
		Bold(true).
		Padding(0, 1)

	return filterStyle.Render(fmt.Sprintf("Filter: %s", filterText))
}
