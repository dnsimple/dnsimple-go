package dnsimple

import (
	"errors"
	"fmt"
)

// DomainsService handles communication with the domain related
// methods of the DNSimple API.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/
type DomainsService struct {
	client *Client
}

type Domain struct {
	Id             int    `json:"id,omitempty"`
	UserId         int    `json:"user_id,omitempty"`
	RegistrantId   int    `json:"registrant_id,omitempty"`
	Name           string `json:"name,omitempty"`
	UnicodeName    string `json:"unicode_name,omitempty"`
	Token          string `json:"token,omitempty"`
	State          string `json:"state,omitempty"`
	Language       string `json:"language,omitempty"`
	Lockable       bool   `json:"lockable,omitempty"`
	AutoRenew      bool   `json:"auto_renew,omitempty"`
	WhoisProtected bool   `json:"whois_protected,omitempty"`
	RecordCount    int    `json:"record_count,omitempty"`
	ServiceCount   int    `json:"service_count,omitempty"`
	ExpiresOn      string `json:"expires_on,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`

	RenewWhoisPrivacy bool `json:"renew_whois_privacy,omitempty"`
}

type domainWrapper struct {
	Domain Domain `json:"domain"`
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

// domainPath generates the resource path for given domain.
func domainPath(domain interface{}) string {
	if domain != nil {
		return fmt.Sprintf("domains/%s", domainIdentifier(domain))
	}
	return "domains"
}

// List the domains for the authenticated user.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/#list
func (s *DomainsService) List() ([]Domain, *Response, error) {
	path := domainPath(nil)
	wrappedDomains := []domainWrapper{}

	res, err := s.client.get(path, &wrappedDomains);
	if err != nil {
		return []Domain{}, res, err
	}

	domains := []Domain{}
	for _, domain := range wrappedDomains {
		domains = append(domains, domain.Domain)
	}

	return domains, res, nil
}

// Create a new domain in the authenticated account.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/#create
func (s *DomainsService) Create(domain Domain) (Domain, *Response, error) {
	path := domainPath(nil)
	wrappedDomain := domainWrapper{Domain: domain}
	returnedDomain := domainWrapper{}

	res, err := s.client.post(path, wrappedDomain, &returnedDomain)
	if err != nil {
		return Domain{}, res, err
	}

	if res.StatusCode == 400 {
		return Domain{}, res, errors.New("Invalid Domain")
	}

	return returnedDomain.Domain, res, nil
}

// Get fetches a domain from the authenticated account.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/#get
func (s *DomainsService) Get(domain interface{}) (Domain, *Response, error) {
	path := domainPath(domain)
	wrappedDomain := domainWrapper{}

	res, err := s.client.get(path, &wrappedDomain)
	if err != nil {
		return Domain{}, res, err
	}

	return wrappedDomain.Domain, res, nil
}

// Delete a domain from the authenticated account.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/#delete
func (s *DomainsService) Delete(domain interface{}) (*Response, error) {
	path := domainPath(domain)

	res, err := s.client.delete(path, nil)
	return res, err
}

func (s *DomainsService) CheckAvailability(domain interface{}) (bool, *Response, error) {
	path := fmt.Sprintf("%s/check", domainPath(domain))

	res, err := s.client.get(path, nil)
	if err != nil && res != nil && res.StatusCode != 404 {
		return false, res, err
	}

	return res.StatusCode == 404, res, nil
}

func (s *DomainsService) Renew(domain string, renewWhoisPrivacy bool) (*Response, error) {
	wrappedDomain := domainWrapper{Domain: Domain{
		Name:              domain,
		RenewWhoisPrivacy: renewWhoisPrivacy}}

	res, err := s.client.post("domain_renewals", wrappedDomain, nil)
	if err != nil {
		return res, err
	}

	if res.StatusCode == 400 {
		return res, errors.New("Failed to Renew")
	}

	return res, nil
}
