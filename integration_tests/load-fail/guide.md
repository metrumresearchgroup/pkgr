# load-fail

tags: load-fail

 ## Description
Environment to help test `pkgr load` command.

 ## Expected Behaviors
* In the pkg-environment stored in preinstalled-library, packages `R6`, `pillar`, `utf8` and `fansi` should be unable to be loaded.
  - `R6` and `fansi` have had their "R" folders deleted.
  - `pillar` is "validly installed", but one of its dependencies (`utf8`) has been removed. **THIS CHECK DOES NOT FAIL -- `library(pillar)` runs successfully even with the missing dependency.**
  - `utf8` has been removed
* `pkgr load` should indicate that `R6` and `pillar` fail to load.
* `pkgr load --all --loglevel=debug` should indicate that `R6`, `pillar`, `utf8`, and `fansi` fail to load. `cli`, `assertthat`, `crayon`, and `rlang` should load successfully.

1. `pkgr load` will indicate that two packages failed to load: R6, fansi, and pillar.
  1. R6 and fansi will fail because their installation is corrupted (the `R` folder has been deleted from both.)
  2. pillar will fail because it will be missing a dependency package, fansi.

Dependency tree as output by `pkgr inspect --deps`
```
{
  "cli": [
    "assertthat",
    "crayon"
  ],
  "pillar": [
    "utf8",
    "assertthat",
    "crayon",
    "fansi",
    "rlang",
    "cli"
  ]
}
```
