// Package dnsimple implements a client for the DNSimple API.
//
// In order to use this package you will need a DNSimple account and your API Token.
package dnsimple

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"strings"
	//"io/ioutil"
)

const (
	// libraryVersion identifies the current library version.
	// This is a pro-forma convention given that Go dependencies
	// tends to be fetched directly from the repo.
	// It is also used in the user-agent identify the client.
	libraryVersion = "0.5.0-dev"

	// defaultBaseURL to the DNSimple production API.
	defaultBaseURL = "https://api.dnsimple.com"

	// userAgent represents the default user agent used
	// when no other user agent is set.
	defaultUserAgent = "dnsimple-go/" + libraryVersion

	apiVersion = "v2"
)

// Client represents a client to the DNSimple API.
type Client struct {
	// HttpClient is the underlying HTTP client
	// used to communicate with the API.
	HttpClient *http.Client

	// Credentials used for accessing the DNSimple API
	Credentials Credentials

	// BaseURL for API requests.
	// Defaults to the public DNSimple API, but can be set to a different endpoint (e.g. the sandbox).
	BaseURL string

	// UserAgent used when communicating with the DNSimple API.
	UserAgent string

	// Services used for talking to different parts of the DNSimple API.
	Auth     *AuthService
	Contacts *ContactsService
	Domains  *DomainsService
	Zones    *ZonesService

	// Set to true to output debugging logs during API calls
	Debug bool
}

// NewClient returns a new DNSimple API client using the given credentials.
func NewClient(credentials Credentials) *Client {
	c := &Client{Credentials: credentials, HttpClient: &http.Client{}, BaseURL: defaultBaseURL, UserAgent: defaultUserAgent}
	c.Auth = &AuthService{client: c}
	c.Contacts = &ContactsService{client: c}
	c.Domains = &DomainsService{client: c}
	c.Zones = &ZonesService{client: c}
	return c
}

// NewRequest creates an API request.
// The path is expected to be a relative path and will be resolved
// according to the BaseURL of the Client. Paths should always be specified without a preceding slash.
func (client *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	url := client.BaseURL + fmt.Sprintf("/%s/%s", apiVersion, strings.Trim(path, "/"))

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
	req.Header.Add(client.Credentials.HttpHeader())

	return req, nil
}

func (c *Client) get(path string, obj interface{}) (*LegacyResponse, error) {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil, obj)
}

func (c *Client) post(path string, payload, obj interface{}) (*LegacyResponse, error) {
	req, err := c.NewRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

func (c *Client) put(path string, payload, obj interface{}) (*LegacyResponse, error) {
	req, err := c.NewRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

func (c *Client) patch(path string, payload, obj interface{}) (*LegacyResponse, error) {
	req, err := c.NewRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

func (c *Client) delete(path string, payload interface{}, obj interface{}) (*LegacyResponse, error) {
	req, err := c.NewRequest("DELETE", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

// Do sends an API request and returns the API response.
//
// The API response is JSON decoded and stored in the value pointed by obj,
// or returned as an error if an API error has occurred.
// If obj implements the io.Writer interface, the raw response body will be written to obj,
// without attempting to decode it.
func (c *Client) Do(req *http.Request, payload, obj interface{}) (*LegacyResponse, error) {
	if c.Debug {
		log.Printf("Executing request (%v): %#v", req.URL, req)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.Debug {
		log.Printf("Response received: %#v", resp)
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	// If obj implements the io.Writer,
	// the response body is decoded into v.
	if obj != nil {
		if w, ok := obj.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(obj)
		}
	}

	response := &LegacyResponse{HttpResponse: resp}
	return response, err
}

// A Response represents an API response.
type LegacyResponse struct {
	HttpResponse *http.Response // HTTP response
}

// A Response represents an API response.
type Response struct {
	HttpResponse *http.Response // HTTP response
}

type ResponseInterface interface {
	RawData() interface{}
}

// An ErrorResponse represents an API response that generated an error.
type ErrorResponse struct {
	HttpResponse *http.Response // HTTP response that caused this error
	Message      string         `json:"message"` // human-readable message
}

// Error implements the error interface.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %v %v",
		r.HttpResponse.Request.Method, r.HttpResponse.Request.URL,
		r.HttpResponse.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if the status code is different than 2xx. Specific requests
// may have additional requirements, but this is sufficient in most of the cases.
func CheckResponse(r *http.Response) error {
	if code := r.StatusCode; 200 <= code && code <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{HttpResponse: r}
	err := json.NewDecoder(r.Body).Decode(errorResponse)
	if err != nil {
		return err
	}

	return errorResponse
}

// Date custom type.
type Date struct {
	time.Time
}

// UnmarshalJSON handles the deserialization of the custom Date type.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("date should be a string, got %s", data)
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date: %v", err)
	}
	d.Time = t
	return nil
}
