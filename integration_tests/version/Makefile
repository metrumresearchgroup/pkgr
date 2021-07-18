TEST_HOME=${PWD}

.PHONY: clean

clean:
	rm -rf test-library
	rm -rf test-cache

update: clean
	go test ./... -update

test: clean
	go test ./... -v

test-json: clean
	go test ./... -json