package dnsimple

import (
	"fmt"
)

// ZonesService handles communication with the zone related
// methods of the DNSimple API.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/zones/
type ZonesService struct {
	client *Client
}

type zoneResponse struct {
	Zone string `json:"zone,omitempty"`
}

// Get downloads the Bind-like zone file.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/zones/#get
func (s *ZonesService) Get(domain interface{}) (string, *Response, error) {
	path := fmt.Sprintf("%s/zone", domainPath(domain))
	zoneResponse := zoneResponse{}

	res, err := s.client.get(path, &zoneResponse)
	if err != nil {
		return "", res, err
	}

	return zoneResponse.Zone, res, nil
}
