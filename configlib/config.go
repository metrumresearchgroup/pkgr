package configlib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"runtime"

	"github.com/spf13/viper"
)

// LoadGlobalConfig loads pkgcheck configuration into the global Viper
func LoadGlobalConfig(configFilename string) error {
	viper.SetConfigName(configFilename)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("pkgr")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			// found config file but couldn't parse it, should error
			panic(fmt.Errorf("unable to parse config file with error (%s)", err))
		}
		fmt.Println("no config file detected, using default settings")
	}

	loadDefaultSettings()
	return nil
}

// LoadConfigFromPath loads pkc configuration into the global Viper
func LoadConfigFromPath(configFilename string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("pkgr")
	configFilename = expand(filepath.Clean(configFilename))
	viper.SetConfigFile(configFilename)
	err := viper.ReadInConfig()
	if err != nil {
		// panic if can't find or parse config as this could be explicit to user expectations
		if _, ok := err.(*os.PathError); ok {
			panic(fmt.Errorf("could not find a config file at path: %s", configFilename))
		}
		if _, ok := err.(viper.ConfigParseError); ok {
			// found config file but couldn't parse it, should error
			panic(fmt.Errorf("unable to parse config file with error (%s)", err))
		}
	}

	loadDefaultSettings()
	return nil
}

func loadDefaultSettings() {
	// the lib paths to use, colon separated list of paths
	viper.SetDefault("library", "")
	viper.SetDefault("debug", false)
	viper.SetDefault("preview", false)
	// should be one of Debug,Info,Warn,Error,Fatal,Panic
	viper.SetDefault("loglevel", "info")
	// path to R on system, defaults to R in path
	viper.SetDefault("rpath", "R")
	viper.SetDefault("threads", runtime.NumCPU())

	// packrat related
	viper.SetDefault("pr_lockfile", "")
	viper.SetDefault("pr_dir", "")

}

func expand(s string) string {
	if strings.HasPrefix(s, "~/") {
		return filepath.Join(os.Getenv("HOME"), s[1:])
	}

	return s
}
