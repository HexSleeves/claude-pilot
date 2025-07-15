package cmd

import (
	"fmt"
	"sort"

	"claude-pilot/core/api"
	"claude-pilot/internal/ui"
	"claude-pilot/shared/components"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active Claude sessions",
	Long: `List all active Claude coding sessions with their details.
Shows session ID, name, status, creation time, last activity, and message count.

Examples:
  claude-pilot list           # List all sessions
  claude-pilot list --all     # Include inactive sessions
  claude-pilot list --sort=name # Sort by name instead of last activity`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		// Get flags
		showAll, _ := cmd.Flags().GetBool("all")
		sortBy, _ := cmd.Flags().GetString("sort")

		// Get all sessions
		sessions, err := ctx.Client.ListSessions()
		if err != nil {
			HandleError(err, "list sessions")
		}

		// Filter sessions if not showing all
		if !showAll {
			// Pre-allocate with estimated capacity (assume most sessions are active)
			activeSessions := make([]*api.Session, 0, len(sessions))
			for _, sess := range sessions {
				if sess.Status == api.StatusActive || sess.Status == api.StatusConnected {
					activeSessions = append(activeSessions, sess)
				}
			}
			sessions = activeSessions
		}

		// Sort sessions
		switch sortBy {
		case "name":
			sort.Slice(sessions, func(i, j int) bool {
				return sessions[i].Name < sessions[j].Name
			})
		case "created":
			sort.Slice(sessions, func(i, j int) bool {
				return sessions[i].CreatedAt.Before(sessions[j].CreatedAt)
			})
		case "status":
			sort.Slice(sessions, func(i, j int) bool {
				return sessions[i].Status < sessions[j].Status
			})
		default: // "activity" or default
			sort.Slice(sessions, func(i, j int) bool {
				return sessions[i].LastActive.After(sessions[j].LastActive)
			})
		}

		// Display header with enhanced styling
		fmt.Println(ui.Header("Claude Pilot Sessions"))
		fmt.Printf("%s Backend: %s\n", ui.InfoMsg("Current"), ui.Highlight(ctx.Client.GetBackend()))
		fmt.Println()

		if len(sessions) == 0 {
			if showAll {
				fmt.Println(ui.Dim("No sessions found."))
			} else {
				fmt.Println(ui.Dim("No active sessions found."))
				fmt.Println(ui.InfoMsg("Use --all to show inactive sessions"))
			}
			fmt.Println()
			fmt.Println(ui.InfoMsg("Create a new session:"))
			fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
			return
		}

		// Convert API sessions to shared table format
		sessionData := convertToSessionData(sessions)

		// Create and configure table for CLI output
		table := components.NewTable(components.TableConfig{
			Width:       120, // Reasonable width for CLI
			ShowHeaders: true,
			Interactive: false,
			MaxRows:     0, // Show all rows
		})

		// Set the session data
		table.SetSessionData(sessionData)

		// Display sessions table using shared component
		fmt.Println(table.RenderCLI())
		fmt.Println()

		// Show enhanced summary
		activeCount := 0
		inactiveCount := 0
		for _, sess := range sessions {
			if sess.Status == api.StatusActive || sess.Status == api.StatusConnected {
				activeCount++
			} else {
				inactiveCount++
			}
		}

		fmt.Println(ui.SessionSummary(len(sessions), activeCount, inactiveCount, showAll))

		// Show helpful commands with enhanced styling
		fmt.Println(ui.AvailableCommands(
			"claude-pilot attach <session-name>",
			"claude-pilot kill <session-name>",
			"claude-pilot create [session-name]",
		))
	},
}

// convertToSessionData converts API sessions to the shared table SessionData format
func convertToSessionData(sessions []*api.Session) []components.SessionData {
	sessionData := make([]components.SessionData, len(sessions))

	for i, sess := range sessions {
		sessionData[i] = components.SessionData{
			ID:          sess.ID,
			Name:        sess.Name,
			Status:      string(sess.Status),
			Backend:     "claude", // Default backend for Claude sessions
			Created:     sess.CreatedAt,
			LastActive:  sess.LastActive,
			Messages:    len(sess.Messages),
			ProjectPath: sess.ProjectPath,
		}
	}

	return sessionData
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolP("all", "a", false, "Show all sessions including inactive ones")
	listCmd.Flags().StringP("sort", "s", "activity", "Sort by: name, created, status, activity")
}
