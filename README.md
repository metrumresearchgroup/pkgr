# pkgr

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

