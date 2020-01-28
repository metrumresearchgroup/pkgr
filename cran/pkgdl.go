package cran

// PkgAndRepoNames returns the strings for package name and repo name
func (pd PkgDl) PkgAndRepoNames() (string, string) {
	return pd.Package.Package, pd.Config.GetOrigin().Name
}
