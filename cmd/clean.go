// Copyright © 2019 John Carlo Salter <juuncaerlum@gmail.com>
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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cleanAll bool
var cleanPkgdbs bool
var pkgdbs string
var cleanCache bool
var srcCaches string
var binaryCaches string

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean up cached information",
	Long:  "clean up cached source files and binaries, as well as the saved package database.",
	RunE:  clean,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("clean called")
	//	},
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanAll, "all", false, "clean all cached items")
	cleanCmd.Flags().BoolVar(&cleanPkgdbs, "pkgdbs", false, "Remove cached package databases.")
	cleanCmd.Flags().StringVar(&pkgdbs, "dbs", "ALL", "Package databases to remove.")
	cleanCmd.Flags().BoolVar(&cleanCache, "cache", false, "Remove cache sources and/or binaries")
	cleanCmd.Flags().StringVar(&srcCaches, "src", "ALL", "Clean src caches in clean --cache")
	cleanCmd.Flags().StringVar(&binaryCaches, "binary", "ALL", "Clean binary caches in clean --cache")

	RootCmd.AddCommand(cleanCmd)
}

func clean(cmd *cobra.Command, args []string) error {

	if !cleanAll && !cleanPkgdbs && !cleanCache {
		fmt.Println("No clean options passed -- not cleaning.")
	}

	if cleanAll {
		fmt.Println("Cleaning all.")
	} else {
		if cleanPkgdbs {
			if pkgdbs == "ALL" {
				fmt.Println("Cleaning all pkgdbs")
			} else {
				fmt.Println(fmt.Sprintf("Cleaning specific package databases: %s", pkgdbs))
			}
		}

		if cleanCache {
			if srcCaches == "ALL" {
				fmt.Println("Cleaning all src caches.")
				clearCaches(nil, nil)
			} else {
				fmt.Println(fmt.Sprintf("Cleaning specific src caches: %s", srcCaches))
			}

			if binaryCaches == "ALL" {
				fmt.Println("Cleaning all binary caches.")
			} else {
				fmt.Println(fmt.Sprintf("Cleaning specific binary caches: %s", binaryCaches))
			}
		}
	}
	fmt.Println("Donezo.")
	return nil
}

func clearCaches(src, binary []string) error {
	cacheDir := userCache(cfg.Cache)
	log.WithField("dir", cacheDir).Info("clearing cache at directory ")

	return nil
}

//When you bind a bool flag, it's basically "if that flag is not set then it's false, if it is set then it's true"
//Default string is going
