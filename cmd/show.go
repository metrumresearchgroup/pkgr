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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd shows information about internal settings
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show with the cli",
	Long: `
	show internal settings
 `,
	RunE: rShow,
}

func rShow(cmd *cobra.Command, args []string) error {

	//AppFs := afero.NewOsFs()
	// can use this to redirect log output
	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	Log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	fmt.Println("------viper-------")
	as := viper.AllSettings()
	fmt.Println(as)
	fmt.Println("subs:")
	fmt.Println(viper.Sub("Customizations"))
	fmt.Println(viper.Sub("Customizations").AllSettings()["packages"])
	fmt.Println(viper.Sub("Customizations").AllSettings()["repos"])
	fmt.Println("------config-------")
	prettyPrint(cfg)
	return nil
}

func init() {
	RootCmd.AddCommand(showCmd)
}
