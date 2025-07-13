package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

Use "claude-pilot [command] --help" for more information about a command.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.claude-pilot.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".claude-pilot")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
