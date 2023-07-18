package login

import (
	"fmt"

	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	"github.com/JulienBreux/mksctl/pkg/validator/url"
	"github.com/spf13/cobra"
)

const (
	cmdName      = "login URL"
	cmdShortDesc = "Authorize Microcks CLI to access the Mickrocks server using user credentials."
)

var (
	cmd = &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDesc,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			url.ValidArg(0),
		),
		PreRunE: preRun,
		RunE:    run,
	}
)

// New returns the login command
func New() *cobra.Command {
	return cmd
}

func preRun(_ *cobra.Command, args []string) error {
	return config.WriteKey(config.URLField, args[0])
}

func run(_ *cobra.Command, _ []string) error {
	fmt.Printf("Server URL: %s\n", config.GetKey(config.URLField))
	return nil
}
