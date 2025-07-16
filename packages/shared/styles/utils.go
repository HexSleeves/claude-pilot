package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
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

// ConfigureBubblesViewportStyles applies styling to a viewport model
func ConfigureBubblesViewportStyles(vp viewport.Model) viewport.Model {
	// Note: Viewport styling is limited in bubbles v0.21.0
	// Most styling is handled through the content
	return vp
}

// ConfigureBubblesHelpStyles applies styling to a help model
func ConfigureBubblesHelpStyles(h help.Model) help.Model {
	// Configure help styles
	h.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(ClaudePrimary).
		Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(TextSecondary)
	h.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(TextMuted)
	h.Styles.FullKey = lipgloss.NewStyle().
		Foreground(ClaudePrimary).
		Bold(true)
	h.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(TextSecondary)
	h.Styles.FullSeparator = lipgloss.NewStyle().
		Foreground(TextMuted)

	return h
}

// ConfigureBubblesTable applies Claude theme to a table model
func ConfigureBubblesTable(t table.Model) table.Model {
	t.SetStyles(GetBubblesTableStyles())
	return t
}

// ConfigureBubblesTextInput applies Claude theme to a textinput model
func ConfigureBubblesTextInput(ti textinput.Model) textinput.Model {
	return ConfigureBubblesTextInputStyles(ti)
}

// ConfigureBubblesViewport applies Claude theme to a viewport model
func ConfigureBubblesViewport(vp viewport.Model) viewport.Model {
	return ConfigureBubblesViewportStyles(vp)
}

// ConfigureBubblesHelp applies Claude theme to a help model
func ConfigureBubblesHelp(h help.Model) help.Model {
	return ConfigureBubblesHelpStyles(h)
}
