# DNSimple Go Client

A Go client for the [DNSimple API v2](https://developer.dnsimple.com/v2/).

[![Build Status](https://github.com/dnsimple/dnsimple-go/actions/workflows/ci.yml/badge.svg)](https://github.com/dnsimple/dnsimple-go/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/dnsimple/dnsimple-go/dnsimple?status.svg)](https://godoc.org/github.com/dnsimple/dnsimple-go/dnsimple)

## Requirements

- Go 1.21+
- An activated DNSimple account

## Installation

```shell
go get github.com/dnsimple/dnsimple-go/v8/dnsimple
```

## Usage

This library is a Go client you can use to interact with the [DNSimple API v2](https://developer.dnsimple.com/v2/). Here are some examples.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strconv"

    "github.com/dnsimple/dnsimple-go/v8/dnsimple"
)

func main() {
    tc := dnsimple.StaticTokenHTTPClient(context.Background(), "your-token")

    // new client
    client := dnsimple.NewClient(tc)

    // get the current authenticated account (if you don't know who you are)
    whoamiResponse, err := client.Identity.Whoami(context.Background())
    if err != nil {
        fmt.Printf("Whoami() returned error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(whoamiResponse.Data.Account)
    fmt.Println(whoamiResponse.Data.User)

    // either assign the account ID or fetch it from the response
    // if you are authenticated with an account token
    accountID := strconv.FormatInt(whoamiResponse.Data.Account.ID, 10)

    // get the list of domains
    domainsResponse, err := client.Domains.ListDomains(context.Background(), accountID, nil)
    if err != nil {
        fmt.Printf("Domains.ListDomains() returned error: %v\n", err)
        os.Exit(1)
    }

    // iterate over all the domains in the
    // paginated response.
    for _, domain := range domainsResponse.Data {
        fmt.Println(domain)
    }

    // List methods support a variety of options to paginate, sort and filter records.
    // Here's a few example:

    // get the list of domains filtered by name and sorted by expiration
    client.Domains.ListDomains(context.Background(), accountID, &dnsimple.DomainListOptions{NameLike: dnsimple.String("com"), ListOptions: {Sort: dnsimple.String("expiration:DESC")}})
}
```

For more complete documentation, see [godoc](https://godoc.org/github.com/dnsimple/dnsimple-go/dnsimple).

## Configuration

### Sandbox Environment

We highly recommend testing against our [sandbox environment](https://developer.dnsimple.com/sandbox/) before using our production environment. This will allow you to avoid real purchases, live charges on your credit card, and reduce the chance of your running up against rate limits.

The client supports both the production and sandbox environment. To switch to sandbox pass the sandbox API host using the `base_url` option when you construct the client:

```go
client := dnsimple.NewClient(tc)
client.BaseURL = "https://api.sandbox.dnsimple.com"
```

You will need to ensure that you are using an access token created in the sandbox environment. Production tokens will *not* work in the sandbox environment.

### Setting a custom `User-Agent` header

You can customize the `User-Agent` header for the calls made to the DNSimple API:

```go
client := dnsimple.NewClient(tc)
client.SetUserAgent("my-app/1.0")
```

The value you provide will be prepended to the default `User-Agent` the client uses. For example, if you use `my-app/1.0`, the final header value will be `my-app/1.0 dnsimple-go/0.14.0` (note that it will vary depending on the client version).

We recommend to customize the user agent. If you are building a library or integration on top of the official client, customizing the client will help us to understand what is this client used for, and allow to contribute back or get in touch.

## Authentication

When creating a new client you are required to provide an `http.Client` to use for authenticating the requests.
Supported authentication mechanisms are OAuth and HTTP Digest. We provide convenient helpers to generate a preconfigured HTTP client.

### Authenticating with OAuth

```go
tc := dnsimple.StaticTokenHTTPClient(context.Background(), "your-token")

// new client
client := dnsimple.NewClient(tc)
```

### Authenticating with HTTP Basic Auth

```go
hc := dnsimple.BasicAuthHTTPClient(context.Background(), "your-user", "your-password")
client := dnsimple.NewClient(hc)
```

For requests made to authorize OAuth access, and to exchange the short lived authorization token for the OAuth token, use an HTTP client with a timeout:

```go
client := dnsimple.NewClient(&http.Client{Timeout: time.Second * 10})
```

For any other custom need you can define your own `http.RoundTripper` implementation and
pass a client that authenticated with the custom round tripper.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for details.

## License

Copyright (c) 2014-2026 DNSimple Corporation. This is Free Software distributed under the [MIT License](LICENSE.txt).
