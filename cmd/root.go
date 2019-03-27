// Copyright © 2018 Devin Pastoor <devin.pastoor@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	. "github.com/metrumresearchgroup/pkgr/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/spf13/cobra"
)

// VERSION is the current pkc version
var VERSION = "0.2.0-alpha.2"

var fs afero.Fs
var cfg configlib.PkgrConfig

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pkgr",
	Short: "package manager",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(build string) {
	if build != "" {
		VERSION = fmt.Sprintf("%s-%s", VERSION, build)
	}
	RootCmd.Long = fmt.Sprintf("pkgr cli version %s", VERSION)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().String("config", "", "config file (default is pkgr.yml)")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	RootCmd.PersistentFlags().String("loglevel", "", "level for logging")
	viper.BindPFlag("loglevel", RootCmd.PersistentFlags().Lookup("loglevel"))

	RootCmd.PersistentFlags().Int("threads", 0, "number of threads to execute with")
	viper.BindPFlag("threads", RootCmd.PersistentFlags().Lookup("threads"))

	RootCmd.PersistentFlags().Bool("preview", false, "preview action, but don't actually run command")
	viper.BindPFlag("preview", RootCmd.PersistentFlags().Lookup("preview"))

	RootCmd.PersistentFlags().Bool("debug", false, "use debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// globals
	RootCmd.PersistentFlags().String("library", "", "library to install packages")
	viper.BindPFlag("library", RootCmd.PersistentFlags().Lookup("library"))

	// packrat related
	// RootCmd.PersistentFlags().String("pr_lockfile", "", "packrat lockfile")
	// viper.BindPFlag("pr_lockfile", RootCmd.PersistentFlags().Lookup("pr_lockfile"))
	// RootCmd.PersistentFlags().String("pr_dir", "", "packrat dir")
	// viper.BindPFlag("pr_dir", RootCmd.PersistentFlags().Lookup("pr_dir"))
}

func setGlobals() {

	fs = afero.NewOsFs()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	switch logLevel := strings.ToLower(viper.GetString("loglevel")); logLevel {
	case "trace":
		Log.SetLevel(logrus.TraceLevel)
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Log.SetLevel(logrus.FatalLevel)
	case "panic":
		Log.SetLevel(logrus.PanicLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" { // enable ability to specify config file via flag
	// 	viper.SetConfigFile(cfgFile)
	// }
	configlib.LoadConfigFromPath(viper.GetString("config"))

	setGlobals()
	if viper.GetBool("debug") {
		viper.Debug()
	}
	viper.Unmarshal(&cfg)
	configFilePath, _ := filepath.Abs(viper.ConfigFileUsed())
	cwd, _ := os.Getwd()
	Log.WithFields(logrus.Fields{
		"cwd": cwd,
		"nwd": filepath.Dir(configFilePath),
	}).Trace("setting directory to configuration file")
	os.Chdir(filepath.Dir(configFilePath))

	if cfg.Logging.All != "" {
		fileHook, err := logger.NewLogrusFileHook(cfg.Logging.All, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err == nil {
			Log.AddHook(fileHook)
		}
	}
}
