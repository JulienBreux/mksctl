package url

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

// ValidArg validates a URL in the arguments according to the position of the argument
func ValidArg(argURLPosition int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if argURLPosition > len(args) || len(args) == 0 {
			return fmt.Errorf("unable to valid URL, no enough arguments")
		}
		uri := args[argURLPosition]
		if _, err := url.ParseRequestURI(uri); err != nil {
			return fmt.Errorf("URL '%s' is not valid", uri)
		}
		return nil
	}
}
