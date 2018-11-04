package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/desc"
	"github.com/dpastoor/rpackagemanager/gpsr"
	"github.com/dpastoor/rpackagemanager/rcmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	appFS := afero.NewOsFs()
	log := logrus.New()
	log.Level = logrus.DebugLevel
	// log.SetFormatter(&logrus.JSONFormatter{})
	// appFS.Remove("logfile.txt")
	// logf, _ := appFS.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// log.SetOutput(logf)
	startTime := time.Now()
	repos := []cran.RepoURL{
		cran.RepoURL{
			Name: "CRAN",
			URL:  "https://cran.rstudio.com",
		},
		cran.RepoURL{
			Name: "gh_releases",
			URL:  "https://metrumresearchgroup.github.io/rpkgs/gh_releases",
		},
	}
	cdb, err := cran.NewPkgDb(repos)
	if err != nil {
		fmt.Println("error getting pkgdb ", err)
		panic(err)
	}
	//PrettyPrint(cdb)
	for _, db := range cdb.Db {
		fmt.Println(fmt.Sprintf("%v packages pulled in for %s from %s", len(db.Db), db.Repo.Name, db.Repo.URL))
	}

	p, repo, _ := cdb.GetPackage("logrrr")
	fmt.Println(repo)
	PrettyPrint(p)

	p, repo, _ = cdb.GetPackage("PKPDmisc")
	fmt.Println(repo)
	PrettyPrint(p)
	// PrettyPrint(dmap["dplyr"])
	// PrettyPrint(dmap["PKPDmisc"])
	//AppFs := afero.NewOsFs()
	// can use this to redirect log output
	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	pkgs := []string{
		// "PKPDmisc",
		// "mrgsolve",
		// "rmarkdown",
		// "bitops",
		// "caTools",
		// "GGally",
		// "knitr",
		// "gridExtra",
		// "htmltools",
		// "xtable",
		// "tidyverse",
		// "shiny",
		// "shinydashboard",
		// "data.table",
		"logrrr",
	}
	ip, err := gpsr.ResolveInstallationReqs(pkgs, cdb)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(time.Since(startTime))
	PrettyPrint(ip)
	fmt.Println("inverted dependencies: ")
	PrettyPrint(ip.InvertDependencies())
	if err != nil {
		log.Fatalf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Info("The dependency graph resolved successfully")
	}
	var toDl []cran.PkgDl
	// starting packages
	for _, p := range ip.StartingPackages {
		pkg, repo, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Repo: repo})
	}
	// all other packages
	for p := range ip.DepDb {
		pkg, repo, _ := cdb.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Repo: repo})
	}
	// // want to download the packages and return the full path of any downloaded package
	dl, err := cran.DownloadPackages(appFS, toDl, cran.Source, "dump")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// // ia := rcmd.NewDefaultInstallArgs()
	// // ia.Library = "../integration_tests/lib"
	// for pn, p := range dl {
	// 	fmt.Println(pn, p)
	// }
	ia := rcmd.NewDefaultInstallArgs()
	ia.Library, _ = filepath.Abs("dump/test3lib/")
	fmt.Println("library set to: ", ia.Library)
	err = rcmd.InstallPackagePlan(appFS, ip, dl, ia, rcmd.NewRSettings(), rcmd.ExecSettings{}, log, 12)
	if err != nil {
		fmt.Println("failed package install")
		fmt.Println(err)
	}
	fmt.Println("duration:", time.Since(startTime))
	// fmt.Println("library: ", viper.GetString("library"))
}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func GobDB() {
	appFS := afero.NewOsFs()
	ok, _ := afero.Exists(appFS, "crandb.gob")
	if !ok {
		startTime := time.Now()
		res, err := http.Get("https://cran.rstudio.com/src/contrib/PACKAGES")
		if err != nil {
			fmt.Println("problem getting packages")
			panic(err)
		}
		file, err := os.Create("crandb.gob")
		if err != nil {
			fmt.Println("problem creating crandb")
			panic(err)
		}
		dmap := make(map[string]desc.Desc)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		cb := bytes.Split(body, []byte("\n\n"))
		for _, p := range cb {
			reader := bytes.NewReader(p)
			d, err := desc.ParseDesc(reader)
			dmap[d.Package] = d
			if err != nil {
				fmt.Println("problem parsing")
				panic(err)
			}
			//PrettyPrint(d)
		}
		fmt.Println("duration:", time.Since(startTime))
		fmt.Println("length: ", len(dmap))

		e := gob.NewEncoder(file)

		// Encoding the map
		err = e.Encode(dmap)
		if err != nil {
			panic(err)
		}
	}
}
