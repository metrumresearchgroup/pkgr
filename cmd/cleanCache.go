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
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/metrumresearchgroup/pkgr/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//var cleanAllFromCache bool
var srcOnly bool
var binariesOnly bool
var reposToClear string

// cacheCmd represents the cache command
var cleanCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Subcommand to clean cached source and binary files.",
	Long: `This command is a subcommand of the "clean" command.

	Using this command deletes cached source and binary files. Use the
	--src and --binary options to specify which repos to clean each
	file type from.

	`,
	RunE: cache,
}

func init() {
	//cleanCacheCmd.Flags().BoolVar(&cleanAllFromCache, "all", true, "Clean both src and binary files from the cache.")
	cleanCacheCmd.Flags().BoolVar(&srcOnly, "src-only", false, "Clean only src files from the cache")
	cleanCacheCmd.Flags().BoolVar(&binariesOnly, "binaries-only", false, "Clean only binary files from the cache")
	cleanCacheCmd.Flags().StringVar(&reposToClear, "repos", "ALL", "Comma separated list of repositories to be cleaned. Defaults to all.")

	CleanCmd.AddCommand(cleanCacheCmd)
}

func cache(cmd *cobra.Command, args []string) error {
	cleanCacheFolders()
	return nil
}

func cleanCacheFolders() error {
	cachePath := userCache(cfg.Cache)
	repos := strings.Split(reposToClear, ",")

	Log.WithFields(logrus.Fields{
		"repos argument": reposToClear,
		"repos parsed":   sliceToString(repos),
		"cache dir":      cachePath,
	}).Trace("cleaning cache")

	if !srcOnly && !binariesOnly {
		Log.Info("clearing source and binary files from the cache")
		deleteAllCacheSubfolders(repos, cachePath)
	} else if srcOnly && binariesOnly {
		err := errors.New("invalid argument combination -- cannot combine srcOnly and binaryOnly flags")
		Log.Error(err)
		return err
	} else if srcOnly {
		Log.Info("clearing source files only from the cache")
		deleteCacheSubfolders(repos, "src", cachePath)
	} else if binariesOnly {
		Log.Info("clearing binary files only from the cache")
		deleteCacheSubfolders(repos, "binary", cachePath)
	} else {
		return errors.New("'what? that's impossible! my logic is flawless!'")
	}

	return nil
}

func deleteAllCacheSubfolders(repos []string, cacheDirectory string) {
	deleteCacheSubfolders(repos, "src", cacheDirectory)
	deleteCacheSubfolders(repos, "binary", cacheDirectory)
}

func deleteCacheSubfolders(repos []string, subfolder string, cacheDirectory string) error {
	var err error

	cacheDirFsObject, err := fs.Open(cacheDirectory)
	if err != nil {
		return err
	}

	repoFolderFsObjects, _ := cacheDirFsObject.Readdir(0)

	Log.WithFields(logrus.Fields{
		"repos argument": reposToClear,
		"repos parsed":   sliceToString(repos),
		"cache dir":      cacheDirectory,
	}).Trace("cleaning cache")

	if repos == nil || len(repos) == 0 || reposToClear == "ALL" {
		for _, repoFolderFsObject := range repoFolderFsObjects {
			subfolderPath := filepath.Join(
				cacheDirectory,
				repoFolderFsObject.Name(),
				subfolder,
			)
			err = fs.RemoveAll(subfolderPath)
			if err != nil {
				Log.Error(err)
			}
		}
	} else {
		for _, repoToClear := range repos {
			for _, repoFolderFsObject := range repoFolderFsObjects {

				if repoToClear == repoFolderFsObject.Name() {

					subfolderPath := filepath.Join(
						cacheDirectory,
						repoFolderFsObject.Name(),
						subfolder,
					)

					Log.WithFields(logrus.Fields{
						"repoToClear":             repoToClear,
						"repoFolderFsObject Name": repoFolderFsObject.Name(),
						"subfolder":               subfolder,
						"subfolder path":          subfolderPath,
					}).Trace("match found")

					err = fs.RemoveAll(subfolderPath)
					if err != nil {
						Log.Error(err)
					}
				}
			}
		}
	}
	return err
}

func sliceToString(str []string) string {
	returnString := "Size: " + strconv.Itoa(len(str)) + " :"
	for _, s := range str {
		returnString += s + "|"
	}
	return returnString
}
