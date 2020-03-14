# outdated-pkgs-no-update
tags: pkg-update, pkg-outdated, command-flags

## Description

Environment to help test the "--update" [flag or config setting]. The environment is configured for the same packages as [simple](../simple/guide.md), plus the base-R package `Matrix`\*.

\* Note: We included the "Matrix" package here because a user had trouble using this feature with Matrix. It's mainly here as a regression test.

** DO NOT MODIFY THE [outdated-library](outdated-library) DIRECTORY OR YOU WILL DESTROY THIS ENVIRONMENT **

## Expected behavior

You can modify the [pkgr.yml](pkgr.yml) file by setting the `Update` value to `true` or `false`. By default, this test case does not specify, so it should default to FALSE.

* **Update True**:
  - `pkgr plan` indicates that four packages are outdated: `R6`, `pillar`, and `Matrix`, plus `crayon`, which is a dependency of pillar.

  - `pkgr install` will install the most recent versions of `R6`, `pillar`, `Matrix`, and `crayon`, as well as several other dependenceis.

* **Update False**
  - `pkgr plan` will warn that `R6`, `pillar`, `Matrix`, and `crayon` are outdated.
  - `pkgr install` will leave `R6`, `pillar`, `Matrix`, and `crayon` alone (their versions will not be updated) and will only install the dependencies not already installed.
  - `pkgr plan --update` and `pkgr install --update` will behave as if the Update flag was set to "true" in the config file.
