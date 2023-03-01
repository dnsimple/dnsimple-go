package dnsimple

import (
	"context"
	"errors"
	"fmt"
)

// DomainPush represents a domain push in DNSimple.
type DomainPush struct {
	ID         int64  `json:"id,omitempty"`
	DomainID   int64  `json:"domain_id,omitempty"`
	ContactID  int64  `json:"contact_id,omitempty"`
	AccountID  int64  `json:"account_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	AcceptedAt string `json:"accepted_at,omitempty"`
}

func domainPushesPath(accountID string, domainIdentifier string) (string, error) {
	basePath, err := domainPath(accountID, domainIdentifier)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/%v/pushes", basePath), nil
}

func pushesPath(accountID string) (string, error) {
	if accountID == "" {
		return "", errors.New("account parameter should not be empty")
	}

	return fmt.Sprintf("/%v/pushes", accountID), nil

}

func pushPath(accountID string, pushID int64) (string, error) {
	basePath, err := pushesPath(accountID)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%v/%v", basePath, pushID), nil
}

// DomainPushResponse represents a response from an API method that returns a DomainPush struct.
type DomainPushResponse struct {
	Response
	Data *DomainPush `json:"data"`
}

// DomainPushesResponse represents a response from an API method that returns a collection of DomainPush struct.
type DomainPushesResponse struct {
	Response
	Data []DomainPush `json:"data"`
}

// DomainPushAttributes represent a domain push payload (see initiate).
type DomainPushAttributes struct {
	NewAccountEmail string `json:"new_account_email,omitempty"`
	ContactID       int64  `json:"contact_id,omitempty"`
}

// InitiatePush initiate a new domain push.
//
// See https://developer.dnsimple.com/v2/domains/pushes/#initiateDomainPush
func (s *DomainsService) InitiatePush(ctx context.Context, accountID, domainID string, pushAttributes DomainPushAttributes) (*DomainPushResponse, error) {
	path, err := domainPushesPath(accountID, domainID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	pushResponse := &DomainPushResponse{}

	resp, err := s.client.post(ctx, path, pushAttributes, pushResponse)
	if err != nil {
		return nil, err
	}

	pushResponse.HTTPResponse = resp
	return pushResponse, nil
}

// ListPushes lists the pushes for an account.
//
// See https://developer.dnsimple.com/v2/domains/pushes/#listPushes
func (s *DomainsService) ListPushes(ctx context.Context, accountID string, options *ListOptions) (*DomainPushesResponse, error) {
	path, err := pushesPath(accountID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	path, err = addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	pushesResponse := &DomainPushesResponse{}

	resp, err := s.client.get(ctx, path, pushesResponse)
	if err != nil {
		return nil, err
	}

	pushesResponse.HTTPResponse = resp
	return pushesResponse, nil
}

// AcceptPush accept a push for a domain.
//
// See https://developer.dnsimple.com/v2/domains/pushes/#acceptPush
func (s *DomainsService) AcceptPush(ctx context.Context, accountID string, pushID int64, pushAttributes DomainPushAttributes) (*DomainPushResponse, error) {
	path, err := pushPath(accountID, pushID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	pushResponse := &DomainPushResponse{}

	resp, err := s.client.post(ctx, path, pushAttributes, nil)
	if err != nil {
		return nil, err
	}

	pushResponse.HTTPResponse = resp
	return pushResponse, nil
}

// RejectPush reject a push for a domain.
//
// See https://developer.dnsimple.com/v2/domains/pushes/#rejectPush
func (s *DomainsService) RejectPush(ctx context.Context, accountID string, pushID int64) (*DomainPushResponse, error) {
	path, err := pushPath(accountID, pushID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	pushResponse := &DomainPushResponse{}

	resp, err := s.client.delete(ctx, path, nil, nil)
	if err != nil {
		return nil, err
	}

	pushResponse.HTTPResponse = resp
	return pushResponse, nil
}
