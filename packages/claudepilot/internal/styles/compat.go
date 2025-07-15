package styles

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Compatibility functions that match the existing UI interface
// These functions provide the same API as the current ui/colors.go but with lipgloss styling

// Message formatting functions (matching ui/colors.go interface)
func SuccessMsg(text string) string {
	return Success(text)
}

func ErrorMsg(text string) string {
	return Error(text)
}

func WarningMsg(text string) string {
	return Warning(text)
}

func InfoMsg(text string) string {
	return Info(text)
}

// Status formatting functions
func FormatStatus(status string) string {
	switch strings.ToLower(status) {
	case "active", "running":
		return StatusActive(status)
	case "inactive", "stopped":
		return StatusInactive(status)
	case "connected", "attached":
		return StatusConnected(status)
	case "error", "failed":
		return StatusError(status)
	default:
		return Dim(status)
	}
}

// Tmux status formatting (matching existing interface)
func FormatTmuxStatus(status string) string {
	switch status {
	case "running":
		return StatusActive(status)
	case "attached":
		return StatusConnected(status)
	case "stopped":
		return StatusInactive(status)
	case "error":
		return StatusError(status)
	default:
		return Dim("? " + status)
	}
}

// Progress indicators
func CheckMark() string {
	return SuccessStyle.Render("✓")
}

func CrossMark() string {
	return ErrorStyle.Render("✗")
}

// Time formatting
func FormatTime(t time.Time) string {
	return Dim(t.Format("2006-01-02 15:04"))
}

// Text truncation with styling
func TruncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + Dim("...")
}

// Enhanced banner for the root command
func RootBanner() string {
	// Create the main title
	title := lipgloss.NewStyle().
		Foreground(ClaudePrimary).
		Bold(true).
		Render("Claude Pilot 🚀")

	// Create the subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(TextSecondary).
		Render("A powerful CLI tool for managing multiple Claude code sessions")

	// Create decorative border
	borderStyle := lipgloss.NewStyle().
		Foreground(ClaudePrimary).
		Bold(true)

	topBorder := borderStyle.Render("╔" + strings.Repeat("═", 58) + "╗")
	bottomBorder := borderStyle.Render("╚" + strings.Repeat("═", 58) + "╝")

	// Center the content
	titleCentered := lipgloss.PlaceHorizontal(58, lipgloss.Center, title)
	subtitleCentered := lipgloss.PlaceHorizontal(58, lipgloss.Center, subtitle)
	emptyCentered := lipgloss.PlaceHorizontal(58, lipgloss.Center, "")

	// Add side borders
	titleLine := borderStyle.Render("║") + titleCentered + borderStyle.Render("║")
	subtitleLine := borderStyle.Render("║") + subtitleCentered + borderStyle.Render("║")
	emptyLine := borderStyle.Render("║") + emptyCentered + borderStyle.Render("║")

	// Join all parts
	banner := lipgloss.JoinVertical(lipgloss.Left,
		topBorder,
		emptyLine,
		titleLine,
		emptyLine,
		subtitleLine,
		emptyLine,
		bottomBorder,
	)

	return banner
}

// Enhanced command list formatting
func CommandList(commands map[string]string) string {
	var lines []string

	// Header
	header := lipgloss.NewStyle().
		Foreground(InfoColor).
		Bold(true).
		Render("Available Commands:")
	lines = append(lines, header)
	lines = append(lines, "")

	// Commands
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

// Session summary formatting
func SessionSummary(total, active, inactive int, showAll bool) string {
	var parts []string

	if total > 0 {
		totalText := fmt.Sprintf("Total: %d", total)
		parts = append(parts, BoldStyle.Background(BackgroundSecondary).Render(totalText))

		if active > 0 {
			activeText := fmt.Sprintf("Active: %d", active)
			parts = append(parts, StatusActiveStyle.Background(BackgroundSecondary).Render(activeText))
		}

		if inactive > 0 {
			inactiveText := fmt.Sprintf("Inactive: %d", inactive)
			parts = append(parts, StatusInactiveStyle.Background(BackgroundSecondary).Render(inactiveText))
		}
	}

	if len(parts) == 0 {
		return Dim("No sessions found")
	}

	// Join the parts with a separator that also has the background
	separator := WithBackground(" | ", BackgroundSecondary)
	summary := strings.Join(parts, separator)

	if !showAll && inactive > 0 {
		hint := Dim("(Use --all to show inactive sessions)")
		summary = fmt.Sprintf("%s\n%s", summary, hint)
	}

	return InfoBox(summary)
}

// Next steps formatting
func NextSteps(commands ...string) string {
	var lines []string

	header := Info("Next steps:")
	lines = append(lines, header)

	for _, cmd := range commands {
		line := fmt.Sprintf("  %s %s", Arrow(), Highlight(cmd))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// Available commands formatting
func AvailableCommands(commands ...string) string {
	var lines []string

	header := Info("Available commands:")
	lines = append(lines, header)

	for _, cmd := range commands {
		line := fmt.Sprintf("  %s %s", Arrow(), Highlight(cmd))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// Session details formatting
func SessionDetails(sessionID, name, status, backend, created, project, description string) string {
	var lines []string

	// Create a consistent width for labels
	labelWidth := 15

	// Format each field
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("ID:"), sessionID))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("Name:"), Title(name)))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("Status:"), FormatStatus(status)))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("Backend:"), backend))
	lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("Created:"), created))

	if project != "" {
		lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("Project:"), project))
	}

	if description != "" {
		lines = append(lines, fmt.Sprintf("%-*s %s", labelWidth, Bold("Description:"), description))
	}

	return strings.Join(lines, "\n")
}

// Available sessions list formatting
func AvailableSessionsList(sessions []SessionInfo) string {
	if len(sessions) == 0 {
		return Dim("  No sessions available") + "\n" +
			fmt.Sprintf("  %s %s", Arrow(), Highlight("claude-pilot create [session-name]"))
	}

	var lines []string
	for _, s := range sessions {
		idDisplay := s.ID
		if len(s.ID) > 8 {
			idDisplay = s.ID[:8]
		}

		line := fmt.Sprintf("  %s %s (%s)", Arrow(), Highlight(s.Name), Dim(idDisplay))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// SessionInfo represents basic session information for display
type SessionInfo struct {
	ID   string
	Name string
}
