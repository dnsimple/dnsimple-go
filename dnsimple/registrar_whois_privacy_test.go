package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrarService_GetWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getWhoisPrivacy/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyResponse, err := client.Registrar.GetWhoisPrivacy(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	privacy := privacyResponse.Data
	wantSingle := &WhoisPrivacy{
		ID:        1,
		DomainID:  2,
		Enabled:   true,
		ExpiresOn: "2017-02-13",
		CreatedAt: "2016-02-13T14:34:50Z",
		UpdatedAt: "2016-02-13T14:34:52Z",
	}
	assert.Equal(t, wantSingle, privacy)
}

func TestRegistrarService_EnableWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/enableWhoisPrivacy/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		// want := map[string]interface{}{}
		// testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyResponse, err := client.Registrar.EnableWhoisPrivacy(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	privacy := privacyResponse.Data
	assert.Equal(t, int64(1), privacy.ID)
}

func TestRegistrarService_DisableWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/disableWhoisPrivacy/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		// want := map[string]interface{}{}
		// testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyResponse, err := client.Registrar.DisableWhoisPrivacy(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	privacy := privacyResponse.Data
	assert.Equal(t, int64(1), privacy.ID)
}

func TestRegistrarService_RenewWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy/renewals", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/renewWhoisPrivacy/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyRenewalResponse, err := client.Registrar.RenewWhoisPrivacy(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	privacyRenewal := privacyRenewalResponse.Data
	assert.Equal(t, int64(1), privacyRenewal.ID)
}
