TEST_HOME=${PWD}

.PHONY: clean

clean:
	go clean -testcache

update: clean
	go test ./... -update

test: clean
	go test ./... -v

test-json: clean
	go test ./... -json