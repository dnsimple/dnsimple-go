package dnsimple

import (
	"fmt"
)

// RegistrarService handles communication with the registrar related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/registrar/
type RegistrarService struct {
	client *Client
}

// RegistrationResponse represents a response from an API method that results in a domain registration.
type RegistrationResponse struct {
	Response
	Data *Domain `json:"data"`
}

// Register a domain name.
//
// TODO: ? Switch to a RegistrationOptions struct for the payload.
//
// See https://developer.dnsimple.com/v2/registrar/#register
func (s *RegistrarService) Register(accountID string, domainName string, domainAttributes Domain) (*RegistrationResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/registration", accountID, domainName))
	registrationResponse := &RegistrationResponse{}

	// TODO: validate mandatory attributes RegistrantID

	resp, err := s.client.post(path, domainAttributes, registrationResponse)
	if err != nil {
		return nil, err
	}

	registrationResponse.HttpResponse = resp
	return registrationResponse, nil
}


// RenewOptions represents the option you can pass to a renew API request.
type RenewOptions struct {
	// The number of years
	Period int `json:"period"`
}

// RenewalResponse represents a response from an API method that results in a domain renewal.
type RenewalResponse struct {
	Response
	Data *Domain `json:"data"`
}

// Renew a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#register
func (s *RegistrarService) Renew(accountID string, domainName string, options *RenewOptions) (*RenewalResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/renewal", accountID, domainName))
	renewalResponse := &RenewalResponse{}

	resp, err := s.client.post(path, options, renewalResponse)
	if err != nil {
		return nil, err
	}

	renewalResponse.HttpResponse = resp
	return renewalResponse, nil
}
