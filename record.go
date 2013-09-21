package dnsimple

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Record struct {
	Id         int
	Content    string
	Name       string
	TTL        int
	RecordType string `json:"record_type"`
	Priority   int    `json:"prio"`
	DomainId   int    `json:"domain_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type recordWrapper struct {
	Record Record
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

func (record *Record) UpdateIP(client *DNSimpleClient, IP string) error {
	// lame, but easy enough for now
	jsonPayload := fmt.Sprintf(`{"record": {"content": "%s"}}`, IP)
	url := fmt.Sprintf("https://dnsimple.com/domains/%d/records/%d", record.DomainId, record.Id)

	_, err := client.sendRequest("PUT", url, strings.NewReader(jsonPayload))
	if err != nil {
		return err
	}

	return nil
}
