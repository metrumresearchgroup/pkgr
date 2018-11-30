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
	"github.com/metrumresearchgroup/pkgr/rcmd"

	"github.com/spf13/cobra"
)

// checkCmd represents the R CMD check command
var runCmd = &cobra.Command{
	Use:   "run R",
	Short: "Run R with the configuration settings used with other R commands",
	Long: `
	allows for interactive use and debugging based on the configuration specified by pkgr
 `,
	RunE: rRun,
}

func rRun(cmd *cobra.Command, args []string) error {

	rs := rcmd.NewRSettings()
	// installation through binary doesn't do this exactly, but its pretty close
	// at least for experimentation for now. If necessary can refactor out the
	// specifics so could be run here exactly.
	rs.LibPaths = append(rs.LibPaths, cfg.Library)
	rcmd.RunR(fs, rs, "", log)
	return nil
}

func init() {
	RootCmd.AddCommand(runCmd)
}
