package auth

import (
	cmdLogin "github.com/JulienBreux/mksctl/internal/mksctl/command/auth/login"
	"github.com/spf13/cobra"
)

const (
	cmdName      = "auth"
	cmdShortDesc = "Manage credentials for the Microcks CLI."
)

// New returns the auth command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDesc,
	}

	cmd.AddCommand(cmdLogin.New())

	return cmd
}
