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

`pkgr plan`  and  `pkgr install` will work the same as [simple](../simple/guide.md),
and the ending environment will be equivalent, with the following differences:

`pkgr plan` will indicate that R6 will be installed from a Tarball.
`pkgr install` will install a R6 from the indicated tarball.
`pkgr clean cache` will remove everything from the local cache. 

Relevant entries in the `logs` directory should be created depending on the command run.
- `pkgr plan` currently doesn't save logs.
- `pkgr install` saves to `logs-install.log`
- `pkgr clean <subcommad>` `saves to logs.log`
