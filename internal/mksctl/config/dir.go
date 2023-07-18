package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func fullDirPath() (p string) {
	p = filepath.Join(userHomeDir(), Path)
	p = filepath.Clean(p)
	return
}

func defaultDirNotExists() bool {
	_, err := os.Stat(fullDirPath())
	return os.IsNotExist(err)
}

func createDefaultDir() {
	err := os.MkdirAll(fullDirPath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
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
