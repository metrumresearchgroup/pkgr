# simple

tags: basic, dependencies, cache-system, local-library, clean-cache, clean-pkgdb, inspect, install-type

 ## Description
Environment to help test basic pkgr functionality, such as the `plan`, `install`, `inspect --deps`

 ## Expected Behaviors
1. `pkgr plan --loglevel debug` will indicate that repositories have been set for packages "R6" and "pillar", as well as their dependencies.
2. `pkgr inspect --deps` will print the following object (order does not matter as long as everything is nested the same way):
```
  {
  "cli": [
    "assertthat",
    "crayon"
  ],
  "pillar": [
    "fansi",
    "rlang",
    "utf8",
    "assertthat",
    "crayon",
    "cli"
  ]
}
```
3. `pkgr install` will install the following packages, using the system default to determine whether those packages are installed through source or binary:
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)

4. After running `pkgr install`, you should see a pkgr cache-folder created in an appropriate temp directory. On Mac, for example, it might be `/Users/<user>/Library/Caches/pkgr`. If you can't find it, look in the install logs for a line such as: `INFO[0004] downloading required packages within directory   dir=/Users/johncarlos/Library/Caches/pkgr` to help figure out where you computer has made a cache.
  - Inside the top-level pkgr cache folder, you should see at least two folders:
    - `CRAN-<HASH>`: Should contain `src` and `binary` subfolders that, after drilling down, contain the source and binary packages that were used to perform the installation.
    - `r_packagedb_caches`:  Should contain a file with a hashed name. This file is the parsed PACKAGES information from one of the repos used during installation.
5. After running `pkgr install` and observing the behavior in 4, run `pkgr clean cache`. Verify that the `CRAN-<HASH>` folder has been removed from the cache.
6. After running `pkgr install` and observing the behavior in 4, run `pkgr clean pkgdbs`. Verify that the file with the hashed name is removed from `r_packagedb_caches`.
7. Re-run `pkgr install`, then run `pkgr clean --all`. Verify that the `CRAN-<HASH>` folder has been removed from the cache and that the file with the hashed name is removed from `r_packagedb_caches`.
