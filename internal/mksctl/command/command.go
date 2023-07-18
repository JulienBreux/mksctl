package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	cmdAuth "github.com/JulienBreux/mksctl/internal/mksctl/command/auth"
	cmdRoot "github.com/JulienBreux/mksctl/internal/mksctl/command/root"
	cmdVersion "github.com/JulienBreux/mksctl/internal/mksctl/command/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName      = "mksctl"
	appShortDesc = "Microcks server CLI."
	appLongDesc  = "CLI for interacting with Microcks server."
	appIssueURL  = "https://github.com/JulienBreux/mksctl/issues/new?labels=bug"
)

var (
	cfgPathFile string
	cfgSubPath  = ".config/mksctl/"
	cfgType     = "yml"
	cfgFile     = "mksctl.yml"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgPathFile != "" {
		viper.SetConfigFile(cfgPathFile)
	} else {
		viper.AddConfigPath(cfgPath())
		viper.SetConfigType(cfgType)
		viper.SetConfigName(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

type IOs struct {
	In       io.Reader
	Out, Err io.Writer
}

// Execute executes command
func New(ios *IOs, args ...string) *cobra.Command {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Improve error message color
			fmt.Println("Internal " + appName + " error")
			// TODO: Add logger at debug level
			// TODO: Add "tips" option
			// TODO: Get URL from outside
			fmt.Println("➡ Please report here: " + appIssueURL)
			os.Exit(1)
		}
	}()

	// Cobra initialization
	cobra.OnInitialize(initConfig)

	// Create root command
	cmd := &cobra.Command{
		Use:   appName,
		Short: appShortDesc,
		Long:  appLongDesc,
		RunE:  cmdRoot.Run,
	}

	cmd.SetIn(ios.In)
	cmd.SetOut(ios.Out)
	cmd.SetErr(ios.Err)
	cmd.SetArgs(args)

	// Add flags
	flags(cmd)

	// Add subcommands
	cmd.AddCommand(cmdVersion.New())
	cmd.AddCommand(cmdAuth.New())

	return cmd
}

func flags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cfgPathFile, "config", "c", filepath.Join(cfgPath(), cfgFile), "configuration file")
	cmd.PersistentFlags().StringP("output", "o", "", "Output format, one of 'yaml', 'json', 'toml' or 'xml'.")
}

func cfgPath() string {
	// Home directory
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(homeDir, cfgSubPath)
}
