## pkgr load

Checks that installed packages can be loaded

### Synopsis

Attempts to load user packages specified in pkgr.yml to validate that each package has been installed
successfully and can be used. Use the --all flag to load all packages in the user-library dependency tree instead of just user-level packages.

```
pkgr load [flags]
```

### Options

```
      --all    load user packages as well as their dependencies
  -h, --help   help for load
      --json   output a JSON object of package info at the end
```

### Options inherited from parent commands

```
      --config string     config file (default is pkgr.yml)
      --debug             use debug mode
      --library string    library to install packages
      --logjson           log as json
      --loglevel string   level for logging
      --no-rollback       Disable rollback
      --no-secure         disable TLS certificate verification
      --no-update         don't update installed packages
      --preview           preview action, but don't actually run command
      --strict            Enable strict mode
      --threads int       number of threads to execute with
      --update            whether to update installed packages
```

### SEE ALSO

* [pkgr](pkgr.md)	 - package manager

