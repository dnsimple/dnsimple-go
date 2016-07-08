package dnsimple

import (
	"fmt"
)

// ServicesService handles communication with the service related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/services/
type ServicesService struct {
	client *Client
}

// ServiceSetting represents a single group of settings for a DNSimple Service.
type ServiceSetting struct {
	Name        string `json:"name,omitempty"`
	Label       string `json:"label,omitempty"`
	Append      string `json:"append,omitempty"`
	Description string `json:"description,omitempty"`
	Example     string `json:"example,omitempty"`
	Password    bool   `json:"password,omitempty"`
}

// Service represents a Service in DNSimple.
type Service struct {
	ID               int              `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	ShortName        string           `json:"short_name,omitempty"`
	Description      string           `json:"description,omitempty"`
	SetupDescription string           `json:"setup_description,omitempty"`
	RequiresSetup    bool             `json:"requires_setup,omitempty"`
	DefaultSubdomain string           `json:"default_subdomain,omitempty"`
	CreatedAt        string           `json:"created_at,omitempty"`
	UpdatedAt        string           `json:"updated_at,omitempty"`
	Settings         []ServiceSetting `json:"settings,omitempty"`
}

func servicePath(service interface{}) string {
	if service != nil {
		return fmt.Sprintf("/services/%v", service)
	}
	return fmt.Sprintf("/services")
}

// ServiceResponse represents a response from an API method that returns a Service struct.
type ServiceResponse struct {
	Response
	Data *Service `json:"data"`
}

// ServicesResponse represents a response from an API method that returns a collection of Service struct.
type ServicesResponse struct {
	Response
	Data []Service `json:"data"`
}

// ListServices list the services for an account.
//
// See https://developer.dnsimple.com/v2/services/#list
func (s *ServicesService) ListServices(options *ListOptions) (*ServicesResponse, error) {
	path := versioned(servicePath(nil))
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

// GetService fetches a service.
//
// See https://developer.dnsimple.com/v2/services/#get
func (s *ServicesService) GetService(serviceID string) (*ServiceResponse, error) {
	path := versioned(servicePath(serviceID))
	serviceResponse := &ServiceResponse{}

	resp, err := s.client.get(path, serviceResponse)
	if err != nil {
		return nil, err
	}

	serviceResponse.HttpResponse = resp
	return serviceResponse, nil
}
