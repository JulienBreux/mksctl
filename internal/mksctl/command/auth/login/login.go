package login

import (
	"fmt"
	"os"

	"github.com/JulienBreux/mksctl/internal/mksctl/api/client"
	"github.com/JulienBreux/mksctl/internal/mksctl/auth/server"
	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	urlValidator "github.com/JulienBreux/mksctl/pkg/validator/url"
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
			urlValidator.ValidArg(0),
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
	// Connect to client
	cli, err := client.New(args[0])
	if err != nil {
		return err
	}

	// Retrieve configuration
	authConfig, err := cli.AuthConfig()
	if err != nil {
		return err
	}

	// Save configuation
	config.Config.AuthEnabled = authConfig.Enabled
	config.Config.APIURL = cli.APIURL()
	config.Config.AuthClientRealm = authConfig.Realm
	config.Config.AuthServerURL = authConfig.AuthServerUrl
	return config.Save()
}

func run(_ *cobra.Command, _ []string) error {
	if !config.Config.AuthEnabled {
		fmt.Printf("Authentication is disabled, you're aleardy connected!\n")
		os.Exit(0)
	}
	return server.New().Run()
}
