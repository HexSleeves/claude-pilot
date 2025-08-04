package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"claude-pilot/internal/cli"
	"claude-pilot/internal/ui"

	"github.com/spf13/cobra"
)

var attachCmd = &cobra.Command{
	Use:   "attach [session-name-or-id]",
	Short: "Attach to a Claude session for interactive communication",
	Long: `Attach to a specific Claude coding session to start interactive communication.
This opens a terminal-based multiplexer session where you can interact with your
coding environment directly.

This command requires an interactive terminal (TTY) and cannot be used in 
non-interactive environments like CI/CD pipelines or automated scripts.

Examples:
  claude-pilot attach my-session            # Attach to session by name  
  claude-pilot attach abc123def             # Attach to session by ID
  claude-pilot attach --id my-session       # Attach using explicit --id flag`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			return fmt.Errorf("failed to initialize command: %w", err)
		}

		// TTY-Only Operation: Check if we're in an interactive terminal
		if !ctx.TTYDetector.IsInteractive() {
			return cli.NewUnsupportedError(
				"attach command requires an interactive terminal",
				"Use 'claude-pilot list' or 'claude-pilot details --id <session>' for non-interactive session inspection. Consider using --read-only flag for future non-interactive access (coming soon).",
			)
		}

		// Argument Standardization: Handle both positional and --id flag
		var identifier string
		if len(args) > 0 {
			identifier = args[0]
			// Show deprecation warning for positional session identifier
			if err := CreateDeprecationWarning(ctx, "positional session identifier", "--id flag"); err != nil {
				return err
			}
		}

		// Get --id flag
		idFlag, _ := cmd.Flags().GetString("id")

		// Prefer --id flag over positional argument
		if idFlag != "" {
			if identifier != "" {
				return cli.NewValidationError(
					"cannot specify both session identifier as positional argument and --id flag",
					"use either 'claude-pilot attach session-name' or 'claude-pilot attach --id session-name'",
				)
			}
			identifier = idFlag
		}

		// Require session identifier
		if identifier == "" {
			return cli.NewValidationError(
				"session identifier is required",
				"specify session name or ID: 'claude-pilot attach my-session' or 'claude-pilot attach --id my-session'",
			)
		}

		// Get the session
		sess, err := ctx.Client.GetSession(identifier)
		if err != nil {
			// Provide helpful session listing when target session is not found
			if cli.IsErrorCategory(err, cli.ErrorCategoryNotFound) {
				// Show available sessions for context
				sessions, listErr := ctx.Client.ListSessions()
				if listErr == nil && len(sessions) > 0 {
					suggestions := []string{
						"claude-pilot list",
						"claude-pilot create <session-name>",
					}
					if len(sessions) <= 5 {
						// Show a few session names as suggestions if list is short
						for i, session := range sessions {
							if i >= 3 {
								break // Limit to first 3 suggestions
							}
							suggestions = append(suggestions, fmt.Sprintf("claude-pilot attach %s", session.Name))
						}
					}

					writeErr := WriteHelpfulMessage(ctx, "Session not found. Available sessions:", suggestions)
					if writeErr != nil {
						return writeErr
					}
				}
			}
			return err
		}

		// Check if session is running - map to appropriate error category
		if !ctx.Client.IsSessionRunning(sess.Name) {
			return cli.NewStructuredError(
				cli.ErrorCodeSessionNotRunning,
				cli.ErrorCategoryConflict,
				fmt.Sprintf("Session '%s' is not running. It may have been terminated.", sess.Name),
				fmt.Sprintf("Start the session with: claude-pilot create %s", sess.Name),
				map[string]string{
					"sessionName": sess.Name,
					"sessionID":   sess.ID,
					"action":      "restart",
				},
			)
		}

		// Interactive Behavior: Show clear attachment instructions
		backend := ctx.Client.GetBackend()
		err = ctx.OutputWriter.WriteString(ui.InfoMsg(fmt.Sprintf("Attaching to session '%s' (%s backend)...", sess.Name, backend)) + "\n")
		if err != nil {
			return err
		}

		// Show multiplexer-specific detach keys
		var detachKey string
		switch backend {
		case "tmux":
			detachKey = "Ctrl+B, then D"
		default:
			detachKey = "your multiplexer's detach key"
		}

		err = ctx.OutputWriter.WriteString(ui.InfoMsg(fmt.Sprintf("Use %s to detach from the session", detachKey)) + "\n\n")
		if err != nil {
			return err
		}

		// Emit attachment event for NDJSON mode
		if ctx.OutputWriter.GetFormat() == cli.OutputFormatNDJSON {
			attachmentEvent := map[string]interface{}{
				"type":        "session_attached",
				"sessionId":   sess.ID,
				"sessionName": sess.Name,
				"timestamp":   time.Now().Format(time.RFC3339),
				"backend":     backend,
			}

			result := cli.OperationResult{
				Success: true,
				Message: "Session attached",
				Data:    attachmentEvent,
			}

			if err := ctx.OutputWriter.WriteOperationResult(result); err != nil {
				// Don't fail attachment for logging errors, just proceed
				// Log error could be handled here if needed
			}
		}

		// Handle SIGINT gracefully during attachment
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Channel to handle attachment completion
		attachDone := make(chan error, 1)

		// Start attachment in goroutine
		go func() {
			attachDone <- ctx.Client.AttachToSession(identifier)
		}()

		// Wait for either attachment completion or signal
		select {
		case err := <-attachDone:
			// Normal attachment completion
			signal.Stop(sigChan)
			if err != nil {
				return cli.WrapError(err, cli.ErrorCodeAttachmentFailed, cli.ErrorCategoryInternal,
					"Check if the session is still running and the multiplexer is accessible.")
			}
		case sig := <-sigChan:
			// Signal received during attachment
			signal.Stop(sigChan)
			return cli.NewStructuredError(
				"attachment_interrupted",
				cli.ErrorCategoryInternal,
				fmt.Sprintf("Attachment interrupted by signal: %v", sig),
				"The attachment was cancelled. The session may still be running.",
				nil,
			)
		}

		// Post-detachment status message
		err = ctx.OutputWriter.WriteString(ui.InfoMsg("Detached from session") + "\n")
		if err != nil {
			return err
		}

		// Emit detachment event for NDJSON mode
		if ctx.OutputWriter.GetFormat() == cli.OutputFormatNDJSON {
			detachmentEvent := map[string]interface{}{
				"type":        "session_detached",
				"sessionId":   sess.ID,
				"sessionName": sess.Name,
				"timestamp":   time.Now().Format(time.RFC3339),
				"backend":     backend,
			}

			result := cli.OperationResult{
				Success: true,
				Message: "Session detached",
				Data:    detachmentEvent,
			}

			if err := ctx.OutputWriter.WriteOperationResult(result); err != nil {
				// Don't fail command for logging errors, just proceed
				// Log error could be handled here if needed
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)

	// Add --id flag for consistent argument handling
	attachCmd.Flags().String("id", "", "Session ID/name to attach to (alternative to positional argument)")
}
