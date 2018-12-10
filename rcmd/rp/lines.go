package rp

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

// ScanLines scans lines from Rscript output and returns an array with
// the line numbers removed and whitespace trimmed
func ScanLines(b []byte) []string {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	output := []string{}
	re := regexp.MustCompile("^\\[\\d+\\]")
	for scanner.Scan() {
		newLine := strings.TrimSpace(re.ReplaceAllString(scanner.Text(), ""))
		if newLine != "" {
			output = append(output, newLine)
		}
	}
	return output
}
