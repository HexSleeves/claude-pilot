package cmd

import (
	"fmt"
	"os"
	"strings"

	"claude-pilot/internal/config"
	"claude-pilot/internal/manager"
	"claude-pilot/internal/ui"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <session-name-or-id>",
	Short: "Terminate a Claude session",
	Long: `Terminate a specific Claude coding session by name or ID.
This will permanently delete the session and all its data.

Examples:
  claude-pilot kill my-session      # Kill session by name
  claude-pilot kill abc123def       # Kill session by ID
  claude-pilot kill --force my-session # Skip confirmation prompt`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		force, _ := cmd.Flags().GetBool("force")

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

		// Get the session to verify it exists
		sess, err := sm.GetSession(identifier)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Session not found: %v", err)))
			os.Exit(1)
		}

		// Show session details
		fmt.Println(ui.WarningMsg(fmt.Sprintf("About to terminate session '%s'", sess.Name)))
		fmt.Println()

		// Show session details
		fmt.Printf("%-15s %s\n", ui.Bold("ID:"), sess.ID)
		fmt.Printf("%-15s %s\n", ui.Bold("Name:"), ui.Title(sess.Name))
		fmt.Printf("%-15s %s\n", ui.Bold("Status:"), ui.FormatStatus(string(sess.Status)))
		fmt.Printf("%-15s %s\n", ui.Bold("Backend:"), cfg.Backend)
		fmt.Printf("%-15s %s\n", ui.Bold("Created:"), sess.CreatedAt.Format("2006-01-02 15:04:05"))
		if sess.ProjectPath != "" {
			fmt.Printf("%-15s %s\n", ui.Bold("Project:"), sess.ProjectPath)
		}
		fmt.Printf("%-15s %d\n", ui.Bold("Messages:"), len(sess.Messages))

		fmt.Println()

		// Confirmation prompt (unless forced)
		if !force {
			fmt.Print(ui.Prompt("Are you sure you want to terminate this session? [y/N]: "))
			var response string
			fmt.Scanln(&response)

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println(ui.InfoMsg("Session termination cancelled."))
				return
			}
		}

		// Delete the session
		if err := sm.DeleteSession(identifier); err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to terminate session: %v", err)))
			os.Exit(1)
		}

		// Success message
		fmt.Println(ui.SuccessMsg(fmt.Sprintf("Session '%s' has been terminated", sess.Name)))

		// Show remaining sessions count
		remainingSessions, err := sm.ListSessions()
		if err != nil {
			fmt.Println(ui.WarningMsg("Failed to check remaining sessions"))
			return
		}

		if len(remainingSessions) > 0 {
			fmt.Printf("%s %d sessions remaining\n", ui.InfoMsg("Status:"), len(remainingSessions))
			fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot list"))
		} else {
			fmt.Println(ui.InfoMsg("No sessions remaining"))
			fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
		}
	},
}

var killAllCmd = &cobra.Command{
	Use:   "kill-all",
	Short: "Terminate all Claude sessions",
	Long: `Terminate all Claude coding sessions.
This will permanently delete all sessions and their data.

Examples:
  claude-pilot kill-all              # Kill all sessions with confirmation
  claude-pilot kill-all --force      # Skip confirmation prompt`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

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

		if len(sessions) == 0 {
			fmt.Println(ui.InfoMsg("No sessions to terminate."))
			return
		}

		// Show sessions to be terminated
		fmt.Println(ui.WarningMsg(fmt.Sprintf("About to terminate %d sessions:", len(sessions))))
		fmt.Println()
		fmt.Println(ui.SessionTable(sessions, sm.GetMultiplexer()))
		fmt.Println()

		// Confirmation prompt (unless forced)
		if !force {
			fmt.Print(ui.Prompt(fmt.Sprintf("Are you sure you want to terminate all %d sessions? [y/N]: ", len(sessions))))
			var response string
			fmt.Scanln(&response)

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println(ui.InfoMsg("Session termination cancelled."))
				return
			}
		}

		// Delete all sessions
		if err := sm.KillAllSessions(); err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to terminate all sessions: %v", err)))
			os.Exit(1)
		}

		// Success message
		fmt.Println(ui.SuccessMsg(fmt.Sprintf("Successfully terminated all %d sessions", len(sessions))))
		fmt.Println(ui.InfoMsg("All sessions have been terminated"))
		fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
	rootCmd.AddCommand(killAllCmd)

	// Add flags
	killCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	killAllCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}
