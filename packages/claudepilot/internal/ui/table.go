package ui

import (
	"claude-pilot/core/api"
	"claude-pilot/shared/styles"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// sessionStatusToMultiplexerDisplay converts session status to multiplexer display format
func sessionStatusToMultiplexerDisplay(status api.SessionStatus) string {
	switch status {
	case api.StatusConnected:
		return FormatTmuxStatus("attached")
	case api.StatusActive:
		return FormatTmuxStatus("running")
	case api.StatusInactive:
		return FormatTmuxStatus("stopped")
	case api.StatusError:
		return FormatTmuxStatus("error")
	default:
		return Dim("unknown")
	}
}

// createTableStyleFunc creates a StyleFunc for the session table using existing color scheme
func createTableStyleFunc() table.StyleFunc {
	// Use existing lipgloss styles from the styles package
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.TextPrimary).
		Background(styles.BackgroundSecondary)

	evenRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary)

	oddRowStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted)

	return func(row, col int) lipgloss.Style {
		switch {
		case row == table.HeaderRow:
			return headerStyle
		case row%2 == 0:
			return evenRowStyle
		default:
			return oddRowStyle
		}
	}
}

// SessionTable creates a formatted table for displaying sessions
func SessionTable(sessions []*api.Session, backend string) string {
	if len(sessions) == 0 {
		return Dim("No active sessions found.")
	}

	// Create table
	t := table.New()

	// Set border and styling using existing color scheme
	t.Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.ClaudePrimary)).
		StyleFunc(createTableStyleFunc())

	// Set headers
	t.Headers(
		Bold("ID"),
		Bold("Name"),
		Bold("Status"),
		Bold("Backend"),
		Bold("Created"),
		Bold("Last Active"),
		Bold("Messages"),
		Bold("Project"),
	)

	// Add rows
	for _, sess := range sessions {
		// Use the session status computed by SessionService
		muxStatus := sessionStatusToMultiplexerDisplay(sess.Status)

		// Truncate and format data to fit column widths
		id := sess.ID
		if len(id) > 11 {
			id = id[:8] + "..."
		}

		name := sess.Name
		if len(name) > 19 {
			name = name[:16] + "..."
		}

		project := formatProjectPath(sess.ProjectPath)
		if len(project) > 29 {
			project = project[:26] + "..."
		}

		t.Row(
			Highlight(id),
			Title(name),
			FormatStatus(string(sess.Status)),
			muxStatus,
			formatTime(sess.CreatedAt),
			formatTimeAgo(sess.LastActive),
			fmt.Sprintf("%d", len(sess.Messages)),
			project,
		)
	}

	return t.String()
}

// SessionDetail creates a detailed view of a single session
func SessionDetail(sess *api.Session, backend string) string {
	var builder strings.Builder

	// Header
	builder.WriteString(Title("Session Details") + "\n")
	builder.WriteString(HorizontalLine(50) + "\n\n")

	// Use the session status computed by SessionService
	muxStatusDisplay := sessionStatusToMultiplexerDisplay(sess.Status)

	// Basic info
	builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("ID:"), sess.ID))
	builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Name:"), Title(sess.Name)))
	builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Status:"), FormatStatus(string(sess.Status))))
	builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Backend:"), muxStatusDisplay))
	builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Created:"), formatTime(sess.CreatedAt)))
	builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Last Active:"), formatTimeAgo(sess.LastActive)))
	builder.WriteString(fmt.Sprintf("%-15s %d\n", Bold("Messages:"), len(sess.Messages)))

	if sess.ProjectPath != "" {
		builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Project:"), sess.ProjectPath))
	}

	if sess.Description != "" {
		builder.WriteString(fmt.Sprintf("%-15s %s\n", Bold("Description:"), sess.Description))
	}

	// Recent messages with enhanced styling
	if len(sess.Messages) > 0 {
		builder.WriteString("\n" + styles.Subtitle("Recent Messages:") + "\n")
		builder.WriteString(styles.HorizontalLine(50) + "\n")

		// Show last 5 messages
		start := len(sess.Messages) - 5
		if start < 0 {
			start = 0
		}

		for i := start; i < len(sess.Messages); i++ {
			msg := sess.Messages[i]
			roleColor := Info
			if msg.Role == "user" {
				roleColor = ClaudePrimary
			}

			builder.WriteString(fmt.Sprintf("%s %s %s\n",
				roleColor.Sprint(fmt.Sprintf("[%s]", msg.Role)),
				Dim(msg.Timestamp.Format("15:04:05")),
				styles.TruncateText(msg.Content, 60),
			))
		}
	}

	return builder.String()
}

// MessageHistory creates a formatted view of session messages
func MessageHistory(messages []api.Message, limit int) string {
	if len(messages) == 0 {
		return Dim("No messages in this session.")
	}

	var builder strings.Builder
	builder.WriteString(Title("Message History") + "\n")
	builder.WriteString(HorizontalLine(80) + "\n\n")

	start := 0
	if limit > 0 && len(messages) > limit {
		start = len(messages) - limit
		builder.WriteString(Dim(fmt.Sprintf("Showing last %d messages (of %d total)\n\n", limit, len(messages))))
	}

	for i := start; i < len(messages); i++ {
		msg := messages[i]

		// Role header
		roleColor := Info
		roleIcon := "🤖"
		if msg.Role == "user" {
			roleColor = ClaudePrimary
			roleIcon = "👤"
		}

		builder.WriteString(fmt.Sprintf("%s %s %s %s\n",
			roleIcon,
			roleColor.Sprint(strings.ToUpper(msg.Role)),
			Dim(msg.Timestamp.Format("2006-01-02 15:04:05")),
			Dim(fmt.Sprintf("(%s)", msg.ID[:8])),
		))

		// Message content
		builder.WriteString(TextPrimary.Sprint(msg.Content) + "\n\n")
		builder.WriteString(Dim(strings.Repeat("─", 80)) + "\n\n")
	}

	return builder.String()
}

// Helper functions
func formatTime(t time.Time) string {
	return Dim(t.Format("2006-01-02 15:04"))
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return Success.Sprint("just now")
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 0 {
			return Success.Sprint("just now")
		}
		unit := "min"
		if minutes != 1 {
			unit = "mins"
		}
		return Info.Sprint(fmt.Sprintf("%d%s ago", minutes, unit))
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 0 {
			return Success.Sprint("just now")
		}
		unit := "hour"
		if hours != 1 {
			unit = "hours"
		}
		return Warning.Sprint(fmt.Sprintf("%d %s ago", hours, unit))
	default:
		days := int(duration.Hours() / 24)
		if days == 0 {
			return Success.Sprint("just now")
		}
		unit := "day"
		if days != 1 {
			unit = "days"
		}
		return Dim(fmt.Sprintf("%d %s ago", days, unit))
	}
}

func formatProjectPath(path string) string {
	if path == "" {
		return Dim("—")
	}

	// Show only the last part of the path if it's too long
	if len(path) > 25 {
		parts := strings.Split(path, "/")
		if len(parts) > 1 {
			return Dim(".../" + parts[len(parts)-1])
		}
	}

	return Dim(path)
}
