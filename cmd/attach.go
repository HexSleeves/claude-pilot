package cmd

import (
	"fmt"
	"os"

	"claude-pilot/internal/config"
	"claude-pilot/internal/manager"
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
		identifier := args[0]

		// Load configuration
		configManager := config.NewConfigManager("")
		cfg, err := configManager.Load()
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to load configuration: %v", err)))
			os.Exit(1)
		}

		// Create session manager
		sm, err := manager.NewSessionManager(cfg)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to initialize session manager: %v", err)))
			os.Exit(1)
		}

		// Get the session
		sess, err := sm.GetSession(identifier)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Session not found: %v", err)))
			fmt.Println()
			fmt.Println(ui.InfoMsg("Available sessions:"))

			// Show available sessions
			sessions, err := sm.ListSessions()
			if err != nil {
				fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to list sessions: %v", err)))
				os.Exit(1)
			}

			if len(sessions) == 0 {
				fmt.Println(ui.Dim("  No sessions available"))
				fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
			} else {
				for _, s := range sessions {
					fmt.Printf("  %s %s (%s)\n", ui.Arrow(), ui.Highlight(s.Name), ui.Dim(s.ID[:8]))
				}
			}
			os.Exit(1)
		}

		// Check if session is running
		if !sm.IsSessionRunning(sess.Name) {
			fmt.Println(ui.WarningMsg(fmt.Sprintf("Session '%s' is not running. It may have been terminated.", sess.Name)))
			fmt.Println(ui.InfoMsg("You can recreate it with: claude-pilot create " + sess.Name))
			os.Exit(1)
		}

		// Show session info
		fmt.Println(ui.InfoMsg(fmt.Sprintf("Attaching to session '%s' (%s backend)...", sess.Name, cfg.Backend)))
		fmt.Println(ui.InfoMsg("Use your multiplexer's detach key to exit (tmux: Ctrl+B,D | zellij: Ctrl+P,D)"))
		fmt.Println()

		// Attach to the session
		if err := sm.AttachToSession(identifier); err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to attach to session: %v", err)))
			os.Exit(1)
		}

		// After detaching, we're back to the CLI
		fmt.Println(ui.InfoMsg("Detached from session"))
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)
}
