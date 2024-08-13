vtdir := $(dir $(lastword $(MAKEFILE_LIST)))
vtdir := $(patsubst %/,%,$(vtdir))
ifeq ($(strip $(vtdir)),)
$(error "bug: vtdir is unexpectedly empty")
endif

ifeq ($(VT_PKG),)
VT_PKG := $(notdir $(CURDIR))
endif

version := $(shell git describe --tags --always HEAD)
version := $(version:v%=%)
name := $(VT_PKG)_$(version)

VT_OUT_DIR ?= $(vtdir)/output/$(name)
prefix := $(VT_OUT_DIR)/$(name)

VT_BIN_DIR ?= $(vtdir)/bin
VT_DOC_DIR ?= docs/commands
VT_DOC_DO_GEN ?= yes
VT_MATRIX ?= docs/validation/matrix.yaml

VT_TEST_ALLOW_SKIPS ?= no
VT_TEST_RUNNERS ?=
ifeq ($(strip $(VT_TEST_RUNNERS)),)
$(error "VT_TEST_RUNNERS must point to space-delimited list of test scripts")
endif

.PHONY: vt-help
vt-help:
	$(info Primary targets:)
	$(info * vt-all: create all validation artifacts under $(VT_OUT_DIR)/)
	$(info * vt-gen-docs: generate command docs under $(VT_DOC_DIR)/)
	$(info )
	$(info Other targets:)
	$(info * vt-cover-unlisted: show Go files that are not in coverage JSON)
	$(info * vt-test: invoke each script listed in VT_TEST_RUNNERS)
	$(info )
	$(info Other targets, triggered by vt-all:)
	$(info * vt-archive: write source archive to $(prefix).tar.gz)
	$(info * vt-bin: install executables for packages under current directory to $(VT_BIN_DIR)/)
	$(info * vt-checkmat: check $(VT_MATRIX) with checkmat)
	$(info * vt-copymat: copy $(VT_MATRIX) to $(prefix).matrix.yaml)
	$(info * vt-cover: invoke each script listed in VT_TEST_RUNNERS with coverage enabled)
	$(info * vt-metadata: write scorecard metadata to $(prefix).metadata.json)
	$(info * vt-pkg: write package name and version to $(prefix).pkg.json)
	$(info * vt-scores: write scorecard scores to $(prefix).scores.json)
	@:

.PHONY: help-valtools
help-valtools: vt-help

.PHONY: vt-all
vt-all: vt-copymat
vt-all: vt-cover
vt-all: vt-scores
vt-all: vt-pkg
vt-all: vt-metadata
vt-all: vt-archive

.PHONY: vt-bin
vt-bin:
	rm -rf '$(VT_BIN_DIR)' && mkdir '$(VT_BIN_DIR)'
	go build -o '$(VT_BIN_DIR)/' ./...

.PHONY: vt-gen-docs
vt-gen-docs:
ifeq ($(VT_DOC_DO_GEN),yes)
	$(MAKE) vt-bin
	@test -f '$(VT_BIN_DIR)/docgen' || \
	  { printf '"make vt-bin" did not generate $(VT_BIN_DIR)/docgen\n'; exit 1; }
	'$(VT_BIN_DIR)/docgen' '$(VT_DOC_DIR)'
else
	@:
endif

$(VT_BIN_DIR)/checkmat: $(vtdir)/checkmat/main.go

.PHONY: vt-checkmat
vt-checkmat: $(VT_BIN_DIR)/checkmat
	'$(VT_BIN_DIR)/checkmat' '$(VT_MATRIX)' '$(VT_DOC_DIR)'

.PHONY: vt-copymat
vt-copymat:
	$(MAKE) vt-gen-docs
	$(MAKE) vt-checkmat
	@test -z "$$(git status -unormal --porcelain -- '$(VT_DOC_DIR)')" || \
	  { printf 'commit changes to $(VT_DOC_DIR) first\n'; exit 1; }
	@mkdir -p '$(VT_OUT_DIR)'
	cp '$(VT_MATRIX)' '$(prefix).matrix.yaml'

$(VT_BIN_DIR)/fmttests: $(vtdir)/fmttests/main.go

.PHONY: vt-test
vt-test: $(VT_BIN_DIR)/fmttests
	@unset GOCOVERDIR; \
	  '$(vtdir)/scripts/run-tests' '$(VT_BIN_DIR)/fmttests' \
	  '$(VT_TEST_ALLOW_SKIPS)' $(VT_TEST_RUNNERS)

$(VT_BIN_DIR)/filecov: $(vtdir)/filecov/main.go

# ATTN: Make coverage directory absolute because we cannot rely on
# test subprocesses to be executed from the same directory.
cov_dir := $(abspath $(VT_OUT_DIR)/.coverage)
cov_prof := $(cov_dir).profile

.PHONY: vt-cover
vt-cover: export GOCOVERDIR=$(cov_dir)
vt-cover: $(VT_BIN_DIR)/filecov
vt-cover: $(VT_BIN_DIR)/fmttests
	@mkdir -p '$(VT_OUT_DIR)'
	rm -rf '$(cov_dir)' && mkdir '$(cov_dir)'
	'$(vtdir)/scripts/run-tests' '$(VT_BIN_DIR)/fmttests' \
	  '$(VT_TEST_ALLOW_SKIPS)' $(VT_TEST_RUNNERS) \
	  >'$(prefix).check.txt'
	go tool covdata textfmt -i '$(cov_dir)' -o '$(cov_prof)'
	'$(VT_BIN_DIR)/filecov' -mod go.mod '$(cov_prof)' \
	  >'$(prefix).coverage.json'

.PHONY: vt-cover-unlisted
vt-cover-unlisted:
	@test -f '$(prefix).coverage.json' || \
	  { printf >&2 'vt-cover-unlisted requires $(prefix).coverage.json\n'; exit 1; }
	@'$(vtdir)/scripts/cover-unlisted' '$(prefix).coverage.json' || :

.PHONY: vt-scores
vt-scores:
	@mkdir -p '$(VT_OUT_DIR)'
	'$(vtdir)/scripts/write-scores' '$(prefix).coverage.json' \
	  >'$(prefix).scores.json'

.PHONY: vt-pkg
vt-pkg:
	@mkdir -p '$(VT_OUT_DIR)'
	jq -n --arg p '$(VT_PKG)' --arg v "$(version)" \
	'{"mpn_scorecard_format": "1.0",'\
	' "pkg_name": $$p, "pkg_version": $$v,'\
	' "scorecard_type": "cli"}' \
	>'$(prefix).pkg.json'

.PHONY: vt-metadata
vt-metadata:
	@mkdir -p '$(VT_OUT_DIR)'
	'$(vtdir)/scripts/metadata' >'$(prefix).metadata.json'

.PHONY: vt-archive
vt-archive:
	@mkdir -p '$(VT_OUT_DIR)'
	@test -z "$(git status --porcelain -unormal --ignore-submodules=none)" || \
	  { printf >&2 'working tree is dirty; commit changes first\n'; exit 1; }
	git archive -o '$(prefix).tar.gz' --format=tar.gz HEAD

$(VT_BIN_DIR)/%: $(vtdir)/%/main.go
	@mkdir -p '$(VT_BIN_DIR)'
	go build -o '$@' '$<'
