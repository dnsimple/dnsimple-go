package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnssecPath(t *testing.T) {
	assert.Equal(t, "/1010/domains/example.com/dnssec", dnssecPath("1010", "example.com"))
}

func TestDomainsService_EnableDnssec(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/dnssec", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/enableDnssec/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	res, err := client.Domains.EnableDnssec(context.Background(), accountID, "example.com")

	assert.NoError(t, err)
	assert.Equal(t, &Dnssec{Enabled: true}, res.Data)
}

func TestDomainsService_DisableDnssec(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/dnssec", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/disableDnssec/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
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
