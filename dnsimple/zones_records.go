package dnsimple

import (
	"fmt"
)

// ZoneRecordResponse represents a response from an API method that returns a ZoneRecord struct.
type ZoneRecordResponse struct {
	Response
	Data *Record `json:"data"`
}

// ZoneRecordsResponse represents a response from an API method that returns a collection of ZoneRecord struct.
type ZoneRecordsResponse struct {
	Response
	Data []Record `json:"data"`
}

type Record struct {
	ID           int    `json:"id,omitempty"`
	ZoneID       string `json:"zone_id,omitempty"`
	ParentID     int    `json:"parent_id,omitempty"`
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	Content      string `json:"content,omitempty"`
	TTL          int    `json:"ttl,omitempty"`
	Priority     int    `json:"priority,omitempty"`
	SystemRecord bool   `json:"system_record,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// recordPath generates the resource path for given record that belongs to a domain.
func recordPath(accountID string, domain interface{}, record interface{}) string {
	path := fmt.Sprintf("/%v/zones/%v/records", accountID, domainIDentifier(domain))

	if record != nil {
		path += fmt.Sprintf("/%v", record)
	}

	return path
}

// List the zone records.
//
// See https://developer.dnsimple.com/v2/zones/#list
func (s *ZonesService) ListRecords(accountID string, domain interface{}) (*ZoneRecordsResponse, error) {
	path := recordPath(accountID, domain, nil)
	recordsResponse := &ZoneRecordsResponse{}

	resp, err := s.client.get(path, recordsResponse)
	if err != nil {
		return nil, err
	}

	recordsResponse.HttpResponse = resp
	return recordsResponse, nil
}

// CreateRecord creates a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#create
func (s *ZonesService) CreateRecord(accountID string, domain interface{}, recordAttributes Record) (*ZoneRecordResponse, error) {
	path := recordPath(accountID, domain, nil)
	recordResponse := &ZoneRecordResponse{}

	resp, err := s.client.post(path, recordAttributes, recordResponse)
	if err != nil {
		return nil, err
	}

	recordResponse.HttpResponse = resp
	return recordResponse, nil
}

// GetRecord gets the zone record.
//
// See https://developer.dnsimple.com/v2/zones/#get
func (s *ZonesService) GetRecord(accountID string, domain interface{}, recordID int) (*ZoneRecordResponse, error) {
	path := recordPath(accountID, domain, recordID)
	recordResponse := &ZoneRecordResponse{}

	resp, err := s.client.get(path, recordResponse)
	if err != nil {
		return nil, err
	}

	recordResponse.HttpResponse = resp
	return recordResponse, nil
}

// UpdateRecord updates a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#update
func (s *ZonesService) UpdateRecord(accountID string, domain interface{}, recordID int, recordAttributes Record) (*ZoneRecordResponse, error) {
	path := recordPath(accountID, domain, recordID)
	recordResponse := &ZoneRecordResponse{}

	resp, err := s.client.patch(path, recordAttributes, recordResponse)
	if err != nil {
		return nil, err
	}

	recordResponse.HttpResponse = resp
	return recordResponse, nil
}

// DeleteRecord deletes a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#delete
func (s *ZonesService) DeleteRecord(accountID string, domain interface{}, recordID int) (*ZoneRecordResponse, error) {
	path := recordPath(accountID, domain, recordID)
	recordResponse := &ZoneRecordResponse{}

	resp, err := s.client.delete(path, nil, nil)
	if err != nil {
		return nil, err
	}

	recordResponse.HttpResponse = resp
	return recordResponse, nil
}
