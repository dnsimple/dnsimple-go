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
    oauthToken := "xxxxxxx"
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oauthToken})
    tc := oauth2.NewClient(context.Background(), ts)

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
Currently supported authentication mechanisms are OAuth and HTTP Digest.

For OAuth we suggest to use the client provided by the `golang.org/x/oauth2` package:

```go
ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "XXX"})
tc := oauth2.NewClient(context.Background(), ts)

// new client
client := dnsimple.NewClient(tc)
```

For HTTP Digest you can use the client provided by this package:

```go
tp := dnsimple.BasicAuthTransport{
    Username: "XXX",
    Password: "XXX",
}

client := dnsimple.NewClient(tp.Client())
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

You can customize the `User-Agent` header for the calls made to the DNSimple API:

```go
client := dnsimple.NewClient(tc)
client.SetUserAgent("my-app/1.0")
```

The value you provide will be prepended to the default `User-Agent` the client uses. For example, if you use `my-app/1.0`, the final header value will be `my-app/1.0 dnsimple-go/0.14.0` (note that it will vary depending on the client version).

We recommend to customize the user agent. If you are building a library or integration on top of the official client, customizing the client will help us to understand what is this client used for, and allow to contribute back or get in touch.


## Contributing

For instructions about contributing and testing, visit the [CONTRIBUTING](CONTRIBUTING.md) file.


## License

Copyright (c) 2014-2020 DNSimple Corporation. This is Free Software distributed under the MIT license.
