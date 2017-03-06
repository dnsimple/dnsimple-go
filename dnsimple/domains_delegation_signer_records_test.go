package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDelegationSignerRecordPath(t *testing.T) {
	if want, got := "/1010/domains/example.com/ds_records", delegationSignerRecordPath("1010", "example.com", 0); want != got {
		t.Errorf("delegationSignerRecordPath(%v) = %v, want %v", "", got, want)
	}

	if want, got := "/1010/domains/example.com/ds_records/2", delegationSignerRecordPath("1010", "example.com", 2); want != got {
		t.Errorf("delegationSignerRecordPath(%v) = %v, want %v", "2", got, want)
	}
}

func TestDomainsService_ListDelegationSignerRecords(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listDelegationSignerRecords/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	dsRecordsResponse, err := client.Domains.ListDelegationSignerRecords("1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Domains.ListDelegationSignerRecords() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 1}), dsRecordsResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Domains.ListDelegationSignerRecords() pagination expected to be %v, got %v", want, got)
	}

	dsRecords := dsRecordsResponse.Data
	if want, got := 1, len(dsRecords); want != got {
		t.Errorf("Domains.ListDelegationSignerRecords() expected to return %v delegation signer records, got %v", want, got)
	}

	if want, got := 24, dsRecords[0].ID; want != got {
		t.Fatalf("Domains.ListDelegationSignerRecords() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "8", dsRecords[0].Algorithm; want != got {
		t.Fatalf("Domains.ListDelegationSignerRecords() returned Algorithm expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_ListDelegationSignerRecords_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/ds_records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listDelegationSignerRecords/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListDelegationSignerRecords("1010", "example.com", &ListOptions{Page: 2, PerPage: 20})
	if err != nil {
		t.Fatalf("Domains.ListDelegationSignerRecords() returned error: %v", err)
	}
}
