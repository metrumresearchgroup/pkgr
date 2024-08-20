## pkgr clean pkgdbs

Subcommand to clean cached pkgdbs

### Synopsis

This command parses the currently-cached pkgdbs and removes all
	of them by default, or specific ones if desired. Identify specific repos using the "repos" argument, i.e.
	pkgr clean pkgdbs --repos="CRAN,r_validated"
	Repo names should match names in the pkgr.yml file.

```
pkgr clean pkgdbs [flags]
```

### Options

```
  -h, --help           help for pkgdbs
      --repos string   Set the repos you wish to clear the pkgdbs for. (default "ALL")
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

* [pkgr clean](pkgr_clean.md)	 - clean up cached information

