## :exclamation::exclamation::exclamation: Development Warning :exclamation::exclamation::exclamation:

This project targets the development of the API client for the [DNSimple API v2](https://developer.dnsimple.com/v2/). If you are looking for the initial version of the client for [DNSimple API v1](https://developer.dnsimple.com/v1/) then use the [`weppos/go-dnsimple`](https://github.com/weppos/go-dnsimple) project.

This version is currently under development, therefore the methods and the implementation should he considered a work-in-progress. Changes in the method naming, method signatures, public or internal APIs may happen at any time.

The code is tested with an automated test suite connected to a continuous integration tool, therefore you should not expect :bomb: bugs to be merged into master. Regardless, use this library at your own risk. :boom:


# DNSimple Go Client

A Go client for the [DNSimple API v2](https://developer.dnsimple.com/v2/).

[![Build Status](https://travis-ci.org/aetrion/dnsimple-go.svg)](https://travis-ci.org/aetrion/dnsimple-go)
[![Coverage Status](https://img.shields.io/coveralls/aetrion/dnsimple-go.svg)](https://coveralls.io/r/aetrion/dnsimple-go?branch=master)
[![GoDoc](https://godoc.org/github.com/aetrion/dnsimple-go/dnsimple?status.svg)](https://godoc.org/github.com/aetrion/dnsimple-go/dnsimple)

## Installation

```
$ go get github.com/aetrion/dnsimple-go/dnsimple
```


## Getting Started

This library is a Go client you can use to interact with the [DNSimple API v2](https://developer.dnsimple.com/v2/). Here are some examples.

```go
package main

import (
  "fmt"
  "github.com/aetrion/dnsimple-go/dnsimple"
)

func main() {
  oauthToken := "xxxxxxx"
  email := "foo@example.com"

  client := dnsimple.NewClient(dnsimple.NewOauthTokenCredentials(oauthToken))

  // Get a list of your domains
  domains, _, _ := client.Domains.List()
  for _, domain := range domains {
      fmt.Printf("Domain: %s (id: %d)\n", domain.Name, domain.Id)
  }

  // Get a list of your domains (with error management)
  domains, _, error := client.Domains.List()
  if error != nil {
      log.Fatalln(error)
  }
  for _, domain := range domains {
      fmt.Printf("Domain: %s (id: %d)\n", domain.Name, domain.Id)
  }

  // Create a new Domain
  newDomain := Domain{Name: "example.com"}
  domain, _, _ := client.Domains.Create(newDomain)
  fmt.Printf("Domain: %s\n (id: %d)", domain.Name, domain.Id)
}
```

For more complete documentation, see [godoc](https://godoc.org/github.com/aetrion/dnsimple-go/dnsimple).


## License

Copyright (c) 2014-2016 Aetrion LLC. This is Free Software distributed under the MIT license.
