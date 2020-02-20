package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Checks that installed packages can be loaded",
	Long: `Attempts to load user packages specified in pkgr.yml to validate that each package has been installed
successfully and can be used. Use the --all flag to load all packages in the user-library dependency tree instead of just user-level packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		all := viper.GetBool("all")
		json := viper.GetBool("json")

		rs := rcmd.NewRSettings(cfg.RPath)

		// Get directory to test package loads from.
		// By default, use the same directory as the config file.
		rDir := viper.GetString("config")
		if rDir == "" {
			rDir, _ = os.Getwd()
		} else {
			rDir = filepath.Dir(rDir)
		}
		rDir, err := filepath.Abs(rDir)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"directory": rDir,
			}).Fatal("error getting absolute path for R directory")
		}

		load(rs, rDir, all, json)
	},
}

func init() {
	RootCmd.AddCommand(loadCmd)
	loadCmd.Flags().Bool("all", false, "load user packages as well as their dependencies")
	viper.BindPFlag("all", loadCmd.LocalFlags().Lookup("all")) //There doesn't seem to be a way to bind local flags.
	loadCmd.Flags().Bool("json", false, "output a JSON object of package info at the end")
	viper.BindPFlag("json", loadCmd.LocalFlags().Lookup("json")) //There doesn't seem to be a way to bind local flags.
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func load(rs rcmd.RSettings, rDir string, all, json bool) {
	log.WithFields(log.Fields{
		"all_arg" : all,
		"json_arg" : json,
	}).Info("attempting to load packages")

	log.Info("getting relevant packages via `pkgr plan`..................")

	rVersion := rcmd.GetRVersion(&rs)
	_, installPlan, _ := planInstall(rVersion, false)

	log.WithFields(log.Fields{
		"all_arg" : all,
		"json_arg" : json,
	}).Info("attempting to load packages")

	log.Info("finished getting packages from `pkgr plan`__________________")

	var toLoad []string
	if all {
		toLoad = installPlan.GetAllPackages()
	} else {
		toLoad = cfg.Packages
	}

	report := MakeLoadReport(getRSessionMetadata(rs, rDir))

	for _, pkg := range toLoad {
		lr := attemptLoad(rs, rDir, pkg)
		report.AddResult(pkg, lr)
	}

	if report.failures == 0 {
		log.WithFields(log.Fields{
			"working_directory": rDir,
			"user_packages_attempted": "true",
			"dependencies_attempted": all,
		}).Info("packages loaded successfully")
	}

	if json {
		//getRSessionMetadata(rs, rDir)
		log.Info("printing load report as JSON")
		logLoadReport(report)
	}

	//Add failure condition.
}

func logLoadReport(rpt loadReport) {
	//var rptPkgs []string
	//for key := range rpt.results {
	//	rptPkgs = append(rptPkgs, key)
	//}
	jsonObj, err := json.Marshal(rpt)
	if err != nil {
		log.WithFields(log.Fields{"err" : err}).Error("encountered problem marshalling load report to JSON")
		return
	}
	log.Infof("%s \n", jsonObj)
	log.Info("done printing load report")
	//fmt.Printf("%s \n", jsonObj)
}

func getRSessionMetadata(rs rcmd.RSettings, rDir string) rSessionMetadata {
	sessionMetadata := rSessionMetadata{
		rPath : rs.Rpath,
		rVersion : fmt.Sprintf("%d.%d.%d", rs.Version.Major, rs.Version.Minor, rs.Version.Patch),
	}


	// We need to get the libPaths that the session uses, not necessarily the ones pkgr would set.
	libPathsCmd := fmt.Sprintf(".libPaths()")

	cmdArgs := []string{
		"-e",
		libPathsCmd,
	}
	cmd := exec.Command(rs.R(runtime.GOOS), cmdArgs...)
	cmd.Dir = rDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.WithFields(log.Fields{"pkg" : pkg, "rDir": rDir,}).Trace("attempting to get R libpaths")
	cmdErr := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if cmdErr != nil || errStr != "" {
		log.WithFields(log.Fields{
			"session_work_dir" : rDir,
			"err" : cmdErr,
		}).Warn("could not get libPaths -- there was a problem running `.libPaths()` in an R session")
		sessionMetadata.libPaths = []string {"could not retrieve libpaths"}
		return sessionMetadata
	}

	sessionMetadata.libPaths = strings.Split(outStr, "\n")
	return sessionMetadata
}



func attemptLoad(rs rcmd.RSettings, rDir, pkg string) loadResult {

	libraryCmd := fmt.Sprintf("library('%s')", pkg)

	cmdArgs := []string{
		//"--vanilla",
		"-e",
		libraryCmd,
	}
	cmd := exec.Command(rs.R(runtime.GOOS), cmdArgs...)
	cmd.Dir = rDir
	//cmd.Env = envVars

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	//cmd.Stdin = os.Stdin

	var exitError *exec.ExitError

	log.WithFields(log.Fields{"pkg" : pkg, "rDir": rDir,}).Trace("attempting to load package.")
	cmdErr := cmd.Run()

	var didSucceed bool

	if cmdErr != nil {
		if errors.As(cmdErr, &exitError) { // If the command had an exit code != 0
			didSucceed = false
		} else {
			log.WithFields(log.Fields{
				"cmd": rs.R(runtime.GOOS),
				"rDir": rDir,
				"pkg": pkg,
			}).Fatal("could not execute R 'library' call for package") //TODO: revisit
		}
	}


	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if errStr == "" {
		didSucceed = true
	} else {
		didSucceed = false
	}

	if didSucceed {
		log.WithFields(log.Fields{
			"pkg": pkg,
			"rDir": rDir,
		}).Debug("Package loaded successfully")
		log.WithFields(log.Fields{
			"pkg": pkg,
			"rDir": rDir,
			"stdOut": string(stdout.Bytes()),
		}).Trace("Package loaded successfully")
	} else {
		log.WithFields(log.Fields{
			"goErr" : cmdErr,
			"stdErr": errStr,
			"pkg": pkg,
			"rDir": rDir,
		}).Error("error loading package via `library(<pkg>)` in R")
	}

	additionalInfo := getAdditionalPkgInfo(rs, rDir, pkg)
	return MakeLoadResult(pkg, additionalInfo.pkgVersion, additionalInfo.pkgPath, outStr, errStr, didSucceed, cmdErr)

	//return MakeLoadResult(outStr, errStr, didSucceed)

}

// Get the path and version for a given package, assuming that that package can be loaded.
func getAdditionalPkgInfo(rs rcmd.RSettings, rDir, pkg string) pkgLoadMetadata {
	infoCmd := fmt.Sprintf("find.package('%s'); packageVersion('%s')", pkg, pkg)

	cmdArgs := []string{
		//"--vanilla",
		"-e",
		infoCmd,
	}
	cmd := exec.Command(rs.R(runtime.GOOS), cmdArgs...)
	cmd.Dir = rDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.WithFields(log.Fields{"pkg" : pkg, "rDir": rDir,}).Trace("attempting to load package.")
	cmdErr := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if cmdErr != nil || errStr != "" {
		log.WithFields(log.Fields{
			"pkg" : pkg,
			"rDir" : rDir,
			"err" : cmdErr,
		}).Warn("could not get package path and info data during load")
		return pkgLoadMetadata{"could not retrieve",	"could not retrieve",}
	}

	outStrSplit := strings.Split(outStr, "\n")
	return pkgLoadMetadata{
		outStrSplit[0], //path
		outStrSplit[1], //version
	}
}



//// Load report struct
type loadReport struct {
	rMetadata rSessionMetadata
	results map[string]loadResult
	failures int
}

func MakeLoadReport(rMetadata rSessionMetadata) loadReport {
	return loadReport {
		rMetadata : rMetadata,
		results : make(map[string]loadResult),
		failures : 0,
	}
}

func (report *loadReport) AddResult(pkg string, result loadResult) {
	report.results[pkg] = result
	if !result.success {
		report.failures = report.failures + 1
	}
}


//// Load result struct
type loadResult struct {
	pkg string
	version string
	path string
	exiterr error // could be equivalent to exit code.
	stdout string
	stderr string
	success bool
	// Can store information for JSON here
}

type pkgLoadMetadata struct { // Used to help create loadResult
	pkgPath string
	pkgVersion string
}

func MakeLoadResult(pkg, version, path, outStr, errStr string, success bool, exitErr error) loadResult {
	return loadResult {
		pkg : pkg,
		version : version,
		path : path,
		exiterr : exitErr,
		stdout : outStr,
		stderr : errStr,
		success : success,
	}
}

//// R Session Info Struct
type rSessionMetadata struct {
	libPaths []string
	rPath string
	rVersion string
}


//// Constructor for load result
//func MakeLoadResult(outStr, errStr string, success bool) loadResult {
//	lr := loadResult{
//		stdout : outStr,
//		stderr : errStr,
//		success : success,
//	}
//	//lr.updateResult()
//	return lr
//}



