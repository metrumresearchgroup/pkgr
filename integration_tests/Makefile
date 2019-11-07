BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.buildTime=${BUILD}"
# MAKE_HOME=${PWD}
# TEST_HOME=${MAKE_HOME}/integration_tests
BUILD_HOME=${PWD}/../cmd/pkgr
TEST_HOME=${PWD}

.PHONY: install test-setup test-mixed-source-reset test-simple-reset test-threads-reset test-logging-config-reset test-outdated-pkgs-reset test-repo-order-reset test-repo-local-reset test-repo-local-and-remote-reset test-rollback-reset test-simple-suggests-reset test-rollback-reset test-create-library-reset test-strict-mode-reset 

install:
	cd ${BUILD_HOME}; go install ${LDFLAGS}

test-setup: test-mixed-source-reset test-simple-reset test-threads-reset test-logging-config-reset test-outdated-pkgs-reset test-outdated-pkgs-stress-reset test-repo-order-reset test-repo-local-reset test-repo-local-and-remote-reset test-simple-suggests-reset test-rollback-reset test-create-library-reset test-strict-mode-reset install


test-mixed-source-reset:
	cd ${TEST_HOME}/mixed-source; rm -rf test-library/
	cd ${TEST_HOME}/mixed-source; rm -rf testsets/
	cd ${TEST_HOME}/mixed-source; mkdir test-library
	cd ${TEST_HOME}/mixed-source; mkdir -p testsets/pkgcache

test-simple-reset:
	cd ${TEST_HOME}/simple; rm -rf test-library/
	cd ${TEST_HOME}/simple; mkdir test-library/

test-threads-reset:
	cd ${TEST_HOME}/threads; rm -rf test-library/
	cd ${TEST_HOME}/threads; mkdir test-library/

test-logging-config-reset:
	cd ${TEST_HOME}/logging-config/install-log; rm -rf logs/
	cd ${TEST_HOME}/logging-config/install-log; rm -rf test-library/
	cd ${TEST_HOME}/logging-config/install-log; mkdir logs
	cd ${TEST_HOME}/logging-config/install-log; mkdir test-library
	cd ${TEST_HOME}/logging-config/default; rm -rf logs/
	cd ${TEST_HOME}/logging-config/default; rm -rf test-library/
	cd ${TEST_HOME}/logging-config/default; mkdir logs
	cd ${TEST_HOME}/logging-config/default; mkdir test-library
	cd ${TEST_HOME}/logging-config/overwrite-setting; rm -rf logs/
	cd ${TEST_HOME}/logging-config/overwrite-setting; rm -rf test-library/
	cd ${TEST_HOME}/logging-config/overwrite-setting; mkdir logs
	cd ${TEST_HOME}/logging-config/overwrite-setting; mkdir test-library
	cd ${TEST_HOME}/logging-config/overwrite-setting; echo "This text should be deleted" > logs/all.log
	cd ${TEST_HOME}/logging-config/overwrite-setting; echo "This text should be deleted" > logs/install.log

test-outdated-pkgs-reset:
	rm -rf ${TEST_HOME}/outdated-pkgs/test-library/
	mkdir ${TEST_HOME}/outdated-pkgs/test-library
	cp -r ${TEST_HOME}/outdated-pkgs/outdated-library/* ${TEST_HOME}/outdated-pkgs/test-library/

test-outdated-pkgs-stress-reset:
	rm -rf ${TEST_HOME}/outdated-pkgs-stress/test-library/
	mkdir ${TEST_HOME}/outdated-pkgs-stress/test-library
	cp -r ${TEST_HOME}/outdated-pkgs-stress/outdated-library/* ${TEST_HOME}/outdated-pkgs-stress/test-library/

test-repo-order-reset:
	rm -rf ${TEST_HOME}/repo-order/test-library/
	mkdir ${TEST_HOME}/repo-order/test-library/

test-repo-local-reset:
	rm -rf ${TEST_HOME}/repo-local/test-library/
	mkdir ${TEST_HOME}/repo-local/test-library/

test-repo-local-and-remote-reset:
	rm -rf ${TEST_HOME}/repo-local-and-remote/test-library/
	mkdir ${TEST_HOME}/repo-local-and-remote/test-library/

test-simple-suggests-reset:
	rm -rf ${TEST_HOME}/simple-suggests/test-library/
	rm -rf ${TEST_HOME}/simple-suggests/pkgcache/
	mkdir ${TEST_HOME}/simple-suggests/test-library
	mkdir ${TEST_HOME}/simple-suggests/pkgcache/

test-rollback-reset:
	cd ${TEST_HOME}/rollback; rm -rf test-library/
	cd ${TEST_HOME}/rollback; mkdir test-library/
	cd ${TEST_HOME}/rollback; cp -r ${TEST_HOME}/rollback/preinstalled-library/* ${TEST_HOME}/rollback/test-library

test-create-library-reset:
	cd ${TEST_HOME}/create-library; rm -rf test-library/

test-strict-mode-reset:
	cd ${TEST_HOME}/strict-mode; rm -rf test-library/
