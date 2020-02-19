# load

tags: load

 ## Description
Environment to help test `pkgr load` command.

 ## Expected Behaviors
* `pkgr load --loglevel=debug` will indicate that `R6` and `pillar` (user packages) load successfully.
* `pkgr load --all --loglevel=debug` will indicate that `R6`, `pillar`, `utf8`, `fansi`, `cli`, `assertthat`, `crayon`, and `rlang` are loaded successfully.
