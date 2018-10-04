package rcmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const defaultFailedCode = 1
const defaultSuccessCode = 0

// RunR launches an interactive R console
func RunR(
	fs afero.Fs,
	rs RSettings,
	rdir string, // this should be put into RSettings
	lg *logrus.Logger,
) error {

	cmdArgs := []string{
		"--no-save",
		"--no-restore-data",
	}

	envVars := os.Environ()
	ok, rLibsSite := rs.LibPathsEnv()
	if ok {
		envVars = append(envVars, rLibsSite, "R_LIBS=''")
	}

	lg.WithFields(
		logrus.Fields{
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       rLibsSite,
		}).Debug("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		rs.R(),
		cmdArgs...,
	)

	if rdir == "" {
		rdir, _ = os.Getwd()
		lg.WithFields(
			logrus.Fields{"rdir": rdir},
		).Debug("launch dir")
	}
	cmd.Dir = rdir
	cmd.Env = envVars
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Install installs a given tarball
// exit code 0 - success, 1 - error
func Install(
	fs afero.Fs,
	tbp string, // tarball path
	args InstallArgs,
	rs RSettings,
	es ExecSettings,
	lg *logrus.Logger,
) (string, error, int) {

	cmdArgs := []string{
		"CMD",
		"install",
	}
	envVars := os.Environ()

	lg.WithFields(
		logrus.Fields{
			"cmd":       "install",
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       envVars,
		}).Debug("command args")
	lg.WithFields(
		logrus.Fields{
			"cmd":  "install",
			"exec": es,
		}).Debug("execution")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		fmt.Sprintf("%s", rs.R()),
		cmdArgs...,
	)

	rdir := es.WorkDir
	if rdir == "" {
		rdir, _ = os.Getwd()
		lg.WithFields(
			logrus.Fields{"rdir": rdir},
		).Debug("launch dir")
	}
	cmd.Dir = rdir
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	exitCode := defaultSuccessCode
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	lg.WithFields(
		logrus.Fields{
			"stdout":   stdout,
			"stderr":   stderr,
			"exitCode": exitCode,
		}).Info("cmd output")
	return stdout, err, exitCode
}
