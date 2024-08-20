## pkgr add

add one or more packages

### Synopsis


	add package/s to the configuration file and optionally install


```
pkgr add [package name1] [package name2] [package name3] ... [flags]
```

### Options

```
  -h, --help      help for add
      --install   install package/s after adding
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

* [pkgr](pkgr.md)	 - package manager

