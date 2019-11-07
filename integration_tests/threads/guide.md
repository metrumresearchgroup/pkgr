# threads

## Description
Very simple pkgr environment meant to easily test the `threads` setting in pkgr.yml.

## Expected Behaviors
* `pkgr plan` will indicate that 5 threads are to be used.
* `pkgr install` will install the following packages, launching five threads to do so:
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
