// Copyright © 2019 Metrum Research Group
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

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// installCmd represents the R CMD install command
var downloadCmd = &cobra.Command{
	Use:   "download-src",
	Short: "download all packages",
	Long: `
	download src version of packages in the plan into a single folder.
	download works by converting the package type to source, and downloading
	the src version from the appropriate repo. It does not do version mismatch
	reconciliation or other more sophisticated analysis to confirm that 
	the src version and binary version are consistent for platforms that
	support binary downloads
	
	This will download packages to the folder 'src' in the cache folder
`,
	RunE: rDownload,
}

func init() {
	downloadCmd.PersistentFlags().String("dir", "", "directory to download")
	viper.BindPFlag("dir", downloadCmd.PersistentFlags().Lookup("dir"))
	RootCmd.AddCommand(downloadCmd)
}

func rDownload(cmd *cobra.Command, args []string) error {
	dir := viper.GetString("dir")
	if dir != "" {
		cfg.Cache = dir
	}
	// set Type to source for everything to make sure we get source tarballs
	// this will brak the abstraction if src and binary version are different for the platforms
	if cfg.Customizations.Packages == nil {
		cfg.Customizations.Packages = make(map[string]configlib.PkgConfig)
	}
	for _, p := range cfg.Packages {
		//lets set every package as having a customization of type source
		pkg, _ := cfg.Customizations.Packages[p]
		pkg.Type = "source"
		cfg.Customizations.Packages[p] = pkg
	}
	// Initialize log and start time.
	rSettings := rcmd.NewRSettings(cfg.RPath)
	rVersion := rcmd.GetRVersion(&rSettings)
	log.Infoln("R Version " + rVersion.ToFullString())
	// most people should know what platform they are on
	log.Debugln("OS Platform " + rSettings.Platform)

	// Get master object containing the packages available in each repository (pkgNexus),
	//  as well as a master install plan to guide our process.
	_, installPlan, _ := planInstall(rVersion, cran.Source, true)
	for _, pdl := range installPlan.PackageDownloads {
		pkg, repo := pdl.PkgAndRepoNames()
		fmt.Println(pkg, repo, pdl.Package.Version, pdl.Config.Type)
	}
	// Retrieve a cache to store any packages we need to download for the install.
	packageCache := rcmd.NewPackageCache(userCache(cfg.Cache), false)
	dlPath := filepath.Join(packageCache.BaseDir, "src")
	err := fs.MkdirAll(dlPath, 0777)
	if err != nil {
		log.Fatal("error creating download cache dir at: ", dlPath)
	}
	// //Create a pkgMap object, which helps us with parallel downloads (?)
	_, err = cran.DownloadPackages(fs, installPlan.PackageDownloads, packageCache.BaseDir, rVersion, false)
	if err != nil {
		log.Fatalf("error downloading packages: %s", err)
	}
	return nil
}
