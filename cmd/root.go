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
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/spf13/cobra"
)

// VERSION is the current pkc version
const VERSION string = "0.0.1-beta.3"

var log *logrus.Logger
var fs afero.Fs
var cfg configlib.PkgrConfig

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pkgr",
	Short: "package manager",
	Long:  fmt.Sprintf("pkgr cli version %s", VERSION),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

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

	log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Hooks:     make(logrus.LevelHooks),
		// Minimum level to log at (5 is most verbose (debug), 0 is panic)
	}
	switch logLevel := strings.ToLower(viper.GetString("loglevel")); logLevel {
	case "trace":
		log.Level = logrus.TraceLevel
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	default:
		log.Level = logrus.InfoLevel
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" { // enable ability to specify config file via flag
	// 	viper.SetConfigFile(cfgFile)
	// }
	if viper.GetString("config") == "" {
		_ = configlib.LoadGlobalConfig("pkgr")
	} else {
		_ = configlib.LoadConfigFromPath(viper.GetString("config"))
	}

	setGlobals()
	if viper.GetBool("debug") {
		viper.Debug()
	}
	viper.Unmarshal(&cfg)
	configDir, _ := filepath.Abs(viper.ConfigFileUsed())
	cwd, _ := os.Getwd()
	log.WithFields(logrus.Fields{
		"cwd": cwd,
		"nwd": filepath.Dir(configDir),
	}).Trace("setting directory to configuration file")
	os.Chdir(filepath.Dir(configDir))

	if cfg.Logging.File != "" {
		fileHook, err := logger.NewLogrusFileHook(cfg.Logging.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err == nil {
			log.Hooks.Add(fileHook)
		}
	}
}
