BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.buildTime=${BUILD}"
# MAKE_HOME=${PWD}
# TEST_HOME=${MAKE_HOME}/integration_tests
BUILD_HOME=${PWD}/../../cmd/pkgr
TEST_HOME=${PWD}

.PHONY: test-setup test-partially-installed-reset test

install:
	cd ${BUILD_HOME}; go install ${LDFLAGS}

test-setup: test-partially-installed-reset install
	pkgr clean cache

test-partially-installed-reset:
	rm -rf ./test-library
	mkdir ./test-library
	cp -r ./preinstalled-library/* ./test-library/

test: test-setup
	pkgr install
