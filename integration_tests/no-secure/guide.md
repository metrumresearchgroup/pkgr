# no-secure

tags: tls-skip

## Description
Environment that points to a repo with an unsigned certificate (tls verification should fail).
Used to help validate the `--no-secure` argument to pkgr.

## Expected Behaviors
1. `pkgr plan --no-secure` will indicate that repositories have been set for the packages in `pkgr.yml`
3. `pkgr install --no-secure` will install the following packages and dependencies as specified in `pkgr.yml`
