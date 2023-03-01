package dnsimple

import (
	"context"
	"fmt"
)

// DelegationSignerRecord represents a delegation signer record for a domain in DNSimple.
type DelegationSignerRecord struct {
	ID         int64  `json:"id,omitempty"`
	DomainID   int64  `json:"domain_id,omitempty"`
	Algorithm  string `json:"algorithm"`
	Digest     string `json:"digest,omitempty"`
	DigestType string `json:"digest_type,omitempty"`
	Keytag     string `json:"keytag,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

func delegationSignerRecordsPath(accountID string, domainIdentifier string) (string, error) {
	basePath, err := domainPath(accountID, domainIdentifier)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v/ds_records", basePath), nil
}

func delegationSignerRecordPath(accountID string, domainIdentifier string, dsRecordID int64) (string, error) {
	basePath, err := delegationSignerRecordsPath(accountID, domainIdentifier)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v/%v", basePath, dsRecordID), nil
}

// DelegationSignerRecordResponse represents a response from an API method that returns a DelegationSignerRecord struct.
type DelegationSignerRecordResponse struct {
	Response
	Data *DelegationSignerRecord `json:"data"`
}

// DelegationSignerRecordsResponse represents a response from an API method that returns a DelegationSignerRecord struct.
type DelegationSignerRecordsResponse struct {
	Response
	Data []DelegationSignerRecord `json:"data"`
}

// ListDelegationSignerRecords lists the delegation signer records for a domain.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#ds-record-list
func (s *DomainsService) ListDelegationSignerRecords(ctx context.Context, accountID string, domainIdentifier string, options *ListOptions) (*DelegationSignerRecordsResponse, error) {
	path, err := delegationSignerRecordsPath(accountID, domainIdentifier)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	path, err = addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	dsRecordsResponse := &DelegationSignerRecordsResponse{}

	resp, err := s.client.get(ctx, path, dsRecordsResponse)
	if err != nil {
		return nil, err
	}

	dsRecordsResponse.HTTPResponse = resp
	return dsRecordsResponse, nil
}

// CreateDelegationSignerRecord creates a new delegation signer record.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#ds-record-create
func (s *DomainsService) CreateDelegationSignerRecord(ctx context.Context, accountID string, domainIdentifier string, dsRecordAttributes DelegationSignerRecord) (*DelegationSignerRecordResponse, error) {
	path, err := delegationSignerRecordsPath(accountID, domainIdentifier)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	dsRecordResponse := &DelegationSignerRecordResponse{}

	resp, err := s.client.post(ctx, path, dsRecordAttributes, dsRecordResponse)
	if err != nil {
		return nil, err
	}

	dsRecordResponse.HTTPResponse = resp
	return dsRecordResponse, nil
}

// GetDelegationSignerRecord fetches a delegation signer record.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#ds-record-get
func (s *DomainsService) GetDelegationSignerRecord(ctx context.Context, accountID string, domainIdentifier string, dsRecordID int64) (*DelegationSignerRecordResponse, error) {
	path, err := delegationSignerRecordPath(accountID, domainIdentifier, dsRecordID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	dsRecordResponse := &DelegationSignerRecordResponse{}

	resp, err := s.client.get(ctx, path, dsRecordResponse)
	if err != nil {
		return nil, err
	}

	dsRecordResponse.HTTPResponse = resp
	return dsRecordResponse, nil
}

// DeleteDelegationSignerRecord PERMANENTLY deletes a delegation signer record
// from the domain.
//
// See https://developer.dnsimple.com/v2/domains/dnssec/#ds-record-delete
func (s *DomainsService) DeleteDelegationSignerRecord(ctx context.Context, accountID string, domainIdentifier string, dsRecordID int64) (*DelegationSignerRecordResponse, error) {
	path, err := delegationSignerRecordPath(accountID, domainIdentifier, dsRecordID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	dsRecordResponse := &DelegationSignerRecordResponse{}

	resp, err := s.client.delete(ctx, path, nil, nil)
	if err != nil {
		return nil, err
	}

	dsRecordResponse.HTTPResponse = resp
	return dsRecordResponse, nil
}
