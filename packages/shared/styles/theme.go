package styles

import "github.com/charmbracelet/lipgloss"

// Responsive breakpoints for adaptive layouts
const (
	BreakpointSmall  = 80  // Small terminal
	BreakpointMedium = 120 // Medium terminal
	BreakpointLarge  = 160 // Large terminal
)

// Color palette inspired by Claude Pilot CLI - unified across TUI and CLI
var (
	// Primary colors
	ClaudePrimary   = lipgloss.Color("#FF6B35") // Claude orange
	ClaudeSecondary = lipgloss.Color("#6BB6FF") // More readable Claude blue

	// Status colors
	SuccessColor = lipgloss.Color("#2ECC71") // Success green
	ErrorColor   = lipgloss.Color("#E74C3C") // Error red
	WarningColor = lipgloss.Color("#F39C12") // Warning amber
	InfoColor    = lipgloss.Color("#5DADE2") // Lighter, more readable blue

	// Text colors
	TextPrimary   = lipgloss.Color("#FFFFFF") // White
	TextSecondary = lipgloss.Color("#D5DBDB") // More readable light gray
	TextMuted     = lipgloss.Color("#AEB6BF") // More readable muted gray
	TextDim       = lipgloss.Color("#495057") // Dark gray

	// Background colors
	BackgroundPrimary   = lipgloss.Color("#2C3E50") // Dark blue-gray
	BackgroundSecondary = lipgloss.Color("#34495E") // Lighter blue-gray
	BackgroundAccent    = lipgloss.Color("#58D68D") // More readable green accent

	// Session status colors
	StatusActiveColor    = SuccessColor // Green
	StatusInactiveColor  = WarningColor // Yellow
	StatusConnectedColor = InfoColor    // Blue
	StatusErrorColor     = ErrorColor   // Red
)

// Typography styles
var (
	// Title and headers
	TitleStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true).
			Underline(true)

	// Text styles
	PrimaryTextStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	SecondaryTextStyle = lipgloss.NewStyle().
				Foreground(TextSecondary)

	MutedTextStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(TextDim)

	BoldStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextPrimary)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)
)

// Status and state styles
var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true)

	// Session status styles
	SessionStatusActiveStyle = lipgloss.NewStyle().
					Foreground(StatusActiveColor).
					Bold(true)

	SessionStatusInactiveStyle = lipgloss.NewStyle().
					Foreground(StatusInactiveColor).
					Bold(true)

	SessionStatusConnectedStyle = lipgloss.NewStyle().
					Foreground(StatusConnectedColor).
					Bold(true)

	SessionStatusErrorStyle = lipgloss.NewStyle().
				Foreground(StatusErrorColor).
				Bold(true)
)

// Interactive element styles
var (
	// Selection and focus states
	SelectedStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(ClaudePrimary).
			Bold(true).
			Padding(0, 1)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1)

	FocusedStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)

	BlurredStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(ClaudePrimary).
			Bold(true).
			Padding(0, 2).
			Margin(0, 1)

	ButtonFocusedStyle = lipgloss.NewStyle().
				Foreground(ClaudePrimary).
				Background(TextPrimary).
				Bold(true).
				Padding(0, 2).
				Margin(0, 1)

	// Form elements
	InputStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(BackgroundSecondary).
			Padding(0, 1).
			Margin(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(BackgroundSecondary).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ClaudePrimary).
				Padding(0, 1).
				Margin(0, 1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Bold(true).
			Margin(0, 1)
)

// Layout and container styles
var (
	// Basic containers
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(TextMuted).
			Padding(1, 2)

	// Box styles
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ClaudePrimary).
			Padding(1, 2).
			Margin(1, 0)

	InfoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(InfoColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	WarningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(WarningColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ErrorColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	// Panel styles for dashboard layout
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(TextMuted).
			Padding(1).
			Margin(0, 1)

	PanelHeaderStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(BackgroundAccent).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)
)

// Table styles
var (
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(BackgroundAccent).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1)

	TableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(ClaudePrimary).
				Bold(true)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)
)

// Footer and navigation
var (
	FooterStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(1, 2).
			Align(lipgloss.Left)

	ArrowStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)
)

// Session-specific styles
var (
	SessionNameStyle = lipgloss.NewStyle().
				Foreground(ClaudePrimary).
				Bold(true)

	SessionIDStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	SessionDescriptionStyle = lipgloss.NewStyle().
				Foreground(TextSecondary)
)

// Progress and loading styles
var (
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	ProgressStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)
)

// Utility functions for status formatting
func FormatSessionStatus(status string) lipgloss.Style {
	switch status {
	case "active":
		return SessionStatusActiveStyle
	case "inactive":
		return SessionStatusInactiveStyle
	case "connected":
		return SessionStatusConnectedStyle
	case "error":
		return SessionStatusErrorStyle
	default:
		return MutedTextStyle
	}
}

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

// Responsive sizing helper
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

	content := lipgloss.JoinVertical(lipgloss.Center, titleRendered, subtitleRendered)

	return MainBoxStyle.
		Width(60).
		Align(lipgloss.Center).
		Render(content)
}

// Helper for background color application
func WithBackground(text string, background lipgloss.Color) string {
	return lipgloss.NewStyle().Background(background).Render(text)
}
