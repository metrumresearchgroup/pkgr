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
	"strings"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pkgdbsToClear string

// pkgdbCmd represents the pkgdb command
var pkgdbCmd = &cobra.Command{
	Use:   "pkgdb",
	Short: "Subcommand to clean cached pkgdbs",
	Long: `This command parses the currently-cached pkgdbs and removes all
	of them by default, or specific ones if desired.`,
	RunE: pkgdb,
}

func init() {

	pkgdbCmd.Flags().StringVar(&pkgdbsToClear, "repos", "ALL", "Set the repos you wish to clear the pkgdbs for.")
	CleanCmd.AddCommand(pkgdbCmd)
}

func pkgdb(cmd *cobra.Command, args []string) error {

	var clearRepos []string

	//db := NewRepoDb(url, dst, cfgdb.Repos[url.Name], rv)
	if pkgdbsToClear == "ALL" {
		for _, repoMap := range cfg.Repos {
			for key := range repoMap {
				log.Info("Key: " + key)
				clearRepos = append(clearRepos, key)
			}
		}
	} else {
		clearRepos = strings.Split(pkgdbsToClear, ",")
	}

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

			log.WithField("repoFromConfig", repoFromConfig).Info("In the for loop.")
			urlString, found := repoFromConfig[clearRepo]
			log.WithFields(logrus.Fields{
				"urlString":      urlString,
				"found":          found,
				"repoFromConfig": repoFromConfig,
				"clearRepo":      clearRepo,
			}).Info("Check it out, yo.")

			urlObject := cran.RepoURL{Name: clearRepo, URL: urlString}

			if found {
				db, _ := cran.NewRepoDb(urlObject, cran.DefaultType(), cic.Repos[clearRepo], rVersion)
				filepathToRemove := db.GetRepoDbCacheFilePath()
				log.Info("Attempting to remove " + filepathToRemove)
				err := fs.Remove(filepathToRemove)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	return nil
}
