
Tools for validating Go executables
===================================


Overview
--------

As part of validating a release of a package, MetrumRG produces two
artifacts, the source code archive and a scorecard report.  For non-R
packages, the [mpn.scorecard][ms] R package takes care of *rendering*
the scorecard report.

[ms]: https://github.com/metrumresearchgroup/mpn.scorecard

This repository includes Go commands, shell scripts, and Makefile
rules that Go modules can use to generate the source code archive and
the inputs for mpn.scorecard.  Its current focus is on validating
functionality exposed via **command-line executables**.


Key components
--------------

 * [rules.mk][]: a Makefile that defines a `vt-all` target that
   generates the source code archive and mpn.scorecard inputs

 * [checkmat/][]: Go package that defines an executable for running
   various checks on a traceability matrix YAML file

 * [filecov/][]: Go package that defines an executable for calculating
    *file*-level coverage from a Go coverage profile

 * [fmttests/][]: Go package that defines an executable for converting
   `go test -json` into a custom format meant for inclusion in the
   scorecard

 * [scripts/][]: various shell scripts invoked by targets in
   [rules.mk][]

[rules.mk]: ./rules.mk
[checkmat/]: ./checkmat
[filecov/]: ./filecov
[fmttests/]: ./fmttests
[scripts/]: ./scripts/


Setup instructions
------------------

High-level summary of the main steps (details in subsections below):

 1. add this repository as a subtree in the main repository

 2. add documentation file for each command

 3. create a traceability matrix YAML

 4. add scripts for running the tests

 5. include this subtree's `rules.mk` into the top-level Makefile

### 1. Add subtree

This repository is designed to be incorporated as a subtree in the
main project.  The suggested location is `internal/valtools/`.

You can use [git subtree][gs] to manage the subtree.

[gs]: https://manpages.debian.org/stable/git-man/git-subtree.1.en.html

Run `go mod tidy` to update `go.mod` with new dependencies, if any,
brought in by the subtree.

### 2. Command documentation

Consider a Go module that defines one executable, `foo`, where the
user-facing functionality is exposed via two subcommands, `foo bar`
and `foo baz`.  In order to define the traceability matrix, there must
be a directory that contains a documentation file for each of these
commands, replacing any spaces in the name with underscores.

    docs/commands/
    |-- foo.md
    |-- foo_bar.md
    `-- foo_baz.md

`docs/commands/` is taken as the directory for the command
documentation unless the `VT_DOC_DIR` variable specifies another
locatino (see [step 5](#step5)).

The Makefile rules expect the main repository to define a `docgen`
package (suggested location: `internal/tools/docgen`) whose executable
generates the documentation.

The `docgen` executable should accept one argument, the directory in
which to write the documentation files.  The executable is responsible
for ensuring that no stale documentation remains in the directory
(e.g., by removing the directory before writing new files).

If a module uses `cobra`, `docgen` can likely be defined as a light
wrapper around the `cobra/doc.GenMarkdownTree` function.

### 3. Define traceability matrix

Create a traceability matrix YAML file that maps each command to the
code file that defines it, the file that documents it, and the main
files that test it.  By default, the Makefile rules look for the
matrix YAML file at `docs/validation/matrix.yaml`.  To change this
location, set the `VT_MATRIX` variable (see [step 5](#step5)).

The file should consist of a sequence of entries with following items:

 * `entrypoint`: the name of the command, as invoked by the user

 * `code`: the path to where the command is defined

 * `doc`: path to the command's main documentation

 * `tests`: list of paths where the command is tested

Example entry:

    - entrypoint: foo bar
      code: cmd/bar.go
      doc: docs/commands/foo_bar.md
      tests:
        - cmd/bar_test.go
        - integration/bar_test.go

The `checkmat` tool will flag any documentation file that does not map
to a command with an entry in the YAML file.  In some cases, you may
not want to include a command in the rendered matrix.  For example,
the top-level `foo` command may serve only as an entry point for
subcommands, making it "uninteresting" to include in the matrix.  For
such cases, you can add a skip entry to the matrix.

    - command: foo
      skip: true

### 4. Add test runners

The main repository must define one or more scripts for running its
tests.  Specify the paths to these scripts via the `VT_TEST_RUNNERS`
variable (see [step 5](#step5)).

A test script must

 * write `go test` JSON records to standard output when passed the
   `-json` argument.  It must not write anything else to standard
   output in this case.

 * check whether the `GOCOVERDIR` environment variable is set and, if
   so, instrument the tests to write coverage data files under it.

How to handle the `GOCOVERDIR` environment variable is determined by
whether tests are integration tests that use a built executable or
unit tests.

 * integration tests: when `GOCOVERDIR` is set, the script should pass
   `-cover` to the `go build` call to build the instrumented
   executable(s) to test

 * unit tests: when `GOCOVERDIR` is set, `go test` should be
   instructed to write coverage data files to `GOCOVERDIR`.  This can
   be done by adding `-args -test.gocoverdir="$GOCOVERDIR"` to the end
   of the call.

   When tallying coverage, Go does not by default count a statement as
   covered if it's only executed via another packge's unit tests.  To
   change that, list all the module's packages of interest by
   specifying `-coverpkg` in the `go test` call.

Notes:

 * Go 1.20 [introduced support][newcov] for `GOCOVERDIR` and
   instrumenting executables.

 * There's a [proposed patch][covarg] to expose `test.gocoverdir` as
   top-level argument to `go test`, although it's currently on hold.

[newcov]: https://go.dev/blog/integration-test-coverage
[covarg]: https://go-review.googlesource.com/c/go/+/456595/14


<a id="step5"></a>

### 5. Wire up Makefile

To wire up the subtree to the main repository, include the subtree's
`rule.mk` file in the repository's top-level Makefile.  Before the
`include` directive, you can specify any Makefile [variables](#vars),
but, at a minimum, you should set `VT_TEST_RUNNERS`.

    VT_TEST_RUNNERS = scripts/run-unit-tests
    VT_TEST_RUNNERS += scripts/run-integration-tests
    include internal/valtools/rules.mk


Running the pipeline
--------------------

The `vt-all` target provides the main entry point.  It generates the
source archive and mpn.scorecard inputs.

    make vt-all

The generated files are written under the directory specified by the
variable `VT_OUT_DIR`.  By default, this points to
`{subtree}/output/{package}_{version}`.

<a id="vars"></a>

### Makefile variables

 * `VT_BIN_DIR`: where to install executables (default:
   `{subtree}/bin`)

 * `VT_DOC_DIR`: tell `docgen` executable to generate documentation
   files under this directory (default: `docs/commands`)

 * `VT_MATRIX`: path to matrix file (default:
   `docs/validation/matrix.yaml`)

 * `VT_OUT_DIR`: where to generate the results (default:
   `{subtree}/output/{package}_{version}`)

 * `VT_PKG`: name of the package (default: the base name of the
   top-level directory).

 * `VT_TEST_ALLOW_SKIPS`: whether to allow skips when running the
   `VT_TEST_RUNNERS` scripts (default: `no`)

 * `VT_TEST_RUNNERS`: a space-delimited list of scripts to invoke to
   run the test suite

### Auxillary targets

In addition to `vt-all`, the following targets can be useful to run
directly:

 * `vt-gen-docs`: invoke the `docgen` exectuable to refresh the
   documentation in `VT_DOC_DIR`

 * `vt-test`: invoke each script in `VT_TEST_RUNNERS` *without*
   coverage enabled

Run `vt-help` target to see a more complete list of targets.
