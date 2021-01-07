# linux-binaries

tags: linux-binaries

## Special Instructions

This test must be run on a Linux machine.

## Description
Test to make sure Linux binaries can be installed.

## Expected Behaviors
1. `pkgr plan --loglevel=debug` will indicate that repositories have been set for package `R6`, `pillar`, and their dependencies. The install type for all packages will be as binaries.
2. `pkgr install --loglevel=debug` will install the following packages, all through binaries:
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
