package dnsimple

import (
//"fmt"
)

// TldsService handles communication with the Tld related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/tlds/
type TldsService struct {
	client *Client
}

// Tld represents a TLD in DNSimple.
type Tld struct {
	Tld           string `json:"tld"`
	TldType       int    `json:"tld_type"`
	WhoisPrivacy  bool   `json:"whois_privacy"`
	AutoRenewOnly bool   `json:"auto_renew_only"`
}

// TldResponse represents a response from an API method that returns a Tld struct.
type TldResponse struct {
	Response
	Data *Tld `json:"data"`
}

// TldsResponse represents a response from an API method that returns a collection of Tld struct.
type TldsResponse struct {
	Response
	Data []Tld `json:"data"`
}

// ListTlds lists the supported TLDs.
//
// See https://developer.dnsimple.com/v2/tlds/#list
func (s *TldsService) ListTlds() (*TldsResponse, error) {
	path := versioned("/tlds")
	tldsResponse := &TldsResponse{}

	resp, err := s.client.get(path, tldsResponse)
	if err != nil {
		return tldsResponse, err
	}

	tldsResponse.HttpResponse = resp
	return tldsResponse, nil
}
