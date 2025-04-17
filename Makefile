all: test

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: lint
lint:
	golangci-lint run
