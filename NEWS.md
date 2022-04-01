
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
