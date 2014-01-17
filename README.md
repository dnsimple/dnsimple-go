# go-dnsimple

A Go library for using the DNSimple API.

**Build Status:** [![Build Status](https://travis-ci.org/rubyist/go-dnsimple.png?branch=master)](https://travis-ci.org/rubyist/go-dnsimple)  
**Test Coverage:** [![Coverage Status](https://coveralls.io/repos/rubyist/go-dnsimple/badge.png)](https://coveralls.io/r/rubyist/go-dnsimple)
## Examples

```go
package main
import (
  "fmt"
  dnsimple "github.com/rubyist/go-dnsimple"
)

func main() {
  apiToken := "xxxxxxx"
  email := "foo@example.com"

  client := dnsimple.NewClient(apiToken, email)

  // Get a list of your domains
  domains, _ := client.Domains.List()
  for _, domain := range domains {
    fmt.Printf("Domain: %s\n", domain.Name)
  }

  // Create a new Domain
  newDomain := Domain{Name: "example.com"}
  domain, _ := client.Domains.Create(newDomain)

  // Get a list of records for a domain
  records, _ := client.Records.List("example.com")
  for _, record := range records {
    fmt.Printf("Record: %s -> %s\n", record.Name, record.Content)
  }

  // Create a new Record
  newRec := Record{Name: "www", Content: "127.0.0.1", RecordType: "A"}
  rec, _ := client.Records.Create("example.com")

  // Update a Record
  rec, _ = rec.Update(client, Record{Content: "192.168.0.1"})

  // Convenience method for updating a Record's IP
  rec.UpdateIP(client, "10.0.0.1")
}
```

For more complete documentation, load up godoc and find the package.

## Development

- Source hosted at [GitHub](https://github.com/rubyist/go-dnsimple)
- Report issues and feature requests to [GitHub Issues](https://github.com/rubyist/go-dnsimple/issues)

Pull requests welcome!

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
