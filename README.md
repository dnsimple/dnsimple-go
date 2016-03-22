# DNSimple Go Client

A Go client for the [DNSimple API v2](https://developer.dnsimple.com/v2/).

[![Build Status](https://travis-ci.org/aetrion/dnsimple-go.svg)](https://travis-ci.org/aetrion/dnsimple-go)
[![GoDoc](https://godoc.org/github.com/aetrion/dnsimple-go/dnsimple?status.svg)](https://godoc.org/github.com/aetrion/dnsimple-go/dnsimple)


## :warning: Beta Warning

This project targets the development of the API client for the [DNSimple API v2](https://developer.dnsimple.com/v2/). If you are looking for the initial version of the client for [DNSimple API v1](https://developer.dnsimple.com/v1/) then use the [`weppos/dnsimple-go`](https://github.com/weppos/dnsimple-go) project.

This library is currently in beta version, the methods and the implementation should be considered a work-in-progress. Changes in the method naming, method signatures, public or internal APIs may happen during the beta period.


## Installation

```
$ go get github.com/aetrion/dnsimple-go/dnsimple
```


## Usage

This library is a Go client you can use to interact with the [DNSimple API v2](https://developer.dnsimple.com/v2/). Here are some examples.

```go
package main

import (
  "fmt"
  "os"

  "github.com/aetrion/dnsimple-go/dnsimple"
)

func main() {
    oauthToken := "xxxxxxx"

    // new client
    client := dnsimple.NewClient(dnsimple.NewOauthTokenCredentials(oauthToken))

    // get the current authenticated account (if you don't know who you are)
    whoamiResponse, err := client.Identity.Whoami()
    if err != nil {
        fmt.Printf("Whoami() returned error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(whoamiResponse.Data.Account)
    fmt.Println(whoamiResponse.Data.User)

    // get the list of domains
    domainsResponse, err := client.Domains.ListDomains(whoami.Account.Id, nil)
    if err != nil {
        fmt.Printf("Domains.ListDomains() returned error: %v\n", err)
        os.Exit(1)
    }

    // iterate over all the domains in the
    // paginated response.
    for _, domain := range domainsResponse.Data {
        fmt.Println(domain)
    }
}
```

For more complete documentation, see [godoc](https://godoc.org/github.com/aetrion/dnsimple-go/dnsimple).


## Contributing

For instructions about contributing and testing, visit the [CONTRIBUTING](CONTRIBUTING.md) file.


## License

Copyright (c) 2014-2016 Aetrion LLC. This is Free Software distributed under the MIT license.
