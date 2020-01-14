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
* R6 and Rcpp do not update
* fansi, bitops, and RCurl are installed
* xml2 fails to install and says `ERRO[0000] installation failed for packages: xml2`     
* flaxtml does not install, but there is no message saying the installation failed.

Output:
```
johncarlos: ~/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/rollback-disabled$: pkgr install
INFO[0000] R Version 3.6.0                              
INFO[0000] found installed packages                      count=4
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:13396 (binary:source) packages available in for CRAN2 from https://cran.microsoft.com/snapshot/2018-11-18
INFO[0000] package will be updated                       installed_version=2.0.0 pkg=R6 update_version=2.3.0
INFO[0000] package will be updated                       installed_version=0.1.0 pkg=Rcpp update_version=1.0.0
INFO[0000] package installation status                   installed=4 not_from_pkgr=0 outdated=2 total_packages_required=8
INFO[0000] package installation sources                  CRAN2=8
INFO[0000] package installation plan                     to_install=7 to_update=2
INFO[0000] to install                                    package=bitops repo=CRAN2 type=source version=1.0-6
INFO[0000] to install                                    package=fansi repo=CRAN2 type=source version=0.4.0
INFO[0000] to install                                    package=xml2 repo=CRAN2 type=source version=1.2.0
INFO[0000] to install                                    package=RCurl repo=CRAN2 type=source version=1.95-4.11
INFO[0000] to install                                    package=flatxml repo=CRAN2 type=source version=0.0.2
INFO[0000] resolution time 133.397087ms                 
INFO[0000] downloading required packages within directory   dir=/Users/johncarlos/Library/Caches/pkgr
INFO[0000] all packages downloaded                       duration="205.169µs"
INFO[0000] starting initial install                     
INFO[0000] Successfully Installed.                       package=bitops remaining=6 repo=CRAN2 version=1.0-6
INFO[0000] Successfully Installed.                       package=fansi remaining=5 repo=CRAN2 version=0.4.0
INFO[0000] Successfully Installed.                       package=RCurl remaining=4 repo=CRAN2 version=1.95-4.11
ERRO[0000] cmd output                                    exitCode=1 package=xml2 stderr="* installing *source* package ‘xml2’ ...\n** package ‘xml2’ successfully unpacked and MD5 sums checked\n** using staged installation\nERROR: configuration failed for package ‘xml2’\n* removing ‘/private/var/folders/kn/ny6x14mj6sj97c050mp06ywc0000gn/T/NBWFQRMELAIP/xml2’\n" stdout="Found pkg-config cflags and libs!\nUsing PKG_CFLAGS=-I/usr/include/libxml2\nUsing PKG_LIBS=-L/usr/lib -lxml2 -lz -lpthread -licucore -lm\n------------------------- ANTICONF ERROR ---------------------------\nConfiguration failed because libxml-2.0 was not found. Try installing:\n * deb: libxml2-dev (Debian, Ubuntu, etc)\n * rpm: libxml2-devel (Fedora, CentOS, RHEL)\n * csw: libxml2_dev (Solaris)\nIf libxml-2.0 is already installed, check that 'pkg-config' is in your\nPATH and PKG_CONFIG_PATH contains a libxml-2.0.pc file. If pkg-config\nis unavailable you can set INCLUDE_DIR and LIB_DIR manually via:\nR CMD INSTALL --configure-vars='INCLUDE_DIR=... LIB_DIR=...'\n--------------------------------------------------------------------\n"
WARN[0000] error installing                              err="exit status 1"
INFO[0000] total install time                            duration=778.171557ms
ERRO[0000] did not install xml2                         
ERRO[0000] did not install flatxml                      
ERRO[0000] installation failed for packages: xml2       
INFO[0000] duration:924.635672ms                        
WARN[0000] failed package install with err, %sfailed installation for packages: xml2 
```

## Expected Behavior:

* `pkgr install` will fail to install xml2, but all other packages will be installed successfully, and no updates will be applied.
* `pkgr install --update` will fail to install xml2, but will install and update other packages.


## Possible test cases not covered here:
* What if an **update** fails to install instead of just a regular installation?
  - Since updates are treated like regular installations, I think it's safe to exclude this.
