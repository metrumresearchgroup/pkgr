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
	"time"

	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/gpsr"
	"github.com/dpastoor/rpackagemanager/rcmd"
	"github.com/sajari/fuzzy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// installCmd represents the R CMD install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a package",
	Long: `
	install a package
 `,
	RunE: rInstall,
}

func rInstall(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	repos := []cran.RepoURL{
		cran.RepoURL{
			Name: "CRAN",
			URL:  "https://cran.rstudio.com",
		},
		cran.RepoURL{
			Name: "gh_releases",
			URL:  "https://metrumresearchgroup.github.io/rpkgs/gh_releases",
		},
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

	pkgs := []string{
		// "PKPDmisc",
		// "mrgsolve",
		// "rmarkdown",
		// "bitops",
		// "caTools",
		// "GGally",
		// "knitr",
		// "gridExtra",
		// "htmltools",
		// "xtable",
		// "ggplot2",
		//"dplyr",
		// "shiny",
		// "shinydashboard",
		// "data.table",
		// "logrr",
		// "crayon",
		// "glue",
		// "rcpp",

		// should cover misspelled packages!!
		//"tidyVerse",
	}
	ap := cdb.GetPackages(pkgs)
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
	ip, err := gpsr.ResolveInstallationReqs(pkgs, make(map[string]gpsr.InstallDeps), cdb)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(time.Since(startTime))
	if err != nil {
		log.Fatalf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Info("The dependency graph resolved successfully")
	}
	var toDl []cran.PkgDl
	// starting packages
	for _, p := range ip.StartingPackages {
		pkg, repo, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Repo: repo})
	}
	// all other packages
	for p := range ip.DepDb {
		pkg, repo, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Repo: repo})
	}
	// // want to download the packages and return the full path of any downloaded package
	dl, err := cran.DownloadPackages(fs, toDl, cran.Source, "dump/cache")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	ia := rcmd.NewDefaultInstallArgs()
	ia.Library, _ = filepath.Abs("/Users/devinp/repos/queue-cancellation/packrat/lib/x86_64-apple-darwin15.6.0/3.5.1")
	err = rcmd.InstallPackagePlan(fs, ip, dl, rcmd.NewPackageCache("dump/cache", false), ia, rcmd.NewRSettings(), rcmd.ExecSettings{}, log, 12)
	if err != nil {
		fmt.Println("failed package install")
		fmt.Println(err)
	}
	fmt.Println("duration:", time.Since(startTime))
	return nil
}

func init() {
	installCmd.PersistentFlags().String("library", "", "library to install packages to")
	viper.BindPFlag("library", installCmd.PersistentFlags().Lookup("library"))
	RootCmd.AddCommand(installCmd)
}
