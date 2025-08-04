package cmd

import (
	"fmt"
	"os"

	"log/slog"

	"claude-pilot/internal/cli"
	"claude-pilot/internal/ui"
	"claude-pilot/shared/styles"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DEFAULT_CONFIG_FILE = "claude-pilot.yaml"
const DEFAULT_CONFIG_DIR = ".config/claude-pilot"

var (
	cfgFile     string
	outputFlag  string
	debugFlag   bool
	traceFlag   bool
	noColorFlag bool
	yesFlag     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "claude-pilot",
	Short: "A CLI tool for managing multiple Claude code sessions",
	Long: styles.Banner("Claude Pilot 🚀", "A powerful CLI tool for managing multiple Claude code sessions") + "\n\n" + ui.CommandList(map[string]string{
		"create":   "Create a new Claude session",
		"list":     "List all active sessions",
		"attach":   "Attach to a specific session",
		"kill":     "Terminate a session",
		"kill-all": "Terminate all sessions",
		"tui":      "Launch interactive terminal UI",
	}) + "\n\n" + getGlobalFlagsHelp() + "\n\nUse \"claude-pilot [command] --help\" for more information about a command.",
	PersistentPreRunE: initializeGlobalConfig,
	SilenceErrors:     true, // We'll handle errors ourselves
	SilenceUsage:      true, // Don't show usage on errors
	Run: func(cmd *cobra.Command, args []string) {
		// Check if UI mode is set to TUI in config
		ctx, err := InitializeCommand()
		if err != nil {
			// If we can't initialize, just show help
			_ = cmd.Help()
			return
		}

		// Auto-launch TUI if ui.mode is set to "tui"
		if ctx.Client.GetConfig().UI.Mode == "tui" {
			fmt.Println("TUI mode is not implemented yet")
			return
		}

		// Default behavior: show help
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		// Try to initialize context for proper error handling
		ctx, initErr := InitializeCommand()
		if initErr != nil {
			// If we can't initialize context, fall back to basic error handling
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		// Use the structured error handling
		exitCode := HandleErrorWithContextAndExit(ctx, err)
		os.Exit(int(exitCode))
	}
	return nil
}

func RootCmd() *cobra.Command {
	return rootCmd
}

// initConfig sets up environment variable handling for Viper
// The actual config loading is handled by ConfigManager in InitializeCommand()
func initConfig() {
	viper.SetEnvPrefix("CLAUDE_PILOT")
	viper.AutomaticEnv()
}

// initializeGlobalConfig validates and initializes global configuration
func initializeGlobalConfig(cmd *cobra.Command, args []string) error {
	// Validate output format
	outputFormat := cli.OutputFormat(outputFlag)
	if !outputFormat.IsValid() {
		return fmt.Errorf("invalid output format '%s'. Valid formats: human, table, json, ndjson, quiet", outputFlag)
	}

	// Validate flag combinations
	if traceFlag && debugFlag {
		fmt.Fprintf(os.Stderr, "Warning: --trace overrides --debug flag\\n")
	}
	if (traceFlag || debugFlag) && viper.GetBool("verbose") {
		fmt.Fprintf(os.Stderr, "Warning: explicit debug/trace flags override --verbose\\n")
	}

	// Set up logger level based on flags
	logLevel := slog.LevelInfo
	if traceFlag {
		logLevel = slog.LevelDebug - 4 // Trace level equivalent
	} else if debugFlag {
		logLevel = slog.LevelDebug
	} else if viper.GetBool("verbose") {
		logLevel = slog.LevelDebug
	}

	// Initialize global logger (this should be done in each command's context)
	// For now, we'll store the configuration for commands to use
	viper.Set("cli.output_format", outputFlag)
	viper.Set("cli.no_color", noColorFlag)
	viper.Set("cli.yes", yesFlag)
	viper.Set("cli.log_level", logLevel)

	return nil
}

// getGlobalFlagsHelp returns help text for global flags
func getGlobalFlagsHelp() string {
	return `Global Flags:
  -o, --output string    Output format (human|table|json|ndjson|quiet) (default "human")
      --debug           Enable debug logging (overrides --verbose)
      --trace           Enable trace logging (overrides --debug and --verbose)
      --no-color        Disable ANSI colors in output
      --yes             Accept defaults for prompts (non-interactive mode)
  -v, --verbose         Verbose output
      --config string   Config file (default is $HOME/.config/claude-pilot/claude-pilot.yaml)`
}

func init() {
	cobra.OnInitialize(initConfig)

	// Existing global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+DEFAULT_CONFIG_DIR+"/"+DEFAULT_CONFIG_FILE+")")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// New standardized global flags
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "human", "output format (human|table|json|ndjson|quiet)")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "enable debug logging (overrides --verbose)")
	rootCmd.PersistentFlags().BoolVar(&traceFlag, "trace", false, "enable trace logging (overrides --debug and --verbose)")
	rootCmd.PersistentFlags().BoolVar(&noColorFlag, "no-color", false, "disable ANSI colors in output")
	rootCmd.PersistentFlags().BoolVar(&yesFlag, "yes", false, "accept defaults for prompts (non-interactive mode)")

	// Bind flags to viper
	err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		fmt.Println("Error binding verbose flag to viper:", err)
	}
	err = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	if err != nil {
		fmt.Println("Error binding output flag to viper:", err)
	}
}
