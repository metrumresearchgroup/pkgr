package desc

// Constraint on version requirement
type Constraint int

func (c Constraint) ToString() string {
	switch c {
	case GT:
		return ">"
	case GTE:
		return ">="
	case LT:
		return "<"
	case LTE:
		return "<="
	case Equals:
		return "=="
	default:
		return "Unknown constraint"
	}
}

// Constraints on package deps
// Least to most constraining
// TODO: Think what the default constraint is, and set it to the zero value.
const (
	None Constraint = iota
	GTE
	GT
	Equals
	LTE
	LT
)

// Version represents a package version
// want to support package version found on CRAN
// pkgstring <- as.character(package_version(
// unique(as.data.frame(available.packages())$Version)))
// max(sapply(stringr::str_split(pkgstring, "\\."), length))
type Version struct {
	Major int
	Minor int
	Patch int
	Dev   int
	// max amount detected on CRAN was 5
	Other int
	// TODO: idiomatically, this is generally called 'raw'
	// So Can Store the Original version in case needed
	String string
}

// Dep represents a dependency
type Dep struct {
	Name       string
	Version    Version
	Constraint Constraint
}

// Desc represents a package description
type Desc struct {
	Package            string
	Source             string
	Version            string
	Maintainer         string
	Description        string
	License            string
	MD5sum             string
	NeedsCompilation   bool
	Path               string
	Priority           string
	Remotes            []string
	OriginalRepository string
	Repository         string
	Imports            map[string]Dep
	Suggests           map[string]Dep
	Depends            map[string]Dep
	LinkingTo          map[string]Dep
	PkgrVersion        string
	PkgrInstallType    string
	PkgrRepositoryURL  string
}

func (d *Desc) GetCombinedDependencies(suggests bool) map[string]Dep {
	combined := map[string]Dep{}

	for key, value := range d.Imports {
		combined[key] = value
	}
	for key, value := range d.LinkingTo {
		combined[key] = value
	}
	if suggests {
		for key, value := range d.Suggests {
			combined[key] = value
		}
	}
	return combined
}

// TODO figure out unmarshalling pattern so can
//  implement that on Desc so don't need intermediate
//  desc struct
type desc struct {
	Package            string
	Source             string
	Version            string
	Maintainer         string
	Description        string
	License            string
	MD5sum             string
	NeedsCompilation   string
	// Path: 4.1.0/Recommended
	// Path: older
	Path               string
	// Priority: recommended
	Priority           string
	Remotes            []string `delim:"," strip:"\n\r\t "`
	OriginalRepository string
	Repository         string
	Imports            []string `delim:"," strip:"\n\r\t "`
	Suggests           []string `delim:"," strip:"\n\r\t "`
	Depends            []string `delim:"," strip:"\n\r\t "`
	LinkingTo          []string `delim:"," strip:"\n\r\t "`
	PkgrVersion        string
	PkgrInstallType    string
	PkgrRepositoryURL  string
}
