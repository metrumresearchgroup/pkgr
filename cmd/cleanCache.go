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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cleanCache bool
var srcCaches string
var binaryCaches string

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
	cleanCacheCmd.Flags().StringVar(&srcCaches, "src", "ALL", "Clean src caches in clean cache")
	cleanCacheCmd.Flags().StringVar(&binaryCaches, "binary", "ALL", "Clean binary caches in clean cache")

	CleanCmd.AddCommand(cleanCacheCmd)
}

func cache(cmd *cobra.Command, args []string) error {
	cleanCacheFolders()
	return nil
}

func cleanCacheFolders() {
	cachePath := userCache(cfg.Cache)

	if srcCaches == "ALL" {
		fmt.Println("Cleaning all src caches.")
		_ = deleteCacheSubfolders(nil, "src", cachePath)
	} else {
		fmt.Println(fmt.Sprintf("Cleaning specific src caches: %s", srcCaches))
		srcRepos := strings.Split(srcCaches, ",")
		_ = deleteCacheSubfolders(srcRepos, "src", cachePath)
	}

	if binaryCaches == "ALL" {
		fmt.Println("Cleaning all binary caches.")
		deleteCacheSubfolders(nil, "binary", cachePath)
	} else {
		fmt.Println(fmt.Sprintf("Cleaning specific binary caches: %s", binaryCaches))
		binaryRepos := strings.Split(binaryCaches, ",")
		_ = deleteCacheSubfolders(binaryRepos, "binary", cachePath)
	}
}

func deleteAllCacheSubfolders(cacheDirectory string) {
	deleteCacheSubfolders(nil, "src", cacheDirectory)
	deleteCacheSubfolders(nil, "binary", cacheDirectory)
}

func deleteCacheSubfolders(repos []string, subfolder string, cacheDirectory string) error {
	cacheDirFsObject, err := fs.Open(cacheDirectory)
	if err != nil {
		return err
	}

	repoFolders, _ := cacheDirFsObject.Readdir(0)

	if repos == nil || len(repos) == 0 {
		for _, repoFolder := range repoFolders {
			subfolderPath := filepath.Join(cacheDirectory, repoFolder.Name(), subfolder)
			fs.RemoveAll(subfolderPath)
		}
	} else {
		for _, repoToClear := range repos {
			for _, repoFolder := range repoFolders {
				if repoToClear == repoFolder.Name() {
					subfolderPath := filepath.Join(cacheDirectory, repoFolder.Name(), subfolder)
					fs.Remove(subfolderPath)
				}
			}
		}
	}
	return nil
}
