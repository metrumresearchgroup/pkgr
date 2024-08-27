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
	"sort"

	"github.com/metrumresearchgroup/pkgr/logger"

	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/metrumresearchgroup/pkgr/pacman"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

// inspectCmd shows the install inspect
var inspectCmd = &cobra.Command{
	Use:   "inspect --deps [flags] [<package>...]",
	Short: "Inspect package dependencies",
	Long: `The inspect subcommand provides an entry point for displaying
information that can be gathered by examining the configuration file, the
associated package database, and the library. The current focus is on
inspecting package dependencies (triggered by passing --deps).

Note: If the configuration file has 'Suggests: true', that does not affect
the set of dependencies listed for any particular package. Instead the set
of suggested packages is included in the top-level package set.`,
	Example: `  # Show all dependencies as a tree
  pkgr --loglevel=fatal inspect --deps --tree
  # Show dependency tree, restricting roots to the named packages
  pkgr --loglevel=fatal inspect --deps --tree processx here

  # Output a JSON record where each item maps a package to its direct
  # and indirect dependencies
  pkgr --loglevel=fatal inspect --deps
  # Do the same, but filter to records for the named packages
  pkgr --loglevel=fatal inspect --deps processx here

  # Output a JSON record where each item maps a package to
  # the packages that have it as a dependency
  pkgr --loglevel=fatal inspect --deps --reverse`,
	RunE: inspect,
}

var reverse bool
var showDeps bool
var toJson bool
var tree bool
var installedFrom bool

func recurseDeps(pkg string, ddb gpsr.InstallPlan, t treeprint.Tree) {
	pkgDeps := ddb.DepDb[pkg]
	sort.Strings(pkgDeps)
	if len(pkgDeps) == 0 {
		return
	}
	for _, d := range pkgDeps {
		recurseDeps(d, ddb, t.AddBranch(d))
	}
}

func inspect(cmd *cobra.Command, args []string) error {

	logger.AddLogFile(cfg.Logging.All, cfg.Logging.Overwrite)

	if toJson {
		// this should suppress all logging from the planning
		logger.SetLogLevel("fatal")
	}

	if installedFrom {
		printInstalledFromPackages()
		return nil
	}

	rs := rcmd.NewRSettings(cfg.RPath)
	rVersion := rcmd.GetRVersion(&rs)
	_, ip, _ := planInstall(rVersion, true)
	if showDeps {
		var allDeps map[string][]string
		keepDeps := make(map[string][]string)
		if reverse {
			allDeps = ip.InvertDependencies()
		} else {
			allDeps = ip.DepDb
		}
		for p := range allDeps {
			sort.Strings(allDeps[p])
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

func printInstalledFromPackages() {
	prettyPrint(pacman.GetPackagesByInstalledFrom(fs, cfg.Library))
}

func init() {
	inspectCmd.Flags().BoolVar(&showDeps, "deps", false, "show dependency tree")
	inspectCmd.Flags().BoolVar(&reverse, "reverse", false, "show reverse dependencies")
	inspectCmd.Flags().BoolVar(&tree, "tree", false, "show full recursive dependency tree")
	inspectCmd.Flags().BoolVar(&toJson, "json", false,
		// Point the user to --loglevel because it also suppresses logging
		// upstream of the inspect() call.
		"suppress non-fatal logging (note: prefer --loglevel=fatal to this flag)")
	inspectCmd.Flags().BoolVar(&installedFrom, "installed-from", false, "show package installation source")

	// Don't advertise this until work is done to improve it.
	inspectCmd.Flags().MarkHidden("installed-from")

	RootCmd.AddCommand(inspectCmd)
}
