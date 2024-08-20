
## Overview

The package library to be created is defined by a configuration file.
`pkgr` reads the configuration from `pkgr.yml` in the current working
directory unless another location is specified via its `--config`
option.

Here's an example of a minimal configuration:

```yaml {filename="pkgr.yml"}
Version: 1
Packages:
  - here
Repos:
  - MPN: https://mpn.metworx.com/snapshots/stable/2024-06-12
Library: lib
```

That instructs `pkgr` to install `here` and its dependencies, where to
popuplate the library (`lib/` under the current working directory),
and where to get the package and its dependencies from (MPN's
2024-06-12 snapshot).

Running `pkgr install` from the directory that contains that
`pkgr.yml` file crates the library with `here` and its dependency
`rprojroot`:

```
lib/
|-- here/
`-- rprojroot/
```

### Path handling

Several values in the configuration may be paths (e.g., `Library` and
`Descriptions`).

 * A `~` in the path is expanded to the home directory.

 * If the path is a relative path, it is interpreted as relative to
   the *configuration file*.  This is typically the same location from
   which `pkgr` was invoked but may differ if an alternative
   configuration file was specified with `--config`.

### Expansion of environment variables

If a dollar sign precedes a word (with optional brackets around the
word), the sequence is replaced by the value of the environment
variable with that name.  If that variable is undefined, it is
replaced by an empty string.

This is most commonly used to set `Rpath` based on an environment
variable (e.g., `${R_EXE_4_3}`).

> [!CAUTION]
> Dollar signs should be avoided in configuration unless they
> reference an environment variable, as there is no support for
> escaping the dollar sign to disable the expansion.

## Primary sections

The sections below are the most commonly used.

### Customizations

Specify package-specific or repo-specific customizations.

#### Package customizations

Package-level customization is specified as a list of items under
`Customizations: Packages`.

 * **Env**: a set of environment variables to set when installing the
   package.  This is most commonly used to set `R_MAKEVARS_USER`.

   ```yaml {filename="Example"}
   Customizations:
     Packages:
        - RCurl:
            Env:
              R_MAKEVARS_USER: ~/.R/Makevars-RCurl
   ```

 * **Repo**: install the package from this repository (must match a
   key under top-level `Repos` section)

   ```yaml {filename="Example"}
   Customizations:
     Packages:
       - here:
           Repo: PPM
   ```

 * **Suggests**: whether to install the suggested dependencies for a
   package

   ```yaml {filename="Example"}
   Customizations:
     Packages:
       - here:
           Suggests: true
   ```

 * **Type**: type of package, "source" or "binary", to install for a
   package

   Installing a binary package is faster, but the CRAN-like repository
   must have a binary available for the package, which is often not
   the case.

   ```yaml {filename="Example"}
   Customizations:
     Packages:
       - here:
           Type: binary
   ```

#### Repo customizations

 * **Type**: type of packages, "source" or "binary", to install from
   the specified repository

   The default type is "source" on Linux and "binary" on macOS and
   Windows.

   > [!NOTE]
   > It is common for a repository to be missing binaries available
   > for your setup (platform, architecture, and version of R).  If
   > binaries are unavailable and you're running macOS or Windows, you
   > must customize the repository to have type "source" and ensure
   > your machine has the necessary system requirements to build the
   > package from source.

   ```yaml {filename="Example"}
   Customizations:
     Repos:
       - MPN:
           Type: source
   ```

<!-- Note: RepoType and RepoSuffix are left undocumented. -->

### Descriptions

Install the dependencies listed in the specified R package
[DESCRIPTION][desc] files.

[desc]: https://r-pkgs.org/description.html

This main use case for this is R package development, as it avoids
duplicating the list of dependencies under the `Packages` section of
the `pkgr` configuration.

```yaml {filename="Example"}
Descriptions:
  - DESCRIPTION
```

### Library

Directory in which to install packages.

A directory specified via the `--library` command-line option
overrides the value specified in the configuration.

Either this section or `Lockfile` must be present.  `Library` takes
precedence over `Lockfile` if both are specified.

```yaml {filename="Example"}
Library: lib
```

### Lockfile

Use the same library location as
[renv](https://rstudio.github.io/renv) or
[packrat](https://github.com/rstudio/packrat).

> [!WARNING]
> Packrat was superseded by renv.  It is in maintenance-only mode.

`pkgr` invokes `renv` to find the library, as its location depends on
a number of settings.  If `renv` has not yet been initialized for a
project, `pkgr` tries to query the `renv`, if any, available in the
default library paths.  However, the recommended approach is to
initialize `renv` *before* calling `pkgr install` (e.g., by invoking
`renv::init(bare = TRUE)`).

Either this section or `Library` must be present.  `Library` takes
precedence over `Lockfile` if both are specified.

```yaml {filename="Example"}
Lockfile:
  Type: renv
```

### Packages

Packages to install.  They must be available from at least one of the
CRAN-like repositories listed under the `Repos` section.

Any of their dependencies are installed automatically.  By default,
suggested dependencies are not included (but see `Suggests` section).

```yaml {filename="Example"}
Packages:
  - here
  - rlang
```

### Repos

CRAN-like repositories from which to retrieve the packages listed in
`Packages`.

> [!CAUTION]
> The order is important. Unless customized in the `Customizations`
> section, a package is downloaded from the first repository that
> contains it.

> [!TIP]
> To create a library with a reproducible set of packages, use
> CRAN-like repositories that snapshot packages at a particular date,
> such as [Metrum Package Network][mpn] and
> [Posit Package Manager][ppm].

[mpn]: https://mpn.metworx.com
[ppm]: https://packagemanager.posit.co

```yaml {filename="Example"}
Repos:
  - MPN: https://mpn.metworx.com/snapshots/stable/2024-06-12
  - PPM: https://packagemanager.posit.co/cran/2024-06-12
```

### RPath

Which R executable to use.  Defaults to the highest priority R
executable in `PATH` (i.e. the same R that would launch if you invoked
`R` from the current shell).

```yaml {filename="Example"}
Rpath: /opt/R/4.2.3/bin/R
```

### Suggests

Install the suggested dependencies for all packages listed in
`Packages`.

To do so only for some packages, use a package-specific customizations
entry (see `Customizations` section).

```yaml {filename="Example"}
Suggests: true
```

### Version

The configuration format version.  This should always be `1`.

```yaml {filename="Example"}
Version: 1
```

## Other sections

These sections are not as commonly used as the ones above.

### Cache

Store downloaded packages and built binaries in this directory rather
than the default platform-specific directory (e.g.,
`$HOME/.cache/pkgr/` on Linux).

```yaml {filename="Example"}
Cache: cache
```

### IgnorePackages

Do not install the specified packages even if they are a required
dependency.

```yaml {filename="Example"}
IgnorePackages:
- foo
```

<!-- Note(km): leaving LibPaths undocumented.  From my quick testing, it doesn't -->
<!-- behave as described in the user manual. -->

### Logging

Configure logging behavior.

 * **all**: path for all commands to write log records to

 * **install**: path for `pkgr install` to write log records to
   instead of the one specified by **all**.

 * **overwrite**: overwrite the above files instead of the default
   behavior of appending to them

<!-- Note(km): Leave `level` undocumented because, as far as I can tell, -->
<!-- it's not wired up. -->

```yaml {filename="Example"}
Logging:
  all: pkgr-all.log
  install: pkgr-install.log
  overwrite: true
```

### NoRecommended

By default `pkgr intall` considers [recommended packages][rp] when
installing and updating a library.  If `NoRecommended` is set to
`true`, recommended packages are ignored unless they are explicitly
listed under the `Packages` section.

[rp]: https://rstudio.github.io/r-manuals/r-admin/Obtaining-R.html#using-subversion-and-rsync

```yaml {filename="Example"}
NoRecommended: true
```

### NoRollback

By default `pkgr install` restores the library to its original state
if the installation fails.  Set `NoRollback` to `true` to disable that
behavior.

```yaml {filename="Example"}
NoRollback: true
```

### NoSecure

Disable the default TLS certificate verification by setting
`NoSecure: true`.

> [!CAUTION]
> Disable TLS certificate verification is not recommended.  If you
> need this in a particular case, consider passing the `--no-secure`
> command-line flag instead.

```yaml {filename="Example"}
NoSecure: true
```

### NoUpdate

Disable the default behavior of updating packages that are already
present in the library.

This can be set for a single call by passing the `--no-update`
command-line flag instead.

```yaml {filename="Example"}
NoUpdate: true
```

### Strict

`pkgr install` creates the library directory if needed.  Set `Strict`
to `true` to make it abort if the library directory does not exist.

This can be set for a single call by passing the `--strict`
command-line flag instead.

```yaml {filename="Example"}
Strict: true
```

### Tarballs

Install packages from the specified source tarball.

```yaml {filename="Example"}
Tarballs:
  - packages/R6_2.5.1.tar.gz
  - packages/glue_1.7.0.tar.gz
```

### Threads

Number of workers to use for installing packages.  The default is to
use up to eight workers, considering the number of workers on the
machine and any configured Linux CPU quotas.

A value specified via the `--threads` command-line option overrides
the value specified in the configuration.

```yaml {filename="Example"}
Threads: 2
```
