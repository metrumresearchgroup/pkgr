# pkgr

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

