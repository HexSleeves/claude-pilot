package cmd

import (
	"fmt"
	"time"

	"claude-pilot/internal/cli"
	"claude-pilot/internal/ui"

	"github.com/spf13/cobra"
)

var detailsCmd = &cobra.Command{
	Use:   "details --id <session-id>",
	Short: "Show details for a specific Claude session",
	Long: `Show details for a specific Claude coding session.

Examples:
  claude-pilot details --id my-project    # Show details for session with ID "my-project"
  claude-pilot details --id abc123        # Show details for session with ID "abc123"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			return fmt.Errorf("failed to initialize command: %w", err)
		}

		// Get flags
		sessionID, _ := cmd.Flags().GetString("id")

		// Handle deprecation warnings for positional arguments
		if len(args) > 0 && sessionID == "" {
			sessionID = args[0]
			if err := CreateDeprecationWarning(ctx, "positional session ID", "--id flag"); err != nil {
				return err
			}
		}

		// Show usage help if no --id provided
		if sessionID == "" {
			// Get available sessions for help message
			sessions, err := ctx.Client.ListSessions()
			if err != nil {
				return fmt.Errorf("failed to list sessions: %w", err)
			}

			if len(sessions) == 0 {
				suggestions := []string{
					"claude-pilot create [session-name]",
				}
				return WriteHelpfulMessage(ctx, "No sessions available. Create a session first.", suggestions)
			} else {
				suggestions := []string{
					"claude-pilot list",
					fmt.Sprintf("claude-pilot details --id %s", sessions[0].Name),
				}
				return WriteHelpfulMessage(ctx, "Session ID is required. Use --id flag to specify which session to show details for.", suggestions)
			}
		}

		// Get the session
		sess, err := ctx.Client.GetSession(sessionID)
		if err != nil {
			// Check if this is a "session not found" error
			if cli.IsErrorCode(err, cli.ErrorCodeSessionNotFound) ||
				cli.IsErrorCategory(err, cli.ErrorCategoryNotFound) {
				// List available sessions in error scenarios
				allSessions, listErr := ctx.Client.ListSessions()
				if listErr == nil && len(allSessions) > 0 {
					suggestions := []string{
						"claude-pilot list",
					}
					for i, session := range allSessions {
						if i < 3 { // Show up to 3 suggestions
							suggestions = append(suggestions, fmt.Sprintf("claude-pilot details --id %s", session.Name))
						}
					}
					err := WriteHelpfulMessage(ctx, fmt.Sprintf("Session '%s' not found.", sessionID), suggestions)
					if err != nil {
						return err
					}
				}

				// Create structured error for exit code 3
				structuredErr := cli.NewNotFoundError("Session", sessionID)
				return structuredErr
			}
			return err
		}

		// Handle different output formats
		switch ctx.OutputWriter.GetFormat() {
		case cli.OutputFormatQuiet:
			// In quiet mode, output only the session ID
			return ctx.OutputWriter.WriteString(sess.ID + "\n")

		case cli.OutputFormatJSON, cli.OutputFormatNDJSON:
			// Convert session to OutputWriter format
			sessionData := cli.SessionData{
				ID:          sess.ID,
				Name:        sess.Name,
				Description: sess.Description,
				Project:     sess.ProjectPath,
				Status:      string(sess.Status),
				CreatedAt:   sess.CreatedAt,
				UpdatedAt:   sess.LastActive,
				PaneCount:   sess.Panes,
			}

			metadata := map[string]string{
				"backend":   ctx.Client.GetBackend(),
				"operation": "details",
				"timestamp": time.Now().Format(time.RFC3339),
			}

			return ctx.OutputWriter.WriteSession(sessionData, metadata)

		default:
			// Human-readable output with session details and next steps
			details := ui.SessionDetailsFormatted(sess, ctx.Client.GetBackend())
			err := ctx.OutputWriter.WriteString(details + "\n\n")
			if err != nil {
				return err
			}

			// Show enhanced next steps
			nextSteps := ui.NextSteps(
				fmt.Sprintf("claude-pilot attach %s", sess.Name),
				"claude-pilot list",
			)
			return ctx.OutputWriter.WriteString(nextSteps + "\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)

	// Add required --id flag
	detailsCmd.Flags().String("id", "", "Session ID or name to show details for (required)")
}
