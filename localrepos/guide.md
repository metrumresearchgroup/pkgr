# Local Repos

These repos are used for unit and integration testing. Below is a listing of the _base packages_ that
these repos contain the tarballs for. Suggests is assumed to be set to "false" for any pkgr commands,
unless otherwise specified.

## simple
Base packages:
* R6
* pillar

## simple-no-R6
Base packages:
* pillar

To be used with `tarballs/R6_2.4.0.tar.gz`

## testthat_deps
Base packages:
* pillar
* R6

To be used with `tarballs/testthat_2.1.1.tar.gz` (contains deps for testthat).

## bad-xml2
Base packages: 
* xml2
* crayon
* R6
* Rcpp
* crayon
* fansi
* flatxml

Note: the xml2 tarball in this repo is intentionally corrupted to help test the rollback feature. 