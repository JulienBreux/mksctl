package command

import (
	"fmt"
	"io"
	"os"

	cmdAuth "github.com/JulienBreux/mksctl/internal/mksctl/command/auth"
	cmdRoot "github.com/JulienBreux/mksctl/internal/mksctl/command/root"
	cmdVersion "github.com/JulienBreux/mksctl/internal/mksctl/command/version"
	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	"github.com/spf13/cobra"
)

const (
	appName      = "mksctl"
	appShortDesc = "Microcks server CLI."
	appLongDesc  = "CLI for interacting with Microcks server."
	appIssueURL  = "https://github.com/JulienBreux/mksctl/issues/new?labels=bug"
)

var (
	cfgPathFile string

	cmd = &cobra.Command{
		Use:   appName,
		Short: appShortDesc,
		Long:  appLongDesc,
		RunE:  cmdRoot.Run,
	}
)

type IOs struct {
	In       io.Reader
	Out, Err io.Writer
}

// Execute executes command
func New(ios *IOs, args ...string) *cobra.Command {
	defer recoverAndExit()

	// Cobra initialization
	cobra.OnInitialize(initConfig)

	// Setters
	cmd.SetIn(ios.In)
	cmd.SetOut(ios.Out)
	cmd.SetErr(ios.Err)
	cmd.SetArgs(args)

	// Add flags
	cmd.Flags().StringVarP(&cfgPathFile, "config", "c", config.FullFilePath(), "configuration file")
	cmd.PersistentFlags().StringP("output", "o", "", "Output format, one of 'yaml', 'json', 'toml' or 'xml'.")

	// Add subcommands
	cmd.AddCommand(cmdVersion.New())
	cmd.AddCommand(cmdAuth.New())

	return cmd
}

func initConfig() {
	config.Init(cfgPathFile)
	_ = config.Read()
}

func recoverAndExit() {
	if r := recover(); r != nil {
		// TODO: Improve error message color
		fmt.Println("Internal " + appName + " error")
		// TODO: Add logger at debug level
		// TODO: Add "tips" option
		// TODO: Get URL from outside
		fmt.Println("➡ Please report here: " + appIssueURL)
		os.Exit(1)
	}
}
