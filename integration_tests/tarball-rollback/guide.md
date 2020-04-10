# tarball-rollback

tags: rollback, tarball-overwrites

## Description
The purpose of this test is to make sure that rollback operations still function
properly when installing packages from tarballs.

## Assumptions:
* Tarball installations should overwrite existing installations of the same packages

## Expected behavior:

### Rollback behavior
`pkgr plan` will indicate that R6 and crayon will be installed from tarballs.
`pkgr install` will install packages listed in `pkgr.yml`. `rlang` and `crayon` will install
properly. Pkgr will attempt to install a R6 from the indicated tarball, but will
fail and reset the R6 installation to its previous state, leaving only the dummy
R6 folder that it started with.

### Overwrite behavior
Since `crayon` is listed as a user package and also a Tarball, the Tarball package
should always be installed in the end. Run `pkgr plan --rollback=false` and check
the DESCRIPTION file in `crayon` afterwards. You should find the line `FromTarball: TRUE`
