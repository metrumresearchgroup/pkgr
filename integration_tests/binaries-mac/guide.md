# binaries-mac

tags: mac-binaries-3, mac-binaries-4

## Description
Environment to help test that pkgr correctly installs Mac binaries on older and newer versions of R.

## Special Instructions

This test must be run on a machine running MacOS.

This test is dependent on what version of R is running on your machine. To fully
complete this test, the Expected Behaviors need to be verified on a system with
a Version 3.X installation of R (any), as well as a Version 4.X installation of R
(any). This is because of a change in how CRAN-like repositories store Mac
binaries, starting as of R Version 4.0.


## Expected Behaviors

Install a 3.X version of R and verify the following. When you are done, install a
4.X version of R and verify the following again.

1. `pkgr plan` will indicate that repositories have been set for packages "R6" and "pillar".
2. `pkgr inspect --deps` will print the following object:
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
