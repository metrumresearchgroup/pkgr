package desc

import "fmt"

// ToString provides a matching pattern needed for outputting the Constraint
// R (>= 3.6)
// <Name> (<constraint> <version>)
func (d Dep) ToString() string {
	return fmt.Sprintf("%s (%s %s)", d.Name, d.Constraint.ToString(), d.Version.String)
}

// R (>= 3.6)
