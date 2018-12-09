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
	"encoding/json"
	"fmt"
	"time"

	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/metrumresearchgroup/pkgr/rcmd/rp"
	"github.com/spf13/cobra"
)

// checkCmd represents the R CMD check command
var experimentCmd = &cobra.Command{
	Use:   "experiment",
	Short: "experiment with the cli",
	Long: `
	JUST FOR EXPERIMENTATION
 `,
	RunE: rExperiment,
}

func rExperiment(cmd *cobra.Command, args []string) error {

	//AppFs := afero.NewOsFs()
	// can use this to redirect log output
	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	rs := rcmd.NewRSettings()
	// installation through binary doesn't do this exactly, but its pretty close
	// at least for experimentation for now. If necessary can refactor out the
	// specifics so could be run here exactly.
	rs.LibPaths = append(rs.LibPaths, cfg.Library)
	res, _ := rcmd.RunR(fs, rs, ".libPaths()", "", log)
	fmt.Println(rp.ScanLines(res))
	startTime := time.Now()
	res, _ = rcmd.RunR(fs, rs, "paste0(R.Version()$major,'.',R.Version()$minor)", "", log)
	fmt.Println(rp.ScanLines(res)[0])
	fmt.Println(time.Since(startTime))
	res, err := rcmd.RunR(fs, rs, "stop('bad')", "", log)
	fmt.Println(res)
	fmt.Println("err: ", err)
	fmt.Println(rp.ScanLines(res))
	fmt.Println(time.Since(startTime))
	return nil
}

func init() {
	RootCmd.AddCommand(experimentCmd)
}

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
