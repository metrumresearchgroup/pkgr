package rcmd

import (
	"math/rand"
	"strings"
)

func sanitizeDirName(n string) string {
	n = strings.TrimSuffix(n, "/")
	n = strings.TrimSuffix(n, "\\\\")
	return n
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}
