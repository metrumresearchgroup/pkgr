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
	"path/filepath"
	"runtime"
	"time"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// installCmd represents the R CMD install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a package",
	Long: `
	install a package
 `,
	RunE: rInstall,
}

func rInstall(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	cdb, ip := planInstall()

	var toDl []cran.PkgDl
	// starting packages
	for _, p := range ip.StartingPackages {
		pkg, cfg, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	// all other packages
	for p := range ip.DepDb {
		pkg, cfg, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	// // want to download the packages and return the full path of any downloaded package
	pc := rcmd.NewPackageCache(userCache(cfg.Cache), false)
	dl, err := cran.DownloadPackages(fs, toDl, pc.BaseDir)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	ia := rcmd.NewDefaultInstallArgs()
	ia.Library, _ = filepath.Abs(cfg.Library)
	nworkers := runtime.NumCPU()
	// leave at least 1 thread open for coordination, given more than 2 threads available.
	// if only 2 available, will let the OS hypervisor coordinate some else would drop the
	// install time too much for the little bit of additional coordination going on.
	if nworkers > 2 {
		nworkers = nworkers - 1
	}
	err = rcmd.InstallPackagePlan(fs, ip, dl, pc, ia, rcmd.NewRSettings(), rcmd.ExecSettings{}, log, nworkers)
	if err != nil {
		fmt.Println("failed package install")
		fmt.Println(err)
	}
	fmt.Println("duration:", time.Since(startTime))
	return nil
}

func init() {
	installCmd.PersistentFlags().String("library", "", "library to install packages to")
	viper.BindPFlag("library", installCmd.PersistentFlags().Lookup("library"))
	RootCmd.AddCommand(installCmd)
}
