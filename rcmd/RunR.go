package rcmd

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const defaultFailedCode = 1
const defaultSuccessCode = 0

// StartR launches an interactive R console given the same
// configuration as a specific package.
func StartR(
	fs afero.Fs,
	pkg string,
	rs RSettings,
	rdir string, // this should be put into RSettings
) error {

	envVars := configureEnv(os.Environ(), rs, pkg)
	cmdArgs := []string{
		"--vanilla",
	}

	log.WithFields(
		log.Fields{
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       envVars,
		}).Trace("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		rs.R(runtime.GOOS),
		cmdArgs...,
	)

	if rdir == "" {
		rdir, _ = os.Getwd()
		log.WithFields(
			log.Fields{"rdir": rdir},
		).Debug("launch dir")
	}
	cmd.Dir = rdir
	cmd.Env = envVars
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// convertToRscript converts an R path like "/path/to/R" to
// "/path/to/Rscript", accounting for a trailing ".exe", if present.
func convertToRscript(rpath string) string {
	var script string

	if strings.HasSuffix(rpath, ".exe") {
		script = rpath[:len(rpath)-4] + "script" + ".exe"
	} else {
		script = rpath + "script"
	}

	return script
}

// RunR launches an interactive R console
func RunR(
	fs afero.Fs,
	pkg string,
	rs RSettings,
	script string,
	rdir string, // this should be put into RSettings
) ([]byte, error) {

	cmdArgs := []string{
		"--vanilla",
		"-e",
		script,
	}

	return RunRscriptWithArgs(fs, pkg, rs, cmdArgs, rdir)
}

// RunRscriptWithArgs invokes Rscript with cmdArgs, returning
// os.exec.Cmd.Output().
func RunRscriptWithArgs(
	fs afero.Fs,
	pkg string,
	rs RSettings,
	cmdArgs []string,
	rdir string,
) ([]byte, error) {
	envVars := configureEnv(os.Environ(), rs, pkg)
	log.WithFields(
		log.Fields{
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       envVars,
		}).Trace("command args")

	prog := convertToRscript(rs.R(runtime.GOOS))
	cmd := exec.Command(prog, cmdArgs...)

	if rdir == "" {
		rdir, _ = os.Getwd()
		log.WithFields(
			log.Fields{"rdir": rdir},
		).Debug("launch dir")
	}
	cmd.Dir = rdir
	cmd.Env = envVars

	return cmd.Output()
}

// RunRBatch runs a non-interactive R command
func RunRBatch(
	fs afero.Fs,
	rs RSettings,
	cmdArgs []string,
) ([]byte, error) {
	envVars := configureEnv(os.Environ(), rs, "")
	rpath := rs.R(runtime.GOOS)
	log.WithFields(
		log.Fields{
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"rpath":     rpath,
		}).Trace("command args")

	cmd := exec.Command(
		rpath,
		cmdArgs...,
	)
	cmd.Env = envVars

	return cmd.CombinedOutput()
}
