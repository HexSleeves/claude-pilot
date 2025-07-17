package main

import (
	"context"
	"os"

	"claude-pilot/cmd"

	"github.com/charmbracelet/fang"
)

func main() {
	// This is where the magic happens.
	if err := fang.Execute(
		context.Background(),
		cmd.RootCmd(),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	); err != nil {
		os.Exit(1)
	}
}
