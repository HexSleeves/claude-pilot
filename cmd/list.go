package cmd

import (
	"fmt"
	"sort"

	"claude-pilot/internal/interfaces"
	"claude-pilot/internal/ui"

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
		sessions, err := ctx.Service.ListSessions()
		if err != nil {
			HandleError(err, "list sessions")
		}

		// Filter sessions if not showing all
		if !showAll {
			// Pre-allocate with estimated capacity (assume most sessions are active)
			activeSessions := make([]*interfaces.Session, 0, len(sessions))
			for _, sess := range sessions {
				if sess.Status == interfaces.StatusActive || sess.Status == interfaces.StatusConnected {
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

		// Display header
		fmt.Println(ui.Title("Claude Pilot Sessions"))
		fmt.Printf("%s Backend: %s\n", ui.InfoMsg("Current"), ctx.Config.Backend)
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

		// Display sessions table
		fmt.Println(ui.SessionTable(sessions, ctx.Multiplexer))
		fmt.Println()

		// Show summary using common function
		activeCount := 0
		inactiveCount := 0
		for _, sess := range sessions {
			if sess.Status == interfaces.StatusActive || sess.Status == interfaces.StatusConnected {
				activeCount++
			} else {
				inactiveCount++
			}
		}

		ui.DisplaySessionSummary(len(sessions), activeCount, inactiveCount, showAll)

		// Show helpful commands using common function
		ui.DisplayAvailableCommands(
			"claude-pilot attach <session-name>",
			"claude-pilot kill <session-name>",
			"claude-pilot create [session-name]",
		)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolP("all", "a", false, "Show all sessions including inactive ones")
	listCmd.Flags().StringP("sort", "s", "activity", "Sort by: name, created, status, activity")
}
