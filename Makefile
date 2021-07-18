BUILD=`date +%FT%T%z`
MAKE_HOME=${PWD}
TEST_HOME=${MAKE_HOME}/integration_tests

.PHONY: all ci clean install simple test-multiple log-test log-test-reset test-master test-master-reset test-mixed test-mixed-reset

install:
	cd cmd/pkgr; go get; go install

build:
	cd cmd/pkgr; goreleaser build --snapshot --rm-dist --single-target
all: install test-simple

ci: install test-mixed

simple: install test-simple

clean: test-master-reset test-mixed-reset test-simple-reset log-test-reset outdated-test-reset
#	osx+go behave badly and won't remove 0444 files. Chown them first
	chmod -R 700 go
	rm -Rf go

test-multiple: install
	cd ${TEST_HOME}

	rm -rf mixed-source/test-library/*
	rm -rf simple/test-library/*
	rm -rf simple-suggests/test-library/*

	cd ${TEST_HOME}/mixed-source; pkgr install
	cd ${TEST_HOME}/simple;	pkgr install

test-master: install test-master-reset
	cd ${TEST_HOME}/master; pkgr install

test-master-reset:
	cd ${TEST_HOME}; rm -rf master/test-library/*

test-mixed: install test-mixed-reset
	cd ${TEST_HOME}/mixed-source; pkgr install

test-mixed-reset:
	cd ${TEST_HOME}; rm -rf mixed-source/test-library/*

test-simple: install test-simple-reset
	cd ${TEST_HOME}/simple;	pkgr install

test-simple-reset:
	cd ${TEST_HOME};rm -rf simple/test-library/*


log-test: install log-test-reset
	cd ${TEST_HOME}/logging-config/install-log; pkgr install
	cd ${TEST_HOME}/logging-config/default; pkgr install
	cd ${TEST_HOME}/logging-config/overwrite-setting; pkgr install
	cd ${TEST_HOME}/logging-config/overwrite-setting; pkgr clean --all

log-test-reset:
	mkdir -p ${TEST_HOME}/logging-config/overwrite-setting/logs
	cd ${TEST_HOME}/logging-config/install-log; rm -rf logs/*
	cd ${TEST_HOME}/logging-config/default; rm -rf logs/*
	cd ${TEST_HOME}/logging-config/overwrite-setting; rm -rf logs/*
	cd ${TEST_HOME}/logging-config/overwrite-setting; echo "This text should be deleted" > logs/all.log
	cd ${TEST_HOME}/logging-config/overwrite-setting; echo "This text should be deleted" > logs/install.log

outdated-test-reset:
	rm -rf ${TEST_HOME}/outdated-pkgs/test-library/*
	mkdir -p ${TEST_HOME}/outdated-pkgs/test-libary
	cp -r ${TEST_HOME}/outdated-pkgs/outdated-library/* ${TEST_HOME}/outdated-pkgs/test-library/

outdated-test: install outdated-test-reset
	cd ${TEST_HOME}/outdated-pkgs; pkgr plan
	cd ${TEST_HOME}/outdated-pkgs; pkgr install
