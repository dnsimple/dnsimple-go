package dnsimple

import (
	"fmt"
)

type Domain struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Language     string `json:"language,omitempty"`
	Lockable     bool   `json:"lockable,omitempty"`
	State        string `json:"state,omitempty"`
	Token        string `json:"token,omitempty"`
	AutoRenew    bool   `json:"auto_renew,omitempty"`
	ExpiresOn    string `json:"expires_on,omitempty"`
	RegistrantId int    `json:"registrant_id,omitempty"`
	UnicodeName  string `json:"unicode_name,omitempty"`
	UserId       int    `json:"user_id,omitempty"`
	RecordCount  int    `json:"record_count,omitempty"`
	ServiceCount int    `json:"service_count,omitempty"`
	PrivateWhois bool   `json:"private_whois?,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type domainWrapper struct {
	Domain Domain
}

func domainURL(domain interface{}) string {
	str := "https://dnsimple.com/domains"
	if domain != nil {
		str += fmt.Sprintf("/%s", domainIdentifier(domain))
	}
	return str
}

func (client *DNSimpleClient) Domains() ([]Domain, error) {
	wrappedDomains := []domainWrapper{}

	if err := client.get(domainURL(nil), &wrappedDomains); err != nil {
		return []Domain{}, err
	}

	domains := []Domain{}
	for _, domain := range wrappedDomains {
		domains = append(domains, domain.Domain)
	}

	return domains, nil
}

func (client *DNSimpleClient) Domain(domain interface{}) (Domain, error) {
	wrappedDomain := domainWrapper{}

	if err := client.get(domainURL(domain), &wrappedDomain); err != nil {
		return Domain{}, err
	}

	return wrappedDomain.Domain, nil
}

func (client *DNSimpleClient) DomainAvailable(domain interface{}) (bool, error) {
	reqStr := fmt.Sprintf("%s/check", domainURL(domain))

	_, status, err := client.sendRequest("GET", reqStr, nil)

	if err != nil {
		return false, err
	}

	return status == 404, nil
}

func (client *DNSimpleClient) SetAutorenew(domain interface{}, autorenew bool) error {
	reqStr := fmt.Sprintf("%s/auto_renewal", domainURL(domain))

	method := ""
	if autorenew {
		method = "POST"
	} else {
		method = "DELETE"
	}
	_, _, err := client.sendRequest(method, reqStr, nil)

	if err != nil {
		return err
	}
	return nil
}
