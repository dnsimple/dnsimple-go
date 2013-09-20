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

type recordList struct {
	Record Record
}

type Domain struct {
	Id   int
	Name string
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

func (client *DNSimpleClient) Record(domain, name string) (Record, error) {
	reqStr := fmt.Sprintf("https://dnsimple.com/domains/%s/records?name=%s", domain, name)

	body, err := client.sendRequest("GET", reqStr, nil)
	if err != nil {
		return Record{}, err
	}

	var records []recordList

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
