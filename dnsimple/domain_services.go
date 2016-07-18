package dnsimple

import (
	"fmt"
)

// DomainServicesService handles communication with the domain services
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/services/domains/
type DomainServicesService struct {
	client *Client
}

func domainServicesPath(accountID string, domainID string, serviceID string) string {
	if serviceID != "" {
		return fmt.Sprintf("/%v/domains/%v/services/%v", accountID, domainID, serviceID)
	}
	return fmt.Sprintf("/%v/domains/%v/services", accountID, domainID)
}

// AppliedServices list the applied services for a domain.
//
// See https://developer.dnsimple.com/v2/services/domains/#applied
func (s *DomainServicesService) AppliedServices(accountID string, domainID string, options *ListOptions) (*ServicesResponse, error) {
	path := versioned(domainServicesPath(accountID, domainID, ""))
	servicesResponse := &ServicesResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, servicesResponse)
	if err != nil {
		return servicesResponse, err
	}

	servicesResponse.HttpResponse = resp
	return servicesResponse, nil
}
