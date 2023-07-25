package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func fullDirPath() string {
	p := filepath.Join(userHomeDir(), Path)
	p = filepath.Clean(p)
	return p
}

func fullFilePath() string {
	p := filepath.Join(fullDirPath(), Name+"."+Type)
	return p
}

func defaultDirExists() bool {
	_, err := os.Stat(fullDirPath())
	return os.IsExist(err)
}

func createDefaultDir() error {
	if defaultDirExists() {
		return nil
	}

	err := os.MkdirAll(fullDirPath(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func createDefaultFile() error {
	_, err := os.Stat(fullFilePath())
	if os.IsNotExist(err) {
		file, err := os.Create(fullFilePath())
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func initConfigFile() error {
	if err := createDefaultDir(); err != nil {
		return err
	}
	return createDefaultFile()
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
