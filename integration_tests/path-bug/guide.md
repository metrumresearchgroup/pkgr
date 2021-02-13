# path-bug

tags: packages-file-path-bug

## Description
Test to make sure we are addressing the new CRAN "feature" that sometimes causes
duplicate packages to be appended, unsorted, to the end of the PACKAGES file
with a `Path` field.

This bug previously caused pkgr to get the wrong version of xml.

## Special Instructions
This test must be run from a computer with an R version >= 4.0,

## Expected Behaviors
1. `pkgr plan --loglevel=debug` will indicate the package XML will be installed from source with version `3.99-0.5`
2. `pkgr install` will install XML from source with version `3.99-0.5`
