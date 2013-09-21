# go-dnsimple

A Go library for using the DNSimple API.

## Examples

```go
package main
import (
  "fmt"
  "github.com/rubyist/g-dnsimple"
)

func main() {
  apiToken := "xxxxxxx"
  email := "foo@example.com"

  client := dnsimple.NewClient(apiToken, email)

  // Get a list of your domains
  domains, _ := client.Domains()
  for _, domain := range domains {
    fmt.Printf("Domain: %s\n", domain.Name)
  }

  // Get a list of records for a domain
  records, _ := client.Records("example.com")
  for _, record := range records {
    fmt.Printf("Record: %s -> %s\n", record.Name, record.Content)
  }

  // Create a new Record
  newRec := Record{Name: "www", Content: "127.0.0.1", RecordType: "A"}
  rec, _ := client.CreateRecord("example.com")

  // Update a Record
  rec, _ = rec.Update(client, Record{Content: "192.168.0.1"})

  // Convnience method for updating a record's IP
  rec.UpdateIP(client, "10.0.0.1")
}
```

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
