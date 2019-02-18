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

	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

// inspectCmd shows the install inspect
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect a full installation",
	Long: `
	see the inspect for an install
 `,
	RunE: inspect,
}

var reverse bool
var showDeps bool
var toJSON bool
var tree bool

func recurseDeps(pkg string, ddb gpsr.InstallPlan, t treeprint.Tree) {
	pkgDeps := ddb.DepDb[pkg]
	if len(pkgDeps) == 0 {
		return
	}
	for _, d := range pkgDeps {
		recurseDeps(d, ddb, t.AddBranch(d))
	}
}
func inspect(cmd *cobra.Command, args []string) error {
	if toJSON {
		// this should suppress all logging from the planning
		log.SetLevel(logrus.FatalLevel)
	}
	_, ip := planInstall()
	if showDeps {
		var allDeps map[string][]string
		keepDeps := make(map[string][]string)
		if reverse {
			allDeps = ip.InvertDependencies()
		} else {
			allDeps = ip.DepDb
		}
		if len(args) > 0 {
			for _, arg := range args {
				keepDeps[arg] = allDeps[arg]
			}
			printDeps(keepDeps, tree, ip)
		} else {
			printDeps(allDeps, tree, ip)
		}
	}
	return nil
}

func printDeps(deps map[string][]string, tree bool, ip gpsr.InstallPlan) {
	if tree {
		depTree := treeprint.New()
		for p := range deps {
			tb := depTree.AddBranch(p)
			recurseDeps(p, ip, tb)
		}
		fmt.Println(depTree.String())
	} else {
		prettyPrint(deps)
	}
}

func init() {
	inspectCmd.Flags().BoolVar(&showDeps, "deps", false, "show dependency tree")
	inspectCmd.Flags().BoolVar(&reverse, "reverse", false, "show reverse dependencies")
	inspectCmd.Flags().BoolVar(&tree, "tree", false, "show full recursive dependency tree")
	inspectCmd.Flags().BoolVar(&toJSON, "json", false, "output as clean json")

	RootCmd.AddCommand(inspectCmd)
}
