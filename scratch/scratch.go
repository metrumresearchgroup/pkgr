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
	GobDB()
	dmap := make(map[string]desc.Desc)
	appFS := afero.NewOsFs()
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.SetFormatter(&logrus.JSONFormatter{})
	appFS.Remove("logfile.txt")
	logf, _ := appFS.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(logf)

	file, err := os.Open("crandb.gob")
	if err != nil {
		fmt.Println("problem creating crandb")
		panic(err)
	}
	d := gob.NewDecoder(file)

	// Decoding the serialized data
	err = d.Decode(&dmap)
	if err != nil {
		panic(err)
	}

	// PrettyPrint(dmap["dplyr"])
	// PrettyPrint(dmap["PKPDmisc"])
	//AppFs := afero.NewOsFs()
	// can use this to redirect log output
	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	// pkgs := []string{
	// 	"PKPDmisc",
	// 	"mrgsolve",
	// 	"rmarkdown",
	// 	"bitops",
	// 	"caTools",
	// 	"GGally",
	// 	"knitr",
	// 	"gridExtra",
	// 	"htmltools",
	// 	"xtable",
	// 	"tidyverse",
	// 	"shiny",
	// 	"shinydashboard",
	// }
	pkgs := []string{
		"ggplot2",
	}
	startTime := time.Now()
	ip, err := gpsr.ResolveInstallationReqs(pkgs, dmap)
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
	var toDl []desc.Desc
	// starting packages
	for _, p := range ip.StartingPackages {
			toDl = append(toDl, dmap[p])
	}
	// all other packages
	for p := range ip.DepDb {
			toDl = append(toDl, dmap[p])
	}
	// // want to download the packages and return the full path of any downloaded package
	dl, err := cran.DownloadPackages(appFS, toDl, "https://cran.rstudio.com", cran.Source, "dump")
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
	ia.Library, _ = filepath.Abs("dump/test1lib/")
	fmt.Println("library set to: ", ia.Library)
	startTime = time.Now()
	err = rcmd.InstallPackagePlan(appFS, ip, dl, ia, rcmd.RSettings{}, rcmd.ExecSettings{}, log, 6)
	if err != nil {
		fmt.Println("failed package install")
		fmt.Println(err)
	}
	fmt.Println("duration:", time.Since(startTime))
	// fmt.Println("library: ", viper.GetString("library"))
	//rcmd.InstallThroughBinary(appFS, "", ia, rcmd.RSettings{}, rcmd.ExecSettings{}, log)
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
