package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DEFAULT_CONFIG_FILE = "claude-pilot.yaml"
const DEFAULT_CONFIG_DIR = ".config/claude-pilot"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "claude-pilot",
	Short: "A CLI tool for managing multiple Claude code sessions",
	Long: color.New(color.FgHiWhite).Sprint(`
╔═══════════════════════════════════════════════════════════════╗
║                        Claude Pilot 🚀                        ║
║                                                               ║
║    A powerful CLI tool for managing multiple Claude code      ║
║    sessions with tmux-inspired commands and beautiful UI.     ║
╚═══════════════════════════════════════════════════════════════╝

`) + color.New(color.FgCyan).Sprint("Available Commands:") + `
  create    Create a new Claude session
  list      List all active sessions
  attach    Attach to a specific session
  detach    Detach from a specific session
  kill      Terminate a session
  kill-all  Terminate all sessions
  tui       Launch interactive terminal UI

Use "claude-pilot [command] --help" for more information about a command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if UI mode is set to TUI in config
		ctx, err := InitializeCommand()
		if err != nil {
			// If we can't initialize, just show help
			cmd.Help()
			return
		}

		// Auto-launch TUI if ui.mode is set to "tui"
		if ctx.Config.UI.Mode == "tui" {
			fmt.Println("TUI mode is not implemented yet")
			return
		}

		// Default behavior: show help
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+DEFAULT_CONFIG_DIR+"/"+DEFAULT_CONFIG_FILE+")")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig sets up environment variable handling for Viper
// The actual config loading is handled by ConfigManager in InitializeCommand()
func initConfig() {
	viper.SetEnvPrefix("CLAUDE_PILOT")
	viper.AutomaticEnv()
}
