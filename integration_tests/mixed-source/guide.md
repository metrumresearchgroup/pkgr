# mixed-source

## Description

Environment to help test that packages can be pulled from multiple repositories. Also helps test repo and type customizations.

## Expected behavior

* `pkgr plan` indicates that packages are matched to repositories like so:
  - `rmarkdown`, `devtools` ---- CRAN-Micro (this should be the first listed repository that these packages are found in)
  - `shiny` ---- CRAN (set in customizations)
  - `mrgsolve` ---- r_validated (source is only available in r_validated, and Customizations demand source)
* `pkgr install` installs the packages listed above
