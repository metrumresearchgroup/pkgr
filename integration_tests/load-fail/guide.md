# load-fail

tags: load-fail

## Description
Environment to help test `pkgr load` command.

## Note
This test will currently only work as written with R version 4.X (any) installed.
If you need to run this test with R Version 3.X, delete `test-library` and re-create it, filling it with the contents of `preinstalled-library-R3`.


## Expected Behaviors
* In the pkg-environment stored in preinstalled-library, packages `R6`, `pillar`\*, `utf8` and `fansi` should be unable to be loaded.
  - `R6` and `fansi` have had their "R" folders deleted.
  - `pillar`\* is "validly installed", but one of its dependencies (`utf8`) has been removed.
  - `utf8` has been removed
* `pkgr load` should indicate that `R6` and `pillar`\* fail to load.
* `pkgr load --all --loglevel=debug` should indicate that `R6`, `pillar`\*, `utf8`, and `fansi` fail to load. `cli`, `assertthat`, `crayon`, and `rlang` should load successfully.

\* **`library(pillar)` runs successfully even with the missing dependency. It does so even in a regular R session not run through pkgr. Therefore, if pillar loads properly, we will still count these tests as successful.**

1. `pkgr load --all` will indicate that four packages failed to load: R6, fansi, utf8, and pillar. Without the --all command, only R6 and pillar should attempt to be loaded, and therefore only those should fail.
  1. R6 and fansi will fail because their installation is corrupted (the `R` folder has been deleted from both.)
  2. pillar will fail because it will be missing a dependency package, fansi (see exception above).
  3. utf8 will fail because it has been removed completely.

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
