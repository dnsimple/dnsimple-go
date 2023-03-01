package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnssecPath(t *testing.T) {
	t.Run("empty account id", func(t *testing.T) {
		path, err := dnssecPath("", "example.com")
		assert.Error(t, err)
		assert.Empty(t, path)
	})

	t.Run("empty domain identifier", func(t *testing.T) {
		path, err := dnssecPath("1010", "")
		assert.Error(t, err)
		assert.Empty(t, path)
	})

	t.Run("success", func(t *testing.T) {
		path, err := dnssecPath("1010", "example.com")
		assert.NoError(t, err)
		assert.Equal(t, "/1010/domains/example.com/dnssec", path)
	})
}

func TestDomainsService_EnableDnssec(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/dnssec", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/enableDnssec/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
	})

	accountID := "1010"

	_, err := client.Domains.EnableDnssec(context.Background(), accountID, "example.com")

	assert.NoError(t, err)
}

func TestDomainsService_DisableDnssec(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/dnssec", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/disableDnssec/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
	})

	accountID := "1010"

	_, err := client.Domains.DisableDnssec(context.Background(), accountID, "example.com")

	assert.NoError(t, err)
}

func TestDomainsService_GetDnssec(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/dnssec", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDnssec/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	dnssecResponse, err := client.Domains.GetDnssec(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	assert.Equal(t, &Dnssec{Enabled: true}, dnssecResponse.Data)
}
