/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Checks that installed packages can be loaded",
	Long: `Attempts to load user packages specified in pkgr.yml to validate that each package has been installed
successfully and can be used. Use the --all flag to load all packages in the user-library dependency tree instead of just user-level packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load called")
	},
}

func init() {
	RootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loadCmd.Flags().Bool("all", false, "load user packages as well as their dependencies")
	viper.BindPFlag("all", loadCmd.LocalFlags().Lookup("all")) //There doesn't seem to be a way to bind local flags.
	loadCmd.Flags().Bool("json", false, "output a JSON object of package info at the end")
	viper.BindPFlag("json", loadCmd.LocalFlags().Lookup("json")) //There doesn't seem to be a way to bind local flags.
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Load() {

}
