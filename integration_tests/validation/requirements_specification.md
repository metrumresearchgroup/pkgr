# Unofficial
# Requirements Specification: pkgr 1.0.0

## Scope
The purpose of this document is to define the requirements for pkgr and document which tests have been created to ensure that individual requirements are met.
---

## 176: Pkgr will behave as closely as possible to the default R Tooling

### Risk: medium

### Summary
* Should install src versus binary in a platform-specific manner
	* Tests:
		* install-type
* Should install packages in proper order of dependency requirements
	* Tests:
		* dependencies

|Test Tag     |Location                       |File     |Link Test                                       |Link Results                                                   |
|:------------|:------------------------------|:--------|:-----------------------------------------------|:--------------------------------------------------------------|
|dependencies |../../integration_tests/simple |guide.md |[Test](../../integration_tests/simple/guide.md) |[Result](../../integration_tests/validation/simple/results.md) |
|install-type |../../integration_tests/simple |guide.md |[Test](../../integration_tests/simple/guide.md) |[Result](../../integration_tests/validation/simple/results.md) |



---
## 177: Pkgr will be a command-line tool

### Risk: low

### Summary
* All of pkgr’s functions should be executable from a command-line context.
pkgr will read its instructions from a combination of a .yml file (pkgr.yml) and command-line arguments.

|Test Tag |Location |File |Link Test |Link Results |
|:--------|:--------|:----|:---------|:------------|



---
## 178: Pkgr behavior will be controlled by a yml file with certain elements overridable by command line flag

### Risk: medium

### Summary
* Packages to install
Repositories to use
Src/binary configuration for packages
Src/binary configuration for repositories
Libpaths [[ unclear on what this does ]]
Rpath
Cache directory
Lockfile settings (should pgkr make itself compatible with a Lockfile such as renv.lock?)
	* Tests:
		* basic, cache-local, pkg-customizations, lockfile, log-file, install-log
* The following options should be configurable statically in the yml file or dynamically via the command-line (with command line arguments always overwriting yml configurations):
* Library (directory) to install packages to
Update settings (should pkgr attempt to update outdated packages?)
Logging level
Installation of “Suggested” packages config
Number of threads to use
Strict mode enabled/disabled
	* Tests:
		* pkg-update, strict-mode, local-library, log-level, thread-count, suggests
* The following settings should be configurable only via command-line flags
* config (path to yaml file, defaults to pkgr.yml)
preview mode [[ unclear on what this does ]]
Debug mode
	* Tests:
		* command-flags
* If nothing is set, these are the default settings:
* Debug mode: disabled
Preview mode: disabled
Strict mode: disabled
Loglevel: “info”
Rpath: “R”
Threads: Number of processors available on system
Cache: System temp directories
Install suggested packages: false
Update outdated packages: false
Src/binary configurations: System defaults (bin for Mac/Windows, src for Linux)
	* Tests:
		* basic, todo
* Minimum required user-settings:
* Version (pkgr)
Library
Packages to install
Repositories to use
	* Tests:
		* todo

|Test Tag            |Location                                        |File     |Link Test                                                        |Link Results                                                                    |
|:-------------------|:-----------------------------------------------|:--------|:----------------------------------------------------------------|:-------------------------------------------------------------------------------|
|log-file            |../../integration_tests/logging-config          |guide.md |[Test](../../integration_tests/logging-config/guide.md)          |[Result](../../integration_tests/validation/logging-config/results.md)          |
|install-log         |../../integration_tests/logging-config          |guide.md |[Test](../../integration_tests/logging-config/guide.md)          |[Result](../../integration_tests/validation/logging-config/results.md)          |
|log-settings        |../../integration_tests/logging-config          |guide.md |[Test](../../integration_tests/logging-config/guide.md)          |[Result](../../integration_tests/validation/logging-config/results.md)          |
|log-level           |../../integration_tests/logging-config          |guide.md |[Test](../../integration_tests/logging-config/guide.md)          |[Result](../../integration_tests/validation/logging-config/results.md)          |
|repo-customizations |../../integration_tests/mixed-source            |guide.md |[Test](../../integration_tests/mixed-source/guide.md)            |[Result](../../integration_tests/validation/mixed-source/results.md)            |
|pkg-customizations  |../../integration_tests/mixed-source            |guide.md |[Test](../../integration_tests/mixed-source/guide.md)            |[Result](../../integration_tests/validation/mixed-source/results.md)            |
|pkg-update          |../../integration_tests/outdated-pkgs-no-update |guide.md |[Test](../../integration_tests/outdated-pkgs-no-update/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs-no-update/results.md) |
|pkg-outdated        |../../integration_tests/outdated-pkgs-no-update |guide.md |[Test](../../integration_tests/outdated-pkgs-no-update/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs-no-update/results.md) |
|command-flags       |../../integration_tests/outdated-pkgs-no-update |guide.md |[Test](../../integration_tests/outdated-pkgs-no-update/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs-no-update/results.md) |
|pkg-update          |../../integration_tests/outdated-pkgs           |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md)           |[Result](../../integration_tests/validation/outdated-pkgs/results.md)           |
|pkg-outdated        |../../integration_tests/outdated-pkgs           |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md)           |[Result](../../integration_tests/validation/outdated-pkgs/results.md)           |
|suggests            |../../integration_tests/simple-suggests         |guide.md |[Test](../../integration_tests/simple-suggests/guide.md)         |[Result](../../integration_tests/validation/simple-suggests/results.md)         |
|cache-local         |../../integration_tests/simple-suggests         |guide.md |[Test](../../integration_tests/simple-suggests/guide.md)         |[Result](../../integration_tests/validation/simple-suggests/results.md)         |
|basic               |../../integration_tests/simple                  |guide.md |[Test](../../integration_tests/simple/guide.md)                  |[Result](../../integration_tests/validation/simple/results.md)                  |
|strict-mode         |../../integration_tests/strict-mode             |guide.md |[Test](../../integration_tests/strict-mode/guide.md)             |[Result](../../integration_tests/validation/strict-mode/results.md)             |



---
## 179:  Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions

### Risk: medium

### Summary
* Pkgr caches downloaded packages and stores them in a temp folder for faster future installations.
	* Tests:
		* cache-local, cache-system
* Pkgr maintains a cache of Package Databases (pkgdb) in a temp folder. Pkgdbs contain the information in the PACKAGES file for all repos listed in a pkgr.yml file. When pkgdbs are available in the cache, pkgr will not have to re-parse the PACKAGES file for those repos.
	* Tests:
		* cache-system
* Users should be able to specify where the cache is created. If users do not specify, the cache should use System temp folders.
	* Tests:
		* cache-local, cache-system

|Test Tag     |Location                                |File     |Link Test                                                |Link Results                                                            |
|:------------|:---------------------------------------|:--------|:--------------------------------------------------------|:-----------------------------------------------------------------------|
|cache-local  |../../integration_tests/simple-suggests |guide.md |[Test](../../integration_tests/simple-suggests/guide.md) |[Result](../../integration_tests/validation/simple-suggests/results.md) |
|cache-system |../../integration_tests/simple          |guide.md |[Test](../../integration_tests/simple/guide.md)          |[Result](../../integration_tests/validation/simple/results.md)          |



---
## 180: Pkgr installation are idempotent

### Risk: high

### Summary
* Multiple executions of pkgr install with the same pkgr.yml file should all end in the same result, no matter how many times pkgr install has been run.
* Caveat, this is assuming nothing is changing on the repository side (e.g. packages being updated)
Tests:
* idempotence

|Test Tag    |Location                     |File     |Link Test                                     |Link Results                                                 |
|:-----------|:----------------------------|:--------|:---------------------------------------------|:------------------------------------------------------------|
|idempotence |../../integration_tests/misc |guide.md |[Test](../../integration_tests/misc/guide.md) |[Result](../../integration_tests/validation/misc/results.md) |



---
## 181: As as user, I can install R packages in parallel

### Risk: medium

### Summary
* In order to more quickly stand-up/change R environments, we want pkgr to be able to perform its installation operations in parallel.
	* Tests:
		* basic, heavy, thread-count

|Test Tag     |Location                             |File     |Link Test                                             |Link Results                                                         |
|:------------|:------------------------------------|:--------|:-----------------------------------------------------|:--------------------------------------------------------------------|
|heavy        |../../integration_tests/mixed-source |guide.md |[Test](../../integration_tests/mixed-source/guide.md) |[Result](../../integration_tests/validation/mixed-source/results.md) |
|basic        |../../integration_tests/simple       |guide.md |[Test](../../integration_tests/simple/guide.md)       |[Result](../../integration_tests/validation/simple/results.md)       |
|thread-count |../../integration_tests/threads      |guide.md |[Test](../../integration_tests/threads/guide.md)      |[Result](../../integration_tests/validation/threads/results.md)      |



---
## 182: As a user I can access R packages from multiple repositories

### Risk: medium

### Summary
* Users should be able to set up their R environment via pkgr with packages installed from as many unique CRAN-like repositories as they need.
	* Tests:
		* multi-repo

|Test Tag   |Location                             |File     |Link Test                                             |Link Results                                                         |
|:----------|:------------------------------------|:--------|:-----------------------------------------------------|:--------------------------------------------------------------------|
|multi-repo |../../integration_tests/mixed-source |guide.md |[Test](../../integration_tests/mixed-source/guide.md) |[Result](../../integration_tests/validation/mixed-source/results.md) |



---
## 183: As a user I can specify which repository to pull a package from

### Risk: medium

### Summary
* Users can use the Customizations field in pkgr.yml to specify which repository a package should come from.
	* Tests:
		* invalid-yml
* Repository must be defined in the “Repos” section of pkgr.yml.

|Test Tag    |Location                                  |File     |Link Test                                                  |Link Results                                                              |
|:-----------|:-----------------------------------------|:--------|:----------------------------------------------------------|:-------------------------------------------------------------------------|
|invalid-yml |../../integration_tests/bad-customization |guide.md |[Test](../../integration_tests/bad-customization/guide.md) |[Result](../../integration_tests/validation/bad-customization/results.md) |



---
## 184: As a user, I can specify where packages will be installed.

### Risk: medium

### Summary
* Users must be able to define in pkgr.yml where they would like packages to be installed to via a “Library” field in pkgr.yml
	* Tests:
		* local-library

|Test Tag      |Location                       |File     |Link Test                                       |Link Results                                                   |
|:-------------|:------------------------------|:--------|:-----------------------------------------------|:--------------------------------------------------------------|
|local-library |../../integration_tests/simple |guide.md |[Test](../../integration_tests/simple/guide.md) |[Result](../../integration_tests/validation/simple/results.md) |



---
## 185: As a user I can specify whether or not to install the Suggested dependencies for a specific package or all packages

### Risk: medium

### Summary
* Default should be to NOT install suggested when nothing specified
	* Tests:
		* basic
* Cmd line argument should be available to overwrite argument in pkgr.yml
	* Tests:
		* suggests
* Pkgr can install suggested dependencies
	* Tests:
		* suggests

|Test Tag |Location                                |File     |Link Test                                                |Link Results                                                            |
|:--------|:---------------------------------------|:--------|:--------------------------------------------------------|:-----------------------------------------------------------------------|
|suggests |../../integration_tests/simple-suggests |guide.md |[Test](../../integration_tests/simple-suggests/guide.md) |[Result](../../integration_tests/validation/simple-suggests/results.md) |
|basic    |../../integration_tests/simple          |guide.md |[Test](../../integration_tests/simple/guide.md)          |[Result](../../integration_tests/validation/simple/results.md)          |



---
## 186: As a user I can specify the type of package to install (binary vs src)

### Risk: medium

### Summary
* If nothing is specified, should use system default.
	* Tests:
		* basic
* Specifiable at a package-level via Customizations
	* Tests:
		* pkg-customizations

|Test Tag           |Location                             |File     |Link Test                                             |Link Results                                                         |
|:------------------|:------------------------------------|:--------|:-----------------------------------------------------|:--------------------------------------------------------------------|
|pkg-customizations |../../integration_tests/mixed-source |guide.md |[Test](../../integration_tests/mixed-source/guide.md) |[Result](../../integration_tests/validation/mixed-source/results.md) |
|basic              |../../integration_tests/simple       |guide.md |[Test](../../integration_tests/simple/guide.md)       |[Result](../../integration_tests/validation/simple/results.md)       |



---
## 187: As a user, I can specify the package type for a given repository   

### Risk: high

### Summary
* Specifiable at a repo-level via Customizations
	* Tests:
		* repo-customizations

|Test Tag            |Location                             |File     |Link Test                                             |Link Results                                                         |
|:-------------------|:------------------------------------|:--------|:-----------------------------------------------------|:--------------------------------------------------------------------|
|repo-customizations |../../integration_tests/mixed-source |guide.md |[Test](../../integration_tests/mixed-source/guide.md) |[Result](../../integration_tests/validation/mixed-source/results.md) |



---
## 188: As a user, I can view the difference between the current state and the state defined in pkgr.yml without altering current state

### Risk: medium

### Summary
* Command must be available showing exactly what pkgr will download and install, as well as give information about the current state, without actually applying these changes.
	* Tests:
		* plan
pkg-outdated
* This includes the construction of dependency trees.
	* Tests:
		* plan

|Test Tag     |Location                                        |File     |Link Test                                                        |Link Results                                                                    |
|:------------|:-----------------------------------------------|:--------|:----------------------------------------------------------------|:-------------------------------------------------------------------------------|
|pkg-outdated |../../integration_tests/outdated-pkgs-no-update |guide.md |[Test](../../integration_tests/outdated-pkgs-no-update/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs-no-update/results.md) |
|pkg-outdated |../../integration_tests/outdated-pkgs           |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md)           |[Result](../../integration_tests/validation/outdated-pkgs/results.md)           |
|plan         |../../integration_tests/outdated-pkgs           |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md)           |[Result](../../integration_tests/validation/outdated-pkgs/results.md)           |



---
## 189: As a user, I can see how installation dependencies were identified

### Risk: low

### Summary
* There must be a command users can use to view the dependency tree constructed by pkgr.
	* Tests:
		* inspect

|Test Tag |Location                                |File     |Link Test                                                |Link Results                                                            |
|:--------|:---------------------------------------|:--------|:--------------------------------------------------------|:-----------------------------------------------------------------------|
|inspect  |../../integration_tests/simple-suggests |guide.md |[Test](../../integration_tests/simple-suggests/guide.md) |[Result](../../integration_tests/validation/simple-suggests/results.md) |
|inspect  |../../integration_tests/simple          |guide.md |[Test](../../integration_tests/simple/guide.md)          |[Result](../../integration_tests/validation/simple/results.md)          |



---
## 190: As a user, I can invalidate pkgr’s caches at will

### Risk: high

### Summary
* Users should be able to invalidate the download cache, the package db cache, or both as needed.
	* Tests:
		* clean-cache, clean-pkgdb
* Pkgr should be able to remove only those items in the package database cache that are relevant to the in-use pkgr.yml file.
	* Tests:
		* clean-pkgdb
* In theory, removing package caches shouldn’t actually change the end result of a pkgr install run. However, whenever caching is involved, it’s possible for irregularities to enter the process. Either way, pkgr should not be reliant on the presence of information in a cache to achieve the state described in pkgr.yml.
	* Tests:
		* cache-partial, cache-extraneous

|Test Tag    |Location                       |File     |Link Test                                       |Link Results                                                   |
|:-----------|:------------------------------|:--------|:-----------------------------------------------|:--------------------------------------------------------------|
|clean-pkgdb |../../integration_tests/misc   |guide.md |[Test](../../integration_tests/misc/guide.md)   |[Result](../../integration_tests/validation/misc/results.md)   |
|clean-cache |../../integration_tests/simple |guide.md |[Test](../../integration_tests/simple/guide.md) |[Result](../../integration_tests/validation/simple/results.md) |
|clean-pkgdb |../../integration_tests/simple |guide.md |[Test](../../integration_tests/simple/guide.md) |[Result](../../integration_tests/validation/simple/results.md) |



---
## 191:  As a user, when installing packages, I can control whether to update previously installed packages to new versions or only install missing packages.

### Risk: medium

### Summary
* In setups where the repositories listed in pkgr.yml can change, it is possible for new versions of packages to become available outside of the user’s environment. In these cases, we need to provide users options on whether or not they wish to update.
	* Tests:
		* pkg-outdated
pkg-update
* Should be specified via pkgr.yml global option or via command line argument (cmd line argument takes priority)
	* Tests:
		* pkg-outdated
pkg-update
* Default behavior should be to not update
	* Tests:
		* pkg-outdated
pkg-update

|Test Tag     |Location                                        |File     |Link Test                                                        |Link Results                                                                    |
|:------------|:-----------------------------------------------|:--------|:----------------------------------------------------------------|:-------------------------------------------------------------------------------|
|pkg-update   |../../integration_tests/outdated-pkgs-no-update |guide.md |[Test](../../integration_tests/outdated-pkgs-no-update/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs-no-update/results.md) |
|pkg-outdated |../../integration_tests/outdated-pkgs-no-update |guide.md |[Test](../../integration_tests/outdated-pkgs-no-update/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs-no-update/results.md) |
|pkg-update   |../../integration_tests/outdated-pkgs           |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md)           |[Result](../../integration_tests/validation/outdated-pkgs/results.md)           |
|pkg-outdated |../../integration_tests/outdated-pkgs           |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md)           |[Result](../../integration_tests/validation/outdated-pkgs/results.md)           |



---
## 192: As a user, if an installation fails, I can be returned to the latest state before the installation started

### Risk: high

### Summary
* While rolling back, Pkgr will not touch pkgs that weren’t installed via pkgr
	* Tests:
		* rollback
* When a rollback occurs, the end state of the Library folder needs to be exactly the same as the start state (all packages and their versions must match completely)
	* Tests:
		* rollback

|Test Tag |Location                         |File     |Link Test                                         |Link Results                                                     |
|:--------|:--------------------------------|:--------|:-------------------------------------------------|:----------------------------------------------------------------|
|rollback |../../integration_tests/rollback |guide.md |[Test](../../integration_tests/rollback/guide.md) |[Result](../../integration_tests/validation/rollback/results.md) |



---
## 193: As a user, I can identify which packages have been installed by pkgr

### Risk: low

### Summary
* Nothing about using pgkr actually prevents users from installing packages the “normal” way (install.packages). Because of this, it is possible for the user to get their environment into a state that pkgr is not configured to replicate. To remedy this, we need to show users which packages were and were not installed by pkgr. This lets users know that their environment may not be totally recreatable with just the pkgr.yml file.
	* Tests:
		* existing-pkgs
* This information should be displayed automatically as part of the “plan” and “install” commands.
	* Tests:
		* existing-pkgs

|Test Tag      |Location                              |File     |Link Test                                              |Link Results                                                          |
|:-------------|:-------------------------------------|:--------|:------------------------------------------------------|:---------------------------------------------------------------------|
|existing-pkgs |../../integration_tests/outdated-pkgs |guide.md |[Test](../../integration_tests/outdated-pkgs/guide.md) |[Result](../../integration_tests/validation/outdated-pkgs/results.md) |



---
