## pkgr clean cache

Subcommand to clean cached source and binary files.

### Synopsis

This command is a subcommand of the "clean" command.

	Using this command deletes cached source and binary files. Use the
	--src and --binary options to specify which repos to clean each
	file type from.

	

```
pkgr clean cache [flags]
```

### Options

```
      --binaries-only   Clean only binary files from the cache
  -h, --help            help for cache
      --repos string    Comma separated list of repositories to be cleaned. Defaults to all. (default "ALL")
      --src-only        Clean only src files from the cache
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

* [pkgr clean](pkgr_clean.md)	 - clean up cached information

