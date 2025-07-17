package main

import (
	"fmt"
	"os"

	"claude-pilot/core/api"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the core API client
	client, err := api.NewDefaultClient(false) // verbose = false for TUI
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing client: %v\n", err)
		os.Exit(1)
	}

	// Create the main TUI model
	model := NewModel(client)

	// Create the Bubbletea program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
