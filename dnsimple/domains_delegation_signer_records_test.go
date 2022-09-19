package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelegationSignerRecordPath(t *testing.T) {
	assert.Equal(t, "/1010/domains/example.com/ds_records", delegationSignerRecordPath("1010", "example.com", 0))
	assert.Equal(t, "/1010/domains/example.com/ds_records/2", delegationSignerRecordPath("1010", "example.com", 2))
}

func TestDomainsService_ListDelegationSignerRecords(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listDelegationSignerRecords/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	dsRecordsResponse, err := client.Domains.ListDelegationSignerRecords(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 1}, dsRecordsResponse.Pagination)
	dsRecords := dsRecordsResponse.Data
	assert.Len(t, dsRecords, 1)
	assert.Equal(t, int64(24), dsRecords[0].ID)
	assert.Equal(t, "8", dsRecords[0].Algorithm)
}

func TestDomainsService_ListDelegationSignerRecords_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listDelegationSignerRecords/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListDelegationSignerRecords(context.Background(), "1010", "example.com", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestDomainsService_CreateDelegationSignerRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createDelegationSignerRecord/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"algorithm": "13", "digest": "ABC123", "digest_type": "2", "keytag": "1234"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	dsRecordAttributes := DelegationSignerRecord{Algorithm: "13", Digest: "ABC123", DigestType: "2", Keytag: "1234"}

	dsRecordResponse, err := client.Domains.CreateDelegationSignerRecord(context.Background(), "1010", "example.com", dsRecordAttributes)

	assert.NoError(t, err)
	dsRecord := dsRecordResponse.Data
	assert.Equal(t, int64(2), dsRecord.ID)
	assert.Equal(t, "13", dsRecord.Algorithm)
}

func TestDomainsService_GetDelegationSignerRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDelegationSignerRecord/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	dsRecordResponse, err := client.Domains.GetDelegationSignerRecord(context.Background(), "1010", "example.com", 2)

	assert.NoError(t, err)
	dsRecord := dsRecordResponse.Data
	wantSingle := &DelegationSignerRecord{
		ID:         24,
		DomainID:   1010,
		Algorithm:  "8",
		DigestType: "2",
		Digest:     "C1F6E04A5A61FBF65BF9DC8294C363CF11C89E802D926BDAB79C55D27BEFA94F",
		Keytag:     "44620",
		PublicKey:  "",
		CreatedAt:  "2017-03-03T13:49:58Z",
		UpdatedAt:  "2017-03-03T13:49:58Z"}
	assert.Equal(t, wantSingle, dsRecord)
}

func TestDomainsService_DeleteDelegationSignerRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteDelegationSignerRecord/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.DeleteDelegationSignerRecord(context.Background(), "1010", "example.com", 2)

	assert.NoError(t, err)
}
