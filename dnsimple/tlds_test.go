package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTldsService_ListTlds(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/tlds", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listTlds/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	tldsResponse, err := client.Tlds.ListTlds(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 2, TotalPages: 98, TotalEntries: 195}, tldsResponse.Pagination)
	tlds := tldsResponse.Data
	assert.Len(t, tlds, 2)
	assert.Equal(t, "ac", tlds[0].Tld)
	assert.Equal(t, 1, tlds[0].MinimumRegistration)
	assert.True(t, tlds[0].RegistrationEnabled)
	assert.True(t, tlds[0].RenewalEnabled)
	assert.False(t, tlds[0].TransferEnabled)
}

func TestTldsService_ListTlds_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/tlds", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listTlds/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Tlds.ListTlds(context.Background(), &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestTldsService_GetTld(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/tlds/com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getTld/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	tldResponse, err := client.Tlds.GetTld(context.Background(), "com")

	assert.NoError(t, err)
	tld := tldResponse.Data
	assert.Equal(t, "com", tld.Tld)
	assert.Equal(t, 1, tld.TldType)
	assert.True(t, tld.WhoisPrivacy)
	assert.False(t, tld.AutoRenewOnly)
	assert.Equal(t, 1, tld.MinimumRegistration)
	assert.True(t, tld.RegistrationEnabled)
	assert.True(t, tld.RenewalEnabled)
	assert.True(t, tld.TransferEnabled)
	assert.Equal(t, "ds", tld.DnssecInterfaceType)
}

func TestTldsService_GetTldExtendedAttributes(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/tlds/com/extended_attributes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getTldExtendedAttributes/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	tldResponse, err := client.Tlds.GetTldExtendedAttributes(context.Background(), "com")

	assert.NoError(t, err)
	attributes := tldResponse.Data
	assert.Len(t, attributes, 4)
	assert.Equal(t, "uk_legal_type", attributes[0].Name)
}
