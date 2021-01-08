# integration-tests

tags: na

## Description
The directories here contain sample `pkgr.yml` files that are already setup to
provide common use-cases of pkgr. You can find more details in the `guide.md`
files within each folder.

For the time being, these environments are meant to be used for manual testing
as a way to help sanity-check any changes.

Unless otherwise specified, it should be assumed that theses tests are to be run
in a 4.X version of R.

## Quick setup
* To quickly set your tests to a basic-state, set your working directory to
`pkgr/integration_tests` and run `make test-setup`. Please note:
  - `make test-setup` will not reset your system pkgr cache.
  - `make-test-setup` *will* reset the local caches of the test environments if they are set.
    - Right now, only [simple-suggests](./simple-suggests) has a local cache set,
    so only [simple-suggests](./simple-suggests) will have its cache reset.

## Things to consider:
* Right now, most of these pkgr.yml files do not specify a local cache, which
means that they share a folder your system cache directory (on MacOS, that means
`~/Library/Caches/pkgr`). If you wish to test the `pkgr clean` commands, use a test
that specifies a local cache, such as [simple-suggests](./simple-suggests).
* You can reset the cache for a test at any time with the `pkgr clean --all` command.
Again, if you are testing `pkgr clean` specifically, use a test with a local cache.
