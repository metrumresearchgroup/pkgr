
# Traceability Matrix: pkgr 1.0.0

## Scope
This traceability matrix links product risk, test-names, and test results to specific user stories
for the proposed software release. User stories, including requirements and test specifications,
are listed in the [Requirements Specification](req_spec.md).


|issue title                                                                                                                                                |risk   |test tags           |tests | passed|date run               |
|:----------------------------------------------------------------------------------------------------------------------------------------------------------|:------|:-------------------|:-----|------:|:----------------------|
|176: Pkgr will behave as closely as possible to the default R Tooling                                                                                      |medium |install-type        |1     |      1|12-02-2019             |
|176: Pkgr will behave as closely as possible to the default R Tooling                                                                                      |medium |dependencies        |1     |      1|12-02-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |repo-customizations |1     |      1|12-03-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |pkg-customizations  |1     |      1|12-03-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |pkg-update          |2     |      2|12-03-2019, 12-03-2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |pkg-outdated        |2     |      2|12-03-2019, 12-03-2019 |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |command-flags       |1     |      1|12-03-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |basic               |1     |      1|12-02-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |suggests            |1     |      1|12-02-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |strict-mode         |1     |      1|12-03-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |cache-local         |1     |      1|12-02-2019             |
|178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag                                                 |medium |command-flags       |1     |      1|12-03-2019             |
|179:  Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions |medium |cache-local         |1     |      1|12-02-2019             |
|179:  Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions |medium |cache-system        |1     |      1|12-02-2019             |
|180: Pkgr installation are idempotent                                                                                                                      |high   |idempotence         |1     |      1|12-03-2019             |
|181: As as user, I can install R packages in parallel                                                                                                      |medium |basic               |1     |      1|12-02-2019             |
|181: As as user, I can install R packages in parallel                                                                                                      |medium |heavy               |1     |      1|12-03-2019             |
|181: As as user, I can install R packages in parallel                                                                                                      |medium |thread-count        |1     |      1|12-03-2019             |
|182: As a user I can access R packages from multiple repositories                                                                                          |medium |multi-repo          |1     |      1|12-03-2019             |
|183: As a user I can specify which repository to pull a package from                                                                                       |medium |invalid-yml         |1     |      1|12-02-2019             |
|184: As a user, I can specify where packages will be installed.                                                                                            |medium |local-library       |1     |      1|12-02-2019             |
|185: As a user I can specify whether or not to install the Suggested dependencies for a specific package or all packages                                   |medium |suggests            |1     |      1|12-02-2019             |
|185: As a user I can specify whether or not to install the Suggested dependencies for a specific package or all packages                                   |medium |basic               |1     |      1|12-02-2019             |
|186: As a user I can specify the type of package to install (binary vs src)                                                                                |medium |basic               |1     |      1|12-02-2019             |
|186: As a user I can specify the type of package to install (binary vs src)                                                                                |medium |pkg-customizations  |1     |      1|12-03-2019             |
|187: As a user, I can specify the package type for a given repository                                                                                      |high   |repo-customizations |1     |      1|12-03-2019             |
|188: As a user, I can view the difference between the current state and the state defined in pkgr.yml without altering current state                       |medium |plan                |1     |      1|12-03-2019             |
|188: As a user, I can view the difference between the current state and the state defined in pkgr.yml without altering current state                       |medium |pkg-outdated        |2     |      2|12-03-2019, 12-03-2019 |
|189: As a user, I can see how installation dependencies were identified                                                                                    |low    |inspect             |2     |      2|12-02-2019, 12-02-2019 |
|190: As a user, I can invalidate pkgr’s caches at will                                                                                                     |high   |clean-cache         |1     |      1|12-02-2019             |
|190: As a user, I can invalidate pkgr’s caches at will                                                                                                     |high   |clean-pkgdb         |2     |      2|12-03-2019, 12-02-2019 |
|190: As a user, I can invalidate pkgr’s caches at will                                                                                                     |high   |automated           |1     |      1|12-20-2019             |
|191:  As a user, when installing packages, I can control whether to update previously installed packages to new versions or only install missing packages. |medium |pkg-outdated        |2     |      2|12-03-2019, 12-03-2019 |
|191:  As a user, when installing packages, I can control whether to update previously installed packages to new versions or only install missing packages. |medium |pkg-update          |2     |      2|12-03-2019, 12-03-2019 |
|191:  As a user, when installing packages, I can control whether to update previously installed packages to new versions or only install missing packages. |medium |automated           |1     |      1|12-20-2019             |
|192: As a user, if an installation fails, I can be returned to the latest state before the installation started                                            |high   |rollback            |1     |      1|12-03-2019             |
|192: As a user, if an installation fails, I can be returned to the latest state before the installation started                                            |high   |automated           |1     |      1|12-20-2019             |
|193: As a user, I can identify which packages have been installed by pkgr                                                                                  |low    |existing-pkgs       |1     |      1|12-03-2019             |
|203: As a user, I can use pkgr to automatically setup my library for Packrat and Renv                                                                      |medium |automated           |1     |      1|12-20-2019             |
|204: As a user, I can add and remove packages from the config file via command line                                                                        |medium |automated           |1     |      1|12-20-2019             |
|extra test                                                                                                                                                 |N/A   |create-library      |N/A  |      1|12-02-2019             |
|extra test                                                                                                                                                 |N/A   |bug-duplicate-repo  |N/A  |      1|12-03-2019             |
|extra test                                                                                                                                                 |N/A   |cache-extraneous    |N/A  |      1|12-03-2019             |
|extra test                                                                                                                                                 |N/A   |cache-partial       |N/A  |      1|12-03-2019             |
|extra test                                                                                                                                                 |N/A   |repo-order          |N/A  |      1|12-03-2019             |
|extra test                                                                                                                                                 |N/A   |repo-remote         |N/A  |      1|12-03-2019             |
|extra test                                                                                                                                                 |N/A   |repo-local          |N/A  |      2|12-03-2019, 12-03-2019 |
