package main

import (
	"context"
	"os"

	"claude-pilot/cmd"

	"github.com/charmbracelet/fang"
)

func main() {
	// if err := cmd.Execute(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	os.Exit(1)
	// }

	// This is where the magic happens.
	if err := fang.Execute(
		context.Background(),
		cmd.RootCmd(),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		os.Exit(1)
	}
}
