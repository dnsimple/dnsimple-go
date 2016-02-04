package dnsimple

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux *http.ServeMux
	client *Client
	server *httptest.Server
)

func setupMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(NewOauthTokenCredentials("dnsimple-token"))
	client.BaseURL = server.URL + "/"
}

func teardownMockServer() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func testRequestJSON(t *testing.T, r *http.Request, values map[string]interface{}) {
	var dat map[string]interface{}

	body, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		t.Errorf("Could not decode json body: %v", err)
	}

	if !reflect.DeepEqual(values, dat) {
		t.Errorf("Request parameters = %v, want %v", dat, values)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if !reflect.DeepEqual(want, r.Form) {
		t.Errorf("Request parameters = %v, want %v", r.Form, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(NewOauthTokenCredentials("mytoken"))

	if c.BaseURL != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL, defaultBaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(NewOauthTokenCredentials("mytoken"))
	c.BaseURL = "https://go.example.com/"

	inURL, outURL := "foo", "https://go.example.com/v2/foo"
	req, _ := c.NewRequest("GET", inURL, nil)

	// test that relative URL was expanded with the proper BaseURL
	if req.URL.String() != outURL {
		t.Errorf("makeRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}

	// test that default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("makeRequest() User-Agent = %v, want %v", userAgent, c.UserAgent)
	}
}

type badObject struct {
}

func (o *badObject) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Bad object is bad")
}

func TestNewRequestWithBody(t *testing.T) {
	c := NewClient(NewOauthTokenCredentials("mytoken"))
	c.BaseURL = "https://go.example.com/"

	inURL, _ := "foo", "https://go.example.com/v2/foo"
	badObject := badObject{}
	_, err := c.NewRequest("GET", inURL, &badObject)

	if err == nil {
		t.Errorf("NewRequest with body expected error with blank string")
	}
}
