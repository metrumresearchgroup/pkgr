## pkgr run

Run R with the configuration settings used with other R commands

### Synopsis


	allows for interactive use and debugging based on the configuration specified by pkgr
 

```
pkgr run R [flags]
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

