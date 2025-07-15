package ui

import (
	compat "claude-pilot/internal/styles"
	"claude-pilot/shared/styles"

	"github.com/fatih/color" // Deprecated: Use lipgloss styles from shared theme
)

// CLI Color Scheme - Migrated to Shared Theme
// This file provides CLI-optimized styling functions using the unified theme
// from packages/shared/styles while maintaining backward compatibility

// Deprecated: Legacy fatih/color variables - Use shared theme styles instead
// These are kept for backward compatibility but should be migrated to lipgloss
var (
	// Primary colors (DEPRECATED - use sharedStyles.ClaudePrimary)
	ClaudePrimary = color.New(color.FgHiRed).Add(color.Bold)    // Use sharedStyles.TitleStyle instead
	Success       = color.New(color.FgHiGreen).Add(color.Bold)  // Use sharedStyles.SuccessStyle instead
	Error         = color.New(color.FgHiRed).Add(color.Bold)    // Use sharedStyles.ErrorStyle instead
	Warning       = color.New(color.FgHiYellow).Add(color.Bold) // Use sharedStyles.WarningStyle instead
	Info          = color.New(color.FgHiCyan).Add(color.Bold)   // Use sharedStyles.InfoStyle instead

	// Text colors (DEPRECATED - use shared theme text styles)
	TextPrimary   = color.New(color.FgHiWhite) // Use sharedStyles.PrimaryTextStyle instead
	TextSecondary = color.New(color.FgWhite)   // Use sharedStyles.SecondaryTextStyle instead
	TextMuted     = color.New(color.FgHiBlack) // Use sharedStyles.MutedTextStyle instead

	// Status colors (DEPRECATED - use shared theme status styles)
	StatusActive    = color.New(color.FgHiGreen)  // Use sharedStyles.SessionStatusActiveStyle instead
	StatusInactive  = color.New(color.FgHiYellow) // Use sharedStyles.SessionStatusInactiveStyle instead
	StatusConnected = color.New(color.FgHiCyan)   // Use sharedStyles.SessionStatusConnectedStyle instead
	StatusError     = color.New(color.FgHiRed)    // Use sharedStyles.SessionStatusErrorStyle instead

	// Accent colors (DEPRECATED - use shared theme accent colors)
	Accent    = color.New(color.FgHiMagenta) // Use sharedStyles.AccentStyle instead
	AccentDim = color.New(color.FgMagenta)   // Use sharedStyles.MutedStyle instead
)

// CLI-optimized styling functions using the shared theme
// These functions provide terminal-friendly output while using the unified color palette

// Style functions for consistent formatting - Migrated to shared theme
// These functions now use the shared theme through lipgloss for consistency

func Title(text string) string {
	// Use shared theme title style instead of fatih/color
	return styles.TitleStyle.Render(text)
}

func Subtitle(text string) string {
	// Use shared theme subtitle style
	return styles.SubtitleStyle.Render(text)
}

func SuccessMsg(text string) string {
	// Use shared theme success style with icon
	return styles.Success(text)
}

func ErrorMsg(text string) string {
	// Use shared theme error style with icon
	return styles.Error(text)
}

func WarningMsg(text string) string {
	// Use shared theme warning style with icon
	return styles.Warning(text)
}

func InfoMsg(text string) string {
	// Use shared theme info style with icon
	return styles.Info(text)
}

func Highlight(text string) string {
	// Use shared theme highlight style
	return styles.Highlight(text)
}

func Dim(text string) string {
	// Use shared theme dim/muted style
	return styles.Dim(text)
}

func Bold(text string) string {
	// Use shared theme bold style
	return styles.Bold(text)
}

// Status formatting - Migrated to shared theme for consistency
func FormatStatus(status string) string {
	switch status {
	case "active":
		return styles.StatusActive(status)
	case "inactive":
		return styles.StatusInactive(status)
	case "connected":
		return styles.StatusConnected(status)
	case "error":
		return styles.StatusError(status)
	default:
		return styles.Dim("● " + status)
	}
}

// Process status formatting - Enhanced with shared theme
func FormatProcessStatus(status string) string {
	switch status {
	case "running":
		return styles.SuccessStyle.Render("▶ " + status)
	case "starting":
		return styles.WarningStyle.Render("⏳ " + status)
	case "stopped":
		return styles.Dim(status)
	case "error":
		return styles.ErrorStyle.Render("✗ " + status)
	default:
		return styles.Dim("? " + status)
	}
}

// Tmux status formatting - Enhanced with shared theme
func FormatTmuxStatus(status string) string {
	switch status {
	case "running":
		return styles.StatusActive(status)
	case "attached":
		return styles.StatusConnected(status)
	case "stopped":
		return styles.Dim("⏸ " + status)
	case "error":
		return styles.StatusError(status)
	default:
		return styles.Dim("? " + status)
	}
}

// Progress indicators - Enhanced with shared theme
func Spinner(text string) string {
	return styles.Spinner(text)
}

func CheckMark() string {
	return styles.SuccessStyle.Render("✓")
}

func CrossMark() string {
	return styles.ErrorStyle.Render("✗")
}

func Arrow() string {
	return styles.Arrow()
}

// Borders and separators - Enhanced with shared theme
func HorizontalLine(length int) string {
	return styles.HorizontalLine(length)
}

func VerticalSeparator() string {
	return styles.Dim("│")
}

// Interactive prompts - Enhanced with shared theme
func Prompt(text string) string {
	return styles.Prompt(text)
}

func Input(text string) string {
	return styles.TitleStyle.Render("> ") + text
}

// Enhanced lipgloss-based functions for better visual presentation
// These provide improved styling while maintaining backward compatibility

// Enhanced session summary with better formatting
func SessionSummary(total, active, inactive int, showAll bool) string {
	return compat.SessionSummary(total, active, inactive, showAll)
}

// Enhanced next steps display
func NextSteps(commands ...string) string {
	return compat.NextSteps(commands...)
}

// Enhanced available commands display
func AvailableCommands(commands ...string) string {
	return compat.AvailableCommands(commands...)
}

// Enhanced session details formatting
func SessionDetailsFormatted(sessionID, name, status, backend, created, project, description string) string {
	return compat.SessionDetails(sessionID, name, status, backend, created, project, description)
}

// Enhanced available sessions list
func AvailableSessionsList(sessions []compat.SessionInfo) string {
	return compat.AvailableSessionsList(sessions)
}

// Enhanced command list formatting
func CommandList(commands map[string]string) string {
	return compat.CommandList(commands)
}

// Enhanced header formatting
func Header(text string) string {
	return styles.Header(text)
}
