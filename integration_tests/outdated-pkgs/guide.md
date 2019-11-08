# outdated-pkgs
tags: pkg-update, pkg-outdated

## Description

Environment to help test the "--update" [flag or config setting]. The environment is configured for the same packages as [simple](../simple/guide.md), plus the base-R package `Matrix`\*.

\* Note: We included the "Matrix" package here because a user had trouble using this feature with Matrix. It's mainly here as a regression test.

** DO NOT MODIFY THE [outdated-library](outdated-library) DIRECTORY OR YOU WILL DESTROY THIS ENVIRONMENT **

## Expected behavior

You can modify the [pkgr.yml](pkgr.yml) file by setting the `Update` value to `true` or `false`.

* **Update True**:
  - `pkgr plan` indicates that four packages will be updated: `R6`, `pillar`, and `Matrix`, plus `crayon`, which is a dependency of pillar.

```
INFO[0001] package will be updated                       installed_version=1.2.1 pkg=crayon update_version=1.3.4
INFO[0001] package will be updated                       installed_version=2.0 pkg=R6 update_version=2.4.0
INFO[0001] package will be updated                       installed_version=1.2.1 pkg=pillar update_version=1.3.1
INFO[0001] package will be updated                       installed_version=1.2-11 pkg=Matrix update_version=1.2-17
```

  - `pkgr install` will install the most recent versions of `R6`, `pillar`, `Matrix`, and `crayon`, as well as several other dependenceis.

* **Update False**
  - `pkgr plan` will warn that `R6`, `pillar`, `Matrix`, and `crayon` are outdated.
  - `pkgr install` will leave `R6`, `pillar`, `Matrix`, and `crayon` alone (their versions will not be updated) and will only install the dependencies not already installed.
  - `pkgr plan --update` and `pkgr install --update` will behave as if the Update flag was set in the config file.
