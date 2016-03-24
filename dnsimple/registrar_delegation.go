package dnsimple

import (
	"fmt"
)

// Domain represents a list of name servers that correspond to a domain delegation.
type Delegation []string

// WhoisPrivacyResponse represents a response from an API method that returns a delegation struct.
type DelegationResponse struct {
	Response
	Data Delegation `json:"data"`
}

// GetDomainDelegation gets the current delegated name servers for the domain.
//
// See https://developer.dnsimple.com/v2/registrar/delegation/#get
func (s *RegistrarService) GetDomainDelegation(accountID string, domainName string) (*DelegationResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/delegation", accountID, domainName))
	delegationResponse := &DelegationResponse{}

	resp, err := s.client.get(path, delegationResponse)
	if err != nil {
		return nil, err
	}

	delegationResponse.HttpResponse = resp
	return delegationResponse, nil
}
