package pkgreqs

// ReqTypes defines the requirement types for a package
type ReqTypes struct {
	Depends   bool
	Imports   bool
	Suggests  bool
	LinkingTo bool
}
type PkgService interface {
	GetReqsByName(nm string) ([]string, error)
	GetAllReqs(nms []string, r ReqTypes) ([]string, error)
}
