package dnsimple

import (
	"fmt"
)

// DomainsService handles communication with the domain related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/domains/
type DomainsService struct {
	client *Client
}

// DomainResponse represents a response from an API method that returns a Domain struct.
type DomainResponse struct {
	Response
	Data *Domain `json:"data"`
}

// DomainsResponse represents a response from an API method that returns a collection of Domain struct.
type DomainsResponse struct {
	Response
	Data []Domain `json:"data"`
}

// Domain represents a Domain in DNSimple.
type Domain struct {
	ID           int    `json:"id,omitempty"`
	AccountID    int    `json:"account_id,omitempty"`
	RegistrantID int    `json:"registrant_id,omitempty"`
	Name         string `json:"name,omitempty"`
	UnicodeName  string `json:"unicode_name,omitempty"`
	Token        string `json:"token,omitempty"`
	State        string `json:"state,omitempty"`
	AutoRenew    bool   `json:"auto_renew,omitempty"`
	PrivateWhois bool   `json:"private_whois,omitempty"`
	ExpiresOn    string `json:"expires_on,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// domainRequest represents a generic wrapper for a domain request,
// when domainWrapper cannot be used because of type constraint on Domain.
type domainRequest struct {
	Domain interface{} `json:"domain"`
}

func domainIdentifier(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return fmt.Sprintf("%d", value)
	}
	return ""
}

func domainPath(accountID string, domain interface{}) string {
	if domain != nil {
		return fmt.Sprintf("/%v/domains/%v", accountID, domainIdentifier(domain))
	}
	return fmt.Sprintf("/%v/domains", accountID)
}

// List the domains.
//
// See https://developer.dnsimple.com/v2/domains/#list
func (s *DomainsService) List(accountID string) (*DomainsResponse, error) {
	path := domainPath(accountID, nil)
	domainsResponse := &DomainsResponse{}

	resp, err := s.client.get(path, domainsResponse)
	if err != nil {
		return nil, err
	}

	domainsResponse.HttpResponse = resp
	return domainsResponse, nil
}

// Create a new domain.
//
// See https://developer.dnsimple.com/v2/domains/#create
func (s *DomainsService) Create(accountID string, domainAttributes Domain) (*DomainResponse, error) {
	path := domainPath(accountID, nil)
	domainResponse := &DomainResponse{}

	resp, err := s.client.post(path, domainAttributes, domainResponse)
	if err != nil {
		return nil, err
	}

	domainResponse.HttpResponse = resp
	return domainResponse, nil
}

// Get a domain.
//
// See https://developer.dnsimple.com/v2/domains/#get
func (s *DomainsService) Get(accountID string, domain interface{}) (*DomainResponse, error) {
	path := domainPath(accountID, domain)
	domainResponse := &DomainResponse{}

	resp, err := s.client.get(path, domainResponse)
	if err != nil {
		return nil, err
	}

	domainResponse.HttpResponse = resp
	return domainResponse, nil
}

// Delete a domain.
//
// See https://developer.dnsimple.com/v2/domains/#delete
func (s *DomainsService) Delete(accountID string, domain interface{}) (*DomainResponse, error) {
	path := domainPath(accountID, domain)
	domainResponse := &DomainResponse{}

	resp, err := s.client.delete(path, nil, nil)
	if err != nil {
		return nil, err
	}

	domainResponse.HttpResponse = resp
	return domainResponse, nil
}
