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
	"time"

	"github.com/dpastoor/rpackagemanager/configlib"
	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/gpsr"
	"github.com/sajari/fuzzy"
	"github.com/sirupsen/logrus"
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
	startTime := time.Now()
	var repos []cran.RepoURL
	for _, r := range cfg.Repos {
		for nm, url := range r {
			repos = append(repos, cran.RepoURL{Name: nm, URL: url})
		}
	}

	cdb, err := cran.NewPkgDb(repos)
	if err != nil {
		fmt.Println("error getting pkgdb ", err)
		panic(err)
	}
	//PrettyPrint(cdb)
	for _, db := range cdb.Db {
		fmt.Println(fmt.Sprintf("%v packages available in for %s from %s", len(db.Db), db.Repo.Name, db.Repo.URL))
	}
	ids := gpsr.NewDefaultInstallDeps()
	if cfg.Suggests {
		for _, pkg := range cfg.Packages {
			// set all top level packages to install suggests
			dp := ids.Default
			dp.Suggests = true
			ids.Deps[pkg] = dp
		}
	}
	as := viper.AllSettings()
	for pkg, v := range cfg.Customizations {
		if configlib.IsCustomizationSet("Suggests", as, pkg) {
			dp := ids.Default
			dp.Suggests = v.Suggests
			ids.Deps[pkg] = dp
		}
		if configlib.IsCustomizationSet("Repo", as, pkg) {
			err := cdb.SetPackageRepo(pkg, v.Repo)
			if err != nil {
				log.WithFields(logrus.Fields{
					"pkg":  pkg,
					"repo": v.Repo,
				}).Fatal("error finding custom repo to set")
			}
		}
	}

	ap := cdb.GetPackages(cfg.Packages)
	if len(ap.Missing) > 0 {
		fmt.Println("missing packages: ", ap.Missing)
		model := fuzzy.NewModel()

		// For testing only, this is not advisable on production
		model.SetThreshold(1)

		// This expands the distance searched, but costs more resources (memory and time).
		// For spell checking, "2" is typically enough, for query suggestions this can be higher
		model.SetDepth(1)
		pkgs := cdb.GetAllPkgsByName()
		model.Train(pkgs)
		for _, mp := range ap.Missing {
			fmt.Println("did you mean one of: ", model.Suggestions(mp, false))
		}
		os.Exit(1)
	}
	for _, pkg := range ap.Packages {
		log.WithFields(logrus.Fields{
			"pkg":  pkg.Package.Package,
			"repo": pkg.Repo.Name,
		}).Trace("package repository set")
	}

	ip, err := gpsr.ResolveInstallationReqs(cfg.Packages, ids, cdb)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	pkgs := ip.StartingPackages
	for pkg := range ip.DepDb {
		pkgs = append(pkgs, pkg)
	}
	fmt.Println("total packages required:", len(ip.StartingPackages)+len(ip.DepDb))
	fmt.Println(time.Since(startTime))
	return nil
}

func init() {
	RootCmd.AddCommand(planCmd)
}
