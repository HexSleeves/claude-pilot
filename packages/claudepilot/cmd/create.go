package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"claude-pilot/core/api"
	"claude-pilot/internal/cli"
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
  claude-pilot create --id my-session-123          # Create session with specific ID
  `,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			return fmt.Errorf("failed to initialize command: %w", err)
		}

		// Get command-specific parameters
		var sessionName string
		if len(args) > 0 {
			sessionName = args[0]
			// Show deprecation warning for positional session name
			if err := CreateDeprecationWarning(ctx, "positional session name", "--id flag"); err != nil {
				return err
			}
		}

		// Get flags
		description, _ := cmd.Flags().GetString("description")
		projectPath, _ := cmd.Flags().GetString("project")
		attachTo, _ := cmd.Flags().GetString("attach-to")
		asPane, _ := cmd.Flags().GetBool("as-pane")
		asWindow, _ := cmd.Flags().GetBool("as-window")
		splitDirection, _ := cmd.Flags().GetString("split")
		idFlag, _ := cmd.Flags().GetString("id")

		// Prefer --id flag over positional argument
		if idFlag != "" {
			if sessionName != "" {
				return cli.NewValidationError(
					"cannot specify both session name as positional argument and --id flag",
					"use either 'claude-pilot create session-name' or 'claude-pilot create --id session-name'",
				)
			}
			sessionName = idFlag
		}

		// Validate attachment flags
		if err := validateAttachmentFlags(attachTo, asPane, asWindow); err != nil {
			return cli.WrapError(err, cli.ErrorCodeInvalidFlagValue, cli.ErrorCategoryValidation,
				"Check the command usage with --help for proper flag combinations.")
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
				return cli.NewValidationError(
					fmt.Sprintf("invalid split direction '%s'", splitDirection),
					"use 'h' (horizontal) or 'v' (vertical) for split direction",
				)
			}
		} else {
			splitDir = interfaces.SplitVertical // Default to vertical split
		}

		// Resolve project path using common function
		projectPath = GetProjectPath(projectPath)

		// Create the session
		var sessionResult *interfaces.Session

		// Show progress spinner for TTY users
		if ctx.TTYDetector.ShowSpinner() {
			fmt.Fprintln(os.Stderr, "Creating session...")
		}

		sessionResult, err = ctx.Client.CreateSession(api.CreateSessionRequest{
			Name:           sessionName,
			Description:    description,
			ProjectPath:    projectPath,
			AttachTo:       attachTo,
			AttachmentType: attachmentType,
			SplitDirection: splitDir,
		})

		// Clear progress spinner
		if ctx.TTYDetector.ShowSpinner() {
			fmt.Fprintf(os.Stderr, "\r\033[K") // Clear line
		}

		if err != nil {
			return err
		}

		// Handle different output formats
		switch ctx.OutputWriter.GetFormat() {
		case cli.OutputFormatQuiet:
			// In quiet mode, output only the session ID
			return ctx.OutputWriter.WriteString(sessionResult.ID + "\n")

		case cli.OutputFormatJSON, cli.OutputFormatNDJSON:
			// Convert session to OutputWriter format
			sessionData := cli.SessionData{
				ID:          sessionResult.ID,
				Name:        sessionResult.Name,
				Description: sessionResult.Description,
				Project:     sessionResult.ProjectPath,
				Status:      string(sessionResult.Status),
				CreatedAt:   sessionResult.CreatedAt,
				UpdatedAt:   sessionResult.LastActive,
				PaneCount:   sessionResult.Panes,
			}

			metadata := map[string]string{
				"backend":   ctx.Client.GetBackend(),
				"operation": "create",
				"timestamp": time.Now().Format(time.RFC3339),
			}

			return ctx.OutputWriter.WriteSession(sessionData, metadata)

		default:
			// Human-readable output with success message and next steps
			err := ctx.OutputWriter.WriteString(ui.SuccessMsg(fmt.Sprintf("Created session '%s'", sessionResult.Name)) + "\n\n")
			if err != nil {
				return err
			}

			// Show enhanced session details
			details := ui.SessionDetailsFormatted(sessionResult, ctx.Client.GetBackend())
			err = ctx.OutputWriter.WriteString(details + "\n\n")
			if err != nil {
				return err
			}

			// Show enhanced next steps
			nextSteps := ui.NextSteps(
				fmt.Sprintf("claude-pilot attach %s", sessionResult.Name),
				"claude-pilot list",
			)
			return ctx.OutputWriter.WriteString(nextSteps + "\n")
		}
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
	createCmd.Flags().String("id", "", "Session ID/name (alternative to positional argument)")

	// Attachment flags
	createCmd.Flags().StringP("attach-to", "a", "", "Attach to existing session (session name)")
	createCmd.Flags().Bool("as-pane", false, "Create as new pane in existing session")
	createCmd.Flags().Bool("as-window", false, "Create as new window/tab in existing session")
	createCmd.Flags().String("split", "v", "Split direction for panes: 'h' (horizontal) or 'v' (vertical)")
}
