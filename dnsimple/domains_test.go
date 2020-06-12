package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDomainPath(t *testing.T) {
	if want, got := "/1010/domains", domainPath("1010", ""); want != got {
		t.Errorf("domainPath(%v) = %v, want %v", "", got, want)
	}

	if want, got := "/1010/domains/example.com", domainPath("1010", "example.com"); want != got {
		t.Errorf("domainPath(%v) = %v, want %v", "example.com", got, want)
	}
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
	if err != nil {
		t.Fatalf("Domains.ListDomains() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}), domainsResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Domains.ListDomains() pagination expected to be %v, got %v", want, got)
	}

	domains := domainsResponse.Data
	if want, got := 2, len(domains); want != got {
		t.Errorf("Domains.ListDomains() expected to return %v contacts, got %v", want, got)
	}

	if want, got := int64(181984), domains[0].ID; want != got {
		t.Fatalf("Domains.ListDomains() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", domains[0].Name; want != got {
		t.Fatalf("Domains.ListDomains() returned Name expected to be `%v`, got `%v`", want, got)
	}

	if want, got := "2021-06-05", domains[0].ExpiresOn; want != got {
		t.Fatalf("Domains.ListDomains() returned ExpiresAt expected to be `%v`, got `%v`", want, got)
	}
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
	if err != nil {
		t.Fatalf("Domains.ListDomains() returned error: %v", err)
	}
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
	if err != nil {
		t.Fatalf("Domains.Create() returned error: %v", err)
	}

	domain := domainResponse.Data
	if want, got := int64(181985), domain.ID; want != got {
		t.Fatalf("Domains.Create() returned ID expected to be `%v`, got `%v`", want, got)
	}
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
	if err != nil {
		t.Errorf("Domains.Get() returned error: %v", err)
	}

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
		ExpiresOn:    "2021-06-05",
		ExpiresAt:    "2021-06-05T02:15:00Z",
		CreatedAt:    "2020-06-04T19:15:14Z",
		UpdatedAt:    "2020-06-04T19:15:21Z"}

	if !reflect.DeepEqual(domain, wantSingle) {
		t.Fatalf("Domains.Get() returned %+v, want %+v", domain, wantSingle)
	}
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
	if err != nil {
		t.Fatalf("Domains.Delete() returned error: %v", err)
	}
}
