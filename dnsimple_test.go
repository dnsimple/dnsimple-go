package dnsimple

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *DNSimpleClient
	server *httptest.Server
)

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

	if err := json.Unmarshal([]byte(body), &dat); err != nil {
		t.Errorf("Could not decode json body: %v", err)
	}

	if !reflect.DeepEqual(values, dat) {
		t.Errorf("Request parameters = %v, want %v", dat, values)
	}
}
