package cmd

import (
	"fmt"

	"claude-pilot/core/api"
	"claude-pilot/internal/ui"
	"claude-pilot/shared/components"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [session-name]",
	Short: "Kill a Claude session",
	Long: `Kill (terminate) a Claude coding session.
If no session name is provided, kills all sessions.

Examples:
  claude-pilot kill my-session    # Kill specific session
  claude-pilot kill --all         # Kill all sessions
  claude-pilot kill --force       # Kill without confirmation`,
	Aliases: []string{"terminate", "stop", "delete", "remove", "del"},
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		// Get flags
		killAll, _ := cmd.Flags().GetBool("all")
		force, _ := cmd.Flags().GetBool("force")

		// Get all sessions
		allSessions, err := ctx.Client.ListSessions()
		if err != nil {
			HandleError(err, "list sessions")
		}

		var sessions []*api.Session

		if killAll {
			// Kill all sessions
			sessions = allSessions
		} else if len(args) == 0 {
			// No session name provided and --all not specified
			fmt.Println(ui.ErrorMsg("No session name provided"))
			fmt.Println()
			fmt.Println(ui.InfoMsg("Available sessions:"))
			ui.DisplayAvailableSessions(allSessions)
			fmt.Println()
			fmt.Println(ui.InfoMsg("Usage:"))
			fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot kill <session-name>"))
			fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot kill --all"))
			return
		} else {
			// Kill specific session
			sessionName := args[0]
			var targetSession *api.Session

			for _, sess := range allSessions {
				if sess.Name == sessionName || sess.ID == sessionName {
					targetSession = sess
					break
				}
			}

			if targetSession == nil {
				fmt.Println(ui.ErrorMsg(fmt.Sprintf("Session '%s' not found", sessionName)))
				fmt.Println()
				fmt.Println(ui.InfoMsg("Available sessions:"))
				ui.DisplayAvailableSessions(allSessions)
				return
			}

			sessions = []*api.Session{targetSession}
		}

		// Check if any sessions to kill
		if len(sessions) == 0 {
			fmt.Println(ui.InfoMsg("No sessions to kill"))
			return
		}

		// Show sessions to be terminated
		fmt.Println(ui.WarningMsg(fmt.Sprintf("About to terminate %d sessions:", len(sessions))))
		fmt.Println()

		// Convert API sessions to shared table format and display
		sessionData := convertToSessionDataForKill(sessions)
		table := components.NewSessionTable(components.TableConfig{
			ShowHeaders: true,
			Interactive: false,
			MaxRows:     0,
		})
		table.SetSessionData(sessionData)
		fmt.Println(table.RenderCLI())
		fmt.Println()

		// Confirmation prompt (unless forced) using common function
		if !force {
			if !ConfirmAction(fmt.Sprintf("Are you sure you want to kill %d session(s)? [y/N]: ", len(sessions))) {
				fmt.Println(ui.InfoMsg("Operation cancelled"))
				return
			}
		}

		// Kill sessions
		var errors []string
		for _, sess := range sessions {
			if err := ctx.Client.KillSession(sess.ID); err != nil {
				errors = append(errors, fmt.Sprintf("Failed to kill session %s: %v", sess.Name, err))
			} else {
				fmt.Printf("%s Session %s killed successfully\n", ui.SuccessMsg(""), ui.Highlight(sess.Name))
			}
		}

		// Report any errors
		if len(errors) > 0 {
			fmt.Println()
			fmt.Println(ui.ErrorMsg("Some sessions could not be killed:"))
			for _, errMsg := range errors {
				fmt.Printf("  %s %s\n", ui.ErrorMsg("✗"), errMsg)
			}
		}

		// Show remaining sessions
		remainingSessions, err := ctx.Client.ListSessions()
		if err != nil {
			fmt.Printf("%s Warning: Could not list remaining sessions: %v\n", ui.WarningMsg("⚠"), err)
		} else {
			fmt.Println()
			ui.DisplayRemainingSessionsInfo(remainingSessions)
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
}
