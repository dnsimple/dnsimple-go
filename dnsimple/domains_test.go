package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainPath(t *testing.T) {
	assert.Equal(t, "/1010/domains", domainPath("1010", ""))
	assert.Equal(t, "/1010/domains/example.com", domainPath("1010", "example.com"))
}

func TestDomainsService_ListDomains(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1385/domains", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listDomains/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	domainsResponse, err := client.Domains.ListDomains(context.Background(), "1385", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, domainsResponse.Pagination)
	domains := domainsResponse.Data
	assert.Len(t, domains, 2)
	assert.Equal(t, int64(181984), domains[0].ID)
	assert.Equal(t, "example-alpha.com", domains[0].Name)
	assert.Equal(t, "2021-06-05T02:15:00Z", domains[0].ExpiresAt)
}

func TestDomainsService_ListDomains_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listDomains/success.http")

		testQuery(t, r, url.Values{
			"page":          []string{"2"},
			"per_page":      []string{"20"},
			"sort":          []string{"name,expiration:desc"},
			"name_like":     []string{"example"},
			"registrant_id": []string{"10"},
		})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListDomains(context.Background(), "1010", &DomainListOptions{NameLike: String("example"), RegistrantID: Int(10), ListOptions: ListOptions{Page: Int(2), PerPage: Int(20), Sort: String("name,expiration:desc")}})

	assert.NoError(t, err)
}

func TestDomainsService_CreateDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1385/domains", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createDomain/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "example-beta.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1385"
	domainAttributes := Domain{Name: "example-beta.com"}

	domainResponse, err := client.Domains.CreateDomain(context.Background(), accountID, domainAttributes)

	assert.NoError(t, err)
	domain := domainResponse.Data
	assert.Equal(t, int64(181985), domain.ID)
}

func TestDomainsService_GetDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example-alpha.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomain/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	domainResponse, err := client.Domains.GetDomain(context.Background(), accountID, "example-alpha.com")

	assert.NoError(t, err)
	domain := domainResponse.Data
	wantSingle := &Domain{
		ID:           181984,
		AccountID:    1385,
		RegistrantID: 2715,
		Name:         "example-alpha.com",
		UnicodeName:  "example-alpha.com",
		Token:        "",
		State:        "registered",
		AutoRenew:    false,
		PrivateWhois: false,
		ExpiresAt:    "2021-06-05T02:15:00Z",
		CreatedAt:    "2020-06-04T19:15:14Z",
		UpdatedAt:    "2020-06-04T19:15:21Z"}
	assert.Equal(t, wantSingle, domain)
}

func TestDomainsService_DeleteDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteDomain/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	_, err := client.Domains.DeleteDomain(context.Background(), accountID, "example.com")

	assert.NoError(t, err)
}
