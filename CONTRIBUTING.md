# Contributing to DNSimple/Go

The main `dnsimple` package is defined in the `/dnsimple` subfolder of the `dnsimple/dnsimple-go` repository. Therefore, please note that you will need to move into the subfolder to run any `go` command that assumes the current directory to be the package root.

## Getting started

### 1. Clone the repository

Clone the repository [in your workspace](https://go.dev/doc/code#Organization) and move into it:

```shell
git clone git@github.com:dnsimple/dnsimple-go.git
cd dnsimple-go
```

### 2. Build and test

[Run the test suite](#testing) to check everything works as expected.

## Testing

Submit unit tests for your changes. You can test your changes on your machine by running the test suite (see below).

When you submit a PR, tests will also be run on the [continuous integration environment via GitHub Actions](https://github.com/dnsimple/dnsimple-go/actions).

### Test Suite

To run the test suite:

```shell
go test -v ./...
```

To run the test suite in a live environment (integration):

```shell
export DNSIMPLE_TOKEN="some-token"
go test -v ./...
```

