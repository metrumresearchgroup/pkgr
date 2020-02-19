# simple

tags: load

 ## Description
Environment to help test `pkgr load` command.

 ## Expected Behaviors
1. `pkgr load` will indicate that two packages failed to load: R6, fansi, and pillar.
  1. R6 and fansi will fail because their installation is corrupted (the `R` folder has been deleted from both.)
  2. pillar will fail because it will be missing a dependency package, fansi.
