# tilde-expansion

tags: tilde-expansion

## Description
This test is meant to ensure that the `~` is expanded to the home directory in all
paths configured by the user. This test was modeled after the Tarball test,
but with all paths changed to include `~`.  Functionally, this should test all
possible path-inputs.

We have also added a repo-customization for good measure. This should not affect the
final results.

## Expected behavior

`pkgr plan`  and  `pkgr install` will install the following packages:
  - R6 (**user tarball**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)	 
  - utf8 (dependency)	  
  - fansi (dependency)	  
  - assertthat (dependency)	  
  - crayon (dependency)

with the following behavior:

`pkgr plan` will indicate that R6 will be installed from a Tarball.
`pkgr install` will install R6 from the indicated tarball.

Also, `pkgr clean cache` will remove everything from the local cache.

Relevant entries in the `logs` directory should be created depending on the command run.
Logging is currently (as of January 9, 2020) in need of re-work, so consider this part of
of the test optional.
- `pkgr plan` currently doesn't save logs.
- `pkgr install` saves to `logs-install.log`
- `pkgr clean <subcommad>` `saves to logs.log`
