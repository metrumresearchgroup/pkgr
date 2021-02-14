package cran

import (
	"bytes"
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"regexp"
	"runtime"
)

type BinaryUriType int

const (
	DefaultUri = 20
	SuffixUri  = 21
)

var osRelease OsRelease

var supportedDistros = map[string]bool{
	"focal": true,
	"bionic": true,
	"xenial": true,
	"centos": true,
	"rhel": true,
	"ubuntu": true,
}

// these are also duplicated in rcmd for now
func binaryName(pkg, version string) string {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("%s_%s.tgz", pkg, version)
	case "linux":
		if osRelease.Id == "rhel" {
			return fmt.Sprintf("%s_%s_R_x86_64-redhat-linux-gnu.tar.gz", pkg, version)
		}
		// checked centos docker container and returned
		// packaged installation of ‘R6’ as ‘R6_2.5.0_R_x86_64-pc-linux-gnu.tar.gz’
		return fmt.Sprintf("%s_%s_R_x86_64-pc-linux-gnu.tar.gz", pkg, version)
	case "windows":
		return fmt.Sprintf("%s_%s.zip", pkg, version)
	default:
		fmt.Println("platform not supported for binary detection")
		return ""
	}
}

// DefaultType provides the default type for the given platform
// runtime
func DefaultType() SourceType {
	switch runtime.GOOS {
	case "darwin":
		return Binary
	case "windows":
		return Binary
	case "linux":
		return Source
	default:
		return Source
	}
}

// SupportsBinary tells if a platform supports binaries
// namely, windows/mac to, but linux does not
func SupportsBinary(rt RepoType) bool {
	switch runtime.GOOS {
	case "darwin":
		return true
	case "windows":
		return true
	case "linux":
		if linuxKnownSupportsBinary() && rt == MPN {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

// linuxKnownSupportsBinary tells if a distro supports binaries
// namely, Ubuntu 16.04 and 18.04
func linuxKnownSupportsBinary() bool {
	err := ReadOsRelease()
	if err != nil {
		log.Warnf("error reading linux binary information: %s\n", err)
		return false
	}
	if osRelease.Id == "" {
		return false
	}
	if supportedDistros[osRelease.VersionCodename] || supportedDistros[osRelease.Id] {
		return true
	}
	log.Info("The running version of Linux might not support binary packages, please contact the pkgr development team")
	return true
}

func getLinuxCodename() string {
	ReadOsRelease()
	return osRelease.VersionCodename
}


func getLinuxBinaryUri() string {
	if !osRelease.checked {
		err := ReadOsRelease()
		if err != nil {
			log.Warnf("could not get derive linux information with error: %s\n", err)
			return ""
		}
	}
	if osRelease.Id == "" {
		return ""
	}
	// ubuntu should follow ubuntu/focal ubuntu/bionic  etc...
	if osRelease.Id == "ubuntu" {
		return fmt.Sprintf("%s/%s", osRelease.Id, osRelease.VersionCodename)
	}

	// centos/redhat should follow centos/<majorversion> eg centos/8 centos/7
	// for centos version_id is just a single digit, eg 7 or 8, but for redhat will be major.minor
	// rstudio package manager seems confident just 7 vs 8 is sufficient as a default therefore will follow
	// given the compat with redhat can normalize to centos
	if osRelease.Id == "centos" || osRelease.Id == "rhel" {
		// reminder that go indices on strings pull their byte representation, not the string value
		// either coerce back to string or work with runes
		version := string(osRelease.VersionId[0])
		return fmt.Sprintf("centos/%s", version)
	}

	// default for other distros in case wants to fall through
	return fmt.Sprintf("%s/%s", osRelease.Id, osRelease.VersionId)
}

func cranBinaryURL(rv RVersion) string {
	switch runtime.GOOS {
	case "darwin":
		if rv.Major == 4 {
			return "macosx"
		}
		return "macosx/el-capitan"
	case "windows":
		return "windows"
	case "linux":
		return fmt.Sprintf("linux/%s", getLinuxBinaryUri())
	default:
		fmt.Println("platform not supported for binary detection")
		return ""
	}
}

// RepoURLHash provides a hash of the repoURL
// given the structure Name-<urlhash>
func RepoURLHash(r RepoURL) string {
	h := md5.New()
	// don't hash everything as still want a reasonable identifier
	io.WriteString(h, r.URL)
	urlHash := fmt.Sprintf("%x", h.Sum(nil))
	return r.Name + "-" + urlHash[:12]
}

func ReadOsRelease() error {

	if osRelease.checked {
		// Already checked
		return nil
	}
	//os-release doesn't consistently quote variables, so we need to manipulate it a little bit
	configData, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		return err
	}
	//Find all unquoted strings and quote them
	re := regexp.MustCompile(`(.*?=)([^"].*)`)
	fixedConfig := re.ReplaceAll(configData, []byte("${1}\"${2}\""))

	//Let viper map it in otherwise it defaults to yaml
	//Don't pollute the global viper instance or we'll segfault later trying to read pkgr.yml
	vp := viper.New()
	vp.SetConfigType("toml")
	err = vp.ReadConfig(bytes.NewReader(fixedConfig))

	if err != nil {
		return err
	}

	err = vp.Unmarshal(&osRelease)

	if err != nil {
	   return err
	}

	// simplify this so it also works on EL distros
	ltsReleaseMatcher := regexp.MustCompile(`^.*?(\d+)\s.*`)

	ltsRelease := ltsReleaseMatcher.ReplaceAllString(osRelease.Version, "$1")
	osRelease.LtsRelease = ltsRelease
	osRelease.checked = false
	return nil
}
