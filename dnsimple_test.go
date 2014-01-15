package dnsimple

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *DNSimpleClient
	server *httptest.Server
)

// This method of testing http client APIs is borrowed from
// Will Norris's work in go-github @ https://github.com/google/go-github
func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient("mytoken", "me@example.com")
	client.BaseURL = server.URL + "/"
}

func teardown() {
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

func testString(t *testing.T, test, value, want string) {
	if value != want {
		t.Errorf("%s returned %+v, want %+v", test, value, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("mytoken", "me@example.com")

	if c.BaseURL != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL, defaultBaseURL)
	}
}

func TestMakeRequest(t *testing.T) {
	c := NewClient("mytoken", "me@example.com")
	c.BaseURL = "https://go.example.com/"

	inURL, outURL := "foo", "https://go.example.com/foo"
	req, _ := c.makeRequest("GET", inURL, nil)

	// test that relative URL was expanded with the proper BaseURL
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}
}
