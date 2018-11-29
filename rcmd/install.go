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
	"sync"
	"syscall"
	"time"
	"github.com/dpastoor/goutils"
	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/gpsr"
	"github.com/fatih/structs"
	"github.com/fatih/structtag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
)

// NewDefaultInstallArgs provides a set of sane default installation args
func NewDefaultInstallArgs() InstallArgs {
	return InstallArgs{
		WithKeepSource: true,
		NoMultiarch:    true,
		InstallTests:   true,
		Build:          true,
	}
}

// CliArgs converts the InstallArgs struct to the proper cli args
// including only returning the relevant args
func (i InstallArgs) CliArgs() []string {
	args := []string{}
	is := structs.New(i)
	nms := structs.Names(i)
	for _, n := range nms {
		fld := is.Field(n)
		if !fld.IsZero() {
			// ... and start using structtag by parsing the tag
			tag, _ := reflect.TypeOf(i).FieldByName(fld.Name())
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

func configureEnv(
	rs RSettings,
	lg *logrus.Logger,
) ([]string){
	envVars := os.Environ()
	envMap := make(map[string]string)
	for _, ev := range envVars {
		evs := strings.SplitN(ev, "=", 1)
		if len(evs) > 1 && evs[1] != "" {
			envMap[evs[0]] = evs[1] 
		}
	}
	 rlu, exists := envMap["R_LIBS_USER"]
	 if exists {
		 // R_LIBS_USER takes precidence over R_LIBS_SITE
		 // so will cause the loading characteristics to
		 // not be representative of the hierarchy specified
		 // in Library/Libpaths in the pkgr configuration
		delete(envMap, "R_LIBS_USER")
		lg.WithField("path", rlu).Debug("deleting set R_LIBS_USER")
	 }
	envVars = []string{}
	for k, v := range rs.EnvVars {
		envMap[k] = v
	}

	ok, lp := rs.LibPathsEnv()
	if ok {
	// if LibPaths set, lets drop R_LIBS_SITE set as an ENV and instead
	// add the generated R_LIBS_SITE from LibPathsEnv
		delete(envMap, "R_LIBS_SITE")
		envVars = append(envVars, lp)
	}
	for k, v := range envMap {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	// double down on overwriting any specification of user customization
	// and set R_LIBS_SITE to the same as the user
	envVars = append(envVars, strings.Replace(lp, "R_LIBS_SITE", "R_LIBS_USER", 1))

	return envVars
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
) (CmdResult, error) {
	rdir := es.WorkDir
	if rdir == "" {
		rdir, _ = os.Getwd()
		lg.WithFields(
			logrus.Fields{"rdir": rdir},
		).Trace("launch dir set to working directory")
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

	envVars := configureEnv(rs, lg) 

	cmdArgs := []string{
		"CMD",
		"INSTALL",
	}
	cmdArgs = append(cmdArgs, args.CliArgs()...)
	cmdArgs = append(cmdArgs, tbp)
	lg.WithFields(
		logrus.Fields{
			"cmd":       "install",
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
		fmt.Sprintf("%s", rs.R()),
		cmdArgs...,
	)
	cmd.Env = envVars
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
			}).Debug("cmd output")
	}
	return cmdResult, err
}

// isInCache notes if package binary already available in cache
// and returns a new installrequest based on the binary path if available
func isInCache(
	fs afero.Fs,
	ir InstallRequest,
	pc PackageCache,
	lg *logrus.Logger,
) (bool, InstallRequest) {
	// if not in cache just pass back
	meta := ir.Metadata
	pkg := ir.Metadata.Metadata.Package
	bpath := filepath.Join(pc.BaseDir, meta.Metadata.Config.Repo.Name, "binary", binaryName(pkg.Package, pkg.Version))
	lg.WithFields(logrus.Fields{
		"path": bpath,
		"package": pkg.Package,
		}).Trace("checking in cache")
	exists, err := goutils.Exists(fs, bpath)
	if !exists || err != nil {
	lg.WithField("package", pkg.Package).Trace("not found in cache")
		return false, ir
	}
	lg.WithField("package", pkg.Package).Trace("found in cache")
	ir.Metadata.Path = bpath
	ir.Metadata.Metadata.Config.Type = cran.Binary
	return true, ir
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
	ir InstallRequest,
	pc PackageCache,
	lg *logrus.Logger,
) (CmdResult, string, error) {
	exists, _ := goutils.DirExists(fs, filepath.Join(ir.InstallArgs.Library, ir.Package))
	if exists {
		lg.WithField("package", ir.Package).Info("package already installed")
		return CmdResult{
			ExitCode: 0,
			Stderr:   fmt.Sprintf("already installed: %s", ir.Package),
		}, "", nil
	}

	inCache, ir := isInCache(fs, ir, pc, lg)
	if inCache {
		lg.WithField("package", ir.Package).Debug("package detected in cache")
		// don't need to build since already a binary
		ir.InstallArgs.Build = false
		res, err := Install(fs,
			ir.Metadata.Path,
			ir.InstallArgs,
			ir.RSettings,
			ir.ExecSettings,
			lg)
		// don't pass binaryball path back since already in cache
		return res, "", err
	}

	tmpdir := os.TempDir()
	origDir := ir.ExecSettings.WorkDir
	if origDir == "" {
		origDir, _ = os.Getwd()
	}
	ir.ExecSettings.WorkDir = tmpdir
	finalLib := ir.InstallArgs.Library
	// since moving directories to tmp for execution,
	// should treat everything as absolute
	if !filepath.IsAbs(ir.Metadata.Path) {
		ir.Metadata.Path = filepath.Clean(filepath.Join(origDir, ir.Metadata.Path))
	}
	if !filepath.IsAbs(finalLib) {
		finalLib = filepath.Clean(filepath.Join(origDir, finalLib))
	}
	// instead install to tmpdir rather than the library
	// for the first pass, then will ultimately install the
	// generated binary to the proper location
	// this will prevent failed installs overwriting existing
	// properly installed packages in the final lib
	// we also need to point to the requested library as a libpath such that
	// when installing to the tmp dir will still have the proper packages
	ir.RSettings.LibPaths = append(ir.RSettings.LibPaths, ir.InstallArgs.Library)
	ir.InstallArgs.Library = tmpdir
	// TODO: check if installing the binary still has relevant installation args that should be
	// propogated, or if simply pointing to the final lib is sufficient
	ib := InstallArgs{
		Library: finalLib,
	}

	// built binaries have the path extension .tgz rather than tar.gz
	// but otherwise have the same name from empirical testing
	// pkg_0.0.1.tar.gz --> pkg_0.0.1.tgz
	lg.WithFields(logrus.Fields{
		"tbp":  ir.Metadata.Path,
		"args": ir.InstallArgs,
	}).Debug("installing tarball")
	res, err := Install(fs,
		ir.Metadata.Path,
		ir.InstallArgs,
		ir.RSettings,
		ir.ExecSettings,
		lg)
	installPath := filepath.Join(tmpdir, ir.Package)
	de, _ := afero.DirExists(fs, installPath)
	if de {
		err := fs.RemoveAll(installPath)
		if err != nil {
			lg.WithFields(logrus.Fields{
				"err": err,
				"path": installPath,
			}).Error("error removing installed package in tmp dir")
		}
	}
	if err == nil && res.ExitCode == 0 {
		bbp := binaryExt(ir.Metadata.Path)
		binaryBall := filepath.Join(tmpdir, bbp)
		lg.WithFields(logrus.Fields{
			"tbp":        ir.Metadata.Path,
			"bbp":        bbp,
			"binaryBall": binaryBall,
		}).Trace("binary location prior to install")
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
			ir.RSettings,
			ir.ExecSettings,
			lg)
		return res, binaryBall, err
	}
	return res, "", err
}


// InstallPackagePlan installs a set of packages by layer
func InstallPackagePlan(
	fs afero.Fs,
	plan gpsr.InstallPlan,
	dl *cran.PkgMap,
	pc PackageCache,
	args InstallArgs,
	rs RSettings,
	es ExecSettings,
	lg *logrus.Logger,
	ncpu int,
) error {
	wg := sync.WaitGroup{}
	startTime := time.Now()
	// for now this will only be updated in the Update function
	// however if it may be concurrently accessed should consider a syncmap implementation
	installedPkgs := make(map[string]bool)
	// if packages ID'd as ready to install signal so can push them only the queue
	shouldInstall := make(chan string)
	anyFailed := false
	iDeps := plan.InvertDependencies()
	failedPkgs := []string{}
	iq := NewInstallQueue(ncpu,
		InstallThroughBinary,
		func(iu InstallUpdate) {
			if iu.Err != nil {
				lg.WithField("err", iu.Err).Warn("error installing")
				anyFailed = true
				failedPkgs = append(failedPkgs, iu.Package)
			} else {
				// set that the package is installed,
				// then check if any of the inverse dependencies are
				// ready to be installed, and if so, signal they should
				// be installed
				pkg, _ := dl.Get(iu.Package)

				lg.WithFields(logrus.Fields{"binary": iu.BinaryPath, "src": pkg.Path}).Debug(iu.Package)
				lg.WithField("package", iu.Package).Info("Successfully Installed")
				installedPkgs[iu.Package] = true
				deps, exists := iDeps[iu.Package]
				if exists {
					for _, maybeInstall := range deps {
						needDeps := plan.DepDb[maybeInstall]
						allInstalled := true
						for _, d := range needDeps {
							_, installed := installedPkgs[d]
							if !installed {
								allInstalled = false
							}
						}
						if allInstalled && !anyFailed {
							wg.Add(1)
							lg.WithFields(logrus.Fields{
								"from": iu.Package,
								"suggested": maybeInstall,
								}).Trace("suggesting installation")
							shouldInstall <- maybeInstall
						}
					}
				}
				if iu.BinaryPath != "" {
					bdir := filepath.Join(
						filepath.Dir(filepath.Dir(pkg.Path)),
						"binary",
					)
					os.Mkdir(bdir, 0777)
					bpath := filepath.Join(
						bdir,
						filepath.Base(iu.BinaryPath),
					)
					_, err := goutils.Copy(iu.BinaryPath, bpath)
					lg.WithFields(logrus.Fields{"from": iu.BinaryPath, "to": bpath}).Trace("copied binary")
					if err != nil {
						lg.WithFields(logrus.Fields{"from": iu.BinaryPath, "to": bpath}).Error("error copying binary")
					}
				}
			}
			wg.Done()
			fmt.Println("installed with message: ", iu.Result.Stderr)
		}, lg,
	)
	go func(c chan string) {
		requestedPkgs := make(map[string]bool)
		for p := range c {
			if anyFailed {
				// stop trying to install any more
				continue
			}

			// wg added from updater before pushing here
			// wg.Add(1)
			 _, seen := requestedPkgs[p] 
			if seen {
				// should only need to request a package once to install
				continue	
			} else {
				requestedPkgs[p] = true
			}
			pkg, _ := dl.Get(p)
			if p == "data.table" {
				nrs := rs
				fmt.Println("setting data.table makevar")
				nev := make(map[string]string)
				for k, v := range rs.EnvVars {
					nev[k] = v
				}
				nev["R_MAKEVARS_USER"] = "~/.R/Makevars_data.table"
				nrs.EnvVars = nev

				iq.Push(InstallRequest{
					Package:      p,
					Metadata:     pkg,
					Cache:        pc,
					InstallArgs:  args,
					RSettings:    nrs,
					ExecSettings: es,
				})
			} else {
				lg.WithField("package", p).Trace("pushing installation to queue")
				iq.Push(InstallRequest{
					Package:      p,
					Metadata:     pkg,
					Cache:        pc,
					InstallArgs:  args,
					RSettings:    rs,
					ExecSettings: es,
				})
			}
		}
	}(shouldInstall)

	lg.WithField("packages", strings.Join(plan.StartingPackages, ", ")).Info("starting initial install")

	for _, p := range plan.StartingPackages {
		wg.Add(1)
		shouldInstall <- p
	}
	wg.Wait()

	lg.WithField("duration", time.Since(startTime)).Info("total install time")
	for pkg := range plan.DepDb {
		_, exists := installedPkgs[pkg]	
		if !exists {
			lg.Errorf("did not install %s", pkg)
		}
	}
	if anyFailed {
		lg.Errorf("installation failed for packages: %s", strings.Join(failedPkgs, ", "))
		return fmt.Errorf("failed installation for packages: %s", strings.Join(failedPkgs, ", "))
	}
	return nil
}
