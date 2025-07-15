package cmd

import (
	"fmt"

	"claude-pilot/core/api"
	"claude-pilot/internal/ui"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [session-name]",
	Short: "Create a new Claude session",
	Long: `Create a new Claude coding session with an optional name.
If no name is provided, a timestamp-based name will be generated.

Examples:
  claude-pilot create                    # Create session with auto-generated name
  claude-pilot create my-project         # Create session named "my-project"
  claude-pilot create --desc "React app" # Create session with description
  claude-pilot create --project ./src    # Create session with project path`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		// Get command-specific parameters
		var sessionName string
		if len(args) > 0 {
			sessionName = args[0]
		}

		// Get flags
		description, _ := cmd.Flags().GetString("description")
		projectPath, _ := cmd.Flags().GetString("project")

		// Resolve project path using common function
		projectPath = GetProjectPath(projectPath)

		// Create the session
		sess, err := ctx.Client.CreateSession(api.CreateSessionRequest{
			Name:        sessionName,
			Description: description,
			ProjectPath: projectPath,
		})
		if err != nil {
			HandleError(err, "create session")
		}

		// Enhanced success message
		fmt.Println(ui.SuccessMsg(fmt.Sprintf("Created session '%s'", sess.Name)))
		fmt.Println()

		// Show enhanced session details
		details := ui.SessionDetailsFormatted(
			sess.ID,
			sess.Name,
			string(sess.Status),
			ctx.Client.GetBackend(),
			sess.CreatedAt.Format("2006-01-02 15:04:05"),
			sess.ProjectPath,
			sess.Description,
		)
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
	rootCmd.AddCommand(createCmd)

	// Add flags
	createCmd.Flags().StringP("description", "d", "", "Description for the session")
	createCmd.Flags().StringP("project", "p", "", "Project path for the session (defaults to current directory)")
}
