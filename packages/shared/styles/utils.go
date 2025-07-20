package styles

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	evertrastable "github.com/evertras/bubble-table/table"
)

type Context string

const (
	ContextStatus      Context = "status"
	ContextAction      Context = "action"
	ContextInteractive Context = "interactive"
	ContextBackend     Context = "backend"
)

type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusInactive SessionStatus = "inactive"
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

// Theme-aware color selection based on context
func GetContextualColor(context Context, state string) lipgloss.Color {
	switch context {
	case ContextStatus:
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

	case ContextAction:
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

	case ContextInteractive:
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

	case ContextBackend:
		switch state {
		case "tmux":
			return TextPrimary
		default:
			return TextSecondary
		}

	default:
		return TextPrimary
	}
}

// Evertras Table Styling Functions

// GetEvertrasTableStyles returns styling configuration for evertras table
func GetEvertrasTableStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ClaudePrimary).
		Foreground(TextPrimary)
}

// ConfigureEvertrasTable applies comprehensive Claude theme styling to an evertras table model
// and row data formatting since evertras table uses a different styling approach
func ConfigureEvertrasTable(t evertrastable.Model) evertrastable.Model {
	return t.WithBaseStyle(GetEvertrasTableStyles())
}

// FormatStatus provides status formatting with status indicators
func FormatStatus(status string) string {
	switch status {
	case "active":
		return "● " + status
	case "inactive":
		return "⏸ " + status
	case "connected":
		return "🔗 " + status
	case "error", "failed":
		return "✗ " + status
	case "starting", "pending":
		return "⏳ " + status
	case "stopped":
		return "⏹ " + status
	default:
		return "? " + status
	}
}

// FormatTime formats timestamps with consistent styling
func FormatTime(t time.Time) string {
	return TableCellTimestampStyle.Render(t.Format("2006-01-02 15:04"))
}

// FormatTimeAgo formats relative time with semantic colors based on recency
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return TableCellSuccessStyle.Render("just now")
	case duration < time.Hour:
		return TableCellInfoStyle.Render(fmt.Sprintf("%dm ago", int(duration.Minutes())))
	case duration < 24*time.Hour:
		return TableCellWarningStyle.Render(fmt.Sprintf("%dh ago", int(duration.Hours())))
	case duration < 7*24*time.Hour:
		return TableCellTimestampStyle.Render(fmt.Sprintf("%dd ago", int(duration.Hours()/24)))
	default:
		return TableCellStyle.Render(fmt.Sprintf("%dd ago", int(duration.Hours()/24)))
	}
}

// FormatProjectPath formats project paths with consistent styling and smart truncation
func FormatProjectPath(path string, maxLen int) string {
	if path == "" {
		return TableCellStyle.Render("—")
	}

	// If it already fits, just render it
	if len(path) <= maxLen {
		return TableCellStyle.Render(path)
	}

	// Split on “/” and build up the tail from the end
	parts := strings.Split(path, "/")
	var tail string

	// We reserve 4 chars for the ".../" prefix
	for i := len(parts) - 1; i >= 1; i-- {
		cand := parts[i]
		if tail != "" {
			cand = parts[i] + "/" + tail
		}

		if len(cand)+4 > maxLen {
			// adding any earlier segment would only make it longer
			break
		}

		tail = cand
	}

	if tail != "" {
		return TableCellStyle.Render(".../" + tail)
	}

	// Fallback if even the final segment is too long
	return TableCellStyle.Render(TruncateText(path, maxLen))
}
