package dnsimple

import (
	"errors"
	"fmt"
	"net/url"
)

type Record struct {
	Id        int    `json:"id,omitempty"`
	DomainId  int    `json:"domain_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Priority  int    `json:"prio,omitempty"`
	Type      string `json:"record_type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type recordWrapper struct {
	Record Record `json:"record"`
}

// recordPath generates the resource path for given record that belongs to a domain.
func recordPath(domain interface{}, record interface{}) string {
	path := fmt.Sprintf("domains/%s/records", domainIdentifier(domain))

	if record != nil {
		path += fmt.Sprintf("/%d", record)
	}

	return path
}

// List the records for a domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#list
func (s *DomainsService) ListRecords(domain interface{}, recordName, recordType string) ([]Record, *Response, error) {
	reqStr := recordPath(domain, nil)
	v := url.Values{}

	if recordName != "" {
		v.Add("name", recordName)
	}
	if recordType != "" {
		v.Add("type", recordType)
	}
	reqStr += "?" + v.Encode()

	wrappedRecords := []recordWrapper{}

	res, err := s.client.get(reqStr, &wrappedRecords)
	if err != nil {
		return []Record{}, res, err
	}

	records := []Record{}
	for _, record := range wrappedRecords {
		records = append(records, record.Record)
	}

	return records, res, nil
}

// Create a new record for a domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#create
func (s *DomainsService) CreateRecord(domain interface{}, record Record) (Record, *Response, error) {
	wrappedRecord := recordWrapper{Record: record}
	returnedRecord := recordWrapper{}

	res, err := s.client.post(recordPath(domain, nil), wrappedRecord, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if res.StatusCode == 400 {
		return Record{}, res, errors.New("Invalid Record")
	}

	return returnedRecord.Record, res, nil
}

// Get fetches a record for a domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#get
func (s *DomainsService) GetRecord(domain interface{}, recordID int) (Record, *Response, error) {
	wrappedRecord := recordWrapper{}

	res, err := s.client.get(recordPath(domain, recordID), &wrappedRecord)
	if err != nil {
		return Record{}, res, err
	}

	return wrappedRecord.Record, res, nil
}

// Delete a record for a domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#delete
func (s *DomainsService) DeleteRecord(domain interface{}, recordID int) (*Response, error) {
	path := recordPath(domain, recordID)

	res, err := s.client.delete(path, nil)
	return res, err
}

func (record *Record) Update(client *Client, recordAttributes Record) (Record, error) {
	// name, content, ttl, prio - only things allowed
	wrappedRecord := recordWrapper{Record: Record{
		Name:     recordAttributes.Name,
		Content:  recordAttributes.Content,
		TTL:      recordAttributes.TTL,
		Priority: recordAttributes.Priority}}

	returnedRecord := recordWrapper{}

	res, err := client.put(recordPath(record.DomainId, record.Id), wrappedRecord, &returnedRecord)
	if err != nil {
		return Record{}, err
	}

	if res.StatusCode == 400 {
		return Record{}, errors.New("Invalid Record")
	}

	return returnedRecord.Record, nil
}

func (record *Record) UpdateIP(client *Client, IP string) error {
	newRecord := Record{Content: IP, Name: record.Name}
	_, err := record.Update(client, newRecord)
	return err
}
