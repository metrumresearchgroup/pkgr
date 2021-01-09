# outdated-pkgs-no-update
tags: pkg-update, pkg-outdated

## Old Description

Environment specifically to verify that the "update" feature is disabled by default.

## Expected behavior

Without any flags passed, `pkgr plan` and `pkgr install` will behave as if --update was set to false. This is all that's needed for this test.

* **Update True**:
  - `pkgr plan` indicates that four packages are outdated: `R6`, `pillar`, and `Matrix`, plus `crayon`, which is a dependency of pillar.

  - `pkgr install` will install the most recent versions of `R6`, `pillar`, `Matrix`, and `crayon`, as well as several other dependenceis.

* **Update False**
  - `pkgr plan` will warn that `R6`, `pillar`, `Matrix`, and `crayon` are outdated.
  - `pkgr install` will leave `R6`, `pillar`, `Matrix`, and `crayon` alone (their versions will not be updated) and will only install the dependencies not already installed.
  - `pkgr plan --update` and `pkgr install --update` will behave as if the Update flag was set to "true" in the config file.
