export GO111MODULE=on

test:
	go test -v ./...

lint:
	golint ./...
