package get

import "github.com/spf13/cobra"

const (
	cmdUsage     = "get TEST_ID"
	cmdShortDesc = "Get results of a Microcks test."
)

var (
	cmd = &cobra.Command{
		Use:   cmdUsage,
		Short: cmdShortDesc,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
		),
		// PreRunE: preRun,
		RunE: run,
	}
)

// New returns the test get sub command
func New() *cobra.Command {
	return cmd
}

func run(_ *cobra.Command, _ []string) error {
	return nil
}
