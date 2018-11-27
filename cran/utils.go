package cran

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

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
	default:
		return false
	}
}
func cranBinaryURL() string {
	switch runtime.GOOS {
	case "darwin":
		return "macosx/el-capitan"
	case "windows":
		return "windows"
	default:
		fmt.Println("platform not supported for binary detection")
		return ""
	}
}
