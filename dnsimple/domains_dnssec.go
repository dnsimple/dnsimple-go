package dnsimple

import (
	"context"
	"fmt"
)

// Dnssec represents the current DNSSEC settings for a domain in DNSimple.
type Dnssec struct {
	Enabled bool `json:"enabled"`
}

func dnssecPath(accountID string, domainIdentifier string) (path string) {
	path = fmt.Sprintf("%v/dnssec", domainPath(accountID, domainIdentifier))
	return
}

// dnssecResponse represents a response from an API method that returns a Dnssec struct.
type dnssecResponse struct {
	Response
	Data *Dnssec `json:"data"`
}

// EnableDnssec enables DNSSEC on the domain.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#enableDomainDnssec
func (s *DomainsService) EnableDnssec(accountID string, domainIdentifier string) (*dnssecResponse, error) {
	path := versioned(dnssecPath(accountID, domainIdentifier))
	dnssecResponse := &dnssecResponse{}

	resp, err := s.client.post(context.TODO(), path, dnssecResponse, nil)
	if err != nil {
		return nil, err
	}

	dnssecResponse.HTTPResponse = resp
	return dnssecResponse, nil
}

// DisableDnssec disables DNSSEC on the domain.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#disableDomainDnssec
func (s *DomainsService) DisableDnssec(accountID string, domainIdentifier string) (*dnssecResponse, error) {
	path := versioned(dnssecPath(accountID, domainIdentifier))
	dnssecResponse := &dnssecResponse{}

	resp, err := s.client.delete(context.TODO(), path, dnssecResponse, nil)
	if err != nil {
		return nil, err
	}

	dnssecResponse.HTTPResponse = resp
	return dnssecResponse, nil
}

// GetDnssec retrieves the current status of DNSSEC on the domain.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#getDomainDnssec
func (s *DomainsService) GetDnssec(accountID string, domainIdentifier string) (*dnssecResponse, error) {
	path := versioned(dnssecPath(accountID, domainIdentifier))
	dnssecResponse := &dnssecResponse{}

	resp, err := s.client.get(context.TODO(), path, dnssecResponse)
	if err != nil {
		return nil, err
	}

	dnssecResponse.HTTPResponse = resp
	return dnssecResponse, nil
}
