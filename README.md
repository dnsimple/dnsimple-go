# DNSimple Go Client

A Go client for the [DNSimple API v2](https://developer.dnsimple.com/v2/).

[![Build Status](https://travis-ci.com/dnsimple/dnsimple-go.svg?branch=master)](https://travis-ci.com/dnsimple/dnsimple-go)
[![GoDoc](https://godoc.org/github.com/dnsimple/dnsimple-go/dnsimple?status.svg)](https://godoc.org/github.com/dnsimple/dnsimple-go/dnsimple)


## Installation

```
go get github.com/dnsimple/dnsimple-go/dnsimple
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

    "github.com/dnsimple/dnsimple-go/dnsimple"
    "golang.org/x/oauth2"
)

func main() {
    tc := dnsimple.StaticTokenClient(context.Background(), "your-token")

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


## Authentication

When creating a new client you are required to provide an `http.Client` to use for authenticating the requests.
Supported authentication mechanisms are OAuth and HTTP Digest. We provide convenient helpers to generate a preconfigured HTTP client.

**Authenticating with OAuth**

```go
tc := dnsimple.StaticTokenClient(context.Background(), "your-token")

// new client
client := dnsimple.NewClient(tc)
```

**Authenticating with HTTP Basic Auth**

```go
hc := dnsimple.BasicAuthClient(context.Background(), "your-user", "your-password")
client := dnsimple.NewClient(hc)
```

For requests made to authorize OAuth access, and to exchange the short lived authorization token for the OAuth token, use an HTTP client with a timeout:

```go
client := dnsimple.NewClient(&http.Client{Timeout: time.Second * 10})
```

For any other custom need you can define your own `http.RoundTripper` implementation and
pass a client that authenticated with the custom round tripper.


## Sandbox Environment

We highly recommend testing against our [sandbox environment](https://developer.dnsimple.com/sandbox/) before using our production environment. This will allow you to avoid real purchases, live charges on your credit card, and reduce the chance of your running up against rate limits.

The client supports both the production and sandbox environment. To switch to sandbox pass the sandbox API host using the `base_url` option when you construct the client:

```go
client := dnsimple.NewClient(tc)
client.BaseURL = "https://api.sandbox.dnsimple.com"
```

You will need to ensure that you are using an access token created in the sandbox environment. Production tokens will *not* work in the sandbox environment.


## Setting a custom `User-Agent` header

You customize the `User-Agent` header for the calls made to the DNSimple API:

```go
client := dnsimple.NewClient(tc)
client.SetUserAgent("my-app")
```

The value you provide will be appended to the default `User-Agent` the client uses. For example, if you use `my-app`, the final header value will be `dnsimple-go/0.14.0 my-app` (note that it will vary depending on the client version).


## Contributing

For instructions about contributing and testing, visit the [CONTRIBUTING](CONTRIBUTING.md) file.


## License

Copyright (c) 2014-2020 DNSimple Corporation. This is Free Software distributed under the MIT license.
