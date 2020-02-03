sources_path <- "get-sources-here/CRAN-eed1668927b6/src/"
repo_path <- "./simple-no-crayon"

toAdd <- list.files(path = sources_path, full.names = TRUE)

purrr::walk(toAdd, function(p) {
  drat::insertPackage(file = p, repodir = repo_path)
})

