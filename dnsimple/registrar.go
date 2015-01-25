package dnsimple

import (
	"fmt"
)

// RegistrarService handles communication with the registrar related
// methods of the DNSimple API.
//
// DNSimple API docs: http://developer.dnsimple.com/registrar/
type RegistrarService struct {
	client *Client
}

// IsAvailable checks if the domain is available or registered.
func (s *RegistrarService) IsAvailable(domain string) (bool, error) {
	path := fmt.Sprintf("%s/check", domainPath(domain))

	res, err := s.client.get(path, nil)
	if err != nil && res != nil && res.StatusCode != 404 {
		return false, err
	}

	return res.StatusCode == 404, nil
}


type dWrapper struct {
	Domain interface{} `json:"domain"`
}

// renewDomain represents the body of a Renew request.
type renewDomain struct {
	Name string `json:"name,omitempty"`
	RenewWhoisPrivacy bool `json:"renew_whois_privacy,omitempty"`
}

func (s *RegistrarService) Renew(domain string, renewWhoisPrivacy bool) (*Response, error) {
	request := dWrapper{Domain: renewDomain{
		Name:              domain,
		RenewWhoisPrivacy: renewWhoisPrivacy}}

	res, err := s.client.post("domain_renewals", &request, nil)
	if err != nil {
		return res, err
	}

	return res, nil
}
