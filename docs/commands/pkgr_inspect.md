## pkgr inspect

Inspect package dependencies

### Synopsis

The inspect subcommand provides an entry point for displaying
information that can be gathered by examining the configuration file, the
associated package database, and the library. The current focus is on
inspecting package dependencies (triggered by passing --deps).

Note: If the configuration file has 'Suggests: true', that does not affect
the set of dependencies listed for any particular package. Instead the set
of suggested packages is included in the top-level package set.

```
pkgr inspect --deps [flags] [<package>...]
```

### Examples

```
  # Show all dependencies as a tree
  pkgr --loglevel=fatal inspect --deps --tree
  # Show dependency tree, restricting roots to the named packages
  pkgr --loglevel=fatal inspect --deps --tree processx here

  # Output a JSON record where each item maps a package to its direct
  # and indirect dependencies
  pkgr --loglevel=fatal inspect --deps
  # Do the same, but filter to records for the named packages
  pkgr --loglevel=fatal inspect --deps processx here

  # Output a JSON record where each item maps a package to
  # the packages that have it as a dependency
  pkgr --loglevel=fatal inspect --deps --reverse
```

### Options

```
      --deps      show dependency tree
  -h, --help      help for inspect
      --json      suppress non-fatal logging (note: prefer --loglevel=fatal to this flag)
      --reverse   show reverse dependencies
      --tree      show full recursive dependency tree
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

