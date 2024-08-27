## pkgr run

Launch R session with config settings

### Synopsis

Start an interactive R session based on the settings defined in the
configuration file.

   * Use the R executable defined by the 'RPath' value, if any.

   * Set the library paths so that packages come from only the
     configuration's library and the library bundled with the R
     installation.

   * If the --pkg option is passed, set the environment variables defined in
     the package's 'Customizations' entry.

```
pkgr run [flags]
```

### Examples

```
  # Launch an R session, setting values based on pkgr.yml
  pkgr run
  # Also setting environment variables specified for dplyr:
  #
  #   Customizations:
  #     Packages:
  #        - dplyr:
  #            Env:
  #              [...]
  pkgr run --pkg=dplyr
```

### Options

```
  -h, --help         help for run
      --pkg string   package environment to set
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

