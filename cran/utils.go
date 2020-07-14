package cran

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BinaryUriType int

const (
	DefaultUri = 20
	SuffixUri  = 21
)

var osRelease *OsRelease

var supportedDistros = map[string]bool{
	"bionic": true,
	"xenial": true,
	"centos": false,
}

// these are also duplicated in rcmd for now
func binaryName(pkg, version string) string {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("%s_%s.tgz", pkg, version)
	case "linux":
		if strings.Contains(osRelease.IdLike, "rhel") {
			return fmt.Sprintf("%s_%s_R_x86_64-redhat-linux-gnu.tar.gz", pkg, version)
		}
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
		if linuxSupportsBinary() && rt == MPN {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

// LinuxSupportsBinary tells if a distro supports binaries
// namely, Ubuntu 16.04 and 18.04
func linuxSupportsBinary() bool {
	ReadOsRelease()
	if supportedDistros[*osRelease.VersionCodename] || supportedDistros[osRelease.Id] {
		return true
	}
	log.Info("The running version of Linux does not support binary packages")
	return false
}

func getLinuxCodename() *string {
	ReadOsRelease()
	return osRelease.VersionCodename
}

func getLinuxLtsRelease() string {
	ReadOsRelease()
	return osRelease.LtsRelease
}

func getLinuxBinaryUri(params ...BinaryUriType) string {
	if len(params) > 0 && params[0] == SuffixUri {
		return fmt.Sprintf("%s/%s", osRelease.Id, *osRelease.VersionCodename)
	}
	if osRelease.VersionCodename != nil && !strings.Contains(osRelease.IdLike, "rhel") {
		return fmt.Sprintf("%s/%s/%s", osRelease.Id, *osRelease.VersionCodename, osRelease.LtsRelease)
	}
	// EL distros keep the right naming as it is
	return fmt.Sprintf("%s/%s", osRelease.Id, osRelease.VersionId)
}

func cranBinaryURL(rv RVersion, params ...BinaryUriType) string {
	switch runtime.GOOS {
	case "darwin":
		if rv.Major == 4 {
			return "macosx"
		}
		return "macosx/el-capitan"
	case "windows":
		return "windows"
	case "linux":
		if len(params) > 0 && params[0] == SuffixUri {
			return fmt.Sprintf("linux/%s", getLinuxBinaryUri(params[0]))
		}
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

func ReadOsRelease() {
	if osRelease != nil {
		// Already cached
		return
	}
	//os-release doesn't consistently quote variables, so we need to manipulate it a little bit
	configData, err := ioutil.ReadFile("/etc/os-release")

	//Find all unquoted strings and quote them
	re := regexp.MustCompile(`(.*?=)([^"].*)`)
	fixedConfig := re.ReplaceAll(configData, []byte("${1}\"${2}\""))

	//Let viper map it in otherwise it defaults to yaml
	//Don't pollute the global viper instance or we'll segfault later trying to read pkgr.yml
	vp := viper.New()
	vp.SetConfigType("toml")
	err = vp.ReadConfig(bytes.NewReader(fixedConfig))

	if err != nil {
		log.Fatal("%v", err)
	}

	err = vp.Unmarshal(&osRelease)

	if err != nil {
		log.Fatal("%v\n", err)
	}

	// simplify this so it also works on EL distros
	ltsReleaseMatcher := regexp.MustCompile(`^.*?(\d+)\s.*`)

	ltsRelease := ltsReleaseMatcher.ReplaceAllString(osRelease.Version, "$1")
	osRelease.LtsRelease = ltsRelease
}
