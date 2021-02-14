# mixed-source
tags: multi-repo, repo-customizations, pkg-customizations

## Description

Environment to help test that packages can be pulled from multiple repositories. Also helps test repo and type customizations.

## Note
This test should be run on Mac or Windows.

## Expected behavior

* `pkgr plan --loglevel debug` indicates that the installation will apply these restrictions:
  - `fansi` should come from `MPN_Secondary` and be installed as binary/source according to your system's default (Mac/Windows are binary, Linux is source).
  - Everything else should come from `MPN_Primary`
  - Everything coming from `MPN_Primary` should be installed as source, EXCEPT:
  - R6 will come from `MPN_Primary` and be installed as a binary

* `pkgr install` installs `fansi`, `R6`, and `pillar`, as well as their dependencies.
