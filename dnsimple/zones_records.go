package dnsimple

import (
	"fmt"
)

type Record struct {
	Id        int    `json:"id,omitempty"`
	ZoneId    string `json:"zone_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Priority  int    `json:"priority,omitempty"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type recordsWrapper struct {
	Records []Record `json:"data"`
}
type recordWrapper struct {
	Record Record `json:"data"`
}

// recordPath generates the resource path for given record that belongs to a domain.
func recordPath(accountId string, domain interface{}, record interface{}) string {
	path := fmt.Sprintf("%s/zones/%s/records", accountId, domainIdentifier(domain))

	if record != nil {
		path += fmt.Sprintf("/%d", record)
	}

	return path
}

// List the zone records.
//
// See https://developer.dnsimple.com/v2/zones/#list
func (s *ZonesService) ListRecords(accountId string, domain interface{}) ([]Record, *Response, error) {
	path := recordPath(accountId, domain, nil)
	data := recordsWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return []Record{}, res, err
	}

	return data.Records, res, nil
}

// CreateRecord creates a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#create
func (s *ZonesService) CreateRecord(accountId string, domain interface{}, recordAttributes Record) (Record, *Response, error) {
	path := recordPath(accountId, domain, nil)
	data := recordWrapper{}

	res, err := s.client.post(path, recordAttributes, &data)
	if err != nil {
		return Record{}, res, err
	}

	return data.Record, res, nil
}

// GetRecord gets the zone record.
//
// See https://developer.dnsimple.com/v2/zones/#get
func (s *ZonesService) GetRecord(accountId string, domain interface{}, recordID int) (Record, *Response, error) {
	path := recordPath(accountId, domain, recordID)
	data := recordWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return Record{}, res, err
	}

	return data.Record, res, nil
}

// UpdateRecord updates a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#update
func (s *ZonesService) UpdateRecord(accountId string, domain interface{}, recordID int, recordAttributes Record) (Record, *Response, error) {
	path := recordPath(accountId, domain, recordID)
	data := recordWrapper{}

	res, err := s.client.patch(path, recordAttributes, &data)
	if err != nil {
		return Record{}, res, err
	}

	return data.Record, res, nil
}

// DeleteRecord deletes a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#delete
func (s *ZonesService) DeleteRecord(accountId string, domain interface{}, recordID int) (*Response, error) {
	path := recordPath(accountId, domain, recordID)

	return s.client.delete(path, nil)
}
