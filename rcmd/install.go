package rcmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"

	"github.com/fatih/structs"
	"github.com/fatih/structtag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
)

// NewDefaultInstallArgs provides a set of sane default installation args
func NewDefaultInstallArgs() *InstallArgs {
	return &InstallArgs{
		WithKeepSource: true,
		NoMultiarch:    true,
		InstallTests:   true,
		Build:          true,
	}
}

// CliArgs converts the InstallArgs struct to the proper cli args
// including only returning the relevant args
func (i *InstallArgs) CliArgs() []string {
	args := []string{}
	is := structs.New(i)
	nms := structs.Names(i)
	for _, n := range nms {
		fld := is.Field(n)
		if !fld.IsZero() {
			// ... and start using structtag by parsing the tag
			tag, _ := reflect.TypeOf(i).Elem().FieldByName(fld.Name())
			// ... and start using structtag by parsing the tag
			tags, err := structtag.Parse(string(tag.Tag))
			if err != nil {
				panic(err)
			}
			rcmd, err := tags.Get("rcmd")
			if fld.Kind() == reflect.String && funk.Contains(rcmd.Options, "fmt") {
				// format the tag name by injecting any value into the tag name
				// for example lib=%s and struct value is some/path -> lib=some/path
				rcmd.Name = fmt.Sprintf(rcmd.Name, fld.Value())
			}
			args = append(args, fmt.Sprintf("--%s", rcmd.Name))
		}
	}
	return args
}

// Install installs a given tarball
// exit code 0 - success, 1 - error
func Install(
	fs afero.Fs,
	tbp string, // tarball path
	args *InstallArgs,
	rs RSettings,
	es ExecSettings,
	lg *logrus.Logger,
) (CmdResult, error) {
	rdir := es.WorkDir
	if rdir == "" {
		rdir, _ = os.Getwd()
		lg.WithFields(
			logrus.Fields{"rdir": rdir},
		).Debug("launch dir")
	} else {
		ok, err := afero.DirExists(fs, rdir)
		if !ok || err != nil {
			// TODO replace with better error
			return CmdResult{
				Stderr:   err.Error(),
				ExitCode: 1,
			}, err
		}
	}
	if !filepath.IsAbs(tbp) {
		tbp = filepath.Clean(filepath.Join(rdir, tbp))
	}

	ok, err := afero.Exists(fs, tbp)
	if !ok || err != nil {
		lg.WithFields(logrus.Fields{
			"path": tbp,
			"ok":   ok,
			"err":  err,
		}).Error("package tarball not found")
		var errs string
		if err != nil {
			errs = err.Error()
		} else {
			// nil error not ok
			errs = fmt.Sprintf("%s does not exist", tbp)
		}
		return CmdResult{
			Stderr:   fmt.Sprintf("err: %s, ok: %v", errs, ok),
			ExitCode: 1,
		}, err
	}

	cmdArgs := []string{
		"CMD",
		"install",
	}
	cmdArgs = append(cmdArgs, args.CliArgs()...)
	cmdArgs = append(cmdArgs, tbp)
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
			"cmd":       "install",
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
		}).Info("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		fmt.Sprintf("%s", rs.R()),
		cmdArgs...,
	)

	cmd.Dir = rdir
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err = cmd.Run()
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

	cmdResult := CmdResult{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
	}
	if exitCode != 0 {
		lg.WithFields(
			logrus.Fields{
				"stdout":   stdout,
				"stderr":   stderr,
				"exitCode": exitCode,
			}).Error("cmd output")
	} else {
		lg.WithFields(
			logrus.Fields{
				"stdout":   stdout,
				"stderr":   stderr,
				"exitCode": exitCode,
			}).Info("cmd output")
	}
	return cmdResult, err
}

// InstallThroughBinary installs in a two pass fashion
// by first installing and generating a binary in
// a tmp dir, then installs the binary to the desired
// library location
// In addition to returning the CmdResult and any errors
// the path to the binary is also provided for
// additional handling of the binary such as caching
func InstallThroughBinary(
	fs afero.Fs,
	tbp string, // tarball path
	args *InstallArgs,
	rs RSettings,
	es ExecSettings,
	lg *logrus.Logger,
) (CmdResult, string, error) {
	tmpdir := os.TempDir()
	origDir := es.WorkDir
	if origDir == "" {
		origDir, _ = os.Getwd()
	}
	es.WorkDir = tmpdir
	finalLib := args.Library

	// since moving directories to tmp for execution,
	// should treat everything as absolute
	if !filepath.IsAbs(tbp) {
		tbp = filepath.Clean(filepath.Join(origDir, tbp))
	}
	if !filepath.IsAbs(finalLib) {
		finalLib = filepath.Clean(filepath.Join(origDir, finalLib))
	}
	// instead install to tmpdir rather than the library
	// for the first pass, then will ultimately install the
	// generated binary to the proper location
	// this will prevent failed installs overwriting existing
	// properly installed packages in the final lib
	args.Library = tmpdir
	ib := &InstallArgs{
		Library: finalLib,
	}

	// built binaries have the path extension .tgz rather than tar.gz
	// but otherwise have the same name from empirical testing
	// pkg_0.0.1.tar.gz --> pkg_0.0.1.tgz
	lg.WithFields(logrus.Fields{
		"tbp":  tbp,
		"args": args,
	}).Debug("installing tarball")
	res, err := Install(fs,
		tbp,
		args,
		rs,
		es,
		lg)
	if err == nil && res.ExitCode == 0 {
		bbp := strings.Replace(filepath.Base(tbp), "tar.gz", "tgz", 1)
		binaryBall := filepath.Join(tmpdir, bbp)
		lg.WithFields(logrus.Fields{
			"tbp":        tbp,
			"bbp":        bbp,
			"binaryBall": binaryBall,
		}).Debug("binary location prior to install")
		ok, _ := afero.Exists(fs, binaryBall)
		if !ok {
			lg.WithFields(logrus.Fields{
				// check previous stderror, which R logs to installation status
				"stderr":     res.Stderr,
				"tmpdir":     tmpdir,
				"binaryPath": binaryBall,
			}).Error("could not find binary")
			// change the exit code in case top level just blindly looks for
			// 0 exit code means good
			// the successful initial install should still bubble up through
			// the stderr/out
			res.ExitCode = 1
			return res, "", errors.New("no binary found")
		}
		res, err = Install(fs,
			binaryBall,
			ib,
			rs,
			es,
			lg)
		return res, binaryBall, err
	}
	return res, "", err
}
