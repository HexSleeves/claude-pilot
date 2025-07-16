package styles

import (
	"claude-pilot/shared/interfaces"
	"claude-pilot/shared/styles"
	"claude-pilot/shared/utils"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SessionInfo represents basic session information for display
type SessionInfo struct {
	ID   string
	Name string
}

// Compatibility functions that match the existing UI interface
// These functions provide the same API as the current ui/colors.go but with lipgloss styling

// Status formatting functions
func FormatStatus(status string) string {
	switch strings.ToLower(status) {
	case "active", "running":
		return styles.StatusActive(status)
	case "inactive", "stopped":
		return styles.StatusInactive(status)
	case "connected", "attached":
		return styles.StatusConnected(status)
	case "error", "failed":
		return styles.StatusError(status)
	default:
		return styles.Dim(status)
	}
}

// Tmux status formatting (matching existing interface)
func FormatTmuxStatus(status string) string {
	switch status {
	case "running":
		return styles.StatusActive(status)
	case "attached":
		return styles.StatusConnected(status)
	case "stopped":
		return styles.StatusInactive(status)
	case "error":
		return styles.StatusError(status)
	default:
		return styles.Dim("? " + status)
	}
}

// Progress indicators
func CheckMark() string {
	return styles.SuccessStyle.Render("✓")
}

func CrossMark() string {
	return styles.ErrorStyle.Render("✗")
}

// Time formatting
func FormatTime(t time.Time) string {
	return styles.Dim(t.Format("2006-01-02 15:04"))
}

// Text truncation with styling
func TruncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + styles.Dim("...")
}

// Enhanced banner for the root command
func RootBanner() string {
	// Use shared theme banner function for consistency
	return styles.Banner("Claude Pilot 🚀", "A powerful CLI tool for managing multiple Claude code sessions") // TODO: Add banner
}

// Enhanced command list formatting
func CommandList(commands map[string]string) string {
	var lines []string

	// Header with enhanced styling
	header := InfoStyle.Render("Available Commands:")
	lines = append(lines, header)
	lines = append(lines, "")

	// Commands with improved formatting
	for cmd, desc := range commands {
		cmdStyled := lipgloss.NewStyle().
			Foreground(ClaudePrimary).
			Bold(true).
			Width(12).
			Render(cmd)

		descStyled := lipgloss.NewStyle().
			Foreground(TextSecondary).
			Render(desc)

		line := fmt.Sprintf("  %s %s", cmdStyled, descStyled)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// Session summary formatting - Enhanced with standardized theme
func SessionSummary(total, active, inactive int) string {
	var parts []string

	if total > 0 {
		// Use enhanced card-style formatting
		totalText := fmt.Sprintf("Total: %d", total)
		parts = append(parts, lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(BackgroundSecondary).
			Bold(true).
			Padding(0, 1).
			Render(totalText))

		if active > 0 {
			activeText := fmt.Sprintf("Active: %d", active)
			parts = append(parts, lipgloss.NewStyle().
				Foreground(SuccessColor).
				Background(BackgroundSecondary).
				Bold(true).
				Padding(0, 1).
				Render(activeText))
		}

		if inactive > 0 {
			inactiveText := fmt.Sprintf("Inactive: %d", inactive)
			parts = append(parts, lipgloss.NewStyle().
				Foreground(WarningColor).
				Background(BackgroundSecondary).
				Bold(true).
				Padding(0, 1).
				Render(inactiveText))
		}
	}

	if len(parts) == 0 {
		return styles.Dim("No sessions found")
	}

	// Join with enhanced separator
	separator := lipgloss.NewStyle().
		Foreground(TextMuted).
		Background(BackgroundSecondary).
		Render(" | ")
	summary := strings.Join(parts, separator)

	return styles.InfoBox(summary)
}

// Next steps formatting
func NextSteps(commands ...string) string {
	var lines []string

	header := styles.Info("Next steps:")
	lines = append(lines, header)

	for _, cmd := range commands {
		line := fmt.Sprintf("  %s %s", styles.Arrow(), styles.Highlight(cmd))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// Available commands formatting
func AvailableCommands(commands ...string) string {
	var lines []string

	header := styles.Info("Available commands:")
	lines = append(lines, header)

	for _, cmd := range commands {
		line := fmt.Sprintf("  %s %s", styles.Arrow(), styles.Highlight(cmd))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// Session details formatting
func SessionDetails(session *interfaces.Session, backend string) string {
	var lines []string

	// Create a consistent width for labels
	labelWidth := 15

	sessionID := session.ID
	name := session.Name
	created := utils.TimeFormat(session.CreatedAt)
	project := session.ProjectPath
	description := session.Description
	status := string(session.Status)

	// Format each field
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("ID:"), sessionID))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("Name:"), styles.Title(name)))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("Status:"), FormatStatus(status)))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("Backend:"), backend))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("Created:"), created))

	if project != "" {
		lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("Project:"), project))
	}

	if description != "" {
		lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, styles.Bold("Description:"), description))
	}

	return strings.Join(lines, "\n")
}

// Available sessions list formatting
func AvailableSessionsList(sessions []SessionInfo) string {
	if len(sessions) == 0 {
		return styles.Dim("  No sessions available") + "\n" +
			fmt.Sprintf("  %s %s", styles.Arrow(), styles.Highlight("claude-pilot create [session-name]"))
	}

	var lines []string
	for _, s := range sessions {
		idDisplay := s.ID
		if len(s.ID) > 8 {
			idDisplay = s.ID[:8]
		}

		line := fmt.Sprintf("  %s %s (%s)", styles.Arrow(), styles.Highlight(s.Name), styles.Dim(idDisplay))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
