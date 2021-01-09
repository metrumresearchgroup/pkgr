# packrat-library
tags: packrat-lockfile, renv-lockfile

## Description

Environment to test functionality of packrat library, i.e., if  "Library" is not defined and Lockfile->Type is "packrat",
then the "Library" path is defined like a packrat library path:
   `packrat/lib/{PLATFORM}/{R VERSION}`

## Expected behavior

Run `pkgr install`.

If the packrat folder does not exist, it will be created automatically, the same way a `Library` would.
Packages will then be installed.

If the "packrat-like" folder does exist, packages are installed to the packrat-folder as required by the pkgr.yml file.

Execute both of these scenarios in this test.

We expect the same behavior for Lockfile->Type: renv. Use `pkgr install --config=pkgr2.yml` and repeat this test for renv folder-structures,
using the equivalent expected behavior as above.
