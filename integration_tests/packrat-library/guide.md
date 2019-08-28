# packrat-library

## Description

Environment to test functionality of packrat library, i.e., if  "Library" is not defined and Lockfile->Type is "packrat",
then the "Library" path is defined like a packrat library path:
   packrat/lib/{PLATFORM}/{R VERSION}

## Expected behavior

If the "packrat-like" folder does not exist, an error is shown:
pkgr install
INFO[0000] R Version 3.6.0                              
FATA[0000] package library not found at given library path  libraryPath=packrat/lib/x86_64-apple-darwin15.6.0/3.6.0


If the "packrat-like" folder does exist, packages are installed to the packrat-folder and the output is similar to:
pkgr install
INFO[0000] R Version 3.6.0                              
INFO[0000] found installed packages                      count=0
INFO[0000] Default package installation type:  binary   
INFO[0000] 14351:14857 (binary:source) packages available in for CRAN from https://cran.rstudio.com 
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=14
INFO[0000] package installation sources                  CRAN=14
INFO[0000] package installation plan                     to_install=14 to_update=0
INFO[0000] resolution time 226.638867ms                 
INFO[0000] downloading required packages within directory   dir=/Users/davidl/Library/Caches/pkgr
INFO[0000] all packages downloaded                       duration=1.233229ms
INFO[0000] starting initial install                     
INFO[0000] Successfully Installed                        package=assertthat repo=CRAN version=0.2.1
INFO[0000] Successfully Installed                        package=R6 repo=CRAN version=2.4.0
INFO[0000] Successfully Installed                        package=backports repo=CRAN version=1.1.4
INFO[0000] Successfully Installed                        package=zeallot repo=CRAN version=0.1.0
INFO[0000] Successfully Installed                        package=fansi repo=CRAN version=0.4.0
INFO[0000] Successfully Installed                        package=glue repo=CRAN version=1.3.1
INFO[0000] Successfully Installed                        package=utf8 repo=CRAN version=1.1.4
INFO[0000] Successfully Installed                        package=crayon repo=CRAN version=1.3.4
INFO[0000] Successfully Installed                        package=digest repo=CRAN version=0.6.20
INFO[0000] Successfully Installed                        package=rlang repo=CRAN version=0.4.0
INFO[0001] Successfully Installed                        package=cli repo=CRAN version=1.1.0
INFO[0001] Successfully Installed                        package=ellipsis repo=CRAN version=0.2.0.1
INFO[0001] Successfully Installed                        package=vctrs repo=CRAN version=0.2.0
INFO[0001] Successfully Installed                        package=pillar repo=CRAN version=1.4.2
INFO[0001] total install time                            duration=1.392967947s
INFO[0001] duration:1.630248569s                        

