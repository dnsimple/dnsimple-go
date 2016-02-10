package dnsimple

import (
	"fmt"
)

// ZonesService handles communication with the zone related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/zones/
type ZonesService struct {
	client *Client
}

// ZoneResponse represents a response from an API method that returns a Zone struct.
type ZoneResponse struct {
	Response
	Data *Contact `json:"data"`
}

// ZonesResponse represents a response from an API method that returns a collection of Zone struct.
type ZonesResponse struct {
	Response
	Data []Zone `json:"data"`
}

// Zone represents a Zone in DNSimple.
type Zone struct {
	ID        int    `json:"id,omitempty"`
	AccountID int    `json:"account_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Reverse   bool   `json:"reverse,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// List the zones.
//
// See https://developer.dnsimple.com/v2/zones/#list
func (s *ZonesService) ListZones(accountID string) (*ZonesResponse, error) {
	path := fmt.Sprintf("/%v/zones", accountID)
	zonesResponse := &ZonesResponse{}

	resp, err := s.client.get(path, zonesResponse)
	if err != nil {
		return zonesResponse, err
	}

	zonesResponse.HttpResponse = resp
	return zonesResponse, nil
}
