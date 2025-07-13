package cmd

import (
	"fmt"
	"os"
	"sort"

	"claude-pilot/internal/config"
	"claude-pilot/internal/interfaces"
	"claude-pilot/internal/manager"
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
		// Get flags
		showAll, _ := cmd.Flags().GetBool("all")
		sortBy, _ := cmd.Flags().GetString("sort")

		// Load configuration
		configManager := config.NewConfigManager("")
		cfg, err := configManager.Load()
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to load configuration: %v", err)))
			os.Exit(1)
		}

		// Create session manager
		sm, err := manager.NewSessionManager(cfg)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to initialize session manager: %v", err)))
			os.Exit(1)
		}

		// Get all sessions
		sessions, err := sm.ListSessions()
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to list sessions: %v", err)))
			os.Exit(1)
		}

		// Filter sessions if not showing all
		if !showAll {
			var activeSessions []*interfaces.Session
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
		fmt.Printf("%s Backend: %s\n", ui.InfoMsg("Current"), sm.GetConfig().Backend)
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
		fmt.Println(ui.SessionTable(sessions, sm.GetMultiplexer()))
		fmt.Println()

		// Show summary
		activeCount := 0
		inactiveCount := 0
		for _, sess := range sessions {
			if sess.Status == interfaces.StatusActive || sess.Status == interfaces.StatusConnected {
				activeCount++
			} else {
				inactiveCount++
			}
		}

		if showAll {
			fmt.Printf("%s Total: %d sessions (%d active, %d inactive)\n",
				ui.InfoMsg("Summary:"), len(sessions), activeCount, inactiveCount)
		} else {
			fmt.Printf("%s Active sessions: %d\n",
				ui.InfoMsg("Summary:"), activeCount)
			if inactiveCount > 0 {
				fmt.Printf("  %s Use --all to show %d inactive sessions\n",
					ui.Dim("Note:"), inactiveCount)
			}
		}

		// Show helpful commands
		fmt.Println()
		fmt.Println(ui.InfoMsg("Available commands:"))
		fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot attach <session-name>"))
		fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot kill <session-name>"))
		fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolP("all", "a", false, "Show all sessions including inactive ones")
	listCmd.Flags().StringP("sort", "s", "activity", "Sort by: name, created, status, activity")
}
