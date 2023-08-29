package create

import (
	"context"
	_ "embed"
	"fmt"
	"text/template"
	"time"

	"github.com/JulienBreux/mksctl/internal/mksctl/api/client"
	api "github.com/JulienBreux/mksctl/internal/mksctl/api/gen"
	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	cmdUsage     = "create API_NAME:API_VERSION API_ENDPOINT RUNNER"
	cmdShortDesc = "Create a test in Microcks."

	minNArgs  = 3
	runnerArg = 2

	defaultEndpointTimeout = 5 * time.Second
)

var (
	argSecretName      string
	argEndpointTimeout time.Duration
	argOperations      []string
	argHeaders         string

	//go:embed output.tmpl
	outputTemplate string

	cmd = &cobra.Command{
		Use:   cmdUsage,
		Short: cmdShortDesc,
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(minNArgs),
			checkRunnerValue(runnerArg),
		),
		RunE: run,
	}
)

// New returns the test create sub command
func New() *cobra.Command {
	cmd.Flags().StringVar(
		&argSecretName,
		"secret-name",
		"",
		"secret to use for connecting test endpoint",
	)
	cmd.Flags().DurationVar(
		&argEndpointTimeout,
		"endpoint-timeout",
		defaultEndpointTimeout,
		"endpoint timeout, used by Microcks internal client for communicating",
	)
	cmd.Flags().StringArrayVar(
		&argOperations,
		"operation",
		[]string{},
		"list of operations to launch a test for (eg: GET /pastry)",
	)
	cmd.Flags().StringVar(
		&argHeaders,
		"headers",
		"",
		"override of operations headers",
	)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	// Connect to client
	cli, err := client.New(config.Config.APIURL)
	if err != nil {
		return err
	}

	// Create a new test
	ctx := context.Background()
	resp, err := cli.Actions().CreateTestWithResponse(
		ctx,
		testRequest(args[0], args[1], args[2]),
	)
	if err != nil {
		return err
	}

	return output(cmd, resp.JSON201)
}

func testRequest(serviceID, testEndpoint string, runnerType string) api.TestRequest {
	var operationHeaders api.OperationHeaders

	return api.TestRequest{
		FilteredOperations: &argOperations,
		OperationHeaders:   &operationHeaders,
		RunnerType:         api.TestRunnerType(runnerType),
		SecretName:         &argSecretName,
		ServiceId:          serviceID,
		TestEndpoint:       testEndpoint,
		Timeout:            int(argEndpointTimeout),
	}
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

func output(cmd *cobra.Command, resp *api.TestResult) error {
	funcs := template.FuncMap{
		"cyan":              color.CyanString,
		"green":             color.GreenString,
		"red":               color.RedString,
		"dateFromTimestamp": func(i *int) time.Time { return time.Unix(int64(*i), 0) },
	}
	tmpl, err := template.
		New("test").
		Funcs(funcs).
		Parse(outputTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(cmd.OutOrStdout(), resp)
}
