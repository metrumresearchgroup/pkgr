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

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VERSION is the current pkgr version
var VERSION = "2.0.2"

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

	rootInit()
}

func rootInit() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().String("config", "", "config file (default is pkgr.yml)")
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	RootCmd.PersistentFlags().String("loglevel", "", "level for logging")
	_ = viper.BindPFlag("loglevel", RootCmd.PersistentFlags().Lookup("loglevel"))

	RootCmd.PersistentFlags().Int("threads", 0, "number of threads to execute with")
	_ = viper.BindPFlag("threads", RootCmd.PersistentFlags().Lookup("threads"))

	RootCmd.PersistentFlags().Bool("preview", false, "preview action, but don't actually run command")
	_ = viper.BindPFlag("preview", RootCmd.PersistentFlags().Lookup("preview"))

	RootCmd.PersistentFlags().Bool("debug", false, "use debug mode")
	_ = viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// globals
	RootCmd.PersistentFlags().String("library", "", "library to install packages")
	_ = viper.BindPFlag("library", RootCmd.PersistentFlags().Lookup("library"))

	RootCmd.PersistentFlags().Bool("update", cfg.Update, "Update packages along with install")
	_ = viper.BindPFlag("update", RootCmd.PersistentFlags().Lookup("update"))

	RootCmd.PersistentFlags().Bool("rollback", cfg.Rollback, "Enable rollback")
	_ = viper.BindPFlag("rollback", RootCmd.PersistentFlags().Lookup("rollback"))

	RootCmd.PersistentFlags().Bool("no-secure", cfg.Rollback, "disable TLS certificate verification")
	_ = viper.BindPFlag("nosecure", RootCmd.PersistentFlags().Lookup("no-secure"))

	RootCmd.PersistentFlags().Bool("strict", cfg.Strict, "Enable strict mode")
	_ = viper.BindPFlag("strict", RootCmd.PersistentFlags().Lookup("strict"))
}

func setGlobals() {
	fs = afero.NewOsFs()
	logger.SetLogLevel(viper.GetString("loglevel"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" { // enable ability to specify config file via flag
	// 	viper.SetConfigFile(cfgFile)
	// }

	setGlobals()

	if viper.GetBool("debug") {
		viper.Debug()
	}

	log.Trace("attempting to load config file")
	configlib.NewConfig(viper.GetString("config"), &cfg)

	configFilePath, _ := filepath.Abs(viper.ConfigFileUsed())
	cwd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"cwd": cwd,
		"nwd": filepath.Dir(configFilePath),
	}).Trace("setting directory to configuration file")
	_ = os.Chdir(filepath.Dir(configFilePath))

}
