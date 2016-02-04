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

type domainsWrapper struct {
	Domains []Domain `json:"data"`
}

type domainWrapper struct {
	Domain *Domain `json:"data"`
}

// domainRequest represents a generic wrapper for a domain request,
// when domainWrapper cannot be used because of type constraint on Domain.
type domainRequest struct {
	Domain interface{} `json:"domain"`
}

func domainIDentifier(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return fmt.Sprintf("%d", value)
	}
	return ""
}

// domainPath generates the resource path for given domain.
func domainPath(accountID string, domain interface{}) string {
	if domain != nil {
		return fmt.Sprintf("%s/domains/%s", accountID, domainIDentifier(domain))
	}
	return fmt.Sprintf("%s/domains", accountID)
}

// List the domains.
//
// See https://developer.dnsimple.com/v2/domains/#list
func (s *DomainsService) List(accountID string) ([]Domain, *Response, error) {
	path := domainPath(accountID, nil)
	data := domainsWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return []Domain{}, res, err
	}

	return data.Domains, res, nil
}

// Create a new domain.
//
// See https://developer.dnsimple.com/v2/domains/#create
func (s *DomainsService) Create(accountID string, domainAttributes Domain) (*Domain, *Response, error) {
	path := domainPath(accountID, nil)
	data := domainWrapper{}

	res, err := s.client.post(path, domainAttributes, &data)
	if err != nil {
		return &Domain{}, res, err
	}

	return data.Domain, res, nil
}

// Get a domain.
//
// See https://developer.dnsimple.com/v2/domains/#get
func (s *DomainsService) Get(accountID string, domain interface{}) (*Domain, *Response, error) {
	path := domainPath(accountID, domain)
	data := domainWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return &Domain{}, res, err
	}

	return data.Domain, res, nil
}

// Delete a domain.
//
// See https://developer.dnsimple.com/v2/domains/#delete
func (s *DomainsService) Delete(accountID string, domain interface{}) (*Response, error) {
	path := domainPath(accountID, domain)

	return s.client.delete(path, nil, nil)
}
