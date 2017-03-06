package dnsimple

import "fmt"

// DelegationSignerRecord represents a delegation signer record for a domain in DNSimple.
type DelegationSignerRecord struct {
	ID         int    `json:"id,omitempty"`
	DomainID   int    `json:"domain_id,omitempty"`
	Algorithm  string `json:"algorithm"`
	Digest     string `json:"digest"`
	DigestType string `json:"digest_type"`
	Keytag     string `json:"keytag"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

func delegationSignerRecordPath(accountID string, domainIdentifier string, dsRecordID int) (path string) {
	path = fmt.Sprintf("%v/ds_records", domainPath(accountID, domainIdentifier))
	if dsRecordID != 0 {
		path += fmt.Sprintf("/%d", dsRecordID)
	}
	return
}

// delegationSignerRecordResponse represents a response from an API method that returns a DelegationSignerRecord struct.
type delegationSignerRecordResponse struct {
	Response
	Data *DelegationSignerRecord `json:"data"`
}

// delegationSignerRecordResponse represents a response from an API method that returns a DelegationSignerRecord struct.
type delegationSignerRecordsResponse struct {
	Response
	Data []DelegationSignerRecord `json:"data"`
}

// ListDelegationSignerRecords lists the delegation signer records for a domain.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#ds-record-list
func (s *DomainsService) ListDelegationSignerRecords(accountID string, domainIdentifier string, options *ListOptions) (*delegationSignerRecordsResponse, error) {
	path := versioned(delegationSignerRecordPath(accountID, domainIdentifier, 0))
	dsRecordsResponse := &delegationSignerRecordsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, dsRecordsResponse)
	if err != nil {
		return nil, err
	}

	dsRecordsResponse.HttpResponse = resp
	return dsRecordsResponse, nil
}
