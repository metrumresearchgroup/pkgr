
# pkgr 3.1.2

Documentation is now available at
<https://metrumresearchgroup.github.io/pkgr/docs/>, including new
documentation on the configuration format and expanded documentation
of many subcommands.  (#422)


# pkgr 3.1.1

* Fixed renv detection when other startup code writes to stdout. (#408)


# pkgr 3.1.0

* For `Lockfile: Type: renv`, pkgr now invokes `renv` to discover the
  library location rather than assuming it is under the current
  directory's `renv/library/`. This change is important for
  compatibility with renv 0.15 and later, where the default behavior
  is now to put a _package_ project library outside of the main
  project directory. (#396)


# pkgr 3.0.0

This release is primarily about adding a more robust test suite, with
minimal user-facing changes. One notable exception is the first point
below, relating to the new `--no-update` flag.

* Reversed default behavior to update to newest available versions of
  packages when `pkgr install` is run, unless `--no-update` flag is
  present (305a8353). Previously `pkgr` did _not_ update installed
  packages unless `--update` was passed.

* Add explicit version flag (6dd9c3ee).

* Added `IgnorePackages` to ignore specific packages from being
  installed even if they're in the dependency tree (104f1c9c).

* Extend integration test coverage and refactored test suite to
  position better for future changes.

* System CPU quotas are now respected when setting the number of CPUs
  that are used if the `--threads` option isn't explicitly passed and
  the `GOMAXPROCS` environment variable isn't set. (#385)
