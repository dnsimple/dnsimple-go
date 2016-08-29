package dnsimple

import (
	"fmt"
)

// DomainPush represents a domain push in DNSimple.
type DomainPush struct {
	ID         int    `json:"id,omitempty"`
	DomainID   int    `json:"domain_id,omitempty"`
	ContactID  int    `json:"contact_id,omitempty"`
	AccountID  int    `json:"account_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	AcceptedAt string `json:"accepted_at,omitempty"`
}

// DomainPushAttributes represent a domain push payload (see initiate).
type DomainPushAttributes struct {
	NewAccountEmail string `json:"new_account_email,omitempty"`
}

// DomainPushResponse represents a response from an API method that returns a DomainPush struct.
type DomainPushResponse struct {
	Response
	Data *DomainPush `json:"data"`
}

func domainPushPath(accountID string, domain interface{}, pushID int) string {
	path := fmt.Sprintf("%v/pushes", domainPath(accountID, domain))

	if pushID != 0 {
		path += fmt.Sprintf("/%d", pushID)
	}

	return path
}

// InitiatePush initiate a new domain push.
//
// See https://developer.dnsimple.com/v2/domains/pushes/#initiate
func (s *DomainsService) InitiatePush(accountID string, domain interface{}, pushAttributes DomainPushAttributes) (*DomainPushResponse, error) {
	path := versioned(domainPushPath(accountID, domain, 0))
	pushResponse := &DomainPushResponse{}

	resp, err := s.client.post(path, pushAttributes, pushResponse)
	if err != nil {
		return nil, err
	}

	pushResponse.HttpResponse = resp
	return pushResponse, nil
}
