BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.buildTime=${BUILD}"

.PHONY: install

install:
	cd cmd/pkgr; go install ${LDFLAGS}