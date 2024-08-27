## pkgr clean

Clean cached information

### Synopsis

This subcommand is an entry point for cleaning two categories of cached
data:

 * source and binary tarballs

   Use the 'cache' subcommand to remove these.

 * package databases with information about the packages available from
   repositories

   Use the 'pkgdbs' subcommand to remove these.

To remove cached data for both categories, pass the --all flag.

```
pkgr clean [flags]
```

### Options

```
      --all    clean all cached items
  -h, --help   help for clean
```

### Options inherited from parent commands

```
      --config string     config file (default is pkgr.yml)
      --debug             use debug mode
      --library string    library to install packages
      --logjson           log as json
      --loglevel string   level for logging
      --no-rollback       disable rollback
      --no-secure         disable TLS certificate verification
      --no-update         don't update installed packages
      --strict            enable strict mode
      --threads int       number of threads to execute with
```

### SEE ALSO

* [pkgr](pkgr.md)	 - A package manager for R
* [pkgr clean cache](pkgr_clean_cache.md)	 - Clean cached package tarballs
* [pkgr clean pkgdbs](pkgr_clean_pkgdbs.md)	 - Clean cached package databases

