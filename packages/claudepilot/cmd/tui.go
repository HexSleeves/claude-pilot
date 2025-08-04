package cmd

import (
	"fmt"

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

Global flags are automatically forwarded to the TUI for consistent behavior:
- Configuration settings (--config) are respected
- Color settings (--no-color) affect TUI initialization
- Debug/trace flags control TUI logging verbosity
- Output format settings are applied to any CLI-style output

Examples:
  claude-pilot tui                           # Launch the TUI
  claude-pilot tui --config ~/my-config.yaml # Launch TUI with custom config
  claude-pilot tui --no-color                # Launch TUI with colors disabled
  claude-pilot tui --debug                   # Launch TUI with debug logging
  claude-pilot tui --trace                   # Launch TUI with trace logging
  claude-pilot tui --help                    # Show TUI help`,
	Aliases: []string{"ui", "interactive", "terminal"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize common command context with flag forwarding
		ctx, err := InitializeCommand()
		if err != nil {
			return fmt.Errorf("failed to initialize TUI command: %w", err)
		}

		// Validate TUI prerequisites before launching
		if ctx.Client == nil {
			return fmt.Errorf("TUI requires a valid API client")
		}

		// Color settings are already handled through the command context
		// and will be passed to the TUI through the API client configuration.
		// The --no-color flag is automatically forwarded via the global config.

		// Launch the TUI with the configured client
		// The client already contains all global flag settings and configuration
		if err := tui.RunTui(ctx.Client); err != nil {
			// Use new error taxonomy for TUI failures
			HandleErrorWithContext(ctx, fmt.Errorf("TUI initialization failed: %w", err))
			return nil // HandleErrorWithContext will exit
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
