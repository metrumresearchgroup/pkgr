# repo-local

## Description
Environment to help test that pkgr can install packages from local repositories.

## Expected Behaviors
* `pkgr plan` will indicate that packages will be pulled from miniCRAN
* `pkgr install` will complete successfully. pkgr install should have the following output:

```
pkgr install
INFO[0000] R Version 3.6.0                              
INFO[0000] found installed packages                      count=0
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:10 (binary:source) packages available in for miniCRAN from /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/minicran
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=9
INFO[0000] package installation sources                  miniCRAN=9
INFO[0000] package installation plan                     to_install=9 to_update=0
INFO[0000] resolution time 715.331µs                    
INFO[0000] downloading required packages within directory   dir=/Users/johncarlos/Library/Caches/pkgr
INFO[0000] downloading package                           package=jsonlite
INFO[0000] downloading package                           package=httr
INFO[0000] downloading package                           package=askpass
INFO[0000] downloading package                           package=openssl
INFO[0000] downloading package                           package=R6
INFO[0000] downloading package                           package=mime
INFO[0000] downloading package                           package=sys
INFO[0000] downloading package                           package=RColorBrewer
INFO[0000] downloading package                           package=curl
INFO[0000] all packages downloaded                       duration=3.725794ms
INFO[0000] starting initial install                     
INFO[0001] Successfully Installed                        package=R6 repo=miniCRAN version=2.4.0
INFO[0002] Successfully Installed                        package=RColorBrewer repo=miniCRAN version=1.1-2
INFO[0002] Successfully Installed                        package=mime repo=miniCRAN version=0.7
INFO[0002] Successfully Installed                        package=sys repo=miniCRAN version=3.2
INFO[0003] Successfully Installed                        package=askpass repo=miniCRAN version=1.1
INFO[0004] Successfully Installed                        package=curl repo=miniCRAN version=4.0
INFO[0005] Successfully Installed                        package=jsonlite repo=miniCRAN version=1.6
INFO[0011] Successfully Installed                        package=openssl repo=miniCRAN version=1.4.1
INFO[0015] Successfully Installed                        package=httr repo=miniCRAN version=1.4.1
INFO[0015] total install time                            duration=14.991598886s
INFO[0015] duration:15.006644989s
```
