vtdir := $(dir $(lastword $(MAKEFILE_LIST)))
vtdir := $(patsubst %/,%,$(vtdir))
ifeq ($(strip $(vtdir)),)
$(error "bug: vtdir is unexpectedly empty")
endif

VT_BIN_DIR ?= $(vtdir)/bin
VT_DOC_DIR ?= docs/commands

.PHONY: vt-bin
vt-bin:
	rm -rf '$(VT_BIN_DIR)' && mkdir '$(VT_BIN_DIR)'
	go build -o '$(VT_BIN_DIR)/' ./...

.PHONY: vt-gen-docs
vt-gen-docs:
	$(MAKE) vt-bin
	@test -f '$(VT_BIN_DIR)/docgen' || \
	  { printf '"make vt-bin" did not generate $(VT_BIN_DIR)/docgen\n'; exit 1; }
	'$(VT_BIN_DIR)/docgen' '$(VT_DOC_DIR)'
