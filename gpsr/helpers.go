package gpsr

import "github.com/dpastoor/rpackagemanager/desc"

func isDefaultPackage(pkg string) bool {
	_, exists := DefaultPackages[pkg]
	return exists
}

func appendToGraph(m Graph, d desc.Desc, dmap map[string]desc.Desc) {
	var reqs []string
	for r := range d.Imports {
		_, ok := dmap[r]
		if ok && !isDefaultPackage(r) {
			reqs = append(reqs, r)
		}
	}
	for r := range d.Depends {
		_, ok := dmap[r]
		if ok && !isDefaultPackage(r) {
			reqs = append(reqs, r)
		}
	}
	for r := range d.LinkingTo {
		_, ok := dmap[r]
		if ok && !isDefaultPackage(r) {
			reqs = append(reqs, r)
		}
	}
	m[d.Package] = NewNode(d.Package, reqs)
	if len(reqs) > 0 {
		for _, pn := range reqs {
			_, ok := m[pn]
			if pn != "R" && !ok {
				appendToGraph(m, dmap[pn], dmap)
			}
		}
	}
}
