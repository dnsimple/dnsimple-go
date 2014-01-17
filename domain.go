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
	client *DNSimpleClient
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
// DNSimple API docs: http://developer.dnsimple.com/domains/#list-domains
func (s *DomainsService) List() ([]Domain, error) {
	wrappedDomains := []domainWrapper{}

	if _, err := s.client.get(domainPath(nil), &wrappedDomains); err != nil {
		return []Domain{}, err
	}

	domains := []Domain{}
	for _, domain := range wrappedDomains {
		domains = append(domains, domain.Domain)
	}

	return domains, nil
}

// Create a new domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/#create-a-domain
func (s *DomainsService) Create(domain Domain) (Domain, error) {
	wrappedDomain := domainWrapper{Domain: domain}
	returnedDomain := domainWrapper{}

	response, err := s.client.post(domainPath(nil), wrappedDomain, &returnedDomain)
	if err != nil {
		return Domain{}, err
	}

	if response.StatusCode == 400 {
		return Domain{}, errors.New("Invalid Domain")
	}

	return returnedDomain.Domain, nil
}

// Get fetches a domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/#get-a-domain
func (s *DomainsService) Get(domain interface{}) (Domain, error) {
	wrappedDomain := domainWrapper{}

	if _, err := s.client.get(domainPath(domain), &wrappedDomain); err != nil {
		return Domain{}, err
	}

	return wrappedDomain.Domain, nil
}

func (s *DomainsService) SetAutoRenewal(domain interface{}, autorenew bool) error {
	path := fmt.Sprintf("%s/auto_renewal", domainPath(domain))
	method := ""
	if autorenew {
		method = "POST"
	} else {
		method = "DELETE"
	}

	if _, err := s.client.sendRequest(method, path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (s *DomainsService) CheckAvailability(domain interface{}) (bool, error) {
	path := fmt.Sprintf("%s/check", domainPath(domain))

	response, err := s.client.get(path, nil)
	if err != nil {
		return false, err
	}

	return response.StatusCode == 404, nil
}

func (s *DomainsService) Renew(domain string, renewWhoisPrivacy bool) error {
	wrappedDomain := domainWrapper{Domain: Domain{
		Name:              domain,
		RenewWhoisPrivacy: renewWhoisPrivacy}}

	response, err := s.client.post("domain_renewals", wrappedDomain, nil)
	if err != nil {
		return err
	}

	if response.StatusCode == 400 {
		return errors.New("Failed to Renew")
	}

	return nil
}
