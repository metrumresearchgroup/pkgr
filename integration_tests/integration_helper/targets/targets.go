package targets

type Target int

const (
	MixedSource Target = iota
	Simple
	Threads
	LoggingConfig
	OutdatedPkgs
	DescriptionRepoBug
	OutdatedPkgsNoUpdate
	OutdatedPkgsStress
	RepoOrder
	RepoLocal
	RepoLocalAndRemote
	SimpleSuggests
	Rollback
	RollbackDisabled
	CreateLibrary
	StrictMode
	Misc
	Recommended
	BadCustomizations
	Load
	LoadFail
	TarballInstall
	TarballRollback
	TildeExpansion
	BinariesMac
	BinariesLinux
	PathBug
	EnvVars
)

var (
	ResetTargets = map[Target]string{
		MixedSource:          "test-mixed-source-reset",
		Simple:               "test-simple-reset",
		Threads:              "test-threads-reset",
		LoggingConfig:        "test-logging-config-reset",
		OutdatedPkgs:         "test-outdated-pkgs-reset",
		DescriptionRepoBug:   "test-description-repo-bug-reset",
		OutdatedPkgsNoUpdate: "test-outdated-pkgs-no-update-reset",
		OutdatedPkgsStress:   "test-outdated-pkgs-stress-reset",
		RepoOrder:            "test-repo-order-reset",
		RepoLocal:            "test-repo-local-reset",
		RepoLocalAndRemote:   "test-repo-local-and-remote-reset",
		SimpleSuggests:       "test-simple-suggests-reset",
		Rollback:             "test-rollback-reset",
		RollbackDisabled:     "test-rollback-disabled-reset",
		CreateLibrary:        "test-create-library-reset",
		StrictMode:           "test-strict-mode-reset",
		Misc:                 "test-misc-reset",
		Recommended:          "test-recommended-reset",
		BadCustomizations:    "test-bad-customizations-reset",
		Load:                 "test-load-reset",
		LoadFail:             "test-load-fail-reset",
		TarballInstall:       "test-tarball-install-reset",
		TarballRollback:      "test-tarball-rollback-reset",
		TildeExpansion:       "test-tilde-expansion-reset",
		BinariesMac:          "test-binaries-mac-reset",
		BinariesLinux:        "test-binaries-linux-reset",
		PathBug:              "test-path-bug-reset",
		EnvVars:              "test-env-vars-reset",
	}
)

var RunTargets = map[Target]string{
	BadCustomizations:    "bad-customization",
	BinariesLinux:        "binaries-linux",
	BinariesMac:          "binaries-mac",
	CreateLibrary:        "create-library",
	DescriptionRepoBug:   "description-repo-bug",
	EnvVars:              "env-vars",
	Load:                 "load",
	LoadFail:             "load-fail",
	LoggingConfig:        "logging-config",
	Misc:                 "misc",
	MixedSource:          "mixed-source",
	OutdatedPkgs:         "outdated-pkgs",
	OutdatedPkgsNoUpdate: "outdated-pkgs-no-update",
	OutdatedPkgsStress:   "outdated-pkgs-stress",
	PathBug:              "path-bug",
	Recommended:          "recommended",
	RepoLocal:            "repo-local",
	RepoLocalAndRemote:   "repo-local-and-remote",
	RepoOrder:            "repo-order",
	Rollback:             "rollback",
	RollbackDisabled:     "rollback-disabled",
	Simple:               "simple",
	SimpleSuggests:       "simple-suggests",
	StrictMode:           "strict-mode",
	TarballInstall:       "tarball-install",
	TarballRollback:      "tarball-rollback",
	Threads:              "threads",
	TildeExpansion:       "tilde-expansion",
}
