## pkgr install

Install packages

### Synopsis

Create the library defined by the configuration file.

See <https://metrumresearchgroup.github.io/pkgr/docs/config> for details on
the configuration file.

```
pkgr install [flags]
```

### Examples

```
  # Create or update library defined by pkgr.yml
  pkgr install
  # Install new packages and dependencies but don't update packages that already
  # exist in the library.
  pkgr install  --no-update
```

### Options

```
  -h, --help   help for install
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

