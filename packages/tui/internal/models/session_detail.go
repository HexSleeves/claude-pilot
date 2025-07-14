package models

import (
	"fmt"
	"strings"
	"time"

	"claude-pilot/core/api"
	"claude-pilot/tui/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
)

// SessionDetailModel handles the session detail view
type SessionDetailModel struct {
	client  *api.Client
	session *api.Session
	width   int
	height  int
}

// NewSessionDetailModel creates a new session detail model
func NewSessionDetailModel(client *api.Client) *SessionDetailModel {
	return &SessionDetailModel{
		client: client,
	}
}

// Init implements tea.Model
func (m *SessionDetailModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *SessionDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Attach to session
			if m.session != nil {
				return m, func() tea.Msg {
					err := m.client.AttachToSession(m.session.Name)
					if err != nil {
						return ErrorMsg{Error: err}
					}
					return nil
				}
			}

		case "k":
			// Kill session
			if m.session != nil {
				return m, func() tea.Msg {
					err := m.client.KillSession(m.session.Name)
					if err != nil {
						return ErrorMsg{Error: err}
					}
					// Go back to session list after killing
					return SessionKilledMsg{SessionID: m.session.ID}
				}
			}
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *SessionDetailModel) View() string {
	if m.session == nil {
		return styles.MutedTextStyle.Render("No session selected")
	}

	var content strings.Builder

	// Header
	content.WriteString(styles.HeaderStyle.Render("Session Details"))
	content.WriteString("\n\n")

	// Session information
	content.WriteString(styles.LabelStyle.Render("Name:"))
	content.WriteString(" ")
	content.WriteString(styles.SessionNameStyle.Render(m.session.Name))
	content.WriteString("\n")

	content.WriteString(styles.LabelStyle.Render("ID:"))
	content.WriteString(" ")
	content.WriteString(styles.SessionIDStyle.Render(m.session.ID))
	content.WriteString("\n")

	content.WriteString(styles.LabelStyle.Render("Status:"))
	content.WriteString(" ")
	statusStyle := styles.FormatSessionStatus(string(m.session.Status))
	content.WriteString(statusStyle.Render(string(m.session.Status)))
	content.WriteString("\n")

	content.WriteString(styles.LabelStyle.Render("Backend:"))
	content.WriteString(" ")
	content.WriteString(styles.PrimaryTextStyle.Render(m.client.GetBackend()))
	content.WriteString("\n")

	content.WriteString(styles.LabelStyle.Render("Created:"))
	content.WriteString(" ")
	content.WriteString(styles.SecondaryTextStyle.Render(m.session.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString("\n")

	content.WriteString(styles.LabelStyle.Render("Last Active:"))
	content.WriteString(" ")
	content.WriteString(styles.SecondaryTextStyle.Render(m.formatTimeAgo(m.session.LastActive)))
	content.WriteString("\n")

	if m.session.Description != "" {
		content.WriteString(styles.LabelStyle.Render("Description:"))
		content.WriteString(" ")
		content.WriteString(styles.PrimaryTextStyle.Render(m.session.Description))
		content.WriteString("\n")
	}

	if m.session.ProjectPath != "" {
		content.WriteString(styles.LabelStyle.Render("Project Path:"))
		content.WriteString(" ")
		content.WriteString(styles.SecondaryTextStyle.Render(m.session.ProjectPath))
		content.WriteString("\n")
	}

	content.WriteString(styles.LabelStyle.Render("Messages:"))
	content.WriteString(" ")
	content.WriteString(styles.PrimaryTextStyle.Render(fmt.Sprintf("%d", len(m.session.Messages))))
	content.WriteString("\n")

	// Message history (if any)
	if len(m.session.Messages) > 0 {
		content.WriteString("\n")
		content.WriteString(styles.HeaderStyle.Render("Recent Messages"))
		content.WriteString("\n\n")

		// Show last 5 messages
		start := 0
		if len(m.session.Messages) > 5 {
			start = len(m.session.Messages) - 5
		}

		for i := start; i < len(m.session.Messages); i++ {
			msg := m.session.Messages[i]
			
			// Message header
			roleStyle := styles.InfoStyle
			if msg.Role == "user" {
				roleStyle = styles.SuccessStyle
			}
			
			content.WriteString(roleStyle.Render(fmt.Sprintf("[%s]", msg.Role)))
			content.WriteString(" ")
			content.WriteString(styles.MutedTextStyle.Render(msg.Timestamp.Format("15:04:05")))
			content.WriteString("\n")

			// Message content (truncated)
			truncatedContent := styles.TruncateText(msg.Content, 100)
			content.WriteString(styles.SecondaryTextStyle.Render(truncatedContent))
			content.WriteString("\n\n")
		}
	}

	// Actions
	content.WriteString("\n")
	content.WriteString(styles.MutedTextStyle.Render("a: attach • k: kill • esc: back"))

	return content.String()
}

// SetSize updates the model's size
func (m *SessionDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetSession sets the session to display
func (m *SessionDetailModel) SetSession(session *api.Session) {
	m.session = session
}

// formatTimeAgo formats a time as "X ago" string
func (m *SessionDetailModel) formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

// SessionKilledMsg is sent when a session is killed
type SessionKilledMsg struct {
	SessionID string
}
