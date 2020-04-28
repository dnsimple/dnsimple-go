package dnsimple

import (
	"context"
)

// AccountsService handles communication with the account related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/accounts/
type AccountsService struct {
	client *Client
}

// Account represents a DNSimple account.
type Account struct {
	ID             int64  `json:"id,omitempty"`
	Email          string `json:"email,omitempty"`
	PlanIdentifier string `json:"plan_identifier,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// accountsResponse represents a response from an API method that returns a collection of Account struct.
type accountsResponse struct {
	Response
	Data []Account `json:"data"`
}

// ListAccounts list the accounts for an user.
//
// See https://developer.dnsimple.com/v2/accounts/#list
func (s *AccountsService) ListAccounts(options *ListOptions) (*accountsResponse, error) {
	path := versioned("/accounts")
	accountsResponse := &accountsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(context.TODO(), path, accountsResponse)
	if err != nil {
		return accountsResponse, err
	}

	accountsResponse.HTTPResponse = resp
	return accountsResponse, nil
}
