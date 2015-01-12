package dnsimple

import (
	"fmt"
)

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

	return res, nil
}
