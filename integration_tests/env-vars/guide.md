# env-vars

tags: rpath-env-var

## Description
Environment to help test that environment variables are properly propagated into
pkgr.

## NOTE:
**Remember to manually reset any environment variables you used while doing this test. `make test-setup` cannot do this for you.**

## Expected Behaviors/Test Instructions

### PKGR_RPATH
1. In your current terminal session, set the PKGR_RPATH environment variable to the location
of an "alternate" R executable.

```
export PKGR_RPATH="<some_path_to_r>"
```

1. `pkgr plan --loglevel=trace` will contain the following line, which will vary
based on what you set `PKGR_RPATH` to:

```
TRAC[0000] command args                                  RSettings="{{0 0 0} [] <some_path_to_r> {[]} map[] }" cmdArgs="[--version]" rpath="<some_path_to_r>"
```
Note that you should see your RPATH in both the RSettings string as well as the `rpath="..."` string.

2. `pkgr install --loglevel=trace` will have the same line as in step 1, and, if
the RPath you provided is valid, will install `R6` to `test-library`.

3. Set `PKGR_RPATH` to something invalid on purpose and verify that `pkgr install` fails.
