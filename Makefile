lint:
	@golangci-lint run

test: lint
	@go test -v ./...

test_all: test
	@codecrafters test
