package cmd

import (
	"fmt"

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
  claude-pilot list --sort=name # Sort by name instead of last activity
	claude-pilot list --filter=active # Filter to only show active sessions
	claude-pilot list --filter=inactive # Filter to only show inactive sessions`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		// Get flags
		sortBy, _ := cmd.Flags().GetString("sort")
		filter, _ := cmd.Flags().GetString("filter")

		var sessions []*api.Session

		// Apply filters
		if filter != "" {
			sessions, err = ctx.Client.ListFilteredSessions(filter)
			if err != nil {
				HandleError(err, "list filtered sessions")
			}

			if len(sessions) == 0 {
				fmt.Println(ui.Dim("No sessions found. Try another filter or create a new session."))
				fmt.Println()
				fmt.Println(ui.InfoMsg("Try another filter:"))
				fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot list --filter=active"))
				fmt.Println()
				fmt.Println(ui.InfoMsg("Create a new session:"))
				fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
				fmt.Println()
				return
			}
		} else {
			// Get all sessions
			sessions, err = ctx.Client.ListSessions()
			if err != nil {
				HandleError(err, "list sessions")
			}
		}

		// Display header with enhanced styling
		fmt.Println(ui.Header("Claude Pilot Sessions"))
		fmt.Printf("%s Backend: %s\n", ui.InfoMsg("Current"), ui.Highlight(ctx.Client.GetBackend()))
		fmt.Println()

		if len(sessions) == 0 {
			fmt.Println(ui.Dim("No sessions found."))
			fmt.Println()
			fmt.Println(ui.InfoMsg("Create a new session:"))
			fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
			return
		}

		// Convert API sessions to shared table format
		sessionData := convertToSessionData(sessions)

		// Create and configure table for CLI output with enhanced features
		table := components.NewSessionTable(components.TableConfig{
			ShowHeaders: true,
			Interactive: false,
			MaxRows:     0, // Show all rows
			SortEnabled: true,
		})

		// Set the session data
		table.SetSessionData(sessionData)

		// Apply CLI sort option using table's built-in sorting
		if sortBy != "" {
			direction := "asc"
			if sortBy == "activity" {
				sortBy = "last_active"
				direction = "desc" // Most recent first for activity
			}
			table.SetSort(sortBy, direction)
		}

		// Apply CLI filter option using table's built-in filtering
		if filter != "" && filter != "active" {
			// TODO: Implement filtering
			// table.SetFilter(filter)
		}

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

		fmt.Println(ui.SessionSummary(len(sessions), activeCount, inactiveCount))

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
			Backend:     sess.Backend,
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
	listCmd.Flags().StringP("sort", "s", "activity", "Sort by: name, created, status, activity")
	listCmd.Flags().StringP("filter", "f", "", "Filter by: active, inactive")
}
