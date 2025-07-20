package cmd

import (
	"claude-pilot/tui"

	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command that launches the interactive terminal interface
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal UI for managing Claude sessions",
	Long: `Launch the interactive Terminal User Interface (TUI) for Claude Pilot.
This provides a full-screen interactive interface for managing your Claude coding sessions
with features like:

- Browse and filter sessions with an interactive table
- Create new sessions with guided forms
- Attach to sessions directly from the interface
- Kill sessions with confirmation dialogs
- Real-time session status updates
- Keyboard shortcuts for efficient navigation

The TUI uses the same backend as the CLI commands but provides a more visual
and interactive experience for session management.

Examples:
  claude-pilot tui                    # Launch the TUI
  claude-pilot tui --help             # Show TUI help`,
	Aliases: []string{"ui", "interactive", "terminal"},
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize command")
		}

		// Launch the TUI directly using the shared client
		if err := tui.RunTui(ctx.Client); err != nil {
			HandleError(err, "run TUI")
		}

	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
