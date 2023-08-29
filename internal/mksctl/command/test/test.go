package test

import (
	cmdCreate "github.com/JulienBreux/mksctl/internal/mksctl/command/test/create"
	cmdGet "github.com/JulienBreux/mksctl/internal/mksctl/command/test/get"
	"github.com/spf13/cobra"
)

const (
	cmdName      = "test API_NAME:API_VERSION API_ENDPOINT RUNNER"
	cmdShortDesc = "Run test in Microcks."
)

var (
	cmd = &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDesc,
	}
)

// New returns a command to test
func New() *cobra.Command {
	cmd.AddCommand(cmdCreate.New())
	cmd.AddCommand(cmdGet.New())
	return cmd
}
