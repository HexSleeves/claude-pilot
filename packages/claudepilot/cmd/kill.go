package cmd

import (
	"fmt"
	"time"

	"claude-pilot/core/api"
	"claude-pilot/internal/cli"
	"claude-pilot/internal/ui"
	"claude-pilot/shared/components"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [session-name]",
	Short: "Kill a Claude session",
	Long: `Kill (terminate) a Claude coding session.
Supports killing specific sessions by ID/name or all sessions at once.

Examples:
  claude-pilot kill my-session           # Kill specific session (deprecated)
  claude-pilot kill --id my-session      # Kill specific session by ID
  claude-pilot kill --all                # Kill all sessions
  claude-pilot kill --all --yes          # Kill all sessions without confirmation
  claude-pilot kill --id my-session --force  # Force kill without confirmation
  claude-pilot kill --output=json        # Output results in JSON format`,
	Aliases: []string{"terminate", "stop", "delete", "remove", "del"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			return fmt.Errorf("failed to initialize command: %w", err)
		}

		// Get command-specific parameters
		var sessionIdentifier string
		if len(args) > 0 {
			sessionIdentifier = args[0]
			// Show deprecation warning for positional session name
			if err := CreateDeprecationWarning(ctx, "positional session name", "--id flag"); err != nil {
				return err
			}
		}

		// Get flags
		killAll, _ := cmd.Flags().GetBool("all")
		force, _ := cmd.Flags().GetBool("force")
		idFlag, _ := cmd.Flags().GetString("id")

		// Prefer --id flag over positional argument
		if idFlag != "" {
			if sessionIdentifier != "" {
				return cli.NewValidationError(
					"cannot specify both session name as positional argument and --id flag",
					"use either 'claude-pilot kill session-name' or 'claude-pilot kill --id session-name'",
				)
			}
			sessionIdentifier = idFlag
		}

		// Get all sessions
		allSessions, err := ctx.Client.ListSessions()
		if err != nil {
			return cli.WrapError(err, cli.ErrorCodeStorageError, cli.ErrorCategoryInternal,
				"Check if the session storage is accessible and try again.")
		}

		var sessions []*api.Session

		if killAll {
			// Kill all sessions
			sessions = allSessions
		} else if sessionIdentifier == "" {
			// No session identifier provided and --all not specified
			return WriteHelpfulMessage(ctx, "No session identifier provided.",
				[]string{
					"claude-pilot kill --id <session-name>",
					"claude-pilot kill --all",
					"claude-pilot list  # to see available sessions",
				})
		} else {
			// Kill specific session
			var targetSession *api.Session

			for _, sess := range allSessions {
				if sess.Name == sessionIdentifier || sess.ID == sessionIdentifier {
					targetSession = sess
					break
				}
			}

			if targetSession == nil {
				// Return session not found error with exit code 3
				return cli.NewNotFoundError("Session", sessionIdentifier)
			}

			sessions = []*api.Session{targetSession}
		}

		// Check if any sessions to kill
		if len(sessions) == 0 {
			// Handle idempotent behavior for --force
			if force {
				// Return success with operation result
				return ctx.OutputWriter.WriteOperationResult(cli.OperationResult{
					Success: true,
					Message: "No sessions to kill",
					Data:    []cli.SessionData{},
					Metadata: map[string]string{
						"operation": "kill",
						"timestamp": time.Now().Format(time.RFC3339),
						"count":     "0",
					},
				})
			}
			return WriteHelpfulMessage(ctx, "No sessions found to kill.",
				[]string{"claude-pilot list  # to see available sessions"})
		}

		// Show sessions to be terminated (only in human-readable formats)
		if ctx.OutputWriter.GetFormat() == cli.OutputFormatHuman || ctx.OutputWriter.GetFormat() == cli.OutputFormatTable {
			err := ctx.OutputWriter.WriteString(ui.WarningMsg(fmt.Sprintf("About to terminate %d sessions:", len(sessions))) + "\n\n")
			if err != nil {
				return err
			}

			// Convert API sessions to shared table format and display
			sessionData := convertToSessionDataForKill(sessions)
			table := components.NewSessionTable(components.TableConfig{
				ShowHeaders: true,
				Interactive: false,
				MaxRows:     0,
			})
			table.SetSessionData(sessionData)
			err = ctx.OutputWriter.WriteString(table.RenderCLI() + "\n\n")
			if err != nil {
				return err
			}
		}

		// TTY-aware confirmation prompt (unless forced)
		if !force {
			confirmed, err := ConfirmActionWithContext(ctx,
				fmt.Sprintf("Are you sure you want to kill %d session(s)?", len(sessions)),
				false, // default to false for destructive operations
			)
			if err != nil {
				return cli.WrapError(err, cli.ErrorCodeNotInteractiveTerminal, cli.ErrorCategoryUnsupported,
					"Use --yes flag for non-interactive operation or --force to bypass confirmation.")
			}
			if !confirmed {
				// Handle cancellation based on output format
				if ctx.OutputWriter.GetFormat() == cli.OutputFormatJSON || ctx.OutputWriter.GetFormat() == cli.OutputFormatNDJSON {
					return ctx.OutputWriter.WriteOperationResult(cli.OperationResult{
						Success: false,
						Message: "Operation cancelled by user",
						Metadata: map[string]string{
							"operation": "kill",
							"timestamp": time.Now().Format(time.RFC3339),
							"cancelled": "true",
						},
					})
				}
				return ctx.OutputWriter.WriteString(ui.InfoMsg("Operation cancelled") + "\n")
			}
		}

		// Kill sessions with proper error aggregation
		var killErrors []cli.ErrorData
		var successfulSessions []cli.SessionData
		var failedSessions []cli.SessionData

		// Show progress for multiple sessions in TTY mode
		showProgress := len(sessions) > 1 && ctx.TTYDetector.ShowSpinner() && ctx.OutputWriter.GetFormat() == cli.OutputFormatHuman
		if showProgress {
			ctx.OutputWriter.WriteString(fmt.Sprintf("Killing %d sessions...\n", len(sessions)))
		}

		for i, sess := range sessions {
			if showProgress {
				ctx.OutputWriter.WriteString(fmt.Sprintf("\rProgress: %d/%d", i+1, len(sessions)))
			}

			sessionData := cli.SessionData{
				ID:        sess.ID,
				Name:      sess.Name,
				Status:    string(sess.Status),
				Project:   sess.ProjectPath,
				CreatedAt: sess.CreatedAt,
				UpdatedAt: sess.LastActive,
				PaneCount: sess.Panes,
			}

			if err := ctx.Client.KillSession(sess.ID); err != nil {
				// Map error and add to error collection
				errorContract := ctx.ErrorHandler.MapError(err)
				killErrors = append(killErrors, errorContract.ToOutputError())
				failedSessions = append(failedSessions, sessionData)
			} else {
				successfulSessions = append(successfulSessions, sessionData)
			}
		}

		// Clear progress line
		if showProgress {
			ctx.OutputWriter.WriteString("\r\033[K")
		}

		// Handle different output formats
		switch ctx.OutputWriter.GetFormat() {
		case cli.OutputFormatJSON, cli.OutputFormatNDJSON:
			// Return structured operation result
			result := cli.OperationResult{
				Success: len(killErrors) == 0,
				Message: fmt.Sprintf("Killed %d/%d sessions successfully", len(successfulSessions), len(sessions)),
				Data: map[string]any{
					"successful": successfulSessions,
					"failed":     failedSessions,
				},
				Errors: killErrors,
				Metadata: map[string]string{
					"operation":        "kill",
					"timestamp":        time.Now().Format(time.RFC3339),
					"total_count":      fmt.Sprintf("%d", len(sessions)),
					"successful_count": fmt.Sprintf("%d", len(successfulSessions)),
					"failed_count":     fmt.Sprintf("%d", len(killErrors)),
				},
			}
			return ctx.OutputWriter.WriteOperationResult(result)

		case cli.OutputFormatQuiet:
			// In quiet mode, only output essential results
			if len(killErrors) > 0 {
				return ctx.OutputWriter.WriteString(fmt.Sprintf("Failed to kill %d sessions\n", len(killErrors)))
			}
			return nil

		default:
			// Human-readable output with unified feedback
			if len(killErrors) > 0 {
				err := ctx.OutputWriter.WriteString("\n" + ui.ErrorMsg("Failed to kill some sessions:") + "\n")
				if err != nil {
					return err
				}
				for _, killErr := range killErrors {
					err := ctx.OutputWriter.WriteString(fmt.Sprintf("  %s %s\n", ui.ErrorMsg("✗"), killErr.Message))
					if err != nil {
						return err
					}
				}
			}

			// Show unified summary
			if len(successfulSessions) > 0 {
				if len(killErrors) > 0 {
					err := ctx.OutputWriter.WriteString(fmt.Sprintf("\n%s Successfully killed %d out of %d sessions\n",
						ui.SuccessMsg("✓"), len(successfulSessions), len(sessions)))
					if err != nil {
						return err
					}
				} else {
					err := ctx.OutputWriter.WriteString(fmt.Sprintf("%s Successfully killed %d session(s)\n",
						ui.SuccessMsg("✓"), len(successfulSessions)))
					if err != nil {
						return err
					}
				}
			}

			// Show remaining sessions for context
			remainingSessions, err := ctx.Client.ListSessions()
			if err != nil {
				err := ctx.OutputWriter.WriteString(fmt.Sprintf("%s Warning: Could not list remaining sessions: %v\n",
					ui.WarningMsg("⚠"), err))
				if err != nil {
					return err
				}
			} else if len(remainingSessions) > 0 {
				err := ctx.OutputWriter.WriteString("\n")
				if err != nil {
					return err
				}
				ui.DisplayRemainingSessionsInfo(remainingSessions)
			}

			// Return error if there were any failures
			if len(killErrors) > 0 {
				return cli.NewStructuredError(
					"partial_kill_failure",
					cli.ErrorCategoryInternal,
					fmt.Sprintf("Failed to kill %d out of %d sessions", len(killErrors), len(sessions)),
					"Check the individual error messages above for remediation steps.",
					map[string]string{
						"successful_count": fmt.Sprintf("%d", len(successfulSessions)),
						"failed_count":     fmt.Sprintf("%d", len(killErrors)),
						"total_count":      fmt.Sprintf("%d", len(sessions)),
					},
				)
			}
			return nil
		}
	},
}

// convertToSessionDataForKill converts API sessions to the shared table SessionData format for kill command
func convertToSessionDataForKill(sessions []*api.Session) []components.SessionData {
	sessionData := make([]components.SessionData, len(sessions))

	for i, sess := range sessions {
		sessionData[i] = components.SessionData{
			ID:          sess.ID,
			Name:        sess.Name,
			Status:      string(sess.Status),
			Backend:     sess.Backend,
			Created:     sess.CreatedAt,
			LastActive:  sess.LastActive,
			ProjectPath: sess.ProjectPath,
			Panes:       sess.Panes,
		}
	}

	return sessionData
}

func init() {
	rootCmd.AddCommand(killCmd)

	// Add flags
	killCmd.Flags().BoolP("all", "a", false, "Kill all sessions")
	killCmd.Flags().BoolP("force", "f", false, "Force kill without confirmation")
	killCmd.Flags().String("id", "", "Session ID/name to kill (alternative to positional argument)")
}
