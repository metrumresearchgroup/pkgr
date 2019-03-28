# pkgr

[![asciicast](https://asciinema.org/a/wgcPBvCMtEwhpdW793MBjgSi2.svg)](https://asciinema.org/a/wgcPBvCMtEwhpdW793MBjgSi2)

# THIS IS CURRENTLY A WIP, however is getting close for user testing. Check back soon for more comprehensive user docs

# What is pkgr?

`pkgr` is a rethinking of the way packages are managed in R. Namely, it embraces
the declarative philosophy of defining _ideal state_ of the entire system, and working
towards achieving that objective. Furthermore, `pkgr` is built with a focus on reproducibility
and auditability of what is going on, a vital component for the pharmaceutical sciences + enterprises.

# Why pkgr?

`install.packages` and friends such as `remotes::install_github` have a subtle weakness --
they are not good at controlling desired global state. There are some knobs that
can be turned, but overall their APIs are generally not what the user _actually_ needs. Rather, they
are the mechanism by which the user can strive towards their needs, in a forceably iterative fashion.

With pkgr, you can, in a **parallel-processed** manner, do things like:
- Install a number of packages from various repositories, when specific packages must be pulled from specific repositories
- Install `Suggested` packages only for a subset of all packages you'd like to install
- Customize the installation behavior of a single package in a documentable and reproducible way
  - Set custom Makevars for a package that persist across system installations
  - Install source versions of some packages but binaries for others
- **Understand how your R environment will be changed _before_ performing an installation or action.**

Today, packages are highly interwoven. Best practices have pushed towards small, well-scoped packages that
do behaviors well. For example, rather than just having plyr, we now use dplyr+purrr to achieve
the same set of responsibilities (dealing with dataframes + dealing with other list/vector objects in an iterative way).
As such, it is becoming increasingly difficult to manage the _set_ of packages in a transparent and robust
way.

# How it works

`pkgr` is a command line utility with several top level commands. The two primary commands are:

```bash
pkgr plan # show what would happen if install is run
pkgr install # install the packages specified in pkgr.config
```

The actions are controlled by a configuration file that specifies the desired global state, namely,
by defining the top level packages a user cares about, as well as specific configuration customizations.

An example pkgr configuration file might look like:

```yaml
Version: 1
# top level packages
Packages:
  - rmarkdown
  - bitops
  - caTools
  - knitr
  - tidyverse
  - shiny
  - logrrr

# any repositories, order matters
Repos:
  - gh_dev: "https://metrumresearchgroup.github.io/rpkgs/gh_dev"
  - CRAN: "https://cran.microsoft.com/snapshot/2018-11-18"

# path to install packages to
Library: "path/to/install/library"

# package specific customizations
Customizations:
  - tidyverse:
      Suggests: true
```

Another such example, is on CRAN, for OSX, the new devtools (v2.x) is currently available as source,
however the binary is still v1.13. To control and say we would prefer the source version of devtools,
while relying on the platform default (binaries) for all other packages, a customization can be set
for devtools as `Type: source`

```yaml
Version: 1
# top level packages
Packages:
  - rmarkdown
  - shiny
  - devtools

# any repositories, order matters
Repos:
   - CRAN: "https://cran.microsoft.com/snapshot/2018-11-18"

Library: "path/to/install/library"

# can cache both the source and installed binary versions of packages
Cache: "path/to/global/cache"

# can log the actions and outcomes to a file for debugging and auditing
Logging:
  File: pkgr-install.log

Customizations:
  - devtools:
      Type: source
```

A configuration that also pulls from bioconductor:

```
Version: 1
# top level packages
Packages:
  - magrittr
  - rlang
  - ggplot2
  - dplyr
  - tidyr
  - plotly
  - VennDiagram
  - aws.s3
  - data.table
  - forcats
  - preprocessCore
  - loomR
  - ggthemes
  - reshape

# any repositories, order matters
Repos:
  - CRAN: "https://cran.microsoft.com/snapshot/2018-11-18"
  - BioCsoft: "https://bioconductor.org/packages/3.8/bioc"
  - BioCann: "https://bioconductor.org/packages/3.8/data/annotation"
  - BioCexp: "https://bioconductor.org/packages/3.8/data/experiment"
  - BioCworkflows: "https://bioconductor.org/packages/3.8/workflows"

# path to install packages to
Library: pkgs

Cache: pkgcache
Logging:
  File: pkgr-install.log
```

TODO:

- how integrates with packrat
- how integrates with rstudio package manager
- performance characteristics (so much faster than install.packages)
- caveats


## API options

package declaration can become nuanced as the user desires to customize
specifically where a package is pulled from.
Given a set of repositories, the default R tooling
will stop after the package is found and use that.
In some cases the user may prefer to explicitly
declare which repository a package may come from, especially when pulling from
an environment where multiple repos are specified.

The remotes API provides a rich experience around customizing the installation
behavior from external repositories such as github, where combinations such as

- tidyverse/dplyr#12345 - install from PR 123
- github::tidyverse/dplyr#12345 - install from github

The intent of this tooling (for now) is to provide a more modular experience,
in which the ways packages can be identified is minimized, and upstream tooling
can coalesce packages from many of the scenarios outlined.

As such, the biggest focus is on targetting packages placed in a specific _repository_.
With this in mind, the question remains whether the remotes API should be followed.
The concern is that the repo::package pattern slightly obfusicates the package.
This is less noticeable when the package is previously declared in the Imports/Depends
statement of a DESCRIPTION file, however when packages become the forefront
of the requirements.

Some of the potential API designs are:

```yaml
Packages:
  - repository::package
  - package@repository
  - package
    repo: repository
```

### Package dependencies:

By default, packages will need Imports/Depends/LinkingTo to make
sure the packages can work successfully.

```yaml
Packages:
  - PKPDmisc
    Suggests: true
  - dplyr
```

The benefit is customizations related to package requirements
are immediately visible. The downside is it "pollutes" the
packages list.

```yaml
Packages:
  - PKPDmisc
  - dplyr

Customizations:
  PKPDmisc:
    Suggests: true
```


## Assumptions (for now)

Making this tool bulletproof will take a significant effort over time. To bring confidence for use day-to-day
we must clearly outline the assumptions we are making internally to provide guidance on the areas
this tool may not be suitable for, or to explain unexpected behavior.

* Package/versions from a given repo will not change over time
  * if pkgx_0.1.0 is downloaded from repoY, we do not need to check each time that pkgx is consistent
  * this allows simple caching without doing hash comparisons (for now)

R package management

## Install Strategy Background

One of the problems with the full layered implementation is the longest install time dictates the entire layer
installation install. Originally, we did not know if this would be a huge problem, however it was quickly
evident that this was not the case.

For example, when look at the installation layers given a request for ggplot2, the following was
the installation timing. For layer 1, the second longest install time was Rcpp (39 seconds), with most
other packages coming in less than 10 seconds.

| layer  |package | duration|
|:-------|:-------|--------:|
|   1    |stringi |   159.39|
|   2    |Matrix  |    69.98|
|   3    |mgcv    |    34.45|
|   4    |tibble  |     2.37|
|   5    |ggplot2 |    12.12|

Furthermore, when looking at subsequent layers, neither Matrix or mgcv and its dependencies have any relation
to stringi, so there is no reason to wait for the layer to complete.

# Development

run all tests with tabular output:

```
go test ./... -json -cover | tparse -all
```
