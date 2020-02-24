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
		r_dir := viper.GetString("config")
		if r_dir == "" {
			r_dir, _ = os.Getwd()
		} else {
			r_dir = filepath.Dir(r_dir)
		}
		r_dir, err := filepath.Abs(r_dir)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"directory": r_dir,
			}).Fatal("error getting absolute Path for R directory")
		}

		load(rs, r_dir, all, json)
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

func load(rs rcmd.RSettings, r_dir string, all, json bool) {
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

	report := MakeLoadReport(getRSessionMetadata(rs, r_dir))

	for _, pkg := range toLoad {
		lr := attemptLoad(rs, r_dir, pkg)
		report.AddResult(pkg, lr)
	}

	if report.Failures == 0 {
		log.WithFields(log.Fields{
			"working_directory": r_dir,
			"user_packages_attempted": "true",
			"dependencies_attempted": all,
		}).Info("packages loaded successfully")
	}

	if json {
		//getRSessionMetadata(rs, r_dir)
		log.Info("printing load report as JSON")
		logLoadReport(report)
	}

	//Add failure condition.
}

func logLoadReport(rpt loadReport) {
	jsonObj, err := json.MarshalIndent(rpt, "", "    ")
	if err != nil {
		log.WithFields(log.Fields{"err" : err}).Error("encountered problem marshalling load report to JSON")
		return
	}
	log.Infof("%s \n", jsonObj)
	log.Info("done printing load report")
	//fmt.Printf("%s \n", jsonObj)
}

func getRSessionMetadata(rs rcmd.RSettings, r_dir string) rSessionMetadata {
	sessionMetadata := rSessionMetadata{
		RPath:    rs.Rpath,
		RVersion: fmt.Sprintf("%d.%d.%d", rs.Version.Major, rs.Version.Minor, rs.Version.Patch),
	}

	sessionMetadata.LibPaths = getRSessionLibPaths(rs, r_dir)

	return sessionMetadata
}

func getRSessionLibPaths(rs rcmd.RSettings, rDir string) []string {
	// We need to get the LibPaths that the session uses, not necessarily the ones pkgr would set.
	libPathsCmd := fmt.Sprintf(".libPaths()")

	cmdArgs := []string{
		"-q",
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
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	outLines := rp.ScanROutput(stdout.Bytes(), false)
	errLines := rp.ScanROutput(stderr.Bytes(), false)

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

func attemptLoad(rs rcmd.RSettings, rDir, pkg string) loadResult {
	libraryCmd := fmt.Sprintf("library('%s')", pkg)

	cmdArgs := []string{
		"-q",
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

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	outLines := rp.ScanLines(stdout.Bytes())
	errLines := rp.ScanLines(stderr.Bytes())

	if cmdErr != nil {
		if errors.As(cmdErr, &exitError) { // If the command had an exit code != 0
			didSucceed = false
		} else {
			log.WithFields(log.Fields{
				"cmd":           rs.R(runtime.GOOS),
				"r_dir":         rDir,
				"pkg":           pkg,
				"command_error": cmdErr,
			}).Fatal("could not execute R 'library' call for package") //TODO: revisit
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
			"std_out": string(stdout.Bytes()),
		}).Trace("Package loaded successfully")
	}

	additionalInfo := getAdditionalPkgInfo(rs, rDir, pkg)
	return MakeLoadResult(pkg, additionalInfo.pkgVersion, additionalInfo.pkgPath, strings.Join(outLines, "\n"), strings.Join(errLines, "\n"), didSucceed, cmdErr)

}

// Get the Path and Version for a given package, assuming that that package can be loaded.
func getAdditionalPkgInfo(rs rcmd.RSettings, rDir, pkg string) pkgLoadMetadata {
	infoCmd := fmt.Sprintf("find.package('%s'); packageVersion('%s')", pkg, pkg)

	cmdArgs := []string{
		"-q",
		"-e",
		infoCmd,
	}
	cmd := exec.Command(rs.R(runtime.GOOS), cmdArgs...)
	cmd.Dir = rDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.WithFields(log.Fields{"pkg" : pkg, "r_dir": rDir,}).Trace("attempting to load package.")
	cmdErr := cmd.Run()
	outLines := rp.ScanROutput(stdout.Bytes(), false)
	errLines := rp.ScanROutput(stderr.Bytes(), false)

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
		}).Warn("could not parse R command output for package Path and Version info")
		return pkgLoadMetadata{"could not retrieve",	"could not retrieve",}
	}

	pkgPath := outLines[0]
	pkgVersion := outLines[1]

	return pkgLoadMetadata{
		pkgPath,
		pkgVersion,
	}
}


//// Load report struct
type loadReport struct {
	RMetadata   rSessionMetadata
	LoadResults map[string]loadResult
	Failures    int
}

func MakeLoadReport(rMetadata rSessionMetadata) loadReport {
	return loadReport {
		RMetadata:   rMetadata,
		LoadResults: make(map[string]loadResult),
		Failures:    0,
	}
}

func (report *loadReport) AddResult(pkg string, result loadResult) {
	report.LoadResults[pkg] = result
	if !result.Success {
		report.Failures = report.Failures + 1
	}
}


//// Load result struct
type loadResult struct {
	Package string
	Version string
	Path    string
	Exiterr string // could be equivalent to exit code.
	Stdout  string
	Stderr  string
	Success bool
	// Can store information for JSON here
}

type pkgLoadMetadata struct { // Used to help create loadResult
	pkgPath string
	pkgVersion string
}

func MakeLoadResult(pkg, version, path, outStr, errStr string, success bool, exitErr error) loadResult {
	exitErrString := ""
	if exitErr != nil {
		exitErrString = exitErr.Error()
	}

	return loadResult {
		Package: pkg,
		Version: version,
		Path:    path,
		Exiterr: exitErrString,
		Stdout:  outStr,
		Stderr:  errStr,
		Success: success,
	}
}

//// R Session Info Struct
type rSessionMetadata struct {
	LibPaths []string
	RPath    string
	RVersion string
}


