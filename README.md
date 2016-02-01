# go-dnsimple

A Go client for the [DNSimple API](https://developer.dnsimple.com/).

[![Build Status](https://travis-ci.org/weppos/dnsimple-go.svg)](https://travis-ci.org/weppos/dnsimple-go)
[![Coverage Status](https://img.shields.io/coveralls/weppos/dnsimple-go.svg)](https://coveralls.io/r/weppos/dnsimple-go?branch=master)
[![GoDoc](https://godoc.org/github.com/weppos/dnsimple-go/dnsimple?status.svg)](https://godoc.org/github.com/weppos/dnsimple-go/dnsimple)

## Installation

```
$ go get github.com/weppos/dnsimple-go/dnsimple
```


## Getting Started

This library is a Go client you can use to interact with the [DNSimple API](https://developer.dnsimple.com/). Here are some examples.


```go
package main

import (
  "fmt"
  "github.com/weppos/dnsimple-go/dnsimple"
)

func main() {
  apiToken := "xxxxxxx"
  email := "foo@example.com"

  client := dnsimple.NewClient(apiToken, email)

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

For more complete documentation, see [godoc](https://godoc.org/github.com/weppos/dnsimple-go/dnsimple).
