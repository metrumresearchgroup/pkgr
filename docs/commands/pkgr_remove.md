## pkgr remove

Remove packages from the configuration file

### Synopsis

Remove the specified packages from the 'Packages' section of the
configuration file.

```
pkgr remove [flags] <package> [<package>...]
```

### Options

```
  -h, --help   help for remove
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

