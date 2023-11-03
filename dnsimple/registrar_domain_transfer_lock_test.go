package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrarService_GetDomainTransferLock(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/101/transfer_lock", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainTransferLock/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.GetDomainTransferLock(context.Background(), "1010", "101")

	assert.NoError(t, err)
	assert.Equal(t, res.Data, &DomainTransferLock{
		Enabled: true,
	})
}

func TestRegistrarService_EnableDomainTransferLock(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/101/transfer_lock", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/enableDomainTransferLock/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.EnableDomainTransferLock(context.Background(), "1010", "101")

	assert.NoError(t, err)
	assert.Equal(t, res.Data, &DomainTransferLock{
		Enabled: true,
	})
}

func TestRegistrarService_DisableDomainTransferLock(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/101/transfer_lock", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/disableDomainTransferLock/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.DisableDomainTransferLock(context.Background(), "1010", "101")

	assert.NoError(t, err)
	assert.Equal(t, res.Data, &DomainTransferLock{
		Enabled: false,
	})
}
