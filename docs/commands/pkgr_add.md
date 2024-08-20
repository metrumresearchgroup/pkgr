## pkgr add

Add packages to the configuration file

### Synopsis

Add the specified packages to the 'Packages' section of the
configuration file.

```
pkgr add [flags] <package> [<package>...]
```

### Examples

```
  # Add mrgsolve and bbr to list of packages
  pkgr add mrgsolve bbr
  # Add rlang and then do installation
  # (same result as following up with 'pkgr install' call)
  pkgr add --install rlang
```

### Options

```
  -h, --help      help for add
      --install   run install after updating config
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

