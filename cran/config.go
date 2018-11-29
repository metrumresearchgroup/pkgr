package cran

// NewPkgConfigDB initializes a PkgConfig map
func NewPkgConfigDB() map[string]PkgConfig {
	return make(map[string]PkgConfig)
}
