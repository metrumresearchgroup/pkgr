TEST_HOME=${PWD}

.PHONY: clean setup update test test-json

clean:
	rm -rf test-library
	rm -rf test-cache
	go clean -testcache

setup:
	cp -r outdated-library test-library

update: clean setup
	go test ./... -update

test: clean setup
	go test ./... -v

test-json: clean setup
	go test ./... -json