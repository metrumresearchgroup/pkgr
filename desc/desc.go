package desc

import (
	"fmt"
	"os"
	"strings"

	"pault.ag/go/debian/control"
)

// ParseDep parses the dep to a struct
// R package names follow .standard_regexps(), which
// shows that it should be at minimum digit+.digit+
// likewise versions must only be separated by . or -
// so 1.1 and 1-1 are fine. Likewise can mix 1.1-2 1.1-3
// however no other characters allowed so semver
// cannot be followed
func ParseDep(d string) Dep {
	var dep Dep
	sp := strings.Split(d, "(")
	dep.Name = strings.TrimSpace(sp[0])
	if len(sp) > 1 {
		pv := strings.Replace(sp[1], ")", "", 1)
		if strings.Contains(pv, "==") {
			dep.Constraint = Equals
			dep.Version = ParseVersion(strings.TrimSpace(strings.Replace(pv, "==", "", 1)))
		} else if strings.Contains(pv, ">=") {
			dep.Constraint = GTE
			dep.Version = ParseVersion(strings.TrimSpace(strings.Replace(pv, ">=", " ", 1)))
		} else if strings.Contains(pv, ">") {
			dep.Constraint = GT
			dep.Version = ParseVersion(strings.TrimSpace(strings.Replace(pv, ">", " ", 1)))
		} else {
			panic(pv)
		}
	}
	return dep
}

// NewDesc creates a new Description
func NewDesc(d desc) Desc {
	dsc := Desc{
		Package:     d.Package,
		Source:      d.Source,
		Version:     d.Version,
		Maintainer:  d.Maintainer,
		Description: d.Description,
		MD5sum:      d.MD5sum,
		Remotes:     d.Remotes,
		Repository:  d.Repository,
		Imports:     make(map[string]Dep),
		Suggests:    make(map[string]Dep),
		Depends:     make(map[string]Dep),
		// the reason linkingTo is not also a map is it
		// will only contain the package name its linking to
		// since the imports/depends field will give more
		// information about the dependency
		LinkingTo: d.LinkingTo,
	}
	if len(d.Imports) > 0 {
		for _, dp := range d.Imports {
			dep := ParseDep(dp)
			dsc.Imports[dep.Name] = dep
		}
	}
	if len(d.Suggests) > 0 {
		for _, dp := range d.Suggests {
			dep := ParseDep(dp)
			dsc.Suggests[dep.Name] = dep
		}

	}
	if len(d.Depends) > 0 {
		for _, dp := range d.Depends {
			dep := ParseDep(dp)
			dsc.Depends[dep.Name] = dep
		}
	}
	return dsc
}

// ReadDesc reads the DESCRIPTION file
func ReadDesc(p string) (Desc, error) {
	var dsc desc
	f, err := os.Open(p)
	if err != nil {
		return NewDesc(dsc), err
	}
	defer f.Close()
	err = control.Unmarshal(&dsc, f)
	fmt.Println(dsc)
	fmt.Println(err)
	fmt.Println("remotes: ", dsc.Remotes)
	if err != nil {
		return NewDesc(dsc), err
	}
	return NewDesc(dsc), nil
}
