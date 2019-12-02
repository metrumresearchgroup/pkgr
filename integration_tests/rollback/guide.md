# rollback
tags: rollback

## Description
Testing area to demonstrate that the entire package environment is rolled back
whenever a package fails to install. In other words, If anything goes wrong
during a run of pkgr install or pkgr install --update, then the user's Library
should revert to exactly how it was prior to running.

## Justification
To make sure that rollback functionality works as intended, we need a repeatable way to check that package environments rollback correctly.

## Default Setup
From this directory:

`(cd .. ; make install ; make test-rollback-reset)`.

Alternatively:

`(cd .. ; make test-setup)`

to set up all tests.

## Note
**Important:**
This test assumes that attemtpting to install **xml2** from source will FAIL on your machine. If this is not the case, please replace xml2 in pkgr.yml with a package that _will_ fail to install on your machine.


"Setup" for this test includes adding four packages to the test-library, each of
which has been configured in some way to test a different case.

* crayon: **Normal, preinstalled** -- A "regular" package that is already installed and doesn't need updates.
* fansi: **Normal, not installed** -- A "regular" package that will install correctly but should not be in the final package environment after error.
* R6: **Outdated** -- An outdated package that successfully installs prior to failure.
* Rcpp: **Outdated** -- An outdated package that does not attempt to install, as it's installation would happen after the failure.
* utf8: **Not installed by pkgr** -- A package that is not included at all in the pkgr file. We may want to remove this, but I figured we shouldn't be corrupting users' environments even if they're using pkgr in an unintended way.

## Expected Behavior:

* `pkgr install` will fail to install xml2, and after running, `test-library` will be identical to `preinstalled-library`
* `pkgr install --update` will fail to install xml2, and after running, `test-library` will be identical to `preinstalled-library`


## Possible test cases not covered here:
* What if an **update** fails to install instead of just a regular installation?
  - Since updates are treated like regular installations, I think it's safe to exclude this.
