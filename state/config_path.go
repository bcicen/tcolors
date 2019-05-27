package state

import (
	"fmt"
	"os"
	"regexp"
)

var DefaultPalettePath = defaultPalettePath()

func defaultPalettePath() string {
	path, err := getConfigPath()
	if err != nil {
		panic(err)
	}
	if err := ensureDir(path); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/default.toml", path)
}

// attempt create dir if not exist
func ensureDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0755)
		}
		return err
	}
	return nil
}

// return config base dir from environment
func getConfigPath() (string, error) {
	userHome, ok := os.LookupEnv("HOME")
	if !ok {
		return "", fmt.Errorf("$HOME not set")
	}

	if !xdgSupport() {
		// default path
		return fmt.Sprintf("%s/.tcolors", userHome), nil
	}

	xdgHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		// use default xdg config home
		xdgHome = fmt.Sprintf("%s/.config", userHome)
	}
	return fmt.Sprintf("%s/tcolors", xdgHome), nil
}

// Test for environemnt supporting XDG spec
func xdgSupport() bool {
	re := regexp.MustCompile("^XDG_*")
	for _, e := range os.Environ() {
		if re.FindAllString(e, 1) != nil {
			return true
		}
	}
	return false
}
