package rcmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dpastoor/goutils"
	"github.com/fatih/structs"
	"github.com/fatih/structtag"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	funk "github.com/thoas/go-funk"
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

// Install installs a given tarball
// exit code 0 - success, 1 - error
func Install(
	fs afero.Fs,
	pkg string,
	tbp string, // tarball path
	args InstallArgs,
	rs RSettings,
	es ExecSettings,
	ir InstallRequest) (CmdResult, error) {
	rdir := es.WorkDir
	if rdir == "" {
		rdir, _ = os.Getwd()
		log.WithFields(
			log.Fields{"rdir": rdir},
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
		log.WithFields(log.Fields{
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

	envVars := configureEnv(os.Environ(), rs, pkg)

	cmdArgs := []string{
		"--vanilla",
		"CMD",
		"INSTALL",
	}
	cmdArgs = append(cmdArgs, args.CliArgs()...)
	cmdArgs = append(cmdArgs, tbp)
	log.WithFields(
		log.Fields{
			"cmd":       "install",
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       envVars,
			"package":   pkg,
		}).Trace("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		fmt.Sprintf("%s", rs.R(runtime.GOOS)),
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
		// update DESCRIPTION file with pkgr metadata
		writeDescriptionInfo(fs, ir, args)

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
		log.WithFields(
			log.Fields{
				"stdout":   stdout,
				"stderr":   stderr,
				"exitCode": exitCode,
				"package":  pkg,
			}).Error("cmd output")
	} else {
		log.WithFields(
			log.Fields{
				"stdout":   stdout,
				"stderr":   stderr,
				"exitCode": exitCode,
				"package":  pkg,
			}).Debug("cmd output")
	}
	return cmdResult, err
}

// isInCache notes if package binary already available in cache
// and returns a new installrequest based on the binary path if available
func isInCache(
	fs afero.Fs,
	ir InstallRequest,
	pc PackageCache) (bool, InstallRequest) {
	// if not in cache just pass back
	meta := ir.Metadata
	pkg := ir.Metadata.Metadata.Package

	repoHash := cran.RepoURLHash(meta.Metadata.Config.GetOrigin())
	bpath := filepath.Join(
		pc.BaseDir,
		repoHash,
		"binary",
		ir.RSettings.Version.ToString(),
		binaryName(pkg.Package, pkg.Version, ir.RSettings.Platform),
	)
	exists, err := goutils.Exists(fs, bpath)
	if !exists || err != nil {
		log.WithFields(log.Fields{
			"path":    bpath,
			"package": pkg.Package,
		}).Trace("not found in cache")
		return false, ir
	}
	log.WithFields(log.Fields{
		"path":    bpath,
		"package": pkg.Package,
	}).Trace("found in cache")
	ir.Metadata.Path = bpath
	ir.Metadata.Metadata.Config.SetSourceType(cran.Binary.String())
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
	pc PackageCache) (CmdResult, string, error) {
	exists, _ := goutils.DirExists(fs, filepath.Join(ir.InstallArgs.Library, ir.Package))
	if exists {
		log.WithFields(log.Fields{
			"package": ir.Package,
			"version": ir.Metadata.Metadata.Package.Version,
		}).Debug("package already installed")
		return CmdResult{
			ExitCode: -999,
			Stderr:   fmt.Sprintf("already installed: %s", ir.Package),
		}, "", nil
	}

	inCache, ir := isInCache(fs, ir, pc)
	if inCache {
		// don't need to build since already a binary
		ir.InstallArgs.Build = false
		res, err := Install(fs,
			ir.Package,
			ir.Metadata.Path,
			ir.InstallArgs,
			ir.RSettings,
			ir.ExecSettings,
			ir)
		// don't pass binaryball path back since already in cache
		return res, "", err
	}

	tmpdir := filepath.Join(
		os.TempDir(),
		randomString(12),
	)
	err := fs.MkdirAll(tmpdir, 0777)
	if err != nil {
		log.Fatalf("could not make tmpdir at: %s to install package", tmpdir)
	}
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
	log.WithFields(log.Fields{
		"tbp":  ir.Metadata.Path,
		"args": ir.InstallArgs,
	}).Debug("installing tarball")
	res, err := Install(fs,
		ir.Package,
		ir.Metadata.Path,
		ir.InstallArgs,
		ir.RSettings,
		ir.ExecSettings,
		ir)
	installPath := filepath.Join(tmpdir, ir.Package)
	de, _ := afero.DirExists(fs, installPath)
	if de {
		err := fs.RemoveAll(installPath)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": installPath,
			}).Error("error removing installed package in tmp dir")
		}
	}
	if err == nil && res.ExitCode == 0 {
		bbp := binaryExt(ir.Metadata.Path, ir.RSettings.Platform)
		binaryBall := filepath.Join(tmpdir, bbp)
		log.WithFields(log.Fields{
			"tbp":        ir.Metadata.Path,
			"bbp":        bbp,
			"binaryBall": binaryBall,
		}).Trace("binary location prior to install")
		ok, _ := afero.Exists(fs, binaryBall)
		if !ok {
			log.WithFields(log.Fields{
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
			ir.Package,
			binaryBall,
			ib,
			ir.RSettings,
			ir.ExecSettings,
			ir)
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
	ncpu int,
) error {

	//var successCounter uint64
	wg := sync.WaitGroup{}

	startTime := time.Now()

	// for now this will only be updated in the Update function
	// however if it may be concurrently accessed should consider a syncmap implementation
	installedPkgs := make(map[string]bool)
	// if packages ID'd as ready to install signal so can push them only the queue
	shouldInstall := make(chan string)
	anyFailed := false

	packagesNeeded := plan.GetNumPackagesToInstall()

	iDeps := plan.InvertDependencies()

	failedPkgs := []string{}

	installQueue := NewInstallQueue(
		ncpu,
		InstallThroughBinary,
		func(iu InstallUpdate) {
			if iu.Err != nil {
				log.WithField("err", iu.Err).Warn("error installing")
				anyFailed = true
				failedPkgs = append(failedPkgs, iu.Package)
			} else {
				// set that the package is installed,
				// then check if any of the inverse dependencies are
				// ready to be installed, and if so, signal they should
				// be installed
				pkg, _ := dl.Get(iu.Package)
				log.WithFields(log.Fields{"binary": iu.BinaryPath, "src": pkg.Path}).Debug(iu.Package)

				if iu.Result.ExitCode != -999 {
					packagesNeeded = packagesNeeded - 1
					log.WithFields(log.Fields{
						"package": iu.Package,
						"version": pkg.Metadata.Package.Version,
						"repo":    pkg.Metadata.Config.GetOrigin().Name,
						"remaining": packagesNeeded,
					}).Info("Successfully Installed.")
				}
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
							log.WithFields(log.Fields{
								"from":      iu.Package,
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
						rs.Version.ToString(),
						filepath.Base(iu.BinaryPath),
					)
					_, err := goutils.Copy(iu.BinaryPath, bpath)
					log.WithFields(log.Fields{"from": iu.BinaryPath, "to": bpath}).Trace("copied binary")
					if err != nil {
						log.WithFields(log.Fields{"from": iu.BinaryPath, "to": bpath}).Error("error copying binary")
						return
					}
					// want to delete binaries from the existing tmpdir
					// so do not carry around two copies. This is especially
					// relevant for containerized environment where layers get snapshotted
					// before tmp dirs are cleaned up, which can result in very large
					// images
					fs.Remove(iu.BinaryPath)

				}
			}
			wg.Done()
		}, // End anonymous function
	) // End NewInstallQueue
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
			log.WithField("package", p).Trace("pushing installation to queue")
			installQueue.Push(InstallRequest{
				Package:      p,
				Metadata:     pkg,
				Cache:        pc,
				InstallArgs:  args,
				RSettings:    rs,
				ExecSettings: es,
			})
		}
	}(shouldInstall)

	log.Info("starting initial install")

	for _, p := range plan.StartingPackages {
		wg.Add(1)
		shouldInstall <- p
	}
	wg.Wait()

	log.WithField("duration", time.Since(startTime)).Info("total install time")
	for pkg := range plan.DepDb {
		_, exists := installedPkgs[pkg]
		if !exists {
			log.Errorf("did not install %s", pkg)
		}
	}

	if anyFailed {
		log.Errorf("installation failed for packages: %s", strings.Join(failedPkgs, ", "))
		return fmt.Errorf("failed installation for packages: %s", strings.Join(failedPkgs, ", "))
	}
	return nil
}

func writeDescriptionInfo(fs afero.Fs, ir InstallRequest, ia InstallArgs) {
	_, err := updateDescriptionInfo(
		fs,
		filepath.Join(ia.Library, ir.Package, "DESCRIPTION"),
		ir.ExecSettings.PkgrVersion,
		ir.Metadata.Metadata.Config.GetSourceType2().String(),
		ir.Metadata.Metadata.Config.GetOrigin().URL,
		ir.Metadata.Metadata.Config.GetOrigin().Name)

	if err != nil {
		log.Warn(fmt.Sprintf("DESCRIPTION file error:%s", err))
	}
}

// function to apply relevant changes to DESCRIPTION file (if applicable) and ultimately write out the updated
// DESCRIPTION file.
func updateDescriptionInfo(fs afero.Fs, filename, version, installType, repoURL, repo string) ([]string, error) {
	descriptionLines, err := goutils.ReadLinesFS(fs, filename)
	if err != nil {
		return nil, err
	}
	descriptionFinal := updateDescriptionInfoByLines(descriptionLines, version, installType, repoURL, repo)
	err = writeUpdatedDescriptionFile(fs, filename, descriptionFinal)

	return descriptionFinal, err
}

// writes a slice of strings out to a file. Adds newline characters to the end of the each line.
func writeUpdatedDescriptionFile(fs afero.Fs, filename string, update []string) error {
	err := goutils.WriteLinesFS(fs, update, filename)

	if err != nil {
		log.WithFields(log.Fields{
			"filename" : filename,
		}).Error("could not update DESCRIPTION file", err)
		return err
	}

	return nil
}

// takes a string slice of lines (without line-endings) representing a DESCRIPTION file and uses
// provided parameters to make appropriate updates to that description file.
// returns the updated string slice, still without new-line characters.
func updateDescriptionInfoByLines(lines []string, version, installType, repoURL, repo string) []string {

	var newLines []string
	for _, line := range lines {
		//log.Info("Starting repo: " + repo)

		// Don't add these lines to final set, if they exist. This way, we can add/"overwrite" them at the end.
		if strings.Contains(line, "PkgrVersion:") ||
			strings.Contains(line, "PkgrInstallType:") ||
			strings.Contains(line, "PkgrRepositoryURL") {
			continue
		}

		if strings.Contains(line, string("Repository:")) && !strings.Contains(line, repo) {
			//log.Info("Got where we needed.")
			originalRepo := strings.Trim(strings.Split(line, ":")[1], " ")
			newLines = append(newLines, "OriginalRepository: " + originalRepo)
			line = "Repository: " + repo
		}
		newLines = append(newLines, line)
	}

	newLines = append(newLines, "PkgrVersion: "+version)
	newLines = append(newLines, "PkgrInstallType: "+installType)
	newLines = append(newLines, "PkgrRepositoryURL: "+repoURL)

	return newLines
}
