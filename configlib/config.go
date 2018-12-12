package configlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// LoadConfigFromPath loads pkc configuration into the global Viper
func LoadConfigFromPath(configFilename string) error {
	if configFilename == "" {
		configFilename = "pkgr.yml"
	}
	viper.SetEnvPrefix("pkgr")
	viper.AutomaticEnv()
	configFilename, _ = homedir.Expand(filepath.Clean(configFilename))
	viper.SetConfigFile(configFilename)
	b, err := ioutil.ReadFile(configFilename)
	expb := []byte(os.ExpandEnv(string(b)))
	err = viper.ReadConfig(bytes.NewReader(expb))
	if err != nil {
		// panic if can't find or parse config as this could be explicit to user expectations
		if _, ok := err.(*os.PathError); ok {
			panic(fmt.Errorf("could not find a config file at path: %s", configFilename))
		}
		if _, ok := err.(viper.ConfigParseError); ok {
			// found config file but couldn't parse it, should error
			panic(fmt.Errorf("unable to parse config file with error (%s)", err))
		}
		// maybe could be more loose on this later, but for now will require a config file
		fmt.Println("Error with pkgr config file:")
		fmt.Println(err)
		os.Exit(1)
	}

	loadDefaultSettings()
	return nil
}

func loadDefaultSettings() {
	viper.SetDefault("debug", false)
	viper.SetDefault("preview", false)
	// should be one of Debug,Info,Warn,Error,Fatal,Panic
	viper.SetDefault("loglevel", "info")
	// path to R on system, defaults to R in path
	viper.SetDefault("rpath", "R")
	viper.SetDefault("threads", runtime.NumCPU())
}

func IsCustomizationSet(key string, elems []interface{}, elem string) bool {
	for _, v := range elems {
		for k, iv := range v.(map[interface{}]interface{}) {
			if k == elem {
				for k2 := range iv.(map[interface{}]interface{}) {
					if k2 == key {
						return true
					}
				}
			}
		}
	}
	return false
}
