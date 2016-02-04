package dnsimple

import ()

// ZonesService handles communication with the zone related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/zones/
type ZonesService struct {
	client *Client
}
