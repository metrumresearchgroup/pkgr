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
	return ScanROutput(b, false)
}

// Scans lines from RScript output and returns an array with
// the line numbers removed, whitespace trimmed, and (optionally)
// with all input-like lines (which start with ">") excluded.
func ScanROutput(b []byte, outputOnly bool) []string {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	output := []string{}
	re := regexp.MustCompile("^\\[\\d+\\]")
	for scanner.Scan() {
		newLine := strings.TrimSpace(re.ReplaceAllString(scanner.Text(), ""))

		var keepLine bool
		if outputOnly {
			keepLine = newLine != "" && !strings.HasPrefix(newLine, ">")
		} else {
			keepLine = newLine != ""
		}

		if keepLine {
			output = append(output, newLine)
		}
	}
	return output
}