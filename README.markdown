# go-dnsimple

A Go client for the [DNSimple API](http://developer.dnsimple.com/).

[![Build Status](https://travis-ci.org/weppos/go-dnsimple.svg)](https://travis-ci.org/weppos/go-dnsimple)
[![Coverage Status](https://img.shields.io/coveralls/weppos/go-dnsimple.svg)](https://coveralls.io/r/weppos/go-dnsimple?branch=master)



## Installation

```
$ go get github.com/weppos/go-dnsimple
```


## Getting Started

This library is a Go client you can use to interact with the [DNSimple API](http://developer.dnsimple.com/).

Here's a short example.


```go
package main

import (
  "fmt"
  "github.com/weppos/go-dnsimple/dnsimple"
)

func main() {
  apiToken := "xxxxxxx"
  email := "foo@example.com"

  client := dnsimple.NewClient(apiToken, email)

  // Get a list of your domains
  domains, error := client.Domains.List()
  for _, domain := range domains {
      fmt.Printf("Domain: %s (id: %d)\n", domain.Name, domain.Id)
  }

  // Get a list of your domains (with error management)
  domains, error := client.Domains.List()
  if error != nil {
      log.Fatalln(error)
  }
  for _, domain := range domains {
      fmt.Printf("Domain: %s (id: %d)\n", domain.Name, domain.Id)
  }

  // Create a new Domain
  newDomain := Domain{Name: "example.com"}
  domain, _ := client.Domains.Create(newDomain)
  fmt.Printf("Domain: %s\n (id: %d)", domain.Name, domain.Id)
}
```

For more complete documentation, load up godoc and find the package.


## License

(The MIT License)

Copyright (c) 2013 Scott Barron

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
