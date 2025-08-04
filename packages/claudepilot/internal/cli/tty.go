package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// TTYDetector provides methods for TTY detection and terminal interaction
type TTYDetector struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// NewTTYDetector creates a new TTYDetector instance
func NewTTYDetector() *TTYDetector {
	return &TTYDetector{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// NewTTYDetectorWithIO creates a TTYDetector with custom IO streams
func NewTTYDetectorWithIO(stdin io.Reader, stdout, stderr io.Writer) *TTYDetector {
	return &TTYDetector{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// IsInteractive detects if the current session is interactive (TTY)
func (td *TTYDetector) IsInteractive() bool {
	// Check environment variable overrides first
	if forceTTY := os.Getenv("FORCE_TTY"); forceTTY == "1" || forceTTY == "true" {
		return true
	}
	if noTTY := os.Getenv("NO_TTY"); noTTY == "1" || noTTY == "true" {
		return false
	}

	// Check if stdin and stdout are terminals
	stdinFd := int(os.Stdin.Fd())
	stdoutFd := int(os.Stdout.Fd())

	return term.IsTerminal(stdinFd) && term.IsTerminal(stdoutFd)
}

// IsStdinTTY checks if stdin is a terminal
func (td *TTYDetector) IsStdinTTY() bool {
	if noTTY := os.Getenv("NO_TTY"); noTTY == "1" || noTTY == "true" {
		return false
	}
	if forceTTY := os.Getenv("FORCE_TTY"); forceTTY == "1" || forceTTY == "true" {
		return true
	}

	return term.IsTerminal(int(os.Stdin.Fd()))
}

// IsStdoutTTY checks if stdout is a terminal
func (td *TTYDetector) IsStdoutTTY() bool {
	if noTTY := os.Getenv("NO_TTY"); noTTY == "1" || noTTY == "true" {
		return false
	}
	if forceTTY := os.Getenv("FORCE_TTY"); forceTTY == "1" || forceTTY == "true" {
		return true
	}

	return term.IsTerminal(int(os.Stdout.Fd()))
}

// GetTerminalSize returns the width and height of the terminal
func (td *TTYDetector) GetTerminalSize() (width, height int, err error) {
	if !td.IsStdoutTTY() {
		// Default to 80x24 for non-TTY environments
		return 80, 24, nil
	}

	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to default size
		return 80, 24, nil
	}

	return width, height, nil
}

// ConfirmationOptions holds options for confirmation prompts
type ConfirmationOptions struct {
	Message      string
	DefaultValue bool
	Timeout      time.Duration
	YesResponses []string
	NoResponses  []string
	AutoYes      bool // Set by --yes flag
}

// DefaultConfirmationOptions returns default confirmation options
func DefaultConfirmationOptions() *ConfirmationOptions {
	return &ConfirmationOptions{
		Message:      "Do you want to continue?",
		DefaultValue: false,
		Timeout:      30 * time.Second,
		YesResponses: []string{"y", "yes", "true", "1"},
		NoResponses:  []string{"n", "no", "false", "0"},
		AutoYes:      false,
	}
}

// ConfirmWithTimeout displays a confirmation prompt with timeout support
func (td *TTYDetector) ConfirmWithTimeout(options *ConfirmationOptions) (bool, error) {
	if options == nil {
		options = DefaultConfirmationOptions()
	}

	if options.YesResponses == nil {
		options.YesResponses = []string{"y", "yes", "true", "1"}
	}
	if options.NoResponses == nil {
		options.NoResponses = []string{"n", "no", "false", "0"}
	}

	// If --yes flag is set, automatically confirm
	if options.AutoYes {
		return true, nil
	}

	// If not interactive, return default value
	if !td.IsInteractive() {
		return options.DefaultValue, nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	// Display prompt
	defaultText := "N"
	if options.DefaultValue {
		defaultText = "Y"
	}
	prompt := fmt.Sprintf("%s [y/N] (default: %s, timeout: %v): ",
		options.Message, defaultText, options.Timeout)

	fmt.Fprint(td.stdout, prompt)

	// Channel to receive user input
	responseChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Goroutine to read user input
	go func() {
		reader := bufio.NewReader(td.stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- strings.TrimSpace(response)
	}()

	// Wait for either response or timeout
	select {
	case response := <-responseChan:
		return td.parseConfirmationResponse(response, options), nil
	case err := <-errorChan:
		return options.DefaultValue, err
	case <-ctx.Done():
		fmt.Fprintf(td.stdout, "\nTimeout reached. Using default value: %t\n", options.DefaultValue)
		return options.DefaultValue, nil
	}
}

// parseConfirmationResponse parses the user's response
func (td *TTYDetector) parseConfirmationResponse(response string, options *ConfirmationOptions) bool {
	response = strings.ToLower(strings.TrimSpace(response))

	// Empty response uses default
	if response == "" {
		return options.DefaultValue
	}

	// Check yes responses
	for _, yes := range options.YesResponses {
		if response == strings.ToLower(yes) {
			return true
		}
	}

	// Check no responses
	for _, no := range options.NoResponses {
		if response == strings.ToLower(no) {
			return false
		}
	}

	// Invalid response, use default
	return options.DefaultValue
}

// ShowSpinner controls whether progress indicators should be shown
func (td *TTYDetector) ShowSpinner() bool {
	return td.IsInteractive()
}

// PromptForInput prompts for user input with optional default value
func (td *TTYDetector) PromptForInput(prompt, defaultValue string) (string, error) {
	if !td.IsInteractive() {
		return defaultValue, nil
	}

	// Display prompt
	if defaultValue != "" {
		fmt.Fprintf(td.stdout, "%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Fprintf(td.stdout, "%s: ", prompt)
	}

	// Read user input
	reader := bufio.NewReader(td.stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

// PromptForPassword prompts for password input (hidden)
func (td *TTYDetector) PromptForPassword(prompt string) (string, error) {
	if !td.IsInteractive() {
		return "", fmt.Errorf("password input requires interactive terminal")
	}

	fmt.Fprint(td.stdout, prompt)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	fmt.Fprintln(td.stdout) // Add newline after password input
	return string(password), nil
}

// ClearLine clears the current line in the terminal
func (td *TTYDetector) ClearLine() {
	if td.IsStdoutTTY() {
		fmt.Fprint(td.stdout, "\r\033[K")
	}
}

// MoveCursorUp moves the cursor up by the specified number of lines
func (td *TTYDetector) MoveCursorUp(lines int) {
	if td.IsStdoutTTY() && lines > 0 {
		fmt.Fprintf(td.stdout, "\033[%dA", lines)
	}
}

// MoveCursorDown moves the cursor down by the specified number of lines
func (td *TTYDetector) MoveCursorDown(lines int) {
	if td.IsStdoutTTY() && lines > 0 {
		fmt.Fprintf(td.stdout, "\033[%dB", lines)
	}
}

// HideCursor hides the terminal cursor
func (td *TTYDetector) HideCursor() {
	if td.IsStdoutTTY() {
		fmt.Fprint(td.stdout, "\033[?25l")
	}
}

// ShowCursor shows the terminal cursor
func (td *TTYDetector) ShowCursor() {
	if td.IsStdoutTTY() {
		fmt.Fprint(td.stdout, "\033[?25h")
	}
}

// WrapWithFallback wraps an interactive function with a non-interactive fallback
func (td *TTYDetector) WrapWithFallback(interactiveFn func() error, fallbackFn func() error) error {
	if td.IsInteractive() {
		return interactiveFn()
	}
	return fallbackFn()
}

// GetColorSupport checks if the terminal supports colors
func (td *TTYDetector) GetColorSupport() bool {
	if !td.IsStdoutTTY() {
		return false
	}

	// Check NO_COLOR environment variable (universal standard)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check FORCE_COLOR environment variable
	if forceColor := os.Getenv("FORCE_COLOR"); forceColor == "1" || forceColor == "true" {
		return true
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	// Common terminal types that support color
	colorTerms := []string{
		"xterm", "xterm-color", "xterm-256color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"linux", "cygwin",
	}

	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	// Check for color capability indicators
	if strings.Contains(term, "color") || strings.Contains(term, "256") {
		return true
	}

	return false
}

// Package-level convenience functions
var defaultTTYDetector = NewTTYDetector()

// IsInteractive is a package-level convenience function
func IsInteractive() bool {
	return defaultTTYDetector.IsInteractive()
}

// IsStdinTTY is a package-level convenience function
func IsStdinTTY() bool {
	return defaultTTYDetector.IsStdinTTY()
}

// IsStdoutTTY is a package-level convenience function
func IsStdoutTTY() bool {
	return defaultTTYDetector.IsStdoutTTY()
}

// GetTerminalSize is a package-level convenience function
func GetTerminalSize() (width, height int, err error) {
	return defaultTTYDetector.GetTerminalSize()
}

// ConfirmWithTimeout is a package-level convenience function
func ConfirmWithTimeout(options *ConfirmationOptions) (bool, error) {
	return defaultTTYDetector.ConfirmWithTimeout(options)
}

// ShowSpinner is a package-level convenience function
func ShowSpinner() bool {
	return defaultTTYDetector.ShowSpinner()
}

// GetColorSupport is a package-level convenience function
func GetColorSupport() bool {
	return defaultTTYDetector.GetColorSupport()
}
