# mixed-source

## Description

Environment to help test that packages can be pulled from multiple repositories. Also helps test repo and type customizations.

## Expected behavior

* `pkgr plan --loglevel debug` indicates that packages are matched to repositories like so:
  - `rmarkdown`, `devtools` ---- CRAN-Micro.
  - `shiny` ---- CRAN (set in customizations)
  - `mrgsolve` ---- r_validated (source is only available in r_validated, and Customizations demand source)
* `pkgr install` installs the packages listed above
