package rcmd

import (
	"os"
	"os/exec"

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
