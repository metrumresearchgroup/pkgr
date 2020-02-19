package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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

	rVersion := rcmd.GetRVersion(&rs)
	_, installPlan, _ := planInstall(rVersion, false)

	var toLoad []string
	if all {
		toLoad = installPlan.GetAllPackages()
	} else {
		toLoad = cfg.Packages
	}

	var report loadReport

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

	//Add failure condition.
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
	err := cmd.Run()

	var didSucceed bool

	if err != nil {
		if errors.As(err, &exitError) { // If the command had an exit code != 0
			didSucceed = false
		} else {
			log.WithFields(log.Fields{
				"cmd": rs.R(runtime.GOOS),
				"rDir": rDir,
				"pkg": pkg,
			}).Fatal("could not execute R 'library' call for package") //TODO: revist
		}
	}


	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if errStr == "" {
		didSucceed = false
	} else {
		didSucceed = true
	}

	if didSucceed {
		log.WithFields(log.Fields{
			"pkg": pkg,
			"rDir": rDir,
		}).Debug("Package loaded successfully")
		log.WithFields(log.Fields{
			"pkg": pkg,
			"rDir": rDir,
			"stdOut": stdout,
		}).Trace("Package loaded successfully")
	} else {
		log.WithFields(log.Fields{
			"goErr" : err,
			"stdErr": errStr,
			"pkg": pkg,
			"rDir": rDir,
		}).Error("error loading package via `library(<pkg>)` in R")
	}

	return MakeLoadResult(outStr, errStr, didSucceed)

}

//// Load report struct
type loadReport struct {
	results map[string]loadResult
	failures int
}

func (report *loadReport) AddResult(pkg string, result loadResult) {
	report.results[pkg] = result
	if !result.success {
		report.failures = report.failures + 1
	}
}


//// Load result struct
type loadResult struct {
	stdout string
	stderr string
	success bool
	// Can store information for JSON here
}

//func (result *loadResult) updateResult() {
//	if result.stderr == "" {
//		result.success = true
//	} else {
//		result.success = false
//	}
//}

// Constructor for load result
func MakeLoadResult(outStr, errStr string, success bool) loadResult {
	lr := loadResult{
		stdout : outStr,
		stderr : errStr,
		success : success,
	}
	//lr.updateResult()
	return lr
}



