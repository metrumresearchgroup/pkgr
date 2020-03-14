# Recommended
tags: recommended 

## Description
Test whether pkgr will ignore or attempt to install recommended packages in the dependency tree


## Expected Behaviors
* `pkgr plan --loglevel=debug` will indicate that survival is one of the dependencies to install
* `pkgr plan --config=pkgr-no-recommended.yml` will not show survival

```
pkgr plan --loglevel=trace --config=pkgr.yml  
TRAC[0000] setting directory to configuration file       cwd=/Users/devinp/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/recommended nwd=/Users/devinp/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/recommended
INFO[0000] Installation would launch 8 workers          
TRAC[0000] command args                                  RSettings="{{0 0 0} [] R {[]} map[] }" cmdArgs="[--version]" rpath=R
INFO[0000] R Version 3.6.2                              
DEBU[0000] OS Platform x86_64-apple-darwin15.6.0        
INFO[0000] Package Library will be created               path=test-library
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:11 (binary:source) packages available in for miniCRAN from ../minicran 
INFO[0000] 833:833 (binary:source) packages available in for MPN from https://mpn.metworx.com/snapshots/stable/2020-02-18 
DEBU[0000] package repository set                        pkg=survivalrec relationship="user package" repo=miniCRAN type=source version=0.0.1
TRAC[0000] dep config                                    config="{true true false true false}" pkg=survivalrec
TRAC[0000] dep config                                    config="{true true false true false}" pkg=magrittr
TRAC[0000] dep config                                    config="{true true false true false}" pkg=survival
TRAC[0000] skipping Depends dep                          dep=R pkg=survival
TRAC[0000] skipping Imports dep                          dep=methods pkg=survival
TRAC[0000] skipping Imports dep                          dep=splines pkg=survival
TRAC[0000] skipping Imports dep                          dep=stats pkg=survival
TRAC[0000] skipping Imports dep                          dep=utils pkg=survival
TRAC[0000] skipping Imports dep                          dep=graphics pkg=survival
TRAC[0000] dep config                                    config="{true true false true false}" pkg=Matrix
TRAC[0000] skipping Depends dep                          dep=R pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=stats pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=utils pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=methods pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=graphics pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=grid pkg=Matrix
TRAC[0000] dep config                                    config="{true true false true false}" pkg=lattice
TRAC[0000] skipping Depends dep                          dep=R pkg=lattice
TRAC[0000] skipping Imports dep                          dep=graphics pkg=lattice
TRAC[0000] skipping Imports dep                          dep=stats pkg=lattice
TRAC[0000] skipping Imports dep                          dep=utils pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grid pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grDevices pkg=lattice
TRAC[0000] dep config                                    config="{true true false true false}" pkg=Matrix
TRAC[0000] skipping Depends dep                          dep=R pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=stats pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=utils pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=methods pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=graphics pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=grid pkg=Matrix
TRAC[0000] dep config                                    config="{true true false true false}" pkg=lattice
TRAC[0000] skipping Depends dep                          dep=R pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grDevices pkg=lattice
TRAC[0000] skipping Imports dep                          dep=graphics pkg=lattice
TRAC[0000] skipping Imports dep                          dep=stats pkg=lattice
TRAC[0000] skipping Imports dep                          dep=utils pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grid pkg=lattice
TRAC[0000] dep config                                    config="{true true false true false}" pkg=survival
TRAC[0000] skipping Depends dep                          dep=R pkg=survival
TRAC[0000] skipping Imports dep                          dep=methods pkg=survival
TRAC[0000] skipping Imports dep                          dep=splines pkg=survival
TRAC[0000] skipping Imports dep                          dep=stats pkg=survival
TRAC[0000] skipping Imports dep                          dep=utils pkg=survival
TRAC[0000] skipping Imports dep                          dep=graphics pkg=survival
TRAC[0000] dep config                                    config="{true true false true false}" pkg=Matrix
TRAC[0000] skipping Depends dep                          dep=R pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=stats pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=utils pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=methods pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=graphics pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=grid pkg=Matrix
TRAC[0000] dep config                                    config="{true true false true false}" pkg=lattice
TRAC[0000] skipping Depends dep                          dep=R pkg=lattice
TRAC[0000] skipping Imports dep                          dep=stats pkg=lattice
TRAC[0000] skipping Imports dep                          dep=utils pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grid pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grDevices pkg=lattice
TRAC[0000] skipping Imports dep                          dep=graphics pkg=lattice
TRAC[0000] dep config                                    config="{true true false true false}" pkg=survivalrec
TRAC[0000] dep config                                    config="{true true false true false}" pkg=magrittr
TRAC[0000] dep config                                    config="{true true false true false}" pkg=survival
TRAC[0000] skipping Depends dep                          dep=R pkg=survival
TRAC[0000] skipping Imports dep                          dep=methods pkg=survival
TRAC[0000] skipping Imports dep                          dep=splines pkg=survival
TRAC[0000] skipping Imports dep                          dep=stats pkg=survival
TRAC[0000] skipping Imports dep                          dep=utils pkg=survival
TRAC[0000] skipping Imports dep                          dep=graphics pkg=survival
TRAC[0000] dep config                                    config="{true true false true false}" pkg=Matrix
TRAC[0000] skipping Depends dep                          dep=R pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=stats pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=utils pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=methods pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=graphics pkg=Matrix
TRAC[0000] skipping Imports dep                          dep=grid pkg=Matrix
TRAC[0000] dep config                                    config="{true true false true false}" pkg=lattice
TRAC[0000] skipping Depends dep                          dep=R pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grid pkg=lattice
TRAC[0000] skipping Imports dep                          dep=grDevices pkg=lattice
TRAC[0000] skipping Imports dep                          dep=graphics pkg=lattice
TRAC[0000] skipping Imports dep                          dep=stats pkg=lattice
TRAC[0000] skipping Imports dep                          dep=utils pkg=lattice
DEBU[0000] package repository set                        pkg=magrittr relationship=dependency repo=MPN type=binary version=1.5
DEBU[0000] package repository set                        pkg=lattice relationship=dependency repo=MPN type=binary version=0.20-38
DEBU[0000] package repository set                        pkg=Matrix relationship=dependency repo=MPN type=binary version=1.2-18
DEBU[0000] package repository set                        pkg=survival relationship=dependency repo=MPN type=binary version=3.1-8
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=5
INFO[0000] package installation sources                  MPN=4 miniCRAN=1
INFO[0000] package installation plan                     to_install=5 to_update=0
INFO[0000] resolution time 19.438319ms    
```


```              
pkgr plan --loglevel=trace --config=pkgr-no-recommended.yml   
TRAC[0000] setting directory to configuration file       cwd=/Users/devinp/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/recommended nwd=/Users/devinp/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/recommended
INFO[0000] Installation would launch 8 workers          
TRAC[0000] command args                                  RSettings="{{0 0 0} [] R {[]} map[] }" cmdArgs="[--version]" rpath=R
INFO[0000] R Version 3.6.2                              
DEBU[0000] OS Platform x86_64-apple-darwin15.6.0        
INFO[0000] Package Library will be created               path=test-library
INFO[0000] Default package installation type:  binary   
INFO[0000] 0:11 (binary:source) packages available in for miniCRAN from ../minicran 
INFO[0000] 833:833 (binary:source) packages available in for MPN from https://mpn.metworx.com/snapshots/stable/2020-02-18 
DEBU[0000] package repository set                        pkg=survivalrec relationship="user package" repo=miniCRAN type=source version=0.0.1
TRAC[0000] dep config                                    config="{true true false true true}" pkg=survivalrec
TRAC[0000] skipping Imports dep                          dep=survival pkg=survivalrec
TRAC[0000] dep config                                    config="{true true false true true}" pkg=magrittr
TRAC[0000] dep config                                    config="{true true false true true}" pkg=survivalrec
TRAC[0000] skipping Imports dep                          dep=survival pkg=survivalrec
TRAC[0000] dep config                                    config="{true true false true true}" pkg=magrittr
DEBU[0000] package repository set                        pkg=magrittr relationship=dependency repo=MPN type=binary version=1.5
INFO[0000] package installation status                   installed=0 not_from_pkgr=0 outdated=0 total_packages_required=2
INFO[0000] package installation sources                  MPN=1 miniCRAN=1
INFO[0000] package installation plan                     to_install=2 to_update=0
INFO[0000] resolution time 16.197121ms 
```
