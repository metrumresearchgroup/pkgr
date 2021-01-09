# simple-suggests
tags: suggests, inspect, cache-local

## Description

Environment to help test the "Suggests" field in a pkgr.yml file, as well as to test that local pkgr caches caches work.

## Expected Behavior

* `pkgr plan --loglevel=debug` should indicate that the package `ellipsis`, its dependencies, and its suggested packages will be installed, as well as any other sub-dependencies.
* `pkgr install` should install `ellipsis`, its dependencies, and its suggested packages (and their subdependencies).
* `pkgr inspect --deps` should return a large object that includes both dependencies and suggested packages of `ellipsis`
* Instead of the default cache (`~/Library/Caches/pkgr/MPN...` on MacOS), pkgr should use a local `pkgcache` folder. You will
see a subfolder prefixed with `MPN` in `pkgcache`, and you will _not_ see such a folder in the regular directory from this test.
  * To verify that this cache is being used, immediately after installing, delete `test-library` and install again. No new packages should be downloaded, as they should already be in the cache.
* pkgr will still utilize the default location for pkgdbs (`~/Library/Caches/pkgr/r_packagedb_caches` on MacOS)


At the time of this writing, these are the packages we expect to be involved:
* ellipsis (**user package**)
* rlang (dependency)
* covr (_suggested_ package)
* testthat (_suggested_ package)
* All of the dependencies for `covr` and `testthat`, but NOT their suggested dependencies.
  - To quickly check that sub-suggested suggestions aren't installed, verify that `shiny`, which is a suggestion of testthat, is not included.
