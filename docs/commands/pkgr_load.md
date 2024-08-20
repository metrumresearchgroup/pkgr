## pkgr load

Check that installed packages can be loaded

### Synopsis

Load packages specified in the configuration file to validate that
each package has been installed successfully and can be used.

**Execution environment**. This subcommand runs R with the same settings
that R would use if you invoked 'R' from the current working directory. It
relies on that environment being configured to find packages in the library
path specified in the configuration file (via 'Library' or 'Lockfile:
Type').  Pass the --json argument to confirm that the package is being
loaded from the expected library.

```
pkgr load [flags]
```

### Examples

```
  # Load packages listed in config file
  pkgr load --json
  # Load the above packages and all their dependencies
  pkgr load --json --all
```

### Options

```
      --all    load all packages in dependency tree
  -h, --help   help for load
      --json   output results as a JSON object
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

