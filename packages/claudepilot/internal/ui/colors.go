package ui

import (
	"claude-pilot/internal/styles"

	"github.com/fatih/color"
)

// Color scheme constants matching the design specification
var (
	// Primary colors
	ClaudePrimary = color.New(color.FgHiRed).Add(color.Bold)    // #FF6B35 (Claude orange)
	Success       = color.New(color.FgHiGreen).Add(color.Bold)  // #2ECC71 (success green)
	Error         = color.New(color.FgHiRed).Add(color.Bold)    // #E74C3C (error red)
	Warning       = color.New(color.FgHiYellow).Add(color.Bold) // #F39C12 (warning amber)
	Info          = color.New(color.FgHiCyan).Add(color.Bold)   // Info blue

	// Text colors
	TextPrimary   = color.New(color.FgHiWhite) // #FFFFFF (white)
	TextSecondary = color.New(color.FgWhite)   // Dimmed white
	TextMuted     = color.New(color.FgHiBlack) // Gray text

	// Status colors
	StatusActive    = color.New(color.FgHiGreen)
	StatusInactive  = color.New(color.FgHiYellow)
	StatusConnected = color.New(color.FgHiCyan)
	StatusError     = color.New(color.FgHiRed)

	// Accent colors
	Accent    = color.New(color.FgHiMagenta)
	AccentDim = color.New(color.FgMagenta)
)

// Style functions for consistent formatting
func Title(text string) string {
	return ClaudePrimary.Sprint(text)
}

func Subtitle(text string) string {
	return Info.Sprint(text)
}

func SuccessMsg(text string) string {
	return Success.Sprint("✓ " + text)
}

func ErrorMsg(text string) string {
	return Error.Sprint("✗ " + text)
}

func WarningMsg(text string) string {
	return Warning.Sprint("⚠ " + text)
}

func InfoMsg(text string) string {
	return Info.Sprint("ℹ " + text)
}

func Highlight(text string) string {
	return ClaudePrimary.Sprint(text)
}

func Dim(text string) string {
	return TextMuted.Sprint(text)
}

func Bold(text string) string {
	return color.New(color.Bold).Sprint(text)
}

// Status formatting
func FormatStatus(status string) string {
	switch status {
	case "active":
		return StatusActive.Sprint("●") + " " + StatusActive.Sprint(status)
	case "inactive":
		return StatusInactive.Sprint("●") + " " + StatusInactive.Sprint(status)
	case "connected":
		return StatusConnected.Sprint("●") + " " + StatusConnected.Sprint(status)
	case "error":
		return StatusError.Sprint("●") + " " + StatusError.Sprint(status)
	default:
		return TextMuted.Sprint("●") + " " + TextMuted.Sprint(status)
	}
}

// Process status formatting
func FormatProcessStatus(status string) string {
	switch status {
	case "running":
		return StatusActive.Sprint("▶") + " " + StatusActive.Sprint(status)
	case "starting":
		return StatusInactive.Sprint("⏳") + " " + StatusInactive.Sprint(status)
	case "stopped":
		return TextMuted.Sprint("⏸") + " " + TextMuted.Sprint(status)
	case "error":
		return StatusError.Sprint("✗") + " " + StatusError.Sprint(status)
	default:
		return TextMuted.Sprint("?") + " " + TextMuted.Sprint(status)
	}
}

// Tmux status formatting
func FormatTmuxStatus(status string) string {
	switch status {
	case "running":
		return StatusActive.Sprint("●") + " " + StatusActive.Sprint(status)
	case "attached":
		return StatusConnected.Sprint("🔗") + " " + StatusConnected.Sprint(status)
	case "stopped":
		return TextMuted.Sprint("⏸") + " " + TextMuted.Sprint(status)
	case "error":
		return StatusError.Sprint("✗") + " " + StatusError.Sprint(status)
	default:
		return TextMuted.Sprint("?") + " " + TextMuted.Sprint(status)
	}
}

// Progress indicators
func Spinner(text string) string {
	return Info.Sprint("⠋ " + text)
}

func CheckMark() string {
	return Success.Sprint("✓")
}

func CrossMark() string {
	return Error.Sprint("✗")
}

func Arrow() string {
	return ClaudePrimary.Sprint("→")
}

// Borders and separators
func HorizontalLine(length int) string {
	line := ""
	for range length {
		line += "─"
	}
	return TextMuted.Sprint(line)
}

func VerticalSeparator() string {
	return TextMuted.Sprint("│")
}

// Interactive prompts
func Prompt(text string) string {
	return ClaudePrimary.Sprint("? ") + TextPrimary.Sprint(text)
}

func Input(text string) string {
	return ClaudePrimary.Sprint("> ") + TextPrimary.Sprint(text)
}

// Enhanced lipgloss-based functions for better visual presentation
// These provide improved styling while maintaining backward compatibility

// Enhanced banner for root command
func RootBanner() string {
	return styles.RootBanner()
}

// Enhanced session summary with better formatting
func SessionSummary(total, active, inactive int, showAll bool) string {
	return styles.SessionSummary(total, active, inactive, showAll)
}

// Enhanced next steps display
func NextSteps(commands ...string) string {
	return styles.NextSteps(commands...)
}

// Enhanced available commands display
func AvailableCommands(commands ...string) string {
	return styles.AvailableCommands(commands...)
}

// Enhanced session details formatting
func SessionDetailsFormatted(sessionID, name, status, backend, created, project, description string) string {
	return styles.SessionDetails(sessionID, name, status, backend, created, project, description)
}

// Enhanced available sessions list
func AvailableSessionsList(sessions []styles.SessionInfo) string {
	return styles.AvailableSessionsList(sessions)
}

// Enhanced command list formatting
func CommandList(commands map[string]string) string {
	return styles.CommandList(commands)
}

// Enhanced header formatting
func Header(text string) string {
	return styles.Header(text)
}
