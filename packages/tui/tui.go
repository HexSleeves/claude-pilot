// Package tui provides a terminal user interface (TUI) for the Claude Pilot
// session management system. This TUI allows users to view, create, attach to,
// and manage development sessions through an interactive terminal interface
// built with the Bubbletea framework.
package tui

import (
	"fmt"

	"claude-pilot/core/api"

	tea "github.com/charmbracelet/bubbletea"
)

// RunTui initializes and runs the TUI application with the provided API client.
// It creates the TUI model, sets up the Bubbletea program, and runs the interactive
// terminal interface with alternate screen buffer and mouse support enabled.
func RunTui(client *api.Client) error {
	if client == nil {
		return fmt.Errorf("API client cannot be nil")
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
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
