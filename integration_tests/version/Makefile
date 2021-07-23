TEST_HOME=${PWD}

.PHONY: clean update test test-json

clean:
	go clean -testcache

update: clean
	go test ./... -update

test: clean
	go test ./... -v

test-json: clean
	go test ./... -json