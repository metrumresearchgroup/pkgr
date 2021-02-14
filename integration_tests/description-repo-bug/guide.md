# outdated-pkgs
tags: bug-duplicate-repo, bug-repo-name

tags: bug-duplicate-repo, bug-repo-name

## Description

Environment to verify that the "duplicate repo in description file" bug and that
repo names are assigned in the DESCRIPTION file correctly.

** DO NOT MODIFY THE [outdated-library](outdated-library) DIRECTORY OR YOU WILL DESTROY THIS ENVIRONMENT **

## How to quickly reset test environment
For convenience, you can use this script to reset your test environment and come back to this directory. Run the following command from a terminal window in this directory:
`cd ../ && make test-setup && cd description-repo-bug`

Please note that this will build and install pkgr using the source code in this project.

## Expected behavior

Run `pkgr install`. Pkgr will install an updated version of R6. Verify that the DESCRIPTION file for test-library/R6 has the following lines:
```
OriginalRepository: CRAN
Repository: CRAND_ALTHOR
```
Make sure that the lines appear exactly as written here, with no additions.
