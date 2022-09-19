package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainServices_domainServicesPath(t *testing.T) {
	assert.Equal(t, "/1010/domains/example.com/services", domainServicesPath("1010", "example.com", ""))
	assert.Equal(t, "/1010/domains/example.com/services/1", domainServicesPath("1010", "example.com", "1"))
}

func TestServicesService_AppliedServices(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/services", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/appliedServices/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	servicesResponse, err := client.Services.AppliedServices(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 1}, servicesResponse.Pagination)
	services := servicesResponse.Data
	assert.Len(t, services, 1)
	assert.Equal(t, int64(1), services[0].ID)
	assert.Equal(t, "wordpress", services[0].SID)
}

func TestServicesService_ApplyService(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/services/service1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/applyService/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	settings := DomainServiceSettings{Settings: map[string]string{"app": "foo"}}

	_, err := client.Services.ApplyService(context.Background(), "1010", "service1", "example.com", settings)

	assert.NoError(t, err)
}

func TestServicesService_UnapplyService(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/services/service1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/unapplyService/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Services.UnapplyService(context.Background(), "1010", "service1", "example.com")

	assert.NoError(t, err)
}
