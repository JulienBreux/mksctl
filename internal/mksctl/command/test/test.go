package test

import (
	"context"
	"fmt"
	"time"

	"github.com/JulienBreux/mksctl/internal/mksctl/api/client"
	api "github.com/JulienBreux/mksctl/internal/mksctl/api/gen"
	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	"github.com/spf13/cobra"
)

const (
	cmdName      = "test API_NAME:API_VERSION API_ENDPOINT RUNNER API_ENDPOINT_TIMEOUT"
	cmdShortDesc = "Run test in Microcks."
)

const (
	defaultEndpointTimeout = 5 * time.Second
	defaultWaitTimeout     = 5 * time.Second

	minNArgs  = 3
	runnerArg = 2
)

var (
	cmd = &cobra.Command{
		Use:     cmdName,
		Short:   cmdShortDesc,
		Aliases: []string{"t"},
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(minNArgs),
			checkRunnerValue(runnerArg),
		),
		PreRunE: preRun,
		RunE:    run,
	}

	secretName string

	endpointTimeout time.Duration
	waitTimeout     time.Duration

	operations []string
	headers    string
)

// New returns a command to import
func New() *cobra.Command {
	// cmd.Flags().String
	cmd.Flags().StringVar(&secretName, "secret-name", "", "secret to use for connecting test endpoint")
	cmd.Flags().DurationVar(&waitTimeout, "wait", waitTimeout, "time to wait for test to finish")
	cmd.Flags().DurationVar(
		&endpointTimeout,
		"endpoint-timeout",
		defaultEndpointTimeout,
		"endpoint timeout, used by Microcks internal client for communicating",
	)
	cmd.Flags().StringArrayVar(
		&operations,
		"operation",
		[]string{},
		"list of operations to launch a test for",
	)
	cmd.Flags().StringVar(
		&headers,
		"headers",
		"",
		"override of operations headers",
	)
	return cmd
}

// preRun helps to check arguments
func preRun(_ *cobra.Command, _ []string) error {
	return nil
}

// run returns the command
func run(cmd *cobra.Command, args []string) error {
	// Connect to client
	cli, err := client.New(config.Config.APIURL)
	if err != nil {
		return err
	}

	// Decode
	operationHeaders, err := headersDecode()
	if err != nil {
		return err
	}

	// Create a new test
	resp, err := cli.Actions().CreateTestWithResponse(
		context.Background(),
		api.TestRequest{
			FilteredOperations: &operations,
			OperationHeaders:   operationHeaders,
			RunnerType:         api.TestRunnerType(args[2]),
			SecretName:         &secretName,
			ServiceId:          args[0],
			TestEndpoint:       args[1],
			Timeout:            int(endpointTimeout),
		},
	)
	if err != nil {
		return err
	}

	cmd.Print(string(resp.Body))

	return nil
}

func checkRunnerValue(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, allowedRunnerType := range client.RunnerTypes {
			if allowedRunnerType == api.TestRunnerType(args[n]) {
				return nil
			}
		}

		return fmt.Errorf(
			"unable to recognize RUNNER argument with \"%s\" value",
			args[n],
		)
	}
}

func headersDecode() (*api.OperationHeaders, error) {
	return nil, nil
}
