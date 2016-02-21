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

// RegisterRequest represents the attributes you can pass to a register API request.
// Some attributes are mandatory.
type RegisterRequest struct {
	RegistrantID int `json:"registrant_id"`
}

// RegistrationResponse represents a response from an API method that results in a domain registration.
type RegistrationResponse struct {
	Response
	Data *Domain `json:"data"`
}

// Register a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#register
func (s *RegistrarService) Register(accountID string, domainName string, request *RegisterRequest) (*RegistrationResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/registration", accountID, domainName))
	registrationResponse := &RegistrationResponse{}

	// TODO: validate mandatory attributes RegistrantID

	resp, err := s.client.post(path, request, registrationResponse)
	if err != nil {
		return nil, err
	}

	registrationResponse.HttpResponse = resp
	return registrationResponse, nil
}

// TransferRequest represents the attributes you can pass to a transfer API request.
// Some attributes are mandatory.
type TransferRequest struct {
	RegistrantID int    `json:"registrant_id"`
	AuthInfo     string `json:"auth_info,omitempty"`
}

// TransferResponse represents a response from an API method that results in a domain transfer.
type TransferResponse struct {
	Response
	Data *Domain `json:"data"`
}

// Transfer a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#transfer
func (s *RegistrarService) Transfer(accountID string, domainName string, request *TransferRequest) (*TransferResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/transfer", accountID, domainName))
	transferResponse := &TransferResponse{}

	// TODO: validate mandatory attributes RegistrantID

	resp, err := s.client.post(path, request, transferResponse)
	if err != nil {
		return nil, err
	}

	transferResponse.HttpResponse = resp
	return transferResponse, nil
}

// TransferOutResponse represents a response from an API method that results in a domain transfer out.
type TransferOutResponse struct {
	Response
	Data *Domain `json:"data"`
}

// Transfer out a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#transfer-out
func (s *RegistrarService) TransferOut(accountID string, domainName string) (*TransferOutResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/transfer_out", accountID, domainName))
	transferResponse := &TransferOutResponse{}

	resp, err := s.client.post(path, nil, nil)
	if err != nil {
		return nil, err
	}

	transferResponse.HttpResponse = resp
	return transferResponse, nil
}

// RenewRequest represents the attributes you can pass to a renew API request.
// Some attributes are mandatory.
type RenewRequest struct {
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
func (s *RegistrarService) Renew(accountID string, domainName string, request *RenewRequest) (*RenewalResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/renewal", accountID, domainName))
	renewalResponse := &RenewalResponse{}

	resp, err := s.client.post(path, request, renewalResponse)
	if err != nil {
		return nil, err
	}

	renewalResponse.HttpResponse = resp
	return renewalResponse, nil
}
