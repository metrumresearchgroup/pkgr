# mixed-source
tags: multi-repo, repo-customizations, pkg-customizations, heavy

## Description

Environment to help test that packages can be pulled from multiple repositories. Also helps test repo and type customizations.

## Expected behavior

* `pkgr plan --loglevel debug` indicates that packages are matched to repositories like so:
  - shiny should pull from CRAN
  - pmplots should pull from r_validated
  - Everything else should pull from CRAN-Micro

```
pkgr plan --loglevel debug
INFO[0000] Installation would launch 11 workers         
INFO[0000] R Version 3.6.0                              
DEBU[0000] OS Platform x86_64-apple-darwin15.6.0        
INFO[0000] found installed packages                      count=0
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:11 (binary:source) packages available in for r_validated from https://metrumresearchgroup.github.io/r_validated
INFO[0000] 0:13396 (binary:source) packages available in for CRAN-Micro from https://cran.microsoft.com/snapshot/2018-11-18
INFO[0000] 14272:14762 (binary:source) packages available in for CRAN from https://cran.rstudio.com
DEBU[0000] package repository set                        pkg=rmarkdown relationship="user package" repo=CRAN-Micro type=source version=1.10
DEBU[0000] package repository set                        pkg=shiny relationship="user package" repo=CRAN type=binary version=1.3.2
DEBU[0000] package repository set                        pkg=devtools relationship="user package" repo=CRAN-Micro type=source version=2.0.1
DEBU[0000] package repository set                        pkg=pmplots relationship="user package" repo=r_validated type=source version=0.2.0
DEBU[0000] package repository set                        pkg=gtable relationship=dependency repo=CRAN-Micro type=source version=0.2.0
DEBU[0000] package repository set                        pkg=ps relationship=dependency repo=CRAN-Micro type=source version=1.2.1
DEBU[0000] package repository set                        pkg=base64enc relationship=dependency repo=CRAN-Micro type=source version=0.1-3
DEBU[0000] package repository set                        pkg=mime relationship=dependency repo=CRAN-Micro type=source version=0.6
DEBU[0000] package repository set                        pkg=assertthat relationship=dependency repo=CRAN-Micro type=source version=0.2.0
DEBU[0000] package repository set                        pkg=fansi relationship=dependency repo=CRAN-Micro type=source version=0.4.0
DEBU[0000] package repository set                        pkg=jsonlite relationship=dependency repo=CRAN-Micro type=source version=1.5
DEBU[0000] package repository set                        pkg=pkgconfig relationship=dependency repo=CRAN-Micro type=source version=2.0.2
DEBU[0000] package repository set                        pkg=clipr relationship=dependency repo=CRAN-Micro type=source version=0.4.1
DEBU[0000] package repository set                        pkg=Rcpp relationship=dependency repo=CRAN-Micro type=source version=1.0.0
DEBU[0000] package repository set                        pkg=backports relationship=dependency repo=CRAN-Micro type=source version=1.1.2
DEBU[0000] package repository set                        pkg=lazyeval relationship=dependency repo=CRAN-Micro type=source version=0.2.1
DEBU[0000] package repository set                        pkg=crayon relationship=dependency repo=CRAN-Micro type=source version=1.3.4
DEBU[0000] package repository set                        pkg=openssl relationship=dependency repo=CRAN-Micro type=source version=1.1
DEBU[0000] package repository set                        pkg=curl relationship=dependency repo=CRAN-Micro type=source version=3.2
DEBU[0000] package repository set                        pkg=ini relationship=dependency repo=CRAN-Micro type=source version=0.3.1
DEBU[0000] package repository set                        pkg=highr relationship=dependency repo=CRAN-Micro type=source version=0.7
DEBU[0000] package repository set                        pkg=rlang relationship=dependency repo=CRAN-Micro type=source version=0.3.0.1
DEBU[0000] package repository set                        pkg=BH relationship=dependency repo=CRAN-Micro type=source version=1.66.0-1
DEBU[0000] package repository set                        pkg=glue relationship=dependency repo=CRAN-Micro type=source version=1.3.0
DEBU[0000] package repository set                        pkg=digest relationship=dependency repo=CRAN-Micro type=source version=0.6.18
DEBU[0000] package repository set                        pkg=withr relationship=dependency repo=CRAN-Micro type=source version=2.1.2
DEBU[0000] package repository set                        pkg=sourcetools relationship=dependency repo=CRAN-Micro type=source version=0.1.7
DEBU[0000] package repository set                        pkg=stringi relationship=dependency repo=CRAN-Micro type=source version=1.2.4
DEBU[0000] package repository set                        pkg=utf8 relationship=dependency repo=CRAN-Micro type=source version=1.1.4
DEBU[0000] package repository set                        pkg=rstudioapi relationship=dependency repo=CRAN-Micro type=source version=0.8
DEBU[0000] package repository set                        pkg=xfun relationship=dependency repo=CRAN-Micro type=source version=0.4
DEBU[0000] package repository set                        pkg=bindr relationship=dependency repo=CRAN-Micro type=source version=0.1.1
DEBU[0000] package repository set                        pkg=RColorBrewer relationship=dependency repo=CRAN-Micro type=source version=1.1-2
DEBU[0000] package repository set                        pkg=git2r relationship=dependency repo=CRAN-Micro type=source version=0.23.0
DEBU[0000] package repository set                        pkg=viridisLite relationship=dependency repo=CRAN-Micro type=source version=0.3.0
DEBU[0000] package repository set                        pkg=magrittr relationship=dependency repo=CRAN-Micro type=source version=1.5
DEBU[0000] package repository set                        pkg=xtable relationship=dependency repo=CRAN-Micro type=source version=1.8-3
DEBU[0000] package repository set                        pkg=R6 relationship=dependency repo=CRAN-Micro type=source version=2.3.0
DEBU[0000] package repository set                        pkg=remotes relationship=dependency repo=CRAN-Micro type=source version=2.0.2
DEBU[0000] package repository set                        pkg=colorspace relationship=dependency repo=CRAN-Micro type=source version=1.3-2
DEBU[0000] package repository set                        pkg=whisker relationship=dependency repo=CRAN-Micro type=source version=0.3-2
DEBU[0000] package repository set                        pkg=labeling relationship=dependency repo=CRAN-Micro type=source version=0.3
DEBU[0000] package repository set                        pkg=yaml relationship=dependency repo=CRAN-Micro type=source version=2.2.0
DEBU[0000] package repository set                        pkg=plogr relationship=dependency repo=CRAN-Micro type=source version=0.2.0
DEBU[0000] package repository set                        pkg=clisymbols relationship=dependency repo=CRAN-Micro type=source version=1.2.0
DEBU[0000] package repository set                        pkg=evaluate relationship=dependency repo=CRAN-Micro type=source version=0.12
DEBU[0000] package repository set                        pkg=pkgbuild relationship=dependency repo=CRAN-Micro type=source version=1.0.2
DEBU[0000] package repository set                        pkg=httr relationship=dependency repo=CRAN-Micro type=source version=1.3.1
DEBU[0000] package repository set                        pkg=munsell relationship=dependency repo=CRAN-Micro type=source version=0.5.0
DEBU[0000] package repository set                        pkg=usethis relationship=dependency repo=CRAN-Micro type=source version=1.4.0
DEBU[0000] package repository set                        pkg=pillar relationship=dependency repo=CRAN-Micro type=source version=1.3.0
DEBU[0000] package repository set                        pkg=tibble relationship=dependency repo=CRAN-Micro type=source version=1.4.2
DEBU[0000] package repository set                        pkg=rcmdcheck relationship=dependency repo=CRAN-Micro type=source version=1.3.2
DEBU[0000] package repository set                        pkg=ggplot2 relationship=dependency repo=CRAN-Micro type=source version=3.1.0
DEBU[0000] package repository set                        pkg=markdown relationship=dependency repo=CRAN-Micro type=source version=0.8
DEBU[0000] package repository set                        pkg=stringr relationship=dependency repo=CRAN-Micro type=source version=1.3.1
DEBU[0000] package repository set                        pkg=bindrcpp relationship=dependency repo=CRAN-Micro type=source version=0.2.2
DEBU[0000] package repository set                        pkg=purrr relationship=dependency repo=CRAN-Micro type=source version=0.2.5
DEBU[0000] package repository set                        pkg=pkgload relationship=dependency repo=CRAN-Micro type=source version=1.0.2
DEBU[0000] package repository set                        pkg=reshape2 relationship=dependency repo=CRAN-Micro type=source version=1.4.3
DEBU[0000] package repository set                        pkg=sessioninfo relationship=dependency repo=CRAN-Micro type=source version=1.1.1
DEBU[0000] package repository set                        pkg=promises relationship=dependency repo=CRAN-Micro type=source version=1.0.1
DEBU[0000] package repository set                        pkg=tidyselect relationship=dependency repo=CRAN-Micro type=source version=0.2.5
DEBU[0000] package repository set                        pkg=prettyunits relationship=dependency repo=CRAN-Micro type=source version=1.0.2
DEBU[0000] package repository set                        pkg=plyr relationship=dependency repo=CRAN-Micro type=source version=1.8.4
DEBU[0000] package repository set                        pkg=cli relationship=dependency repo=CRAN-Micro type=source version=1.0.1
DEBU[0000] package repository set                        pkg=desc relationship=dependency repo=CRAN-Micro type=source version=1.2.0
DEBU[0000] package repository set                        pkg=dplyr relationship=dependency repo=CRAN-Micro type=source version=0.7.8
DEBU[0000] package repository set                        pkg=later relationship=dependency repo=CRAN-Micro type=source version=0.7.5
DEBU[0000] package repository set                        pkg=fs relationship=dependency repo=CRAN-Micro type=source version=1.2.6
DEBU[0000] package repository set                        pkg=gh relationship=dependency repo=CRAN-Micro type=source version=1.0.1
DEBU[0000] package repository set                        pkg=httpuv relationship=dependency repo=CRAN-Micro type=source version=1.4.5
DEBU[0000] package repository set                        pkg=tidyr relationship=dependency repo=CRAN-Micro type=source version=0.8.2
DEBU[0000] package repository set                        pkg=rprojroot relationship=dependency repo=CRAN-Micro type=source version=1.3-2
DEBU[0000] package repository set                        pkg=scales relationship=dependency repo=CRAN-Micro type=source version=1.0.0
DEBU[0000] package repository set                        pkg=callr relationship=dependency repo=CRAN-Micro type=source version=3.0.0
DEBU[0000] package repository set                        pkg=tinytex relationship=dependency repo=CRAN-Micro type=source version=0.9
DEBU[0000] package repository set                        pkg=xopen relationship=dependency repo=CRAN-Micro type=source version=1.0.0
DEBU[0000] package repository set                        pkg=knitr relationship=dependency repo=CRAN-Micro type=source version=1.20
DEBU[0000] package repository set                        pkg=memoise relationship=dependency repo=CRAN-Micro type=source version=1.1.0
DEBU[0000] package repository set                        pkg=processx relationship=dependency repo=CRAN-Micro type=source version=3.2.0
DEBU[0000] package repository set                        pkg=htmltools relationship=dependency repo=CRAN-Micro type=source version=0.3.6
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=82
INFO[0000] package installation sources                  CRAN=1 CRAN-Micro=80 r_validated=1
INFO[0000] package installation plan                     to_install=82 to_update=0
INFO[0000] resolution time 225.196953ms                 
```

* `pkgr install` installs the packages listed above according to the plan listed above.
