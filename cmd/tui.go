package cmd

import (
	"claude-pilot/internal/tui"

	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal UI",
	Long: `Launch the interactive terminal user interface (TUI) for managing Claude sessions.

The TUI provides a full-screen interactive experience with:
- Real-time session list with status indicators
- Vim-style navigation (hjkl) and standard key bindings
- Session creation, attachment, and deletion
- Help screen with key bindings

Key bindings:
- q: quit
- h/j/k/l: navigate (vim-style)
- Enter: attach to session
- c: create new session
- d: delete session
- r: refresh
- ?: show help`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "initialize")
		}

		// Launch TUI
		if err := tui.RunTUI(ctx.SessionManager); err != nil {
			HandleError(err, "run TUI")
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
