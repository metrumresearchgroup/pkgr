# create-library

## Description
Very simple pkgr environment meant to easily test that a library is created when
it does not exist.

## Expected Behaviors
* `pkgr plan` will indicate test-library will be created
* `pkgr install` will create test-library and install the following packages
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
