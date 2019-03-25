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

var pkgdbsToClearArgument string

// pkgdbCmd represents the pkgdb command
var pkgdbsCmd = &cobra.Command{
	Use:   "pkgdbs",
	Short: "Subcommand to clean cached pkgdbs",
	Long: `This command parses the currently-cached pkgdbs and removes all
	of them by default, or specific ones if desired. Identify specific repos using the "repos" argument, i.e.
	pkgr clean pkgdbs --repos="CRAN,r_validated"
	Repo names should match names in the pkgr.yml file.`,
	RunE: executeCommand,
}

func init() {

	pkgdbsCmd.Flags().StringVar(&pkgdbsToClearArgument, "repos", "ALL", "Set the repos you wish to clear the pkgdbs for.")
	CleanCmd.AddCommand(pkgdbsCmd)
}

// Using a "middle man" function here allows us to call the cleanPackageDatabases function from the "clean"
// command when "pkgr clean --all" is run. Without the middle man, passing in arguments becomes strange.
func executeCommand(cmd *cobra.Command, args []string) error {
	return cleanPackageDatabases(pkgdbsToClearArgument)
}

// cleanPackageDatabases function to remove the cached package databases for each package listed in pkgdbs.
// pkdgbs should be a comma-separated string, e.g. "CRAN,r_validated"
// To remove all package dbs associated with a pkgr.yml file, use pkgdbs = "ALL"
func cleanPackageDatabases(pkgdbs string) error {

	var pkgdbsToClear []string

	//db := NewRepoDb(url, dst, cfgdb.Repos[url.Name], rv)
	if pkgdbs == "ALL" {
		for _, repoMap := range cfg.Repos {
			for key := range repoMap {
				log.Info("Key: " + key)
				pkgdbsToClear = append(pkgdbsToClear, key)
			}
		}
	} else {
		pkgdbsToClear = strings.Split(pkgdbs, ",")
	}

	totalPackageDbsProvided := len(pkgdbsToClear)
	totalPackageDbsDeleted := removePackageDatabases(pkgdbsToClear)

	log.WithFields(logrus.Fields{
		"Packages specified": totalPackageDbsProvided,
		"Packages removed":   totalPackageDbsDeleted,
	}).Info("finished cleaning package dbs.")

	return nil
}

func removePackageDatabases(pkgdbsToClear []string) int {
	var err error
	filesRemoved := 0

	//TODO: This is duplicate code from a another function, see if we can pull this out somewhere.
	//We need to include thise code block because we'll need this object to make a RepoDb object
	//later. We need to make a RepoDb object in order to call the "GetRepoDbCacheFilePath" function
	//on that object. In order for that function to work properly, the RepoDb needs to be constructed
	//the same way it would on a pkgr plan command.
	installConfig := cran.NewInstallConfig()
	for rn, val := range cfg.Customizations.Repos {
		if strings.EqualFold(val.Type, "binary") {
			installConfig.Repos[rn] = cran.RepoConfig{DefaultSourceType: cran.Binary}
		}
		if strings.EqualFold(val.Type, "source") {
			installConfig.Repos[rn] = cran.RepoConfig{DefaultSourceType: cran.Source}
		}
	}

	rs := rcmd.NewRSettings()
	rVersion := rcmd.GetRVersion(&rs)

	for _, dbToClear := range pkgdbsToClear {
		for _, repoFromConfig := range cfg.Repos {
			urlString, foundInConfig := repoFromConfig[dbToClear]
			urlObject := cran.RepoURL{Name: dbToClear, URL: urlString}
			if foundInConfig {
				db, _ := cran.NewRepoDb(urlObject, cran.DefaultType(), installConfig.Repos[dbToClear], rVersion)
				filepathToRemove := db.GetRepoDbCacheFilePath()

				_, err = fs.Stat(filepathToRemove)
				fileExists := err == nil

				if fileExists {
					log.Trace("attempting to remove " + filepathToRemove)
					err = fs.Remove(filepathToRemove)
					if err != nil {
						log.Error(err)
					} else {
						filesRemoved++
					}
				} else {
					log.WithField("pkgdb", filepathToRemove).Warn("pkgdb was not found")
				}
			} else {
				log.WithField("repo", dbToClear).Warn("could not find repository in config file")
			}
		}
	}
	return filesRemoved
}
