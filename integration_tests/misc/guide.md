# Misc

tags: idempotence

 ## Description
Environment to help test miscellaneous pkgr functionality such as idempotence.

## Test steps
1. Run `pkgr install`
2. `pkgr install` will install the following packages:
  - R6 (**user package**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)
  - utf8 (dependency)
  - fansi (dependency)
  - assertthat (dependency)
  - crayon (dependency)
3. Delete the `fansi` folder in `test-library`.
3. Run `pkgr install` install again.
4. Verify that the environment looks the same as it did in step 2 (this tests pkgr for idempotence)
5.

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
