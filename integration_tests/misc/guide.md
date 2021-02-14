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
  - digest (dependency)
  - ellipses (dependency)
  - lifecycle (dependency)
  - vctrs (dependency)
  - glue (dependency)
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
  - Rcpp (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
  - digest (dependency)
  - ellipses (dependency)
  - lifecycle (dependency)
  - vctrs (dependency)
  - glue (dependency)
8. The output of the install from #6 should indicate that Rcpp came from CRAN2 and was redownloaded, but the other packages came from CRAN_Earlier and were NOT redownloaded.
9. Run `pkgr clean --all --config=pkgr2.yml`. Verify that, in localtmp, all repos are removed, and in your system's pkgdb folder (see step 1), only the file from Step 2 remains.
