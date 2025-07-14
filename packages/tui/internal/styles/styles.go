package styles

import "github.com/charmbracelet/lipgloss"

// Color palette inspired by Claude Pilot CLI
var (
	// Primary colors
	ClaudePrimary = lipgloss.Color("#FF6B35") // Claude orange
	Success       = lipgloss.Color("#2ECC71") // Success green
	Error         = lipgloss.Color("#E74C3C") // Error red
	Warning       = lipgloss.Color("#F39C12") // Warning amber
	Info          = lipgloss.Color("#3498DB") // Info blue

	// Text colors
	TextPrimary   = lipgloss.Color("#FFFFFF") // White
	TextSecondary = lipgloss.Color("#BDC3C7") // Light gray
	TextMuted     = lipgloss.Color("#7F8C8D") // Muted gray

	// Status colors
	StatusActive    = lipgloss.Color("#2ECC71") // Green
	StatusInactive  = lipgloss.Color("#F39C12") // Yellow
	StatusConnected = lipgloss.Color("#3498DB") // Blue
	StatusError     = lipgloss.Color("#E74C3C") // Red

	// Background colors
	BackgroundPrimary   = lipgloss.Color("#2C3E50") // Dark blue-gray
	BackgroundSecondary = lipgloss.Color("#34495E") // Lighter blue-gray
	BackgroundAccent    = lipgloss.Color("#1ABC9C") // Teal accent
)

// Base styles
var (
	// Title and headers
	TitleStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Info).
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

	SubtleStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info).
			Bold(true)

	// Interactive elements
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

	// Layout styles
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 0)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(TextMuted).
			Padding(1, 2)

	// Footer and navigation
	FooterStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(1, 2).
			Align(lipgloss.Left)

	// Session-specific styles
	SessionNameStyle = lipgloss.NewStyle().
				Foreground(ClaudePrimary).
				Bold(true)

	SessionIDStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	SessionStatusActiveStyle = lipgloss.NewStyle().
					Foreground(StatusActive).
					Bold(true)

	SessionStatusInactiveStyle = lipgloss.NewStyle().
					Foreground(StatusInactive).
					Bold(true)

	SessionStatusConnectedStyle = lipgloss.NewStyle().
					Foreground(StatusConnected).
					Bold(true)

	SessionStatusErrorStyle = lipgloss.NewStyle().
				Foreground(StatusError).
				Bold(true)

	// Form styles
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
)

// Helper functions for status formatting
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

// Helper function to truncate text with ellipsis
func TruncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}
