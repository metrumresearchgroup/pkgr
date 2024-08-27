// Copyright © 2019 John Carlo Salter <juuncaerlum@gmail.com>
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
	"github.com/metrumresearchgroup/pkgr/logger"
	"os"

	"github.com/spf13/cobra"
)

var cleanAll bool

// CleanCmd represents the clean command
var CleanCmd = &cobra.Command{
	Use:   `clean [flags]`,
	Short: "Clean cached information",
	Long: `This subcommand is an entry point for cleaning two categories of cached
data:

 * source and binary tarballs

   Use the 'cache' subcommand to remove these.

 * package databases with information about the packages available from
   repositories

   Use the 'pkgdbs' subcommand to remove these.

To remove cached data for both categories, pass the --all flag.`,
	RunE: clean,
}

func init() {
	CleanCmd.Flags().BoolVar(&cleanAll, "all", false, "clean all cached items")

	RootCmd.AddCommand(CleanCmd)
}

func clean(cmd *cobra.Command, args []string) error {

	logger.AddLogFile(cfg.Logging.All, cfg.Logging.Overwrite)
	// set the cache control to 0 to make sure will clear out existing caches to re-initialize to get fresh results
	// before purging
	os.Setenv("R_AVAILABLE_PACKAGES_CACHE_CONTROL_MAX_AGE", "0")
	var err error

	if !cleanAll {
		fmt.Println("No clean options passed -- not cleaning.")
	}
	if cleanAll {
		fmt.Println("Cleaning all.")
		err = cleanCacheFolders()
		if err != nil {
			return err
		}
		err = cleanPackageDatabases("ALL")
		if err != nil {
			return err
		}
	}
	return nil
}
