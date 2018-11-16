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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// checkCmd represents the R CMD check command
var experimentCmd = &cobra.Command{
	Use:   "experiment",
	Short: "experiment with the cli",
	Long: `
	JUST FOR EXPERIMENTATION
 `,
	RunE: rExperiment,
}

func rExperiment(cmd *cobra.Command, args []string) error {

	//AppFs := afero.NewOsFs()
	// can use this to redirect log output
	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	prettyPrint(viper.AllSettings())
	prettyPrint(cfg)
	configDir, _ := filepath.Abs(viper.ConfigFileUsed())
	fmt.Println(os.Getwd())
	os.Chdir(filepath.Dir(configDir))
	fmt.Println(configDir)
	fmt.Println(os.Getwd())
	return nil
}

func init() {
	RootCmd.AddCommand(experimentCmd)
}

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
