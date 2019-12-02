# IN PROGRESS
# Traceability Matrix: pkgr 1.0.0

## Scope
This traceability matrix links product risk, test-names, and test results to specific user stories
for the proposed software release. User stories, including requirements and test specifications,
are listed in the [Requirements Specification](req_spec.md).


|issue title                                                                                                                                                |risk |test tags           |tests |passed |date run                 |
|:----------------------------------------------------------------------------------------------------------------------------------------------------------|:----|:-------------------|:-----|:------|:------------------------|
|176: Pkgr will behave as closely as possible to the default R Tooling                                                                                      |TODO |install-type        |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|176: Pkgr will behave as closely as possible to the default R Tooling                                                                                      |TODO |dependencies        |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|177: Pkgr will be a command-line tool                                                                                                                      |TODO |na                  |0     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |log-file            |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |install-log         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |log-settings        |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |log-level           |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |repo-customizations |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |pkg-customizations  |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |pkg-update          |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |pkg-outdated        |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |command-flags       |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |basic               |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |suggests            |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |strict-mode         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |cache-local         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |command-flags       |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |TODO |todo                |0     |TODO   |Tue Nov 26 15:38:01 2019 |
|179:  Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions |TODO |cache-local         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|179:  Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions |TODO |cache-system        |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|180: Pkgr installation are idempotent                                                                                                                      |TODO |idempotence         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|181: As as user, I can install R packages in parallel                                                                                                      |TODO |basic               |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|181: As as user, I can install R packages in parallel                                                                                                      |TODO |heavy               |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|181: As as user, I can install R packages in parallel                                                                                                      |TODO |thread-count        |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|182: As a user I can access R packages from multiple repositories                                                                                          |TODO |multi-repo          |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|183: As a user I can specify which repository to pull a package from                                                                                       |TODO |invalid-yml         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|184: As a user, I can specify where packages will be installed.                                                                                            |TODO |local-library       |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|185: As a user I can specify whether or not to install the Suggested dependencies for a specific package or all packages                                   |TODO |suggests            |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|185: As a user I can specify whether or not to install the Suggested dependencies for a specific package or all packages                                   |TODO |basic               |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|186: As a user I can specify the type of package to install (binary vs src)                                                                                |TODO |basic               |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|186: As a user I can specify the type of package to install (binary vs src)                                                                                |TODO |pkg-customizations  |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|187: As a user, I can specify the package type for a given repository                                                                                      |TODO |repo-customizations |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|188: As a user, I can view the difference between the current state and the state defined in pkgr.yml without altering current state                       |TODO |plan                |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|188: As a user, I can view the difference between the current state and the state defined in pkgr.yml without altering current state                       |TODO |pkg-outdated        |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|189: As a user, I can see how installation dependencies were identified                                                                                    |TODO |inspect             |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|190: As a user, I can invalidate pkgr’s caches at will                                                                                                     |TODO |clean-cache         |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|190: As a user, I can invalidate pkgr’s caches at will                                                                                                     |TODO |clean-pkgdb         |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|191:  As a user, when installing packages, I can control whether to update previously installed packages to new versions or only install missing packages. |TODO |pkg-outdated        |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|191:  As a user, when installing packages, I can control whether to update previously installed packages to new versions or only install missing packages. |TODO |pkg-update          |2     |TODO   |Tue Nov 26 15:38:01 2019 |
|192: As a user, if an installation fails, I can be returned to the latest state before the installation started                                            |TODO |rollback            |1     |TODO   |Tue Nov 26 15:38:01 2019 |
|193: As a user, I can identify which packages have been installed by pkgr                                                                                  |TODO |existing-pkgs       |1     |TODO   |Tue Nov 26 15:38:01 2019 |
