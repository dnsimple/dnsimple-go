package dnsimple

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Record struct {
	Id       int
	Name     string
	Content  string
	DomainId int `json:"domain_id"`
}

type Domain struct {
	Record Record
}

type DNSimpleClient struct {
	DomainToken string
	HttpClient  *http.Client
}

func NewClient(domainToken string) *DNSimpleClient {
	return &DNSimpleClient{DomainToken: domainToken, HttpClient: &http.Client{}}
}

func (client *DNSimpleClient) makeRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("X-DNSimple-Domain-Token", client.DomainToken)
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

func (client *DNSimpleClient) FindDomain(domain, name string) (Record, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s/records?name=%s", domain, name)

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return Record{}, err
	}

	var domains []Domain

	if err = json.Unmarshal([]byte(body), &domains); err != nil {
		return Record{}, err
	}

	return domains[0].Record, nil
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
