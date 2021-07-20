# rollback
tags: rollback

## Description
Testing area to demonstrate that the entire package environment is rolled back
whenever a package fails to install. In other words, If anything goes wrong
during a run of pkgr install or pkgr install --update, then the user's Library
should revert to exactly how it was prior to running.

## Expected Behavior:

* `pkgr install` will fail to install xml2, and after running, `test-library` will be identical to `preinstalled-library`
* `pkgr install --update` will fail to install xml2, and after running, `test-library` will be identical to `preinstalled-library`


## Possible test cases not covered here:
* What if an **update** fails to install instead of just a regular installation?
  - Since updates are treated like regular installations, I think it's safe to exclude this.
