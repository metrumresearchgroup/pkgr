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
	"strconv"

	"github.com/spf13/cobra"
)

var cleanAll bool
var cleanPkgdbs string

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean up cached information",
	Long:  "clean up cached source files and binaries, as well as the cached package database.",
	RunE:  clean,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("clean called")
	//	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cleanCmd.Flags().Bool("all", true, "clean all cached items")
	cleanCmd.Flags().StringVar(&cleanPkgdbs, "pkgdb", "ALL", "Remove cached package databases.")
	//viper.BindFlag("pkgdb")

	RootCmd.AddCommand(cleanCmd)
}

func clean(cmd *cobra.Command, args []string) error {
	fmt.Println(fmt.Sprintf("cleanPkgdbs: %s", cleanPkgdbs))
	fmt.Println(fmt.Sprintf("cleanAll: %s", strconv.FormatBool(cleanAll)))
	return nil
}
