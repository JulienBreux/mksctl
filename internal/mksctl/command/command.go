package command

import (
	"io"
	"os"

	cmdAuth "github.com/JulienBreux/mksctl/internal/mksctl/command/auth"
	cmdImp "github.com/JulienBreux/mksctl/internal/mksctl/command/imp"
	cmdRoot "github.com/JulienBreux/mksctl/internal/mksctl/command/root"
	cmdTest "github.com/JulienBreux/mksctl/internal/mksctl/command/test"
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

	cobra.OnInitialize(initConfig)

	cmd = initSetters(cmd, ios, args...)
	cmd = initFlags(cmd)
	cmd = initSubCommands(cmd)

	return cmd
}

func initSetters(cmd *cobra.Command, ios *IOs, args ...string) *cobra.Command {
	cmd.SetIn(ios.In)
	cmd.SetOut(ios.Out)
	cmd.SetErr(ios.Err)
	cmd.SetArgs(args)
	return cmd
}

func initFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&config.OverrideConfigFile, "config", "c", config.FullFilePath(), "configuration file")
	cmd.PersistentFlags().StringP("output", "o", "", "Output format, one of 'yaml', 'json', 'toml' or 'xml'.")
	return cmd
}

func initSubCommands(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(cmdVersion.New())
	cmd.AddCommand(cmdAuth.New())
	cmd.AddCommand(cmdImp.New())
	cmd.AddCommand(cmdTest.New())
	return cmd
}

func initConfig() {
	if err := config.Init(); err != nil {
		cmd.Printf("Error:\n%v", err)
		os.Exit(1)
	}
}

func recoverAndExit() {
	if r := recover(); r != nil {
		cmd.Println("Internal " + appName + " error")
		cmd.Println("➡ Please report here: " + appIssueURL)
		os.Exit(1)
	}
}
