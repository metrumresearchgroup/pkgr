# tarball-install

tags: tarball-install

## Expected behavior:

Relevant packages:
  - R6 (**user tarball**)
  - pillar (**user package**)
  - rlang (dependency)
  - cli (dependency)	 
  - utf8 (dependency)	  
  - fansi (dependency)	  
  - assertthat (dependency)	  
  - crayon (dependency)

`pkgr plan` and  `pkgr install` will plan for/install the above packages, and the following will be true:

`pkgr plan --loglevel=debug` will indicate that R6 will be installed from a Tarball.
`pkgr install` will install a R6 from the indicated tarball.
`pkgr clean cache` will remove everything from the local cache.
