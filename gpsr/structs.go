package gpsr

//InstallPlan provides metadata around an installation plan
type InstallPlan struct {
	StartingPackages []string
	DepDb            map[string][]string
}
