package dnsimple

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setupMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(NewOauthTokenCredentials("dnsimple-token"))
	client.BaseURL = server.URL
}

func teardownMockServer() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; want != got {
		t.Errorf("Request METHOD expected to be `%v`, got `%v`", want, got)
	}
}

func testHeader(t *testing.T, r *http.Request, name, want string) {
	if got := r.Header.Get(name); want != got {
		t.Errorf("Request() %v expected to be `%v`, got `%v`", name, want, got)
	}
}

func testHeaders(t *testing.T, r *http.Request) {
	testHeader(t, r, "Accept", "application/json")
	testHeader(t, r, "User-Agent", defaultUserAgent)
}

func testRequestJSON(t *testing.T, r *http.Request, values map[string]interface{}) {
	var data map[string]interface{}

	body, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		t.Errorf("Could not decode json body: %v", err)
	}

	if !reflect.DeepEqual(values, data) {
		t.Errorf("Request parameters = %v, want %v", data, values)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(NewOauthTokenCredentials("dnsimple-token"))

	if c.BaseURL != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL, defaultBaseURL)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient(NewOauthTokenCredentials("dnsimple-token"))
	c.BaseURL = "https://go.example.com"

	inURL, outURL := "/foo", "https://go.example.com/v2/foo"
	req, _ := c.NewRequest("GET", inURL, nil)

	// test that relative URL was expanded with the proper BaseURL
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}

	// test that default user-agent is attached to the request
	ua := req.Header.Get("User-Agent")
	if ua != c.UserAgent {
		t.Errorf("NewRequest() User-Agent = %v, want %v", ua, c.UserAgent)
	}
}

type badObject struct {
}

func (o *badObject) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Bad object is bad")
}

func TestClient_NewRequest_WithBody(t *testing.T) {
	c := NewClient(NewOauthTokenCredentials("dnsimple-token"))
	c.BaseURL = "https://go.example.com/"

	inURL, _ := "foo", "https://go.example.com/v2/foo"
	badObject := badObject{}
	_, err := c.NewRequest("GET", inURL, &badObject)

	if err == nil {
		t.Errorf("NewRequest with body expected error with blank string")
	}
}
