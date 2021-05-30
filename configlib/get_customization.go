package configlib

// GetRepoCustomizationByName will return a Repository customization by name as well
// as whether it exists. This is helpful given the customization config of array of maps
// makes it more annoying than it should be to get one.
func GetRepoCustomizationByName(nm string, c Customizations) (RepoConfig, bool){
	for _, rc := range c.Repos {
		rv, ok := rc[nm]
		if !ok {
			continue
		}
		return rv, true
	}
	return RepoConfig{}, false
}

// GetPackageCustomizationByName will return a Package customization by name as well
// as whether it exists. This is helpful given the customization config of array of maps
// makes it more annoying than it should be to get one.
func GetPackageCustomizationByName(nm string, c Customizations) (PkgConfig, bool){
	for _, rc := range c.Packages {
		rv, ok := rc[nm]
		if !ok {
			continue
		}
		return rv, true
	}
	return PkgConfig{}, false
}
