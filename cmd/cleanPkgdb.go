// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"strings"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/cobra"
)

var pkgdbsToClear string

// pkgdbCmd represents the pkgdb command
var pkgdbCmd = &cobra.Command{
	Use:   "pkgdb",
	Short: "Subcommand to clean cached pkgdbs",
	Long: `This command parses the currently-cached pkgdbs and removes all
	of them by default, or specific ones if desired.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pkgdb called")
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pkgdbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pkgdbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pkgdbCmd.Flags().StringVar(&pkgdbsToClear, "repos", "ALL", "Set the repos you wish to clear the pkgdbs for.")
	CleanCmd.AddCommand(pkgdbCmd)
}

func pkgdb(cmd *cobra.Command, args []string) error {

	//db := NewRepoDb(url, dst, cfgdb.Repos[url.Name], rv)
	clearRepos := strings.Split(pkgdbsToClear, ",")

	//TODO: This is duplicate code from a another function, see if we can pull this out somewhere.
	cic := cran.NewInstallConfig()
	for rn, val := range cfg.Customizations.Repos {
		if strings.EqualFold(val.Type, "binary") {
			cic.Repos[rn] = cran.RepoConfig{DefaultSourceType: cran.Binary}
		}
		if strings.EqualFold(val.Type, "source") {
			cic.Repos[rn] = cran.RepoConfig{DefaultSourceType: cran.Source}
		}
	}

	for _, clearRepo := range clearRepos {

		//Taken from plan.go
		rs := rcmd.NewRSettings()
		rVersion := rcmd.GetRVersion(&rs)

		for _, repoFromConfig := range cfg.Repos {
			urlString, found := repoFromConfig[clearRepo]

			urlObject := cran.RepoURL{Name: clearRepo, URL: urlString}

			if found {
				db, _ := cran.NewRepoDb(urlObject, cran.DefaultType(), cic.Repos[clearRepo], rVersion)
				fs.Remove(db.GetRepoDbCacheFilePath())
			}
		}

	}

	return nil
}
