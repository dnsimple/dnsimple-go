package dnsimple

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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

	body, _ := io.ReadAll(r.Body)

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

	body, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(body, &data)

	assert.NoError(t, err)
	assert.Equal(t, data, values)
}

func readHTTPFixture(t *testing.T, filename string) string {
	data, err := os.ReadFile("../fixtures.http" + filename)
	assert.NoError(t, err)

	// Terrible hack
	// Some fixtures have \n and not \r\n

	// Terrible hack
	s := string(data)
	s = strings.ReplaceAll(s, "Transfer-Encoding: chunked\n", "")
	s = strings.ReplaceAll(s, "Transfer-Encoding: chunked\r\n", "")

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

type badObject struct{}

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

func TestClient_NotFound(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/notfound-certificate.http")

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.makeRequest(context.Background(), "POST", "/", nil, nil, nil)

	var got *ErrorResponse
	assert.ErrorAs(t, err, &got)
	assert.Empty(t, got.AttributeErrors)
}

func TestClient_ValidationError(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/validation-error.http")

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.makeRequest(context.Background(), "POST", "/", nil, nil, nil)

	var got *ErrorResponse
	assert.ErrorAs(t, err, &got)
	want := map[string][]string{
		"address1":       {"can't be blank"},
		"city":           {"can't be blank"},
		"country":        {"can't be blank"},
		"email":          {"can't be blank", "is an invalid email address"},
		"first_name":     {"can't be blank"},
		"last_name":      {"can't be blank"},
		"phone":          {"can't be blank", "is probably not a phone number"},
		"postal_code":    {"can't be blank"},
		"state_province": {"can't be blank"},
	}
	assert.Equal(t, want, got.AttributeErrors)
}
