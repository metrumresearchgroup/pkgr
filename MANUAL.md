# Pkgr Manual

## Commands and Flags
All `pkgr` commands follow the format `pkgr <command> <flags>`

### Commands and Subcommands
`pkgr` or `pkgr --help`
* Prints a list of top level pkgr Commands

#### Universal Flags (can be applied to any command)
* `--config [path]`: Path to the `pkgr.yml` file, if not `./pkgr.yml`.
* `--preview`: Enable preview mode [[ unclear on what this does ]]
* `--debug`: Enable debug mode (developer only).
* `--loglevel [info/debug/trace]`: Set the log level to "info", "debug", or "trace".
* `--update` or `--update=false`: Set pkgr to update outdated packages it finds in the user's library.
* `--threads [number]`: Set pkgr to use up to the specified number of threads.
* `--strict` or `--strict=false`: Enable or disable strict mode (see notes)
* `--help`: List information about a command.

`pkgr plan`
- Displays information about what pkgr would do if you were to run `pkgr install`

`pkgr install`
- Installs packages according to the specifications set in the configuration file.

`pkgr clean`
- Top level command for "cleaning" operations.
- Subcommands:
  - `pkgr clean cache`
    - Removes the cache of downloaded packages (source/binary files) for the repos listed in the configuration file.
  - `pkgr clean pkgdbs`
    - Removes the saved package databases for repos listed in the configuration file
  - `pkgr clean --all`
    - Pass the `--all` flag to `pkgr clean` as a shortcut to run both `pkgr clean cache` and `pkgr clean pkgdbs`

`pkgr inspect --deps`
- Prints detailed information about the installation plan pkgr would execute on a `pkgr install` -- provides insight into how pkgr determines dependencies. Must pass the `--deps` flag, or this becomes equivalent to `pkgr install`

`pkgr add [package name]`
- Adds the given package to the config file. Does not check if package is available on any repos, and does not set any customizations.

`pkgr remove [package name]`
- Removes the given package from the config file, if it is there.

## Strict Mode
When strict mode is enabled, pkgr will refuse to run if unless the following conditions are satisfied:
* The library defined in pkgr.yml must exist.

## Default Settings
* Debug mode: disabled
* Preview mode: disabled
* Strict mode: disabled
* Loglevel: “info”
* Rpath: “R”
* Threads: Number of processors available on system
* Cache: System temp directories
* Install suggested packages: false
* Update outdated packages: false
* Src/binary configurations: System defaults (bin for Mac/Windows, src for Linux)






## Notes from issue
The following options should be configurable in only the yml file:

Packages to install
Repositories to use
Src/binary configuration for packages
Src/binary configuration for repositories
Libpaths [[ unclear on what this does ]]
Rpath
Cache directory
Lockfile settings (should pgkr make itself compatible with a Lockfile such as renv.lock?)
Tests:
basic, cache-local, pkg-customizations, lockfile, log-file, install-log

The following options should be configurable statically in the yml file or dynamically via the command-line (with command line arguments always overwriting yml configurations):

Library (directory) to install packages to
Update settings (should pkgr attempt to update outdated packages?)
Logging level
Installation of “Suggested” packages config
Number of threads to use
Strict mode enabled/disabled

The following settings should be configurable only via command-line flags

config (path to yaml file, defaults to pkgr.yml)
preview mode [[ unclear on what this does ]]
Debug mode
Tests:
command-flags


If nothing is set, these are the default settings:

Debug mode: disabled
Preview mode: disabled
Strict mode: disabled
Loglevel: “info”
Rpath: “R”
Threads: Number of processors available on system
Cache: System temp directories
Install suggested packages: false
Update outdated packages: false
Src/binary configurations: System defaults (bin for Mac/Windows, src for Linux)

Tests:
basic, todo
Minimum required user-settings:

Version (pkgr)
Library
Packages to install
Repositories to use
Tests:
todo
