package login

import (
	"fmt"

	"github.com/JulienBreux/mksctl/pkg/validator/url"
	"github.com/spf13/cobra"
)

const (
	cmdName      = "login server_url"
	cmdShortDesc = "Authorize Microcks CLI to access the Mickrocks server using user credentials."
)

// New returns the login command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDesc,
		Args:  cobra.MatchAll(cobra.ExactArgs(1), url.ValidArg(0)),
		RunE:  run,
	}

	return cmd
}

func run(_ *cobra.Command, args []string) error {
	serverURL := args[0]
	fmt.Printf("Server URL: %s\n", serverURL)
	return nil
}
