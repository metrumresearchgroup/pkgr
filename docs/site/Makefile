-include local.mk

HTTP_PORT ?= 1313

version := $(shell git describe --tags --always HEAD)
export HUGO_PARAMS_PKGRVERSION ?= $(version)

.PHONY: help
help:
	$(info Primary targets:)
	$(info * build: build site under public/)
	$(info * serve: build and serve the site locally with 'hugo server')
	@:

.PHONY: build
build: prep
	hugo --gc --minify

.PHONY: serve
serve: prep
	hugo server -p '$(HTTP_PORT)'

copied_files  = content/_index.md
copied_files += content/docs/config.md
copied_files += content/news.md
copied_files += static/windows.png

.PHONY: prep
prep: commands
prep: $(copied_files)
	rm -rf public resources

command_dir = content/docs/commands

.PHONY: commands
commands:
	rm -rf '$(command_dir)'
	./scripts/ingest-command-docs '$(command_dir)'

content/_index.md: ../../README.md
	sed 's|docs/images/windows_tarball_picture\.png|windows.png|' '$<' >'$@'

content/docs/config.md: ../config.md
	printf -- '---\ntitle: "Configuration"\nweight: 10\n---\n\n' >'$@'
	cat '$<' >>'$@'

content/news.md: ../../NEWS.md
	printf -- '---\ntitle: "News"\n---\n\n' >'$@'
	sed 's/^# pkgr/##/' '$<' \
	  >>'$@'

static/windows.png: ../images/windows_tarball_picture.png
	cp '$<' '$@'

.PHONY: clean
clean:
	rm -rf public resources $(command_dir)
	rm -f $(copied_files)
