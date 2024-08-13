#!/bin/sh

# Set up a local CRAN-like repository. In addition to source packages,
# it has binary packages under
# bin/linux/ubuntu/{codename}/pkgr/test/contrib/{Rmajor}.{Rminor}/
#
# Write the repository location to stdout.

set -eu

id=$(lsb_release -si | tr '[:upper:]' '[:lower:]')
test "$id" = ubuntu || {
    printf 'this script is only compatible with Ubuntu, not %s\n' "$id"
    exit 1
}

root=$(mktemp -d "${TMPDIR:-/tmp}"/pkgr-mixed-source-test-XXXXX)

cd "$root"
cat <<EOF >pkgr.yml
Version: 1
Repos:
- MPN: https://mpn.metworx.com/snapshots/stable/2020-07-19
Library: lib
Cache: pkgcache
Packages:
- R6
- ellipsis
- digest
- yaml
Customizations:
  Packages:
    - yaml:
        Suggests: true
EOF

pkgr install >&2
bin=$(ls -d pkgcache/*/binary/*)

nl='
'
case "$bin" in
    *"$nl"*)
        printf >&2 'found more than one binary subdirectory:\n%s\n' "$bin"
        exit 1
        ;;
    '')
        printf >&2 'failed to find binary subdirectory\n'
        exit 1
        ;;
esac

rvers=${bin##*/}
repo=$root/repo

codename=$(lsb_release -sc)
repo_bin=$repo/bin/linux/ubuntu/$codename/pkgr/test/contrib/$rvers
mkdir -p "$repo_bin"
mv "$bin"/*.tar.gz "$repo_bin/"

src=${bin%/binary/*}/src
repo_src=$repo/src/contrib
mkdir -p "$repo_src"
mv "$src"/*.tar.gz "$repo_src"

write_packages () {
    # Note: Despite the confusing value name, 'type = "source"' is the
    # value that's needed for Linux binaries.
    Rscript -e 'tools::write_PACKAGES("'"$1"'", type = "source")' >&2
}

write_packages "$repo_src"
write_packages "$repo_bin"

printf '%s\n' "$repo"
