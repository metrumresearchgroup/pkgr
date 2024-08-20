## pkgr clean pkgdbs

Clean cached package databases

### Synopsis

Delete cached package databases. By default, remove cached databases for
every repository listed in the active configuration file. If the --repos option
is passed, remove only the cached databases for those repositories. Repo names
should match the names in the configuration file.

```
pkgr clean pkgdbs [flags]
```

### Examples

```
  # Clean package databases for CRAN and MPN
  pkgr clean pkgdbs --repos=CRAN,MPN
```

### Options

```
  -h, --help           help for pkgdbs
      --repos string   clear databases for these repos (default "ALL")
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

* [pkgr clean](pkgr_clean.md)	 - Clean cached information

