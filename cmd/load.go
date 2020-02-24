package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/metrumresearchgroup/pkgr/rcmd/rp"
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
				"error":     err,
				"directory": rDir,
			}).Fatal("error getting absolute Path for R directory")
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

	report := InitLoadReport(getRSessionMetadata(rs, rDir))

	//var waitGroup sync.WaitGroup
	resultsChannel := make(chan LoadResult, cfg.Threads * 2)

	//goAttemptLoad()

	go func() {
		for _, pkg := range toLoad {
			//loadResult := attemptLoad(rs, rDir, pkg)
			//report.AddResult(pkg, loadResult)
			//waitGroup.Add(1)
			go goAttemptLoad2(rs, rDir, pkg, resultsChannel)//, waitGroup)

		}
		//close(resultsChannel)
	}()
	for r := range resultsChannel {
		report.AddResult(r.Package, r)
	}

	//for _, pkg := range toLoad {
	//	//loadResult := attemptLoad(rs, rDir, pkg)
	//	//report.AddResult(pkg, loadResult)
	//	//waitGroup.Add(1)
	//	go goAttemptLoad(rs, rDir, pkg, resultsChannel)//, waitGroup)
	//
	//}

	//waitGroup.Wait()

	if report.Failures == 0 {
		log.WithFields(log.Fields{
			"working_directory":       rDir,
			"user_packages_attempted": "true",
			"dependencies_attempted":  all,
		}).Info("packages loaded successfully")
	}

	if json {
		//getRSessionMetadata(rs, rDir)
		log.Info("printing load report as JSON")
		logLoadReport(report)
	}

	//Add failure condition.
}

func getRSessionMetadata(rs rcmd.RSettings, r_dir string) RSessionMetadata {
	sessionMetadata := RSessionMetadata{
		RPath:    rs.Rpath,
		RVersion: fmt.Sprintf("%d.%d.%d", rs.Version.Major, rs.Version.Minor, rs.Version.Patch),
	}
	sessionMetadata.LibPaths = getRSessionLibPaths(rs, r_dir)

	return sessionMetadata
}

func getRSessionLibPaths(rs rcmd.RSettings, rDir string) []string {

	outLines, errLines, cmdErr := runRCmd(".libPaths()", rs, rDir, true)

	if cmdErr != nil || len(errLines) > 0 {
		log.WithFields(log.Fields{
			"sessionWorkingDir" : rDir,
			"cmdErr" :            cmdErr,
			"std_erro":           errLines,
		}).Warn("could not get LibPaths -- there was a problem running `.LibPaths()` in an R session")
		return []string {"could not retrieve libpaths"}
	}

	return outLines
}

func logLoadReport(rpt LoadReport) {
	jsonObj, err := json.MarshalIndent(rpt, "", "    ")
	if err != nil {
		log.WithFields(log.Fields{"err" : err}).Error("encountered problem marshalling load report to JSON")
		return
	}
	log.Infof("%s \n", jsonObj)
	log.Info("done printing load report")
	//fmt.Printf("%s \n", jsonObj)
}

func goAttemptLoad(rs rcmd.RSettings, rDir string, pkgs []string, c chan<- LoadResult) {
	//defer wg.Done()
	for _, pkg := range pkgs {
		c <- attemptLoad(rs, rDir, pkg)
	}
	//result := attemptLoad(rs, rDir, pkg)
	//c <- result
}

func goAttemptLoad2(rs rcmd.RSettings, rDir string, pkg string, c chan<- LoadResult) {
	//defer wg.Done()
	result := attemptLoad(rs, rDir, pkg)
	c <- result
}

func attemptLoad(rs rcmd.RSettings, rDir, pkg string) LoadResult {
	log.WithFields(log.Fields{"pkg": pkg, "rDir": rDir,}).Trace("attempting to load package.")
	outLines, errLines, cmdErr := runRCmd(fmt.Sprintf("library('%s')", pkg), rs, rDir, false)

	var exitError *exec.ExitError
	var didSucceed bool

	if cmdErr != nil {
		didSucceed = false
		if errors.As(cmdErr, &exitError) { // If the command had an exit code != 0
			log.WithFields(log.Fields{
				"go_error" : cmdErr,
				"std_error":  errLines,
				"pkg":       pkg,
				"rDir":     rDir,
			}).Errorf("error loading package via `library(%s)` in R -- received non-zero exit code", pkg)
		} else {
			log.WithFields(log.Fields{
				"cmd":           rs.R(runtime.GOOS),
				"r_dir":         rDir,
				"pkg":           pkg,
				"command_error": cmdErr,
			}).Fatal("could not execute R 'library' call for package") //TODO: revisit -- can we consider this a fatal error?
		}
	} else if len(errLines) > 0 { // Should be impossible, as any errors in stderr should cause exit code to be non-zero.
		didSucceed = false
		log.WithFields(log.Fields{
			"go_error" : cmdErr,
			"std_erro":  errLines,
			"pkg":       pkg,
			"rDir":     rDir,
		}).Errorf("error loading package via `library(%s)` in R", pkg)
	} else {
		didSucceed = true
		log.WithFields(log.Fields{
			"pkg":   pkg,
			"r_dir": rDir,
		}).Debug("Package loaded successfully")
		log.WithFields(log.Fields{
			"pkg":    pkg,
			"r_dir":  rDir,
			"std_out": outLines,
		}).Trace("Package loaded successfully")
	}

	additionalInfo := getAdditionalPkgInfo(rs, rDir, pkg)
	return MakeLoadResult(
		pkg,
		additionalInfo.pkgVersion,
		additionalInfo.pkgPath,
		strings.Join(outLines, "\n"),
		strings.Join(errLines, "\n"),
		didSucceed,
		cmdErr,
	)

}

// Get the Path and Version for a given package, assuming that that package can be loaded.
func getAdditionalPkgInfo(rs rcmd.RSettings, rDir, pkg string) pkgLoadMetadata {
	outLines, errLines, cmdErr := runRCmd(fmt.Sprintf("find.package('%s'); packageVersion('%s')", pkg, pkg), rs, rDir, true)

	if cmdErr != nil || len(errLines) > 0 {
		log.WithFields(log.Fields{
			"pkg" :   pkg,
			"r_dir" : rDir,
			"err" :   cmdErr,
		}).Warn("could not get package Path and Version info data during load")
		return pkgLoadMetadata{"could not retrieve",	"could not retrieve",}
	} else if len(outLines) != 2 {
		log.WithFields(log.Fields{
			"pkg" :    pkg,
			"r_dir" :  rDir,
			"output" : outLines,
		}).Warn("could not parse R command output for package Path and Version info -- expected exactly two lines of output.")
		return pkgLoadMetadata{"could not retrieve",	"could not retrieve",}
	}

	pkgPath := outLines[0]
	pkgVersion := outLines[1]

	return pkgLoadMetadata{
		pkgPath,
		pkgVersion,
	}
}

//// If we can find it, return the exit code and a bool indicating whether or not it was parsed.
//// If there's an error but we can't parse an exit code, return a non-zero exit code and "false".
//// If there's no error, return an exit code of 0 and "false", just to clarify that it wasn't parsed.
//func ParseExitCode(cmdErr error) (int, bool) {
//	if cmdErr != nil {
//		if exitError, ok := cmdErr.(*exec.ExitError); ok {
//			return exitError.ExitCode(), true
//		} else {
//			return -999, false
//		}
//	}
//	return 0, false // Assume exit code zero if cmdErr is nil.
//}

func runRCmd(rExpression string, rs rcmd.RSettings, rDir string, reducedOutput bool) ([]string, []string, error) {
	cmdArgs := []string{
		"-q",
		"-e",
		rExpression,
	}
	cmd := exec.Command(rs.R(runtime.GOOS), cmdArgs...)
	cmd.Dir = rDir
	//cmd.Env = envVars
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	//cmd.Stdin = os.Stdin

	log.WithFields(log.Fields{

	})
	cmdErr := cmd.Run()

	//exitCode, isAuthenticExitCode := ParseExitCode(cmdErr)
	//if !isAuthenticExitCode {
	//	log.WithFields(log.Fields{
	//		"r_dir": rDir,
	//		"cmd_error": cmdErr,
	//		"exit_code": exitCode,
	//		"r_code": rExpression,
	//	}).Warn("exit code for R expression could not be determined -- the exit code provided is our best guess.")
	//}


	outLines := rp.ScanROutput(stdout.Bytes(), reducedOutput)
	errLines := rp.ScanLines(stderr.Bytes())
	return outLines, errLines, cmdErr
}



