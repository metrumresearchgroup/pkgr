# outdated-pkgs
tags: pkg-update, pkg-outdated, plan, existing-pkgs

## Description

Environment to help test the "--update" [flag or config setting]. The environment is configured for the same packages as [simple](../simple/guide.md), plus the base-R package `Matrix`\*.

\* Note: We included the "Matrix" package here because a user had trouble using this feature with Matrix. It's mainly here as a regression test.

** DO NOT MODIFY THE [outdated-library](outdated-library) DIRECTORY OR YOU WILL DESTROY THIS ENVIRONMENT **

## How to quickly reset test environment
For convenience, you can use this script to reset your test environment and come back to this directory. Run the following command from a terminal window in this directory:
`cd ../ && make test-setup && cd outdated-pkgs`

Please note that this will build and install pkgr using the source code in this project.

## Expected behavior

You can modify the [pkgr.yml](pkgr.yml) file by setting the `Update` value to `true` or `false`. For this test case, try both settings. You can also override the settings in pkgr.yml by passing the `--update` flag. Do that as well for these test cases. Remember to reset the test environment each time you run `pkgr install`

* **All cases**
  - Output for `pkgr plan` will contain the following two lines:
```
INFO[0000] found installed packages                      count=4
WARN[0000] Packages not installed by pkgr                packages="[crayon Matrix]"
```

* **Update True**:
  - `pkgr plan` indicates that four packages will be updated: `R6`, `pillar`, and `Matrix`, plus `crayon`, which is a dependency of pillar.
```
INFO[0001] package will be updated                       installed_version=1.2.1 pkg=crayon update_version=1.3.4
INFO[0001] package will be updated                       installed_version=2.0 pkg=R6 update_version=2.4.0
INFO[0001] package will be updated                       installed_version=1.2.1 pkg=pillar update_version=1.3.1
INFO[0001] package will be updated                       installed_version=1.2-11 pkg=Matrix update_version=1.2-17
```

  - `pkgr install` will install the most recent versions (as of 2019-08-14) of `R6`, `pillar`, `Matrix`, and `crayon`, as well as several other dependencies.

* **Update False**
  - `pkgr plan` will warn that `R6`, `pillar`, `Matrix`, and `crayon` are outdated.
  - `pkgr install` will leave `R6`, `pillar`, `Matrix`, and `crayon` alone (their versions will not be updated) and will only install the dependencies not already installed.
  - `pkgr plan --update` and `pkgr install --update` will behave as if the Update flag was set in the config file.
