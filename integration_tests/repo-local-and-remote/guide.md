# repo-local
tags: repo-local, repo-remote

## Description
Environment to help test that pkgr can install packages from a mix of local and remote repositories.

## Expected Behaviors
* `pkgr plan --loglevel debug` will indicate that packages will be pulled according to the following output:

```
pkgr plan --loglevel debug
INFO[0000] Installation would launch 11 workers         
INFO[0000] R Version 3.6.0                              
DEBU[0000] OS Platform x86_64-apple-darwin15.6.0        
INFO[0000] found installed packages                      count=0
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:10 (binary:source) packages available in for miniCRAN from /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/minicran
INFO[0000] 0:13396 (binary:source) packages available in for CRAN from https://cran.microsoft.com/snapshot/2018-11-18
DEBU[0000] package repository set                        pkg=RColorBrewer relationship="user package" repo=miniCRAN type=source version=1.1-2
DEBU[0000] package repository set                        pkg=httr relationship="user package" repo=miniCRAN type=source version=1.4.1
DEBU[0000] package repository set                        pkg=promises relationship="user package" repo=CRAN type=source version=1.0.1
DEBU[0000] package repository set                        pkg=Rcpp relationship=dependency repo=CRAN type=source version=1.0.0
DEBU[0000] package repository set                        pkg=rlang relationship=dependency repo=CRAN type=source version=0.3.0.1
DEBU[0000] package repository set                        pkg=curl relationship=dependency repo=miniCRAN type=source version=4.0
DEBU[0000] package repository set                        pkg=mime relationship=dependency repo=miniCRAN type=source version=0.7
DEBU[0000] package repository set                        pkg=magrittr relationship=dependency repo=CRAN type=source version=1.5
DEBU[0000] package repository set                        pkg=sys relationship=dependency repo=miniCRAN type=source version=3.2
DEBU[0000] package repository set                        pkg=BH relationship=dependency repo=CRAN type=source version=1.66.0-1
DEBU[0000] package repository set                        pkg=jsonlite relationship=dependency repo=miniCRAN type=source version=1.6
DEBU[0000] package repository set                        pkg=R6 relationship=dependency repo=miniCRAN type=source version=2.4.0
DEBU[0000] package repository set                        pkg=openssl relationship=dependency repo=miniCRAN type=source version=1.4.1
DEBU[0000] package repository set                        pkg=later relationship=dependency repo=CRAN type=source version=0.7.5
DEBU[0000] package repository set                        pkg=askpass relationship=dependency repo=miniCRAN type=source version=1.1
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=15
INFO[0000] package installation sources                  CRAN=6 miniCRAN=9
INFO[0000] package installation plan                     to_install=15 to_update=0               

```

* `pkgr install` will successfully install the packages according to the plan specified above. Output should resemble below:

```
pkgr install
INFO[0000] R Version 3.6.0                              
INFO[0000] found installed packages                      count=0
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:10 (binary:source) packages available in for miniCRAN from /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/minicran
INFO[0000] 0:13396 (binary:source) packages available in for CRAN from https://cran.microsoft.com/snapshot/2018-11-18
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=15
INFO[0000] package installation sources                  CRAN=6 miniCRAN=9
INFO[0000] package installation plan                     to_install=15 to_update=0
INFO[0000] resolution time 126.967307ms                 
INFO[0000] downloading required packages within directory   dir=/Users/johncarlos/Library/Caches/pkgr
INFO[0000] downloading package                           package=rlang
INFO[0000] downloading package                           package=promises
INFO[0000] downloading package                           package=magrittr
INFO[0000] downloading package                           package=Rcpp
INFO[0000] downloading package                           package=BH
INFO[0000] downloading package                           package=later
INFO[0007] all packages downloaded                       duration=7.425887275s
INFO[0007] starting initial install                     
INFO[0010] Successfully Installed                        package=magrittr repo=CRAN version=1.5
INFO[0010] Successfully Installed                        package=R6 repo=miniCRAN version=2.4.0
INFO[0010] Successfully Installed                        package=mime repo=miniCRAN version=0.7
INFO[0010] Successfully Installed                        package=sys repo=miniCRAN version=3.2
INFO[0010] Successfully Installed                        package=RColorBrewer repo=miniCRAN version=1.1-2
INFO[0012] Successfully Installed                        package=askpass repo=miniCRAN version=1.1
INFO[0013] Successfully Installed                        package=curl repo=miniCRAN version=4.0
INFO[0014] Successfully Installed                        package=jsonlite repo=miniCRAN version=1.6
INFO[0017] Successfully Installed                        package=rlang repo=CRAN version=0.3.0.1
INFO[0020] Successfully Installed                        package=openssl repo=miniCRAN version=1.4.1
INFO[0025] Successfully Installed                        package=httr repo=miniCRAN version=1.4.1
INFO[0037] Successfully Installed                        package=Rcpp repo=CRAN version=1.0.0
INFO[0056] Successfully Installed                        package=BH repo=CRAN version=1.66.0-1
INFO[0074] Successfully Installed                        package=later repo=CRAN version=0.7.5
INFO[0081] Successfully Installed                        package=promises repo=CRAN version=1.0.1
INFO[0081] total install time                            duration=1m13.94399499s
INFO[0081] duration:1m21.507678675s  
```
