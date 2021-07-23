# rollback
tags: rollback

## Description
Testing area to demonstrate that the entire package environment is rolled back
whenever a package fails to install. In other words, If anything goes wrong
during a run of pkgr install or pkgr install --update, then the user's Library
should revert to exactly how it was prior to running.

Note, the package failinstall was built to purposefully fail during install
so this can be quickly and easily tested. To do so, the onLoad hook was
intercepted and an error introduce, such that during the load check phase of
a normal install that will fail.