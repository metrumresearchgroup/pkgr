# tarball-install

tags: tarball-install

## Expected behavior:

`pkgr plan`  and  `pkgr install` will work the same as [simple](../simple/guide.md),
and the ending environment will be equivalent, with the following differences:

`pkgr plan` will indicate that R6 will be installed from a Tarball.
`pkgr install` will install a R6 from the indicated tarball.
`pkgr clean cache` will remove everything from the local cache. **Current observed behavior**: Clean removes most entries, but does not touch hashed tarball installations living in the cache.