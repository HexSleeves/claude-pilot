package ui

import (
	"claude-pilot/core/api"
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

// SessionTable creates a formatted table for displaying sessions
func SessionTable(sessions []*api.Session, backend string) string {
	if len(sessions) == 0 {
		return Dim("No active sessions found.")
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)

	// Customize colors to match our theme
	t.Style().Color.Header = text.Colors{text.FgHiWhite, text.Bold}
	t.Style().Color.Row = text.Colors{text.FgWhite}
	t.Style().Color.RowAlternate = text.Colors{text.FgHiBlack}

	// Set headers
	t.AppendHeader(table.Row{
		Bold("ID"),
		Bold("Name"),
		Bold("Status"),
		Bold("Backend"),
		Bold("Created"),
		Bold("Last Active"),
		Bold("Messages"),
		Bold("Project"),
	})

	// Add rows
	for _, sess := range sessions {
		// Use the session status computed by SessionService
		muxStatus := sessionStatusToMultiplexerDisplay(sess.Status)

		t.AppendRow(table.Row{
			Highlight(sess.ID[:8] + "..."), // Truncate ID for readability
			Title(sess.Name),
			FormatStatus(string(sess.Status)),
			muxStatus,
			formatTime(sess.CreatedAt),
			formatTimeAgo(sess.LastActive),
			fmt.Sprintf("%d", len(sess.Messages)),
			formatProjectPath(sess.ProjectPath),
		})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 12}, // ID
		{Number: 2, WidthMax: 20}, // Name
		{Number: 3, WidthMax: 12}, // Status
		{Number: 4, WidthMax: 12}, // Backend
		{Number: 5, WidthMax: 12}, // Created
		{Number: 6, WidthMax: 12}, // Last Active
		{Number: 7, WidthMax: 8},  // Messages
		{Number: 8, WidthMax: 30}, // Project
	})

	return t.Render()
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

	// Recent messages
	if len(sess.Messages) > 0 {
		builder.WriteString("\n" + Subtitle("Recent Messages:") + "\n")
		builder.WriteString(HorizontalLine(50) + "\n")

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
				truncateText(msg.Content, 60),
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

	if duration < time.Minute {
		return Success.Sprint("just now")
	} else if duration < time.Hour {
		return Info.Sprint(fmt.Sprintf("%dm ago", int(duration.Minutes())))
	} else if duration < 24*time.Hour {
		return Warning.Sprint(fmt.Sprintf("%dh ago", int(duration.Hours())))
	} else {
		return Dim(fmt.Sprintf("%dd ago", int(duration.Hours()/24)))
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

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	return text[:maxLen-3] + "..."
}
