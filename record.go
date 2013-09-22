package dnsimple

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	Id         int    `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	Name       string `json:"name,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	Priority   int    `json:"prio,omitempty"`
	DomainId   int    `json:"domain_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type recordWrapper struct {
	Record Record `json:"record"`
}

func recordURL(domain interface{}, record *Record) string {
	str := fmt.Sprintf("https://dnsimple.com/domains/%s/records", domainIdentifier(domain))
	if record != nil {
		str += fmt.Sprintf("/%d", record.Id)
	}
	return str
}

func (client *DNSimpleClient) Records(domain interface{}) ([]Record, error) {
	body, _, err := client.sendRequest("GET", recordURL(domain, nil), nil)
	if err != nil {
		return []Record{}, err
	}

	var recordList []recordWrapper

	if err = json.Unmarshal([]byte(body), &recordList); err != nil {
		return []Record{}, err
	}

	records := []Record{}
	for _, record := range recordList {
		records = append(records, record.Record)
	}

	return records, nil
}

func (client *DNSimpleClient) Record(domain interface{}, name string) (Record, error) {
	reqStr := fmt.Sprintf("%s?name=%s", recordURL(domain, nil), name)

	body, _, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return Record{}, err
	}

	var records []recordWrapper

	if err = json.Unmarshal([]byte(body), &records); err != nil {
		return Record{}, err
	}

	if len(records) == 0 {
		return Record{}, errors.New("Domain not found")
	}

	return records[0].Record, nil
}

func (client *DNSimpleClient) CreateRecord(domain interface{}, record Record) (Record, error) {
	// pre-validate the Record?
	wrappedRecord := recordWrapper{Record: record}
	jsonPayload, err := json.Marshal(wrappedRecord)
	if err != nil {
		return Record{}, err
	}

	resp, status, err := client.sendRequest("POST", recordURL(domain, nil), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return Record{}, err
	}

	if status == 400 {
		// 400: bad request, validation failed
		return Record{}, errors.New("Invalid Record")
	}

	if err = json.Unmarshal([]byte(resp), &wrappedRecord); err != nil {
		return Record{}, err
	}

	return wrappedRecord.Record, nil
}

func (record *Record) Update(client *DNSimpleClient, recordAttributes Record) (Record, error) {
	// pre-validate the Record?
	// name, content, ttl, prio - only things allowed
	wrappedRecord := recordWrapper{Record: Record{
		Name:     recordAttributes.Name,
		Content:  recordAttributes.Content,
		TTL:      recordAttributes.TTL,
		Priority: recordAttributes.Priority}}

	jsonPayload, err := json.Marshal(wrappedRecord)
	if err != nil {
		return Record{}, err
	}

	resp, status, err := client.sendRequest("PUT", recordURL(record.DomainId, record), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return Record{}, err
	}

	if status == 400 {
		// 400: bad request, validation failed
		return Record{}, errors.New("Invalid Record")
	}

	if err = json.Unmarshal([]byte(resp), &wrappedRecord); err != nil {
		return Record{}, err
	}

	return wrappedRecord.Record, nil
}

func (record *Record) Delete(client *DNSimpleClient) error {
	_, status, err := client.sendRequest("DELETE", recordURL(record.DomainId, record), nil)
	if err != nil {
		return err
	}

	if status == 200 {
		return nil
	}
	return errors.New("Failed to delete domain")
}

func (record *Record) UpdateIP(client *DNSimpleClient, IP string) error {
	newRecord := Record{Content: IP}
	_, err := record.Update(client, newRecord)
	return err
}
