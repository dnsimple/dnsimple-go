package dnsimple

import (
	"context"
	"fmt"
)

// EnableDomainAutoRenewal enables auto-renewal for the domain.
//
// See https://developer.dnsimple.com/v2/registrar/auto-renewal/#enable
func (s *RegistrarService) EnableDomainAutoRenewal(accountID string, domainName string) (*domainResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/auto_renewal", accountID, domainName))
	domainResponse := &domainResponse{}

	resp, err := s.client.put(context.TODO(), path, nil, nil)
	if err != nil {
		return nil, err
	}

	domainResponse.HTTPResponse = resp
	return domainResponse, nil
}

// DisableDomainAutoRenewal disables auto-renewal for the domain.
//
// See https://developer.dnsimple.com/v2/registrar/auto-renewal/#enable
func (s *RegistrarService) DisableDomainAutoRenewal(accountID string, domainName string) (*domainResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/auto_renewal", accountID, domainName))
	domainResponse := &domainResponse{}

	resp, err := s.client.delete(context.TODO(), path, nil, nil)
	if err != nil {
		return nil, err
	}

	domainResponse.HTTPResponse = resp
	return domainResponse, nil
}
