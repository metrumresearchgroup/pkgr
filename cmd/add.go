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
	"os"
	"time"

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// installCmd represents the R CMD install command
var addCmd = &cobra.Command{
	Use:   "add [flags] <package> [<package>...]",
	Short: "Add packages to the configuration file",
	Long: `Add the specified packages to the 'Packages' section of the
configuration file.`,
	Example: `  # Add mrgsolve and bbr to list of packages
  pkgr add mrgsolve bbr
  # Add rlang and then do installation
  # (same result as following up with 'pkgr install' call)
  pkgr add --install rlang`,
	RunE: rAdd,
}

var install bool

func init() {
	addCmd.Flags().BoolVar(&install, "install", false, "run install after updating config")
	RootCmd.AddCommand(addCmd)
}

func rAdd(ccmd *cobra.Command, args []string) error {

	// Initialize log and start time.
	initAddLog()
	startTime := time.Now()

	if len(args) == 0 {
		ccmd.Help()
		os.Exit(0)
	}

    err := configlib.AddPackages(args)
	if err != nil {
		log.Fatal(err)
	}
	if install {
		// if installing now, must call initConfig again for cobra to read in the yml file changes and see the new package/s
		initConfig()
		rInstall(nil, nil)
	}

	log.Info("duration:", time.Since(startTime))
	return nil
}

func initAddLog() {
	logger.AddLogFile(cfg.Logging.All, cfg.Logging.Overwrite)
}
