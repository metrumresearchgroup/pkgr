package rcmd

import (
	"math/rand"
	"strings"
	"time"
)

// TODO: not called
func sanitizeDirName(n string) string {
	n = strings.TrimSuffix(n, "/")
	n = strings.TrimSuffix(n, "\\\\")
	return n
}

func randomString(len int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + r1.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}
