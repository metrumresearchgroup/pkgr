package cran

import "fmt"

// ToFullString provides a string representation of the Rversion
func (rv RVersion) ToFullString() string {
	return fmt.Sprintf("%v.%v.%v", rv.Major, rv.Minor, rv.Patch)
}

// ToString provides the major/minor version of R
// TODO: So long as your canonical output is this string, you may simply call this "String()"
func (rv RVersion) ToString() string {
	return fmt.Sprintf("%v.%v", rv.Major, rv.Minor)
}
