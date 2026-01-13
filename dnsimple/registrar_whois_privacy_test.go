package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
