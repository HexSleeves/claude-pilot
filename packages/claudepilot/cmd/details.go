package cmd

import (
	"fmt"

	"claude-pilot/internal/ui"
	"claude-pilot/shared/interfaces"

	"github.com/spf13/cobra"
)

var detailsCmd = &cobra.Command{
	Use:   "details [session-name]",
	Short: "Show details for a specific Claude session",
	Long: `Show details for a specific Claude coding session.

Examples:
  claude-pilot details my-project         # Show details for session named "my-project"`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		// List every session if no identifier is provided
		if len(args) == 0 {
			// Get the session
			sessions, err := ctx.Client.ListSessions()
			if err != nil {
				HandleError(err, "list sessions")
			}

			for _, session := range sessions {
				listDetails(ctx, session)
			}

			// Show enhanced next steps
			fmt.Println(ui.NextSteps(
				"claude-pilot list",
				"claude-pilot details [session-name]",
				"claude-pilot create [session-name]",
				"claude-pilot kill [session-name]",
				"claude-pilot attach [session-name]",
			))

			return
		}

		identifier := args[0]

		// Get the session
		sess, err := ctx.Client.GetSession(identifier)
		if err != nil {
			HandleError(err, "get session")
		}

		// Show enhanced session details
		listDetails(ctx, sess)

		// Show enhanced next steps
		fmt.Println(ui.NextSteps(
			fmt.Sprintf("claude-pilot attach %s", sess.Name),
			"claude-pilot list",
		))
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)
}

func listDetails(ctx *CommandContext, sess *interfaces.Session) {
	// Show enhanced session details
	details := ui.SessionDetailsFormatted(sess, ctx.Client.GetBackend())
	fmt.Println(details)
	fmt.Println()
}
