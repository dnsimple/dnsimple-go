package dnsimple

import (
	"fmt"
)

// WhoisPrivacy represents a whois privacy in DNSimple.
type WhoisPrivacy struct {
	ID        int    `json:"id,omitempty"`
	DomainID  int    `json:"domain_id,omitempty"`
	Enabled   bool   `json:"enabled,omitempty"`
	ExpiresOn string `json:"expires_on,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// WhoisPrivacyResponse represents a response from an API method that returns a WhoisPrivacy struct.
type WhoisPrivacyResponse struct {
	Response
	Data *WhoisPrivacy `json:"data"`
}

// GetWhoisPrivacy gets the whois privacy for the domain.
//
// See https://developer.dnsimple.com/v2/registrar/whois-privacy/#get
func (s *RegistrarService) GetWhoisPrivacy(accountID string, domainName string) (*WhoisPrivacyResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/whois_privacy", accountID, domainName))
	privacyResponse := &WhoisPrivacyResponse{}

	resp, err := s.client.get(path, privacyResponse)
	if err != nil {
		return nil, err
	}

	privacyResponse.HttpResponse = resp
	return privacyResponse, nil
}
