package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVanityNameServers_vanityNameServerPath(t *testing.T) {
	assert.Equal(t, "/1010/vanity/example.com", vanityNameServerPath("1010", "example.com"))
}

func TestVanityNameServersService_EnableVanityNameServers(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/vanity/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/enableVanityNameServers/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	vanityNameServerResponse, err := client.VanityNameServers.EnableVanityNameServers(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	delegation := vanityNameServerResponse.Data[0].Name
	wantSingle := "ns1.example.com"
	assert.Equal(t, wantSingle, delegation)
}

func TestVanityNameServersService_DisableVanityNameServers(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/vanity/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/disableVanityNameServers/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.VanityNameServers.DisableVanityNameServers(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
}
