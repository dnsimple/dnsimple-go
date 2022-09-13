package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainPushPath(t *testing.T) {
	assert.Equal(t, "/1010/pushes", domainPushPath("1010", 0))
	assert.Equal(t, "/1010/pushes/1", domainPushPath("1010", 1))
}

func TestDomainsService_InitiatePush(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/pushes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/initiatePush/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"new_account_email": "admin@target-account.test"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	pushAttributes := DomainPushAttributes{NewAccountEmail: "admin@target-account.test"}

	pushResponse, err := client.Domains.InitiatePush(context.Background(), "1010", "example.com", pushAttributes)

	assert.NoError(t, err)
	push := pushResponse.Data
	assert.Equal(t, int64(1), push.ID)
	assert.Equal(t, int64(2020), push.AccountID)
}

func TestDomainsService_DomainsPushesList(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listPushes/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	pushesResponse, err := client.Domains.ListPushes(context.Background(), "2020", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, pushesResponse.Pagination)
	pushes := pushesResponse.Data
	assert.Len(t, pushes, 2)
	assert.Equal(t, int64(1), pushes[0].ID)
	assert.Equal(t, int64(2020), pushes[0].AccountID)
}

func TestDomainsService_DomainsPushesList_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listPushes/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListPushes(context.Background(), "2020", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestDomainsService_AcceptPush(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/acceptPush/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"contact_id": float64(2)}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	pushAttributes := DomainPushAttributes{ContactID: 2}

	_, err := client.Domains.AcceptPush(context.Background(), "2020", 1, pushAttributes)

	assert.NoError(t, err)
}

func TestDomainsService_RejectPush(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/rejectPush/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.RejectPush(context.Background(), "2020", 1)

	assert.NoError(t, err)
}
