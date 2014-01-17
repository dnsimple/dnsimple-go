// Package dnsimple implements a client for the DNSimple API.
//
// In order to use this package you will need a DNSimple account and your API Token.
package dnsimple

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://api.dnsimple.com/"
	userAgent      = "go-dnsimple/" + libraryVersion

	apiVersion = "v1"
)

type Client struct {
	// HTTP client used to communicate with the API.
	HttpClient *http.Client

	// API Token for the DNSimple account you want to use.
	ApiToken string

	// Email associated with the provided DNSimple API Token.
	Email string

	// Domain Token to be used for authentication
	// as an alternative to the DNSimple API Token for some domain-scoped operations.
	DomainToken string

	// Base URL for API requests.
	// Defaults to the public DNSimple API, but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	// User agent used when communicating with the DNSimple API.
	UserAgent string

	// Services used for talking to different parts of the GitHub API.
	Domains *DomainsService
	Records *RecordsService
}

// NewClient returns a new GitHub API client.
func NewClient(apiToken, email string) *Client {
	c := &Client{ApiToken: apiToken, Email: email, HttpClient: &http.Client{}, BaseURL: defaultBaseURL, UserAgent: userAgent}
	c.Domains = &DomainsService{client: c}
	c.Records = &RecordsService{client: c}
	return c
}

// NewRequest creates an API request.
// The path is expected to be a relative path and will be resolved
// according to the BaseURL of the Client. Paths should always be specified without a preceding slash.
func (client *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	url := client.BaseURL + fmt.Sprintf("%s/%s", apiVersion, path)

	body := new(bytes.Buffer)
	if payload != nil {
		err := json.NewEncoder(body).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", client.UserAgent)
	req.Header.Add("X-DNSimple-Token", fmt.Sprintf("%s:%s", client.Email, client.ApiToken))

	return req, nil
}

func (client *Client) get(path string, val interface{}) (*http.Response, error) {
	return client.sendRequest("GET", path, nil, val)
}

func (client *Client) post(path string, payload, val interface{}) (*http.Response, error) {
	return client.sendRequest("POST", path, payload, val)
}

func (client *Client) put(path string, payload, val interface{}) (*http.Response, error) {
	return client.sendRequest("PUT", path, payload, val)
}

func (client *Client) delete(path string, payload interface{}) (*http.Response, error) {
	return client.sendRequest("DELETE", path, payload, nil)
}

func (client *Client) sendRequest(method, path string, payload, value interface{}) (*http.Response, error) {
	req, err := client.NewRequest(method, path, payload)
	if err != nil {
		return nil, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if value != nil {
		err = json.NewDecoder(resp.Body).Decode(value)
	}

	return resp, err
}
