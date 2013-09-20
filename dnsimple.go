package dnsimple

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Record struct {
	Id         int
	Content    string
	Name       string
	TTL        int
	RecordType string    `json:"record_type"`
	Priority   int       `json:"prio"`
	DomainId   int       `json:"domain_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type recordWrapper struct {
	Record Record
}

type Domain struct {
	Id           int
	Name         string
	Language     string
	Lockable     bool
	State        string
	Token        string
	AutoRenew    bool      `json:"auto_renew"`
	ExpiresOn    string    `json:"expires_on"`
	RegistrantId int       `json:"registrant_id"`
	UnicodeName  string    `json:"unicode_name"`
	UserId       int       `json:"user_id"`
	RecordCount  int       `json:"record_count"`
	ServiceCount int       `json:"service_count"`
	PrivateWhois bool      `json:"private_whois?"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type domainWrapper struct {
	Domain Domain
}

type DNSimpleClient struct {
	ApiToken    string
	Email       string
	DomainToken string
	HttpClient  *http.Client
}

func NewClient(apiToken, email string) *DNSimpleClient {
	return &DNSimpleClient{ApiToken: apiToken, Email: email, HttpClient: &http.Client{}}
}

func (client *DNSimpleClient) makeRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("X-DNSimple-Token", fmt.Sprintf("%s:%s", client.Email, client.ApiToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (client *DNSimpleClient) sendRequest(method, url string, body io.Reader) (string, error) {
	req, err := client.makeRequest(method, url, body)
	if err != nil {
		return "", err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
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

func (client *DNSimpleClient) Domains() ([]Domain, error) {
	reqStr := "https://dnsimple.com/domains"

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return []Domain{}, err
	}

	var domainList []domainWrapper

	if err = json.Unmarshal([]byte(body), &domainList); err != nil {
		return []Domain{}, err
	}

	domains := []Domain{}
	for _, domain := range domainList {
		domains = append(domains, domain.Domain)
	}

	return domains, nil
}

func (client *DNSimpleClient) Domain(domain string) (Domain, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s", domain)

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return Domain{}, err
	}

	wrappedDomain := domainWrapper{}

	if err = json.Unmarshal([]byte(body), &wrappedDomain); err != nil {
		return Domain{}, err
	}
	return wrappedDomain.Domain, nil
}

func (client *DNSimpleClient) DomainAvailable(domain string) (bool, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s/check", domain)

	req, err := client.makeRequest("GET", reqStr, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 404, nil
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
