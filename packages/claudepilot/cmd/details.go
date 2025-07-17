package cmd

import (
	"fmt"

	"claude-pilot/internal/ui"

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

		identifier := args[0]

		// Get the session
		sess, err := ctx.Client.GetSession(identifier)
		if err != nil {
			HandleError(err, "get session")
		}

		// Show enhanced session details
		details := ui.SessionDetailsFormatted(sess, ctx.Client.GetBackend())
		fmt.Println(details)
		fmt.Println()

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
