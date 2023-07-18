package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	Name      = "config"
	Type      = "yaml"
	Path      = ".config/mksctl/"
	EnvPrefix = "MICROCKS"
)

// Init helps to initialize configuration system based on Viper
func Init(cfgPathFile string) {
	if defaultDirNotExists() {
		createDefaultDir()
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetConfigName(Name)
	viper.SetConfigType(Type)
	viper.AddConfigPath(fullDirPath())
	if cfgPathFile != "" {
		viper.SetConfigFile(cfgPathFile)
	}
}

// Read helps to read the configuration file
func Read() error {
	return viper.ReadInConfig()
}

// Write helps to write the configuration file
func Write() error {
	return viper.WriteConfig()
}

// SetKey helps to set a key
func SetKey(key string, value any) {
	viper.Set(key, value)
}

// WriteKey helps to write a key directly to the configuration file
func WriteKey(key string, value any) error {
	SetKey(key, value)
	return Write()
}

// GetKey helps to get a key
func GetKey(key string) any {
	return viper.Get(key)
}

// FullFilePath returns the full path of config file
func FullFilePath() string {
	return filepath.Join(fullDirPath(), strings.Join([]string{Name, Type}, "."))
}
