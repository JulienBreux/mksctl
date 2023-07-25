package config

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	Name      = "config"
	Type      = "yaml"
	Path      = ".mksctl/"
	EnvPrefix = "MICROCKS"
)

var (
	Config *config

	OverrideConfigFile string
)

type config struct {
	Initialized bool

	APIURL     string        `mapstructure:"api_url"`
	APITimeout time.Duration `mapstructure:"api_timeout" default:"5s"`

	AuthEnabled bool `mapstructure:"auth_enabled" default:"false"`

	AuthServerURL string `mapstructure:"auth_server_url"`

	AuthCallbackPort int `mapstructure:"auth_callback_port" default:"8283"`

	AuthClientRealm  string `mapstructure:"auth_client_id" default:"mickrocks"`
	AuthClientID     string `mapstructure:"auth_client_id" default:"mksctl"`
	AuthClientSecret string `mapstructure:"auth_client_secret"` // TODO: Move to github.com/zalando/go-keyring

	AuthAccessToken  string `mapstructure:"auth_access_token"`  // TODO: Move to github.com/zalando/go-keyring
	AuthRefreshToken string `mapstructure:"auth_refresh_token"` // TODO: Move to github.com/zalando/go-keyring
}

// Init helps to initialize configuration system based on Viper
func Init() error {
	// Set Viper behaviour
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetConfigName(Name)
	viper.SetConfigType(Type)
	viper.AddConfigPath(fullDirPath())

	// Use specific configuration file
	if OverrideConfigFile != "" {
		viper.SetConfigFile(OverrideConfigFile)
	}

	// Create configuration file if not exists
	if err := initConfigFile(); err != nil {
		return err
	}

	// Load configuration values
	return loadConfiguration()
}

// Save helps to save the configuration
func Save() error {
	cfg := make(map[string]any)
	if err := mapstructure.Decode(Config, &cfg); err != nil {
		return err
	}

	if err := viper.MergeConfigMap(cfg); err != nil {
		return err
	}

	return viper.WriteConfig()
}

// FullFilePath returns the full path of config file
func FullFilePath() string {
	return filepath.Join(fullDirPath(), strings.Join([]string{Name, Type}, "."))
}

func loadConfiguration() error {
	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Create a default configuration
	Config = new(config)
	defaults.SetDefaults(Config)

	// Unmarshal values in configuration
	if err := viper.Unmarshal(&Config); err != nil {
		return err
	}

	Config.Initialized = true

	return nil
}
