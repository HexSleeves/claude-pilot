package cmd

import (
	"fmt"

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
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		identifier := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// Get the session to verify it exists
		sess, err := ctx.Service.GetSession(identifier)
		if err != nil {
			HandleError(err, "find session")
		}

		// Show session details
		fmt.Println(ui.WarningMsg(fmt.Sprintf("About to terminate session '%s'", sess.Name)))
		fmt.Println()

		// Show session details using common function (with messages)
		ui.DisplaySessionDetailsWithMessages(sess, ctx.Config.Backend)
		fmt.Println()

		// Confirmation prompt (unless forced) using common function
		if !force {
			if !ConfirmAction("Are you sure you want to terminate this session? [y/N]: ") {
				fmt.Println(ui.InfoMsg("Session termination cancelled."))
				return
			}
		}

		// Delete the session
		if err := ctx.Service.DeleteSession(identifier); err != nil {
			HandleError(err, "terminate session")
		}

		// Success message
		fmt.Println(ui.SuccessMsg(fmt.Sprintf("Session '%s' has been terminated", sess.Name)))

		// Show remaining sessions count using common function
		remainingSessions, err := ctx.Service.ListSessions()
		if err != nil {
			fmt.Println(ui.WarningMsg("Failed to check remaining sessions"))
			return
		}

		ui.DisplayRemainingSessionsInfo(remainingSessions)
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
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		force, _ := cmd.Flags().GetBool("force")

		// Get all sessions
		sessions, err := ctx.Service.ListSessions()
		if err != nil {
			HandleError(err, "list sessions")
		}

		if len(sessions) == 0 {
			fmt.Println(ui.InfoMsg("No sessions to terminate."))
			return
		}

		// Show sessions to be terminated
		fmt.Println(ui.WarningMsg(fmt.Sprintf("About to terminate %d sessions:", len(sessions))))
		fmt.Println()
		fmt.Println(ui.SessionTable(sessions, ctx.Multiplexer))
		fmt.Println()

		// Confirmation prompt (unless forced) using common function
		if !force {
			if !ConfirmAction(fmt.Sprintf("Are you sure you want to terminate all %d sessions? [y/N]: ", len(sessions))) {
				fmt.Println(ui.InfoMsg("Session termination cancelled."))
				return
			}
		}

		// Delete all sessions
		if err := ctx.Service.KillAllSessions(); err != nil {
			HandleError(err, "terminate all sessions")
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
