package auth

import (
	cmdLogin "github.com/JulienBreux/mksctl/internal/mksctl/command/auth/login"
	"github.com/spf13/cobra"
)

const (
	cmdUse       = "auth"
	cmdShortDesc = "Manage credentials for the Microcks CLI."
)

var (
	cmd = &cobra.Command{
		Use:   cmdUse,
		Short: cmdShortDesc,
	}
)

// New returns the auth command
func New() *cobra.Command {
	cmd.AddCommand(cmdLogin.New())
	return cmd
}
