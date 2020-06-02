package cran

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

)

type OsRelease struct {
	Name            string `mapstructure:"NAME"`
	Version         string `mapstructure:"VERSION"`
	Id              string `mapstructure:"ID"`
	IdLike          string `mapstructure:"ID_LIKE"`
	LtsRelease      string
	PrettyName      string `mapstructure:"PRETTY_NAME"`
	VersionId       string `mapstructure:"VERSION_ID"`
	VersionCodename string `mapstructure:"VERSION_CODENAME"`
	UbuntuCodename  string `mapstructure:"UBUNTU_CODENAME"`
}

var osRelease *OsRelease

var supportedDistros = map[string]bool {
	"bionic": true,
	"xenial": true,
}

// these are also duplicated in rcmd for now
func binaryName(pkg, version string) string {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("%s_%s.tgz", pkg, version)
	case "linux":
		return fmt.Sprintf("%s_%s_R_x86_64-pc-linux-gnu.tar.gz", pkg, version)
	case "windows":
		return fmt.Sprintf("%s_%s.zip", pkg, version)
	default:
		fmt.Println("platform not supported for binary detection")
		return ""
	}
}

func binaryExt(p string) string {
	switch runtime.GOOS {
	case "darwin":
		return strings.Replace(filepath.Base(p), "tar.gz", "tgz", 1)
	case "linux":
		return strings.Replace(filepath.Base(p), ".tar.gz", "_R_x86_64-pc-linux-gnu.tar.gz", 1)
	case "windows":
		return strings.Replace(filepath.Base(p), "tar.gz", "zip", 1)
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

// SupportsCranBinary tells if a platform supports binaries
// namely, windows/mac to, but linux does not
func SupportsCranBinary() bool {
	switch runtime.GOOS {
	case "darwin":
		return true
	case "windows":
		return true
	case "linux":
		if linuxSupportsBinary() {
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
	codename := getLinuxCodename()
	if supportedDistros[codename] {
		return true
	}
	log.Info("The running version of Linux does not support binary packages")
	return false
}

func getLinuxCodename() string {
	ReadOsRelease()
	return osRelease.VersionCodename
}

func getLinuxLtsRelease() string {
	ReadOsRelease()
	return osRelease.LtsRelease
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
		return fmt.Sprintf("linux/ubuntu/%s/%s", getLinuxCodename(), getLinuxLtsRelease())
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
	viper.SetConfigType("toml")
	err = viper.ReadConfig(bytes.NewReader(fixedConfig))

	if err != nil {
		log.Fatal("%v", err)
	}

	err = viper.Unmarshal(osRelease)

	if err != nil {
		log.Fatal("%v\n", err)
	}

	ltsReleaseMatcher := regexp.MustCompile(`^\d.*?\.(\d+)\s.*`)
	ltsRelease := ltsReleaseMatcher.ReplaceAllString(osRelease.Version, "$1")
	osRelease.LtsRelease = ltsRelease
}