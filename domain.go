package dnsimple

import (
	"encoding/json"
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
	body, err := client.sendRequest("GET", domainURL(nil), nil)
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

func (client *DNSimpleClient) Domain(domain interface{}) (Domain, error) {
	body, err := client.sendRequest("GET", domainURL(domain), nil)
	if err != nil {
		return Domain{}, err
	}

	wrappedDomain := domainWrapper{}

	if err = json.Unmarshal([]byte(body), &wrappedDomain); err != nil {
		return Domain{}, err
	}
	return wrappedDomain.Domain, nil
}

func (client *DNSimpleClient) DomainAvailable(domain interface{}) (bool, error) {
	reqStr := fmt.Sprintf("%s/check", domainURL(domain))

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
