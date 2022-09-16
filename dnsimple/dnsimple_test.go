package dnsimple

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setupMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(http.DefaultClient)
	client.BaseURL = server.URL
}

func teardownMockServer() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func testQuery(t *testing.T, r *http.Request, want url.Values) {
	assert.Equal(t, want, r.URL.Query())
}

func testHeaders(t *testing.T, r *http.Request) {
	assert.Equal(t, "application/json", r.Header.Get("Accept"))
	assert.Equal(t, defaultUserAgent, r.Header.Get("User-Agent"))
}

func getRequestJSON(r *http.Request) (map[string]interface{}, error) {
	var data map[string]interface{}

	body, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func testRequestJSON(t *testing.T, r *http.Request, values map[string]interface{}) {
	data, err := getRequestJSON(r)

	assert.NoError(t, err)
	assert.Equal(t, data, values)
}

func testRequestJSONArray(t *testing.T, r *http.Request, values []interface{}) {
	var data []interface{}

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &data)

	assert.NoError(t, err)
	assert.Equal(t, data, values)
}

func readHTTPFixture(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile("../fixtures.http" + filename)
	assert.NoError(t, err)

	// Terrible hack
	// Some fixtures have \n and not \r\n

	// Terrible hack
	s := string(data[:])
	s = strings.Replace(s, "Transfer-Encoding: chunked\n", "", -1)
	s = strings.Replace(s, "Transfer-Encoding: chunked\r\n", "", -1)

	return s
}

func httpResponseFixture(t *testing.T, filename string) *http.Response {
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(readHTTPFixture(t, filename))), nil)
	assert.NoError(t, err)
	// resp.Body.Close()
	return resp
}

func TestNewClient(t *testing.T) {
	c := NewClient(http.DefaultClient)

	assert.Equal(t, defaultBaseURL, c.BaseURL)
}

func TestClient_SetUserAgent(t *testing.T) {
	c := NewClient(http.DefaultClient)
	customAgent := "custom-agent/0.1"

	c.SetUserAgent(customAgent)
	assert.Equal(t, "custom-agent/0.1", c.UserAgent)

	req, _ := c.newRequest("GET", "/foo", nil)

	assert.Equal(t, "custom-agent/0.1 "+defaultUserAgent, req.Header.Get("User-Agent"))
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient(http.DefaultClient)
	c.BaseURL = "https://go.example.com"

	inURL, outURL := "/foo", "https://go.example.com/foo"
	req, _ := c.newRequest("GET", inURL, nil)

	assert.Equal(t, outURL, req.URL.String())
	assert.Equal(t, defaultUserAgent, req.Header.Get("User-Agent"))
}

func TestClient_NewRequest_CustomUserAgent(t *testing.T) {
	c := NewClient(http.DefaultClient)
	c.UserAgent = "AwesomeClient"

	req, _ := c.newRequest("GET", "/", nil)

	assert.Equal(t, fmt.Sprintf("AwesomeClient %s", defaultUserAgent), req.Header.Get("User-Agent"))
}

type badObject struct {
}

func (o *badObject) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Bad object is bad")
}

func TestClient_NewRequest_WithBody(t *testing.T) {
	c := NewClient(http.DefaultClient)
	c.BaseURL = "https://go.example.com/"

	inURL, _ := "foo", "https://go.example.com/v2/foo"
	badObject := badObject{}
	_, err := c.newRequest("GET", inURL, &badObject)

	assert.Error(t, err)
}
