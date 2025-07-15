package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette inspired by Claude's branding
var (
	// Primary colors
	ClaudePrimary   = lipgloss.Color("#FF6B35") // Claude orange
	ClaudeSecondary = lipgloss.Color("#4A90E2") // Claude blue

	// Status colors
	SuccessColor = lipgloss.Color("#28A745") // Green
	ErrorColor   = lipgloss.Color("#DC3545") // Red
	WarningColor = lipgloss.Color("#FFC107") // Yellow
	InfoColor    = lipgloss.Color("#17A2B8") // Cyan

	// Neutral colors
	TextPrimary   = lipgloss.Color("#FFFFFF") // White
	TextSecondary = lipgloss.Color("#E5E5E5") // Light gray
	TextMuted     = lipgloss.Color("#6C757D") // Gray
	TextDim       = lipgloss.Color("#495057") // Dark gray

	// Background colors
	BackgroundPrimary   = lipgloss.Color("#1A1A1A") // Dark
	BackgroundSecondary = lipgloss.Color("#2D2D2D") // Lighter dark
	BackgroundAccent    = lipgloss.Color("#3A3A3A") // Accent dark
)

// Base styles
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	// Header styles with borders
	HeaderStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(BackgroundSecondary).
			Bold(true).
			Padding(0, 2).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ClaudePrimary)

	// Message styles
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

	// Text styles
	BoldStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextPrimary)

	DimStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)

	// Interactive elements
	ArrowStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	// Code and structured data
	CodeBlockStyle = lipgloss.NewStyle().
			Background(BackgroundSecondary).
			Foreground(TextSecondary).
			Padding(1, 2).
			Margin(1, 0).
			Border(lipgloss.NormalBorder()).
			BorderForeground(TextMuted)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(TextPrimary).
				Background(BackgroundAccent).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 1)

	// Status indicators
	StatusActiveStyle = lipgloss.NewStyle().
				Foreground(SuccessColor).
				Bold(true)

	StatusInactiveStyle = lipgloss.NewStyle().
				Foreground(WarningColor).
				Bold(true)

	StatusConnectedStyle = lipgloss.NewStyle().
				Foreground(InfoColor).
				Bold(true)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(ErrorColor).
				Bold(true)

	// Progress and loading
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(ClaudeSecondary).
			Bold(true)

	ProgressStyle = lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true)
)

// Box styles for containers
var (
	// Main container for CLI output
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ClaudePrimary).
			Padding(1, 2).
			Margin(1, 0)

	// Info box for help text and descriptions
	InfoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(InfoColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	// Warning box for important notices
	WarningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(WarningColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)

	// Error box for error messages
	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ErrorColor).
			Background(BackgroundSecondary).
			Padding(1, 2).
			Margin(1, 0)
)

// Utility function to apply a background color to a string
func WithBackground(text string, background lipgloss.Color) string {
	return lipgloss.NewStyle().Background(background).Render(text)
}

// Utility functions for common styling patterns
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

func BoldWithBackground(text string, background lipgloss.Color) string {
	return BoldStyle.Background(background).Render(text)
}

func Dim(text string) string {
	return DimStyle.Render(text)
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

func CodeBlock(text string) string {
	return CodeBlockStyle.Render(text)
}

// Status formatting functions
func StatusActive(text string) string {
	return StatusActiveStyle.Render("● " + text)
}

func StatusInactive(text string) string {
	return StatusInactiveStyle.Render("⏸ " + text)
}

func StatusConnected(text string) string {
	return StatusConnectedStyle.Render("🔗 " + text)
}

func StatusError(text string) string {
	return StatusErrorStyle.Render("✗ " + text)
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
