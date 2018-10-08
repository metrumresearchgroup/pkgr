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
	"github.com/dpastoor/rpackagemanager/gpsr"
	"github.com/dpastoor/rpackagemanager/packrat"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// planCmd shows the install plan
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "plan a full installation",
	Long: `
	see the plan for an install
 `,
	RunE: plan,
}

func plan(cmd *cobra.Command, args []string) error {

	//AppFs := afero.NewOsFs()
	// can use this to redirect log output
	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	appFS := afero.NewOsFs()
	lf, _ := afero.ReadFile(appFS, viper.GetString("pr_lockfile"))
	pm := packrat.ChunkLockfile(lf)
	var workingGraph gpsr.Graph
	for _, p := range pm.CRANlike {
		workingGraph = append(workingGraph, gpsr.NewNode(p.Package, p.Requires))
	}
	for _, p := range pm.Github {
		workingGraph = append(workingGraph, gpsr.NewNode(p.Reqs.Package, p.Reqs.Requires))
	}

	if viper.GetBool("preview") {
		gpsr.DisplayGraph(workingGraph)
	}

	resolved, err := gpsr.ResolveGraph(workingGraph)
	if err != nil {
		log.Fatalf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Info("The dependency graph resolved successfully")
	}

	for i, pkglayer := range resolved {
		log.WithFields(
			logrus.Fields{
				"layer": i + 1,
				"npkgs": len(pkglayer),
			},
		).Info(pkglayer)
	}
	return nil
}

func init() {
	planCmd.PersistentFlags().String("library", "", "library to plan packages to")
	viper.BindPFlag("library", planCmd.PersistentFlags().Lookup("library"))
	RootCmd.AddCommand(planCmd)
}
