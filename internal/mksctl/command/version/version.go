package version

import (
	"github.com/spf13/cobra"

	ver "github.com/JulienBreux/mksctl/internal/mksctl/version"
)

const (
	cmdName      = "version"
	cmdShortDesc = "Show the Microcks client and server version information."
)

var (
	clientOnly = false
)

// New returns a command to print version
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDesc,
		RunE:  run,
	}

	cmd.Flags().BoolVar(&clientOnly, "client", false, "If true, shows client version only (no server required).")

	return cmd
}

// run returns the command
func run(cmd *cobra.Command, _ []string) error {
	o, err := cmd.Parent().PersistentFlags().GetString("output")
	if err != nil {
		return err
	}
	ver.Print(cmd.OutOrStdout(), clientOnly, o)

	return nil
}
