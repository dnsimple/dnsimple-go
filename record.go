package dnsimple

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func (client *DNSimpleClient) Records(domain string) ([]Record, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s/records", domain)

	body, err := client.sendRequest("GET", reqStr, nil)
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

func (client *DNSimpleClient) Record(domain, name string) (Record, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s/records?name=%s", domain, name)

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return Record{}, err
	}

	var records []recordWrapper

	if err = json.Unmarshal([]byte(body), &records); err != nil {
		return Record{}, err
	}

	return records[0].Record, nil
}

func (client *DNSimpleClient) CreateRecord(domain string, record Record) (Record, error) {
	// pre-validate the Record?
	wrappedRecord := recordWrapper{Record: record}
	jsonPayload, err := json.Marshal(wrappedRecord)
	if err != nil {
		return Record{}, err
	}

	url := fmt.Sprintf("https://dnsimple.com/domains/%s/records", domain)

	resp, err := client.sendRequestResponse("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return Record{}, err
	}

	if resp.StatusCode == 400 {
		// 400: bad request, validation failed
		return Record{}, errors.New("Invalid Record")
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Record{}, err
	}

	if err = json.Unmarshal(responseBody, &wrappedRecord); err != nil {
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

	url := fmt.Sprintf("https://dnsimple.com/domains/%d/records/%d", record.DomainId, record.Id)
	fmt.Println(string(jsonPayload))

	resp, err := client.sendRequestResponse("PUT", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return Record{}, err
	}

	if resp.StatusCode == 400 {
		// 400: bad request, validation failed
		return Record{}, errors.New("Invalid Record")
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Record{}, err
	}

	if err = json.Unmarshal(responseBody, &wrappedRecord); err != nil {
		return Record{}, err
	}

	return wrappedRecord.Record, nil
}

func (record *Record) UpdateIP(client *DNSimpleClient, IP string) error {
	newRecord := Record{Content: IP}
	_, err := record.Update(client, newRecord)
	return err
}
