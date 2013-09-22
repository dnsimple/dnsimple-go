// Package dnsimple implements a client for the DNSimple API.
//
// In order to use this package you will need a DNSimple account and your API Token.
package dnsimple

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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

func (client *DNSimpleClient) sendRequestResponse(method, url string, body io.Reader) (*http.Response, error) {
	req, err := client.makeRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *DNSimpleClient) sendRequest(method, url string, body io.Reader) (string, error) {
	resp, err := client.sendRequestResponse(method, url, body)
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

func domainIdentifier(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return fmt.Sprintf("%d", value)
	case Domain:
		return fmt.Sprintf("%d", value.Id)
	case Record:
		return fmt.Sprintf("%d", value.DomainId)
	}
	return ""
}
