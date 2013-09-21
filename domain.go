package dnsimple

import (
	"encoding/json"
	"fmt"
)

type Domain struct {
	Id           int
	Name         string
	Language     string
	Lockable     bool
	State        string
	Token        string
	AutoRenew    bool   `json:"auto_renew"`
	ExpiresOn    string `json:"expires_on"`
	RegistrantId int    `json:"registrant_id"`
	UnicodeName  string `json:"unicode_name"`
	UserId       int    `json:"user_id"`
	RecordCount  int    `json:"record_count"`
	ServiceCount int    `json:"service_count"`
	PrivateWhois bool   `json:"private_whois?"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type domainWrapper struct {
	Domain Domain
}

func (client *DNSimpleClient) Domains() ([]Domain, error) {
	reqStr := "https://dnsimple.com/domains"

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return []Domain{}, err
	}

	var domainList []domainWrapper

	if err = json.Unmarshal([]byte(body), &domainList); err != nil {
		return []Domain{}, err
	}

	domains := []Domain{}
	for _, domain := range domainList {
		domains = append(domains, domain.Domain)
	}

	return domains, nil
}

func (client *DNSimpleClient) Domain(domain string) (Domain, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s", domain)

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return Domain{}, err
	}

	wrappedDomain := domainWrapper{}

	if err = json.Unmarshal([]byte(body), &wrappedDomain); err != nil {
		return Domain{}, err
	}
	return wrappedDomain.Domain, nil
}

func (client *DNSimpleClient) DomainAvailable(domain string) (bool, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s/check", domain)

	req, err := client.makeRequest("GET", reqStr, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 404, nil
}
