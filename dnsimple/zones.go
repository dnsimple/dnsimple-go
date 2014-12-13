package dnsimple

import (
	"bytes"
	"fmt"
)

// ZonesService handles communication with the zone related
// methods of the DNSimple API.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/zones/
type ZonesService struct {
	client *Client
}

// Get downloads the Bind-like zone file.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/zones/#export
func (s *ZonesService) Get(domain interface{}) (string, *Response, error) {
	var body bytes.Buffer
	path := fmt.Sprintf("%s/zone", domainPath(domain))

	res, err := s.client.get(path, &body)
	if err != nil {
		return "", res, err
	}

	return body.String(), res, nil
}
