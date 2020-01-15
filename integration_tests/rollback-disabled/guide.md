# rollback-disabled
tags: rollback-disabled

## Description
Testing area to demonstrate that the entire package environment is rolled back
whenever a package fails to install. In other words, If anything goes wrong
during a run of pkgr install or pkgr install --update, then the user's Library
should revert to exactly how it was prior to running.

## Justification
To make sure that rollback functionality works as intended, we need a repeatable way to check that package environments rollback correctly.

## Default Setup
From this directory:

`(cd .. ; make install ; make test-rollback-disabled-reset)`.

Alternatively:

`(cd .. ; make test-setup)`

to set up all tests.

## Note
**Important:**
This test assumes that attemtpting to install **xml2** from source will FAIL on your machine. If this is not the case, please replace xml2 in pkgr.yml with a package that _will_ fail to install on your machine.


"Setup" for this test includes adding four packages to the test-library, each of
which has been configured in some way to test a different case.

* crayon: **Normal, preinstalled** -- A "regular" package that is already installed and doesn't need updates.
* fansi: **Normal, not installed** -- A "regular" package that will install correctly.
* R6: **Outdated** -- An outdated package that successfully installs prior to error.
* Rcpp: **Outdated dependency of failure** -- An outdated package that the "fail" package relies on and will update.
* utf8: **Not installed by pkgr** -- A package that is not included at all in the pkgr file. We may want to remove this, but I figured we shouldn't be corrupting users' environments even if they're using pkgr in an unintended way.
* flaxtml: **Depends on fail package** -- A package that depends on xml2 (the package that fails to install). Dependencies for flaxtml should install, but flaxtml itself should fail to install because xml2 failed to install.
* RCurl: **Is a dependency of flaxtml** RCurl should still be installed despite the fact that xml2 -> flaxtml could not be installed.
* bitops: **Is a dependency of RCurl** Bitops should still be installed, because flaxtml should still be installed.

Dependency tree (via `pkgr install --deps`):
```
{
  "RCurl": [
    "bitops"
  ],
  "flatxml": [
    "Rcpp",
    "bitops",
    "RCurl",
    "xml2"
  ],
  "xml2": [
    "Rcpp"
  ]
}
```

## Current behavior observed:
`pkgr install`:
* R6 and Rcpp update successfully
* fansi, bitops, and RCurl are installed
* xml2 fails to install and says `ERRO[0000] installation failed for packages: xml2`     
* flaxtml does not install, but there is no message saying the installation failed.


## Expected Behavior:

* `pkgr install` will fail to install xml2, but all other packages will be installed successfully, and no updates will be applied.
* `pkgr install --update` will fail to install xml2, but will install and update other packages.


## Possible test cases not covered here:
* What if an **update** fails to install instead of just a regular installation?
  - Since updates are treated like regular installations, I think it's safe to exclude this.
