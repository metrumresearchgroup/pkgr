package desc

import (
	"regexp"
	"strconv"
)

// ParseVersion parses a package string to a version
func ParseVersion(v string) Version {
	ver := Version{String: v}
	parts := regexp.MustCompile(`[\.-]`).Split(v, 4)
	ver.Major, _ = strconv.Atoi(parts[0])
	ver.Minor, _ = strconv.Atoi(parts[1])
	if len(parts) > 2 {
		ver.Patch, _ = strconv.Atoi(parts[2])
	}
	if len(parts) > 3 {
		ver.Dev, _ = strconv.Atoi(parts[3])

	}
	if len(parts) > 4 {
		ver.Other, _ = strconv.Atoi(parts[4])

	}
	return ver
}

// CompareVersionStrings compares version strings
func CompareVersionStrings(v1 string, v2 string) int {
	return CompareVersions(ParseVersion(v1), ParseVersion(v2))
}

// CompareVersions compares two version numbers
// similar to bytes.Compare, where
// if v1 < v2 returns - 1, equal 0, v1 > v2 1
func CompareVersions(v1 Version, v2 Version) int {
	switch {
	case v1.Major < v2.Major:
		return -1
	case v1.Major > v2.Major:
		return 1
	}
	switch {
	case v1.Minor < v2.Minor:
		return -1
	case v1.Minor > v2.Minor:
		return 1
	}
	switch {
	case v1.Patch < v2.Patch:
		return -1
	case v1.Patch > v2.Patch:
		return 1
	}
	switch {
	case v1.Dev < v2.Dev:
		return -1
	case v1.Dev > v2.Dev:
		return 1
	}
	switch {
	case v1.Other < v2.Other:
		return -1
	case v1.Other > v2.Other:
		return 1
	}

	return 0
}

// Versions represents an array of Version numbers
type Versions []Version

func isLowerVersion(v1 Version, v2 Version) bool {
	if v1.Major < v2.Major {
		return true
	}
	if v1.Minor < v2.Minor {
		return true
	}
	if v1.Dev < v2.Dev {
		return true
	}
	if v1.Other < v2.Other {
		return true
	}
	return false
}
func (v Versions) Len() int           { return len(v) }
func (v Versions) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Versions) Less(i, j int) bool { return isLowerVersion(v[i], v[j]) }
