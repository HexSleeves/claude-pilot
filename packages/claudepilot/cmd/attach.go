package cmd

import (
	"fmt"
	"os"

	"claude-pilot/internal/ui"

	"github.com/spf13/cobra"
)

var attachCmd = &cobra.Command{
	Use:   "attach <session-name-or-id>",
	Short: "Attach to a Claude session for interactive communication",
	Long: `Attach to a specific Claude coding session to start interactive communication.
This opens a terminal-based chat interface where you can communicate with Claude
in real-time within the context of your coding session.

Examples:
  claude-pilot attach my-session     # Attach to session by name
  claude-pilot attach abc123def      # Attach to session by ID`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		identifier := args[0]

		// Get the session
		sess, err := ctx.Client.GetSession(identifier)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Session not found: %v", err)))
			fmt.Println()
			fmt.Println(ui.InfoMsg("Available sessions:"))

			// Show available sessions using common function
			sessions, err := ctx.Client.ListSessions()
			if err != nil {
				HandleError(err, "list sessions")
			}

			ui.DisplayAvailableSessions(sessions)
			os.Exit(1)
		}

		// Check if session is running
		if !ctx.Client.IsSessionRunning(sess.Name) {
			fmt.Println(ui.WarningMsg(fmt.Sprintf("Session '%s' is not running. It may have been terminated.", sess.Name)))
			fmt.Println(ui.InfoMsg("You can recreate it with: claude-pilot create " + sess.Name))
			os.Exit(1)
		}

		// Show session info
		fmt.Println(ui.InfoMsg(fmt.Sprintf("Attaching to session '%s' (%s backend)...", sess.Name, ctx.Client.GetBackend())))
		fmt.Println(ui.InfoMsg("Use your multiplexer's detach key to exit (tmux: Ctrl+B,D | zellij: Ctrl+P,D)"))
		fmt.Println()

		// Attach to the session
		if err := ctx.Client.AttachToSession(identifier); err != nil {
			HandleError(err, "attach to session")
		}

		// After detaching, we're back to the CLI
		fmt.Println(ui.InfoMsg("Detached from session"))
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)
}
