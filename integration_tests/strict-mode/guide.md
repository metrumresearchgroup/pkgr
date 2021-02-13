# strict-mode
tags: strict-mode

## Description
Verify that strict-mode enforces pre-created library folder.

## Expected Behaviors
* `pkgr plan` will indicate via error message that the library must exist before running in strict-mode.
* `pkgr install` will display a fatal error message saying that the library must exist before running in strict-mode. No library will be created nor packages installed.

* Repeat both of these tests again with the argument `--config=pkgr-renv.yml`. equivalent behavior should happen for config files with `Lockfile: Type: renv`
