package dnsimple

import (
	"fmt"
	"time"
)

type Record struct {
	Id        int        `json:"id,omitempty"`
	ZoneId    int        `json:"zone_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Content   string     `json:"content,omitempty"`
	TTL       int        `json:"ttl,omitempty"`
	Priority  int        `json:"prio,omitempty"`
	Type      string     `json:"record_type,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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

// List the domain records.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#list
func (s *DomainsService) ListRecords(accountId string, domain interface{}) ([]Record, *Response, error) {
	path := recordPath(accountId, domain, nil)
	data := recordsWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return []Record{}, res, err
	}

	return data.Records, res, nil
}

// CreateRecord creates a domain record.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#create
func (s *DomainsService) CreateRecord(accountId string, domain interface{}, recordAttributes Record) (Record, *Response, error) {
	path := recordPath(accountId, domain, nil)
	data := recordWrapper{}

	res, err := s.client.post(path, recordAttributes, &data)
	if err != nil {
		return Record{}, res, err
	}

	return data.Record, res, nil
}

// GetRecord fetches the domain record.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#get
func (s *DomainsService) GetRecord(accountId string, domain interface{}, recordID int) (Record, *Response, error) {
	path := recordPath(accountId, domain, recordID)
	data := recordWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return Record{}, res, err
	}

	return data.Record, res, nil
}

// UpdateRecord updates a domain record.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#update
func (s *DomainsService) UpdateRecord(accountId string, domain interface{}, recordID int, recordAttributes Record) (Record, *Response, error) {
	path := recordPath(accountId, domain, recordID)
	// name, content, ttl, priority
	record := Record{
		Name:     recordAttributes.Name,
		Content:  recordAttributes.Content,
		TTL:      recordAttributes.TTL,
		Priority: recordAttributes.Priority}
	data := recordWrapper{}

	res, err := s.client.put(path, record, &data)
	if err != nil {
		return Record{}, res, err
	}

	return data.Record, res, nil
}

// DeleteRecord deletes a domain record.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#delete
func (s *DomainsService) DeleteRecord(accountId string, domain interface{}, recordID int) (*Response, error) {
	path := recordPath(accountId, domain, recordID)

	return s.client.delete(path, nil)
}

// UpdateIP updates the IP of specific A record.
//
// This is not part of the standard API. However,
// this is useful for Dynamic DNS (DDNS or DynDNS).
func (record *Record) UpdateIP(client *Client, IP, accountId string) error {
	newRecord := Record{Content: IP, Name: record.Name}
	_, _, err := client.Domains.UpdateRecord(accountId, record.ZoneId, record.Id, newRecord)
	return err
}
