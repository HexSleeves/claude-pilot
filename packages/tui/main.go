// Package main provides a terminal user interface (TUI) for the Claude Pilot
// session management system. This TUI allows users to view, create, attach to,
// and manage development sessions through an interactive terminal interface
// built with the Bubbletea framework.
package main

import (
	"fmt"
	"os"

	"claude-pilot/core/api"

	tea "github.com/charmbracelet/bubbletea"
)

// main initializes the TUI application and starts the Bubbletea program.
// It creates an API client, initializes the TUI model, and runs the interactive
// terminal interface with alternate screen buffer and mouse support enabled.
func main() {
	// Initialize the core API client
	client, err := api.NewDefaultClient(false) // verbose = false for TUI
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing client: %v\n", err)
		os.Exit(1)
	}

	// Create the main TUI model
	model := NewModel(client)

	// Create the Bubbletea program with proper cleanup handling
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program and ensure proper cleanup on exit
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		// Ensure terminal is restored before exit
		program.Kill()
		os.Exit(1)
	}
}
