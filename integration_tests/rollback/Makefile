TEST_HOME=${PWD}

.PHONY: clean update test test-json

clean:
	rm -rf test-library
	rm -rf test-cache
	rm -rf initial-cache
	go clean -testcache

update: clean
	go test ./... -update

test: clean
	go test ./... -v

test-json: clean
	go test ./... -json