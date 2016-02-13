package dnsimple

import (
//"fmt"
)

// EmailForwardResponse represents a response from an API method that returns an EmailForward struct.
type EmailForwardResponse struct {
	Response
	Data *EmailForward `json:"data"`
}

// EmailForwardsResponse represents a response from an API method that returns a collection of EmailForward struct.
type EmailForwardsResponse struct {
	Response
	Data []EmailForward `json:"data"`
}

// EmailForward represents an email forward in DNSimple.
type EmailForward struct {
	ID        int    `json:"id,omitempty"`
	DomainID  int    `json:"domain_id,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ListEmailForwards lists the email forwards for a domain.
//
// See https://developer.dnsimple.com/v2/domains/email-forwards/#list
func (s *DomainsService) ListEmailForwards(accountID string, domain interface{}) (*EmailForwardsResponse, error) {
	path := versioned(domainPath(accountID, domain) + "/email_forwards")
	forwardsResponse := &EmailForwardsResponse{}

	resp, err := s.client.get(path, forwardsResponse)
	if err != nil {
		return nil, err
	}

	forwardsResponse.HttpResponse = resp
	return forwardsResponse, nil
}
