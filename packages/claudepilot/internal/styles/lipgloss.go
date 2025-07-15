package styles

import (
	"claude-pilot/shared/styles"
	"github.com/charmbracelet/lipgloss"
)

// Re-export colors from shared styles for backward compatibility
var (
	// Primary colors
	ClaudePrimary   = styles.ClaudePrimary
	ClaudeSecondary = styles.ClaudeSecondary

	// Status colors
	SuccessColor = styles.SuccessColor
	ErrorColor   = styles.ErrorColor
	WarningColor = styles.WarningColor
	InfoColor    = styles.InfoColor

	// Neutral colors
	TextPrimary   = styles.TextPrimary
	TextSecondary = styles.TextSecondary
	TextMuted     = styles.TextMuted
	TextDim       = styles.TextDim

	// Background colors
	BackgroundPrimary   = styles.BackgroundPrimary
	BackgroundSecondary = styles.BackgroundSecondary
	BackgroundAccent    = styles.BackgroundAccent
)

// Re-export styles from shared package for backward compatibility
var (
	// Title styles
	TitleStyle    = styles.TitleStyle
	SubtitleStyle = styles.SubtitleStyle
	HeaderStyle   = styles.HeaderStyle

	// Message styles
	SuccessStyle = styles.SuccessStyle
	ErrorStyle   = styles.ErrorStyle
	WarningStyle = styles.WarningStyle
	InfoStyle    = styles.InfoStyle

	// Text styles
	BoldStyle      = styles.BoldStyle
	DimStyle       = styles.MutedTextStyle
	HighlightStyle = styles.HighlightStyle

	// Interactive elements
	ArrowStyle  = styles.ArrowStyle
	PromptStyle = styles.PromptStyle

	// Table styles
	TableHeaderStyle = styles.TableHeaderStyle
	TableCellStyle   = styles.TableCellStyle

	// Status indicators
	StatusActiveStyle    = styles.SessionStatusActiveStyle
	StatusInactiveStyle  = styles.SessionStatusInactiveStyle
	StatusConnectedStyle = styles.SessionStatusConnectedStyle
	StatusErrorStyle     = styles.SessionStatusErrorStyle

	// Progress and loading
	SpinnerStyle  = styles.SpinnerStyle
	ProgressStyle = styles.ProgressStyle
)

// Re-export box styles from shared package
var (
	MainBoxStyle    = styles.MainBoxStyle
	InfoBoxStyle    = styles.InfoBoxStyle
	WarningBoxStyle = styles.WarningBoxStyle
	ErrorBoxStyle   = styles.ErrorBoxStyle
)

// Re-export utility functions from shared package for backward compatibility
func WithBackground(text string, background lipgloss.Color) string {
	return styles.WithBackground(text, background)
}

func Title(text string) string {
	return styles.Title(text)
}

func Subtitle(text string) string {
	return styles.Subtitle(text)
}

func Header(text string) string {
	return styles.Header(text)
}

func Success(text string) string {
	return styles.Success(text)
}

func Error(text string) string {
	return styles.Error(text)
}

func Warning(text string) string {
	return styles.Warning(text)
}

func Info(text string) string {
	return styles.Info(text)
}

func Bold(text string) string {
	return styles.Bold(text)
}

func BoldWithBackground(text string, background lipgloss.Color) string {
	return styles.WithBackground(styles.Bold(text), background)
}

func Dim(text string) string {
	return styles.Dim(text)
}

func Highlight(text string) string {
	return styles.Highlight(text)
}

func Arrow() string {
	return styles.Arrow()
}

func Prompt(text string) string {
	return styles.Prompt(text)
}

// Status formatting functions
func StatusActive(text string) string {
	return styles.StatusActive(text)
}

func StatusInactive(text string) string {
	return styles.StatusInactive(text)
}

func StatusConnected(text string) string {
	return styles.StatusConnected(text)
}

func StatusError(text string) string {
	return styles.StatusError(text)
}

func Spinner(text string) string {
	return styles.Spinner(text)
}

// Box rendering functions
func MainBox(content string) string {
	return styles.MainBox(content)
}

func InfoBox(content string) string {
	return styles.InfoBox(content)
}

func WarningBox(content string) string {
	return styles.WarningBox(content)
}

func ErrorBox(content string) string {
	return styles.ErrorBox(content)
}

// Horizontal line with styling
func HorizontalLine(width int) string {
	return styles.HorizontalLine(width)
}

// Create a banner with title and subtitle
func Banner(title, subtitle string) string {
	return styles.Banner(title, subtitle)
}
