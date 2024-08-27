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
	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/rcmd"

	"github.com/spf13/cobra"
)

var pkg string

// checkCmd represents the R CMD check command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch R session with config settings",
	Long: `Start an interactive R session based on the settings defined in the
configuration file.

   * Use the R executable defined by the 'RPath' value, if any.

   * Set the library paths so that packages come from only the
     configuration's library and the library bundled with the R
     installation.

   * If the --pkg option is passed, set the environment variables defined in
     the package's 'Customizations' entry.`,
	Example: `  # Launch an R session, setting values based on pkgr.yml
  pkgr run
  # Also setting environment variables specified for dplyr:
  #
  #   Customizations:
  #     Packages:
  #        - dplyr:
  #            Env:
  #              [...]
  pkgr run --pkg=dplyr`,
	RunE: rRun,
}

func rRun(cmd *cobra.Command, args []string) error {

	rs := rcmd.NewRSettings(cfg.RPath)
	// installation through binary doesn't do this exactly, but its pretty close
	// at least for experimentation for now. If necessary can refactor out the
	// specifics so could be run here exactly.
	rs.LibPaths = append(rs.LibPaths, cfg.Library)
	pc, exists := configlib.GetPackageCustomizationByName(pkg, cfg.Customizations)
	if exists {
		rs.PkgEnvVars[pkg] = pc.Env
	}
	rcmd.StartR(fs, pkg, rs, "")
	return nil
}

func init() {
	runCmd.Flags().StringVar(&pkg, "pkg", "", "package environment to set")
	RootCmd.AddCommand(runCmd)
}
