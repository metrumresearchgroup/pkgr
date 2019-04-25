# simple-suggests

## Description

Environment to help test the "Suggests" field in a pkgr.yml file. This environment is configured for the same packages as [simple](../simple/guide.md).

## Expected Behavior

* `pkgr plan` should indicate that ~40 packages need to be installed.
* `pkgr install` should install all of the same packages as [simple](../simple/guide.md), plus ~32 other packages.
*  `pkgr inspect --deps` should return a large object that reflects the "Suggested" packages for R6 and pillar.

```
{
  "cli": [
    "assertthat",
    "crayon"
  ],
  "ggplot2": [
    "pkgconfig",
    "withr",
    "RColorBrewer",
    "utf8",
    "labeling",
    "gtable",
    "lazyeval",
    "stringi",
    "R6",
    "glue",
    "rlang",
    "magrittr",
    "assertthat",
    "crayon",
    "viridisLite",
    "digest",
    "colorspace",
    "Rcpp",
    "fansi",
    "stringr",
    "plyr",
    "cli",
    "munsell",
    "pillar",
    "scales",
    "reshape2",
    "tibble"
  ],
  "knitr": [
    "magrittr",
    "highr",
    "mime",
    "evaluate",
    "stringi",
    "xfun",
    "glue",
    "yaml",
    "markdown",
    "stringr"
  ],
  "lubridate": [
    "magrittr",
    "stringi",
    "glue",
    "Rcpp",
    "stringr"
  ],
  "markdown": [
    "mime"
  ],
  "munsell": [
    "colorspace"
  ],
  "pillar": [
    "utf8",
    "assertthat",
    "crayon",
    "fansi",
    "rlang",
    "cli"
  ],
  "plyr": [
    "Rcpp"
  ],
  "pryr": [
    "Rcpp",
    "magrittr",
    "stringi",
    "glue",
    "stringr"
  ],
  "reshape2": [
    "glue",
    "magrittr",
    "stringi",
    "Rcpp",
    "plyr",
    "stringr"
  ],
  "scales": [
    "Rcpp",
    "viridisLite",
    "labeling",
    "colorspace",
    "R6",
    "RColorBrewer",
    "munsell"
  ],
  "stringr": [
    "stringi",
    "glue",
    "magrittr"
  ],
  "testthat": [
    "praise",
    "digest",
    "R6",
    "crayon",
    "rlang",
    "assertthat",
    "withr",
    "magrittr",
    "cli"
  ],
  "tibble": [
    "crayon",
    "rlang",
    "utf8",
    "assertthat",
    "pkgconfig",
    "fansi",
    "cli",
    "pillar"
  ]
}

```
