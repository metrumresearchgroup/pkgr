package rcmd

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func binaryNameOs(os, pkg, version, platform string) string {
	switch os {
	case "darwin":
		return fmt.Sprintf("%s_%s.tgz", pkg, version)
	case "linux":
		return fmt.Sprintf("%s_%s_R_%s.tar.gz", pkg, version, platform)
	case "windows":
		return fmt.Sprintf("%s_%s.zip", pkg, version)
	default:
		log.Fatal("platform not supported for binary detection")
		return ("")
	}
}

func binaryName(pkg, version, platform string) string {
	return binaryNameOs(runtime.GOOS, pkg, version, platform)
}

func binaryExtOs(os, p, platform string) string {
	switch os {
	case "darwin":
		return strings.Replace(filepath.Base(p), "tar.gz", "tgz", 1)
	case "linux":
		pf := "_R_" + platform + ".tar.gz"
		return strings.Replace(filepath.Base(p), ".tar.gz", pf, 1)
	case "windows":
		return strings.Replace(filepath.Base(p), "tar.gz", "zip", 1)
	default:
		log.Fatal("platform not supported for binary detection")
		return ("")
	}
}

func binaryExt(p, platform string) string {
	return binaryExtOs(runtime.GOOS, p, platform)
}
