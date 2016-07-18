package dnsimple

import (
	"fmt"
)

// DomainServicesService handles communication with the domain one-click services
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

// DomainServiceSettings represents optional settings when applying a DNSimple one-click service to a domain.
type DomainServiceSettings struct {
	Settings map[string]string `url:"settings,omitempty"`
}

// AppliedServices list the applied one-click services for a domain.
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

// ApplyService apply a one-click services to a domain.
//
// See https://developer.dnsimple.com/v2/services/domains/#apply
func (s *DomainServicesService) ApplyService(accountID string, domainID string, serviceID string, settings DomainServiceSettings) (*ServiceResponse, error) {
	path := versioned(domainServicesPath(accountID, domainID, serviceID))
	serviceResponse := &ServiceResponse{}

	resp, err := s.client.post(path, settings, nil)
	if err != nil {
		return nil, err
	}

	serviceResponse.HttpResponse = resp
	return serviceResponse, nil
}

// UnapplyService unapply a one-click services from a domain.
//
// See https://developer.dnsimple.com/v2/services/domains/#unapply
func (s *DomainServicesService) UnapplyService(accountID string, domainID string, serviceID string) (*ServiceResponse, error) {
	path := versioned(domainServicesPath(accountID, domainID, serviceID))
	serviceResponse := &ServiceResponse{}

	resp, err := s.client.delete(path, nil, nil)
	if err != nil {
		return nil, err
	}

	serviceResponse.HttpResponse = resp
	return serviceResponse, nil
}
