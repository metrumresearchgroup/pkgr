## pkgr plan

Display plan for installation

### Synopsis

Preview an installation with the current configuration. This subcommand
is commonly invoked before running 'pkgr install' to confirm that the
configuration is behaving as intended.

The output includes details about which repositories particular packages would
be retrieved from, the library that packages would be installed into, and which
packages would be installed or updated.

```
pkgr plan [flags]
```

### Options

```
  -h, --help        help for plan
      --show-deps   show the (required) dependencies for each package
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

