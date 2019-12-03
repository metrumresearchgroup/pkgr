# Misc

tags: idempotence, cache-partial, cache-extraneous, clean-pkgdb

 ## Description
Environment to help test miscellaneous pkgr functionality such as idempotence.

## Test steps

### Test 1: tests: idempotence
1. Run `pkgr install`
2. `pkgr install` will install the following packages to `test-library`:
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
3. Delete the `fansi` folder in `test-library`.
3. Run `pkgr install` install again.
4. Verify that the environment looks the same as it did in step 2 (this tests pkgr for idempotence)

### Test 2: tests: cache-partial, cache-extraneous
1. Remove all items from your system's pkgdb folder (on Mac, it should be `/Users/<user>/Library/Caches/pkgr/r_packagedb_caches`)
2. Run `pkgr install` (this ensures that the irrelevant entries are added to the cache.) Make a note of what is created in your system's pkgdb folder.
3. Run `pkgr install --config=pkgr2.yml`
4. Remove from `localtmp` the `CRAN2...` directory (not `CRAN_Earlier`).
5. Remove the contents of `test-library2`
6. Rerun `pkgr install --config=pkgr2.yml`
7. Verify that the following packages are installed in `test-library2`
  - R6 (**user package**)
  - pillar (**user package**)
  - Rcpp (user package)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
8. The output of the install from #6 should contain lines equivalent to these:
```
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=9
INFO[0000] package installation sources                  CRAN=1 CRAN_Earlier=8
INFO[0000] package installation plan                     to_install=9 to_update=0
INFO[0000] resolution time 259.670507ms                 
INFO[0000] downloading required packages within directory   dir=/Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/misc/localtmp
INFO[0000] downloading package                           package=Rcpp
INFO[0001] all packages downloaded                       duration=1.646435638s
INFO[0001] starting initial install                     
INFO[0002] Successfully Installed.                       package=assertthat remaining=8 repo=CRAN_Earlier version=0.2.1
INFO[0002] Successfully Installed.                       package=R6 remaining=7 repo=CRAN_Earlier version=2.4.0
INFO[0002] Successfully Installed.                       package=crayon remaining=6 repo=CRAN_Earlier version=1.3.4
INFO[0002] Successfully Installed.                       package=utf8 remaining=5 repo=CRAN_Earlier version=1.1.4
INFO[0002] Successfully Installed.                       package=fansi remaining=4 repo=CRAN_Earlier version=0.4.0
INFO[0002] Successfully Installed.                       package=rlang remaining=3 repo=CRAN_Earlier version=0.3.4
INFO[0002] Successfully Installed.                       package=cli remaining=2 repo=CRAN_Earlier version=1.1.0
INFO[0002] Successfully Installed.                       package=pillar remaining=1 repo=CRAN_Earlier version=1.3.1
INFO[0003] Successfully Installed.                       package=Rcpp remaining=0 repo=CRAN2 version=1.0.1
```
(The important thing here is that Rcpp comes from CRAN2 and must be redownloaded, but the other packages are not.)
9. Run `pkgr clean --all --config=pkgr2.yml`. Verify that, in localtmp, all repos are removed, and in your system's pkgdb folder (see step 1), only the file from Step 2 remains.
