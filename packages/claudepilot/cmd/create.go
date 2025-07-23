package cmd

import (
	"fmt"
	"strings"

	"claude-pilot/core/api"
	"claude-pilot/internal/ui"
	"claude-pilot/shared/interfaces"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [session-name]",
	Short: "Create a new Claude session",
	Long: `Create a new Claude coding session with an optional name.
If no name is provided, a timestamp-based name will be generated.

You can also attach to existing sessions as new panes or windows/tabs.

Examples:
  claude-pilot create                              # Create session with auto-generated name
  claude-pilot create my-project                   # Create session named "my-project"
  claude-pilot create --desc "React app"           # Create session with description
  claude-pilot create --project ./src              # Create session with project path
  claude-pilot create --attach-to main --as-pane   # Create as new pane in 'main' session
  claude-pilot create --attach-to main --as-window # Create as new window in 'main' session
  claude-pilot create debug --attach-to main --as-pane --split h  # Create horizontal pane split
  `,
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
		attachTo, _ := cmd.Flags().GetString("attach-to")
		asPane, _ := cmd.Flags().GetBool("as-pane")
		asWindow, _ := cmd.Flags().GetBool("as-window")
		splitDirection, _ := cmd.Flags().GetString("split")

		// Validate attachment flags
		if err := validateAttachmentFlags(attachTo, asPane, asWindow); err != nil {
			HandleError(err, "validate attachment flags")
		}

		// Determine attachment type
		var attachmentType interfaces.AttachmentType
		if attachTo != "" {
			if asPane {
				attachmentType = interfaces.AttachmentPane
			} else if asWindow {
				attachmentType = interfaces.AttachmentWindow
			} else {
				// Default to pane if attach-to is specified but no type given
				attachmentType = interfaces.AttachmentPane
			}
		}

		// Parse split direction
		var splitDir interfaces.SplitDirection
		if splitDirection != "" {
			switch strings.ToLower(splitDirection) {
			case "h", "horizontal":
				splitDir = interfaces.SplitHorizontal
			case "v", "vertical":
				splitDir = interfaces.SplitVertical
			default:
				HandleError(fmt.Errorf("invalid split direction '%s', use 'h' or 'v'", splitDirection), "parse split direction")
			}
		} else {
			splitDir = interfaces.SplitVertical // Default to vertical split
		}

		// Resolve project path using common function
		projectPath = GetProjectPath(projectPath)

		// Create the session
		sess, err := ctx.Client.CreateSession(api.CreateSessionRequest{
			Name:           sessionName,
			Description:    description,
			ProjectPath:    projectPath,
			AttachTo:       attachTo,
			AttachmentType: attachmentType,
			SplitDirection: splitDir,
		})
		if err != nil {
			HandleError(err, "create session")
		}

		// Enhanced success message
		fmt.Println(ui.SuccessMsg(fmt.Sprintf("Created session '%s'", sess.Name)))
		fmt.Println()

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

// validateAttachmentFlags validates the attachment-related flags
func validateAttachmentFlags(attachTo string, asPane, asWindow bool) error {
	if attachTo == "" && (asPane || asWindow) {
		return fmt.Errorf("--as-pane or --as-window requires --attach-to to be specified")
	}

	if asPane && asWindow {
		return fmt.Errorf("cannot specify both --as-pane and --as-window")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Add flags
	createCmd.Flags().StringP("description", "d", "", "Description for the session")
	createCmd.Flags().StringP("project", "p", "", "Project path for the session (defaults to current directory)")

	// Attachment flags
	createCmd.Flags().StringP("attach-to", "a", "", "Attach to existing session (session name)")
	createCmd.Flags().Bool("as-pane", false, "Create as new pane in existing session")
	createCmd.Flags().Bool("as-window", false, "Create as new window/tab in existing session")
	createCmd.Flags().String("split", "v", "Split direction for panes: 'h' (horizontal) or 'v' (vertical)")
}
