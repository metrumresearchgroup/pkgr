## pkgr clean

clean up cached information

### Synopsis

clean up cached source files and binaries, as well as the saved package database.

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
      --preview           preview action, but don't actually run command
      --strict            enable strict mode
      --threads int       number of threads to execute with
      --update            whether to update installed packages
```

### SEE ALSO

* [pkgr](pkgr.md)	 - package manager
* [pkgr clean cache](pkgr_clean_cache.md)	 - Subcommand to clean cached source and binary files.
* [pkgr clean pkgdbs](pkgr_clean_pkgdbs.md)	 - Subcommand to clean cached pkgdbs

