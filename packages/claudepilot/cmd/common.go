package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"claude-pilot/core/api"
	"claude-pilot/internal/cli"
	"claude-pilot/internal/ui"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// CommandContext holds common dependencies for all commands
type CommandContext struct {
	Client       *api.Client
	OutputWriter cli.OutputWriter
	ErrorHandler *cli.ErrorMapper
	TTYDetector  *cli.TTYDetector
	RequestID    string
	AutoYes      bool
}

// InitializeCommand handles common initialization for all commands
// This creates an API client that provides access to all functionality
func InitializeCommand() (*CommandContext, error) {
	return InitializeCommandWithContext("")
}

// InitializeCommandWithContext handles initialization with request context
func InitializeCommandWithContext(requestID string) (*CommandContext, error) {
	if requestID == "" {
		requestID = uuid.New().String()[:8]
	}

	// Get configuration from viper (set by root command PersistentPreRunE)
	outputFormat := cli.OutputFormat(viper.GetString("cli.output_format"))
	noColor := viper.GetBool("cli.no_color")
	yesFlag := viper.GetBool("cli.yes")

	// Initialize TTY detector
	ttyDetector := cli.NewTTYDetector()
	isTTY := ttyDetector.IsInteractive()

	// Determine color support
	colorEnabled := !noColor && ttyDetector.GetColorSupport()

	// Create output writer
	outputWriter := cli.NewOutputWriter(outputFormat, isTTY, colorEnabled)

	// Create error handler
	errorHandler := cli.NewErrorMapper(requestID)

	// Get verbose flag from viper for backward compatibility
	verbose := viper.GetBool("verbose")

	// Create API client with configuration
	client, err := api.NewClient(api.ClientConfig{
		ConfigFile: cfgFile,
		Verbose:    verbose,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	return &CommandContext{
		Client:       client,
		OutputWriter: outputWriter,
		ErrorHandler: errorHandler,
		TTYDetector:  ttyDetector,
		RequestID:    requestID,
		AutoYes:      yesFlag,
	}, nil
}

// HandleError provides consistent error handling and exit across all commands
// This eliminates the duplicated error handling pattern that appears in every command
// Deprecated: Use HandleErrorWithContext instead
func HandleError(err error, action string) {
	fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to %s: %v", action, err)))
	os.Exit(1)
}

// HandleErrorWithContext handles errors using the new error taxonomy and output system
func HandleErrorWithContext(ctx *CommandContext, err error) {
	if err == nil {
		return
	}

	exitCode := cli.HandleError(err, ctx.OutputWriter, ctx.RequestID)
	os.Exit(int(exitCode))
}

// HandleErrorWithContextAndExit handles errors and returns the exit code without exiting
// This is useful for testing or when you want to handle the exit yourself
func HandleErrorWithContextAndExit(ctx *CommandContext, err error) cli.ExitCode {
	if err == nil {
		return cli.ExitCodeSuccess
	}

	return cli.HandleError(err, ctx.OutputWriter, ctx.RequestID)
}

// ConfirmAction handles user confirmation prompts consistently
// This eliminates the duplicated confirmation logic in kill commands
// Deprecated: Use ConfirmActionWithContext instead
func ConfirmAction(message string) bool {
	fmt.Print(ui.Prompt(message))
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error as empty input is valid
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// ConfirmActionWithContext handles user confirmation with TTY detection and timeout
func ConfirmActionWithContext(ctx *CommandContext, message string, defaultValue bool) (bool, error) {
	options := &cli.ConfirmationOptions{
		Message:      message,
		DefaultValue: defaultValue,
		Timeout:      30 * time.Second,
		AutoYes:      ctx.AutoYes,
	}

	return ctx.TTYDetector.ConfirmWithTimeout(options)
}

// ConfirmActionWithOptions handles user confirmation with custom options
func ConfirmActionWithOptions(ctx *CommandContext, options *cli.ConfirmationOptions) (bool, error) {
	if !options.AutoYes {
		options.AutoYes = ctx.AutoYes
	}
	return ctx.TTYDetector.ConfirmWithTimeout(options)
}

// GetProjectPath handles project path resolution with fallback to current directory
// This eliminates the duplicated project path logic in create command
// GetProjectPath handles project path resolution with fallback to current directory
// This eliminates the duplicated project path logic in create command
func GetProjectPath(projectPath string) string {
	if projectPath == "" {
		cwd, err := os.Getwd()
		if err == nil {
			return cwd
		}
		return ""
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err == nil {
		return absPath
	}
	return projectPath
}

// NoSessionsFoundMessageForFilter displays helpful messages when no sessions are found
// Deprecated: Use WriteHelpfulMessage with context instead
func NoSessionsFoundMessageForFilter() {
	fmt.Println(ui.Dim("No sessions found. Try another filter or create a new session."))
	fmt.Println()
	fmt.Println(ui.InfoMsg("Show active sessions:"))
	fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot list --active"))
	fmt.Println()
	fmt.Println(ui.InfoMsg("Show inactive sessions:"))
	fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot list --inactive"))
	fmt.Println()
	fmt.Println(ui.InfoMsg("Create a new session:"))
	fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot create [session-name]"))
	fmt.Println()
}

// WriteHelpfulMessage writes helpful messages using the output writer
func WriteHelpfulMessage(ctx *CommandContext, message string, suggestions []string) error {
	if ctx.OutputWriter.GetFormat() == cli.OutputFormatQuiet {
		return nil // Don't output help in quiet mode
	}

	// For JSON/NDJSON formats, don't output help messages
	if ctx.OutputWriter.GetFormat() == cli.OutputFormatJSON || ctx.OutputWriter.GetFormat() == cli.OutputFormatNDJSON {
		return nil
	}

	// Write the main message
	if err := ctx.OutputWriter.WriteString(ui.Dim(message) + "\n\n"); err != nil {
		return err
	}

	// Write suggestions
	for _, suggestion := range suggestions {
		// Add title for suggestions
		if err := ctx.OutputWriter.WriteString(ui.Dim("Suggestions:") + "\n"); err != nil {
			return err
		}

		// Add arrow and suggestion
		if err := ctx.OutputWriter.WriteString(fmt.Sprintf("  %s %s\n", ui.Arrow(), ui.Highlight(suggestion))); err != nil {
			return err
		}
	}

	return ctx.OutputWriter.WriteString("\n")
}

// CreateDeprecationWarning shows deprecation warnings for legacy usage patterns
func CreateDeprecationWarning(ctx *CommandContext, oldUsage, newUsage string) error {
	if ctx.OutputWriter.GetFormat() == cli.OutputFormatQuiet {
		return nil
	}

	warning := fmt.Sprintf("Warning: '%s' is deprecated. Use '%s' instead.", oldUsage, newUsage)
	return ctx.OutputWriter.WriteString(ui.WarningMsg(warning) + "\n")
}
