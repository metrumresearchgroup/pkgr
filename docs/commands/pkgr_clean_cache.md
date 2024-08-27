## pkgr clean cache

Clean cached package tarballs

### Synopsis

Delete cached tarballs for source and binary packages. Both source and
binary files are deleted if neither the --src-only or --binaries-only flags
is specified.

By default, files for all repositories are deleted unless specific
repositories are specified via the --repos option. Note that the value must
match the directory name in the cache, including the unique ID that is
appended to the repository name.

```
pkgr clean cache [flags]
```

### Examples

```
  # Clean binary files for all repos
  pkgr clean cache --binaries-only
  # Clean binaries files for MPN-889df4238bae repo
  pkgr clean cache --repos=MPN-889df4238bae --binaries-only
```

### Options

```
      --binaries-only   clean only binary files from the cache
  -h, --help            help for cache
      --repos string    comma-separated list of repositories to be cleaned. Defaults to all. (default "ALL")
      --src-only        clean only source files from the cache
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

* [pkgr clean](pkgr_clean.md)	 - Clean cached information

