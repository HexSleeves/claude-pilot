// Package main provides the standalone TUI binary for Claude Pilot
package main

import (
	"fmt"
	"os"

	"claude-pilot/core/api"
	"claude-pilot/tui"
)

func main() {
	// Initialize the core API client
	client, err := api.NewDefaultClient(false) // verbose = false for TUI
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing client: %v\n", err)
		os.Exit(1)
	}

	// Run the TUI
	if err := tui.RunTui(client); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
