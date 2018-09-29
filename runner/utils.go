package runner

import "strings"

func sanitizeDirName(n string) string {
	n = strings.TrimSuffix(n, "/")
	n = strings.TrimSuffix(n, "\\\\")
	return n
}
