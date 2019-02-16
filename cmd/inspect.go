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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// inspectCmd shows the install inspect
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect a full installation",
	Long: `
	see the inspect for an install
 `,
	RunE: inspect,
}

func inspect(cmd *cobra.Command, args []string) error {
	log.Infof("Installation would launch %v workers\n", getWorkerCount())
	pdb, ip := planInstall()
	fmt.Println(pdb)
	if viper.GetBool("show-deps") {
		if viper.GetBool("invert") {
			prettyPrint(ip.InvertDependencies())
		} else {
			prettyPrint(ip.DepDb)
		}
	}
	return nil
}

func init() {
	inspectCmd.PersistentFlags().Bool("show-deps", false, "show the (required) dependencies for each package")
	viper.BindPFlag("show-deps", inspectCmd.PersistentFlags().Lookup("show-deps"))
	viper.BindPFlag("invert", inspectCmd.PersistentFlags().Lookup("invert"))
	RootCmd.AddCommand(inspectCmd)
}
