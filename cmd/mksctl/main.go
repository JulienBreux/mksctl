package main

import (
	"os"

	"github.com/JulienBreux/mksctl/internal/mksctl/command"
)

func main() {
	// ios :=
	cmd := command.New(
		&command.IOs{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
	)
	if err := cmd.Execute(); err != nil {
		command.PrintError(os.Stderr, err)
		os.Exit(1)
	}
}
