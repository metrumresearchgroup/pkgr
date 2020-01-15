sources_path <- "get-sources-here/CRAN-739227e5b53e/src/"
repo_path <- "bad-xml2/"

toAdd <- list.files(path = sources_path, full.names = TRUE)

purrr::walk(toAdd, function(p) {
  drat::insertPackage(file = p, repodir = repo_path)
})

