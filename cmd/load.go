package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/metrumresearchgroup/pkgr/rcmd/rp"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Check that installed packages can be loaded",
	Long: `Load packages specified in the configuration file to validate that
each package has been installed successfully and can be used.

**Execution environment**. This subcommand runs R with the same settings
that R would use if you invoked 'R' from the current working directory. It
relies on that environment being configured to find packages in the library
path specified in the configuration file (via 'Library' or 'Lockfile:
Type').  Pass the --json argument to confirm that the package is being
loaded from the expected library.`,
	Example: `  # Load packages listed in config file
  pkgr load --json
  # Load the above packages and all their dependencies
  pkgr load --json --all`,
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

		load(cfg.Packages, rs, rDir, cfg.Threads, all, json)
	},
}

func init() {
	RootCmd.AddCommand(loadCmd)
	loadCmd.Flags().Bool("all", false, "load all packages in dependency tree")
	viper.BindPFlag("all", loadCmd.LocalFlags().Lookup("all")) //There doesn't seem to be a way to bind local flags.
	loadCmd.Flags().Bool("json", false, "output results as a JSON object")
	viper.BindPFlag("json", loadCmd.LocalFlags().Lookup("json")) //There doesn't seem to be a way to bind local flags.
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Attempts to load all of the packages in userPackages, and (optionally) the dependencies of those packages, in an R
// session launched from rDir.
func load(userPackages []string, rs rcmd.RSettings, rDir string, threads int, all, toJson bool) {
	if toJson {
		logger.SetLogLevel("fatal") // disable logging
	}

	log.WithFields(log.Fields{
		"all_arg":  all,
		"json_arg": toJson,
	}).Info("attempting to load packages")

	log.Info("getting relevant packages via `pkgr plan`..................")

	rVersion := rcmd.GetRVersion(&rs)
	_, installPlan, _ := planInstall(rVersion, false)

	log.Info("finished getting packages from `pkgr plan`__________________")

	var toLoad []string
	if all {
		toLoad = installPlan.GetAllPackages()
	} else {
		toLoad = userPackages
	}

	report := InitLoadReport(getRSessionMetadata(rs, rDir))

	resultsChannel := make(chan LoadResult, len(toLoad))
	sem := make(chan int, threads*2)

	log.WithFields(
		log.Fields{
			"num_to_load": len(toLoad),
			"threads":     threads * 2,
			"all_arg":     all,
			"json_arg":    toJson,
		}).Info("attempting to load packages")

	// Kick off every load request as a goroutine, each of which will wait on the availability of a semaphore wait group.
	for _, pkg := range toLoad {
		go attemptLoadConcurrent(
			LoadRequest{
				rs,
				rDir,
				pkg,
			},
			sem,
			resultsChannel,
		)
	}

	for i := 0; i < len(toLoad); i++ {
		lr := <-resultsChannel
		log.WithFields(
			log.Fields{
				"pkg":       lr.Package,
				"succeeded": lr.Success,
			}).Trace("completed load attempt and adding to load report.")
		report.AddResult(lr.Package, lr)
	}
	close(resultsChannel)
	close(sem)

	if report.Failures == 0 {
		log.WithFields(log.Fields{
			"working_directory":       rDir,
			"user_packages_attempted": "true",
			"dependencies_attempted":  all,
		}).Info("all packages loaded successfully")
	} else {
		log.WithFields(log.Fields{
			"working_directory": rDir,
			"failures":          report.Failures,
		}).Error("some packages failed to load.")
	}

	if toJson {
		printJsonLoadReport(report)
	}
}

// Get top-level information about an R session launched from rDir.
func getRSessionMetadata(rs rcmd.RSettings, rDir string) RSessionMetadata {
	sessionMetadata := RSessionMetadata{
		RPath:    rs.Rpath,
		RVersion: fmt.Sprintf("%d.%d.%d", rs.Version.Major, rs.Version.Minor, rs.Version.Patch),
	}

	sessionMetadata.LibPaths = getRSessionLibPaths(rs, rDir)

	return sessionMetadata
}

// Get the libPaths that will load for an R session launched in rDir.
func getRSessionLibPaths(rs rcmd.RSettings, rDir string) []string {
	log.Trace("attempting to get session lib paths")
	outLines, errLines, cmdErr := runRCmd(".libPaths()", rs, rDir, true)

	if cmdErr != nil || len(errLines) > 0 {
		log.WithFields(log.Fields{
			"session_working_dir": rDir,
			"cmd_err":             cmdErr,
			"std_error":           errLines,
		}).Warn("could not get LibPaths -- there was a problem running `.LibPaths()` in an R session")
		return []string{"could not retrieve libpaths"}
	}
	var trimmedOutLines []string
	for _, line := range outLines {
		trimmedOutLines = append(trimmedOutLines, strings.ReplaceAll(line, "\"", ""))
	}

	return trimmedOutLines
}

// Print a LoadReport as a JSON object.
func printJsonLoadReport(rpt LoadReport) {
	jsonObj, err := JsonMarshal(rpt)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("encountered problem marshalling load report to JSON")
		return
	}
	fmt.Printf("%s \n", jsonObj)
}

// Worker function to perform work specified in LoadRequest.
// Multiple workers may be launched concurrently to parallelize the process.
func attemptLoadConcurrent(request LoadRequest, sem chan int, out chan LoadResult) {
	sem <- 1 // write to semaphore, which is capped at threads*2. If semaphore is full, block until you can write.
	out <- attemptLoad(request.Rs, request.RDir, request.Pkg)
	<-sem // read from sempaphore to indicate that you are done running, thus opening the "slot" taken earlier.
}

// Try to load the given R package in the rDir specified. Bundle the results in a LoadResults object and return.
func attemptLoad(rs rcmd.RSettings, rDir, pkg string) LoadResult {
	log.WithFields(log.Fields{"pkg": pkg, "rDir": rDir}).Trace("attempting to load package.")
	outLines, errLines, cmdErr := runRCmd(fmt.Sprintf("library('%s')", pkg), rs, rDir, false)

	var exitError *exec.ExitError
	var didSucceed bool

	if cmdErr != nil {
		didSucceed = false
		if errors.As(cmdErr, &exitError) { // If the command had an exit code != 0
			log.WithFields(log.Fields{
				"go_error":  cmdErr,
				"std_error": errLines,
				"pkg":       pkg,
				"rDir":      rDir,
			}).Errorf("error loading package via `library(%s)` in R -- received non-zero exit code", pkg)
		} else {
			log.WithFields(log.Fields{
				"cmd":           rs.R(runtime.GOOS),
				"r_dir":         rDir,
				"pkg":           pkg,
				"command_error": cmdErr,
			}).Fatal("could not execute R 'library' call for package")
		}
	} else {
		didSucceed = true
		log.WithFields(log.Fields{
			"pkg":   pkg,
			"r_dir": rDir,
		}).Debug("Package loaded successfully")
		log.WithFields(log.Fields{
			"pkg":     pkg,
			"r_dir":   rDir,
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
	outLines, errLines, cmdErr := runRCmd(fmt.Sprintf("find.package('%s'); as.character(packageVersion('%s'))", pkg, pkg), rs, rDir, true)

	if cmdErr != nil { // || len(errLines) > 0 {
		log.WithFields(log.Fields{
			"pkg":    pkg,
			"r_dir":  rDir,
			"err":    cmdErr,
			"stdout": outLines,
			"stderr": errLines,
		}).Warn("could not get package Path and Version info data during load")
		return pkgLoadMetadata{"could not retrieve", "could not retrieve"}
	} else if len(outLines) != 2 {
		log.WithFields(log.Fields{
			"pkg":    pkg,
			"r_dir":  rDir,
			"output": outLines,
			"stderr": errLines,
		}).Warn("could not parse R command output for package Path and Version info -- expected exactly two lines of output.")
		return pkgLoadMetadata{"could not retrieve", "could not retrieve"}
	}

	pkgPath := strings.ReplaceAll(outLines[0], "\"", "")
	pkgVersion := strings.ReplaceAll(outLines[1], "\"", "")

	return pkgLoadMetadata{
		pkgPath,
		pkgVersion,
	}
}

// Helper function to run R commands and gather the outputs, using the set of arguments that `pkgr load` internal
// operations require.
func runRCmd(rExpression string, rs rcmd.RSettings, rDir string, reducedOutput bool) ([]string, []string, error) {
	cmdArgs := []string{
		"--slave",
		"--no-restore",
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

	log.WithFields(log.Fields{})
	cmdErr := cmd.Run()

	outLines := rp.ScanROutput(stdout.Bytes(), reducedOutput)
	errLines := rp.ScanLines(stderr.Bytes())

	return outLines, errLines, cmdErr
}
