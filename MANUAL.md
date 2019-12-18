# Pkgr Manual

## Commands and Flags
All `pkgr` commands follow the format `pkgr <command> <flags>`

### Commands and Subcommands
`pkgr` or `pkgr --help`
* Prints a list of top level pkgr Commands

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

#### Universal Flags (can be applied to any command)
* `--config [path]`: Path to the `pkgr.yml` file, if not `./pkgr.yml`.
* `--preview`: Enable preview mode [[ unclear on what this does ]]
* `--library [path_to_directory]`: Set location for packages to be installed.
* `--debug`: Enable debug mode (developer only).
* `--loglevel [info/debug/trace]`: Set the log level to "info", "debug", or "trace".
* `--update` or `--update=false`: Set pkgr to update outdated packages it finds in the user's library.
* `--threads [number]`: Set pkgr to use up to the specified number of threads.
* `--strict` or `--strict=false`: Enable or disable strict mode (see notes)
* `--help`: List information about a command.

## pkgr.yml Configuration Options
**Important Note:** Any flags supplied to the command line will overwrite equivalent settings set in the configuration file.

### The following item MUST be set at the top of your configuration file:
```
Version: 1
```

### The following items can be set ONLY in the configuration file:
* Packages to install (required)
  - Yaml syntax:

  ```
  Packages:
    - <package_name1>
    - <package_name2>
    - dplyr
  ```
* Repositories to use (required)
  - Repos must be named and a URL to the repository must be provided.
  - Unless otherwise configured, pkgr will use Repos in the order listed when searching for packages installation files for an individual package.
  - Yaml syntax:

  ```
  Repos:
    - <repo_name>: "<repo_url>"
    - CRAN_Micro: "https://cran.microsoft.com/snapshot/2018-11-18"
  ```
* Customizations
  - Package-specific customizations
    - Repository to install a package from
    - Installation type for a package
    - Whether or not to install Suggested dependencies for the package.
    - Yaml syntax:
    ```
    Customizations:
      Packages:
        - <package_name1>:
            Suggests: <true/false>
        - <package_name2>:
            Type: <source/binary>
        - <package_name3>:
            Repo: <repo_name as listed in this config file>
        - data.table:
            Suggests: true
            Type: source
            Repo: CRAN_Micro
    ```
  - Repository-specific customizations
    - Installation type for packages pulled from a repo
    - Whether or not to install Suggested dependencies for packages pulled from a repo
    - Yaml syntax:
    ```
    Customizations:
      Repos:
        - <repo_name as listed in this config file>:
            Type: <source/binary>
        - <repo_name as listed in this config file>:
            Suggests: <true/false>
        - CRAN_Micro:
            Type: source
            Suggests: true
    ```
* Libpaths [[ I don't know what this does]]
* Rpath [[I don't know what this does in terms of pkgr]]
* Cache directory
  - Directory to use for storing cached files files
  - Folder must exist
  - Yaml syntax:
  ```
  Cache: <path_to_desired_folder
  ```
* Logging Output
  - Settings for logging output.
  - Yaml syntax:
  ```
  Logging:
    all: <path_to_file>
    install: <path_to_file>
    overwrite: <true/false>
  ```
  - These configuration options are slated for rework, but for now, they behave as described below:
    * Setting: `all: x.log`
      * `pkgr install` logs to x.log when `install=y.log` is not set.
      * `pkgr clean --all` logs to x.log
      * `pkgr inspect` logs to x.log
      * `pkgr clean cache` and `pkgr clean pkgdbs` do NOT log to x.log
      * `pkgr plan` does NOT log to x.log
    * Setting: `install: y.log`
      * `pkgr install` logs to y.log (and only to y.log)
    * Setting: overwrite: `true`
      * Logs are overwritten every time a new pkgr command is run.
      * If false, logs are appended to every time a new pkgr command is run.

* Lockfile settings
  - Lockfile to look for, specified by its type. Pkgr supports Packrat and Renv lockfiles.
  - Yaml syntax:
  ```
  Lockfile:
    Type: [packrat/renv]
  ```

### The following items are configurable statically in the yml file or dynamically via the command-line.
* **Note:** Command-line flags will always overwrite the configuration in the yaml file.

* Library (directory) to install packages to
  - Yaml syntax
  ```
  Library: <path_to_directory>
  ```
* Update settings (should pkgr attempt to update outdated packages?)
  - Yaml syntax
  ```
  Update: <true/false>
  ```
* Logging level
  - Yaml syntax
  ```
  Loglevel: <info/debug/trace>
  ```
* Installation of “Suggested” packages config
  - Yaml syntax
  ```
  Suggested: <true/false>
  ```
* Number of threads to use
  - Yaml syntax
  ```
  Threads: <number_of_threads>
  ```
* Strict mode enabled/disabled
  - Yaml syntax
  ```
  Strict: <true/false>
  ```

## Default Settings
Below are the default settings for non-required pkgr configurations.
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

## Strict Mode
When strict mode is enabled, pkgr will refuse to run if unless the following conditions are satisfied:
* The library defined in pkgr.yml must exist.

## Rollback Feature
If any packages fail to install during a `pkgr install` command, pkgr will attempt to rollback
the user's specified `Library` to its previous state, removing any installed packages and reverting
any updated packages.
