# simple

tags: basic, sanity-check, dependencies, cache-system

 ## Description
Environment to help test basic pkgr functionality, such as the `plan`, `install`, `inspect --deps`

 ## Expected Behaviors
* `pkgr plan` will indicate that repositories have been set for packages "R6" and "pillar".
* `pkgr install` will install the following packages:
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
* `pkgr inspect --deps` will print the following object:
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
