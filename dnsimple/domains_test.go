package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDomains_domainPath(t *testing.T) {
	actual := domainPath("1", nil)
	expected := "/1/domains"

	if actual != expected {
		t.Errorf("domainPath(\"1\", nil): actual %s, expected %s", actual, expected)
	}

	actual = domainPath("1", "example.com")
	expected = "/1/domains/example.com"

	if actual != expected {
		t.Errorf("domainPath(\"1\", \"example.com\", nil): actual %s, expected %s", actual, expected)
	}

	actual = domainPath("1", 1)
	expected = "/1/domains/1"

	if actual != expected {
		t.Errorf("domainPath(\"1\", 1, nil): actual %s, expected %s", actual, expected)
	}
}

func TestDomainsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":[{"id":1,"account_id":1010,"registrant_id":null,"name":"example-alpha.com","unicode_name":"example-alpha.com","token":"domain-token","state":"hosted","auto_renew":false,"private_whois":false,"expires_on":null,"created_at":"2014-12-06T15:56:55.573Z","updated_at":"2015-12-09T00:20:56.056Z"},{"id":2,"account_id":1010,"registrant_id":21,"name":"example-beta.com","unicode_name":"example-beta.com","token":"domain-token","state":"registered","auto_renew":false,"private_whois":false,"expires_on":"2015-12-06","created_at":"2014-12-06T15:46:52.411Z","updated_at":"2015-12-09T00:20:53.572Z"}],"pagination":{"current_page":1,"per_page":30,"total_entries":2,"total_pages":1}}
		`)
	})

	accountID := "1010"
	domains, _, err := client.Domains.List(accountID)

	if err != nil {
		t.Fatalf("Domains.List() returned error: %v", err)
	}

	if want, got := 2, len(domains); want != got {
		t.Errorf("Domains.List() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 1, domains[0].ID; want != got {
		t.Fatalf("Domains.List() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", domains[0].Name; want != got {
		t.Fatalf("Domains.List() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "example.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `
			{"data":{"id":1,"account_id":1010,"registrant_id":null,"name":"example-alpha.com","unicode_name":"example-alpha.com","token":"domain-token","state":"hosted","auto_renew":false,"private_whois":false,"expires_on":null,"created_at":"2014-12-06T15:56:55.573Z","updated_at":"2015-12-09T00:20:56.056Z"}}
		`)
	})

	accountID := "1"
	domainAttributes := Domain{Name: "example.com"}
	domain, _, err := client.Domains.Create(accountID, domainAttributes)

	if err != nil {
		t.Fatalf("Domains.Create() returned error: %v", err)
	}

	if want, got := 1, domain.ID; want != got {
		t.Fatalf("Domains.Create() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", domain.Name; want != got {
		t.Fatalf("Domains.Create() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":{"id":1,"account_id":1010,"registrant_id":null,"name":"example-alpha.com","unicode_name":"example-alpha.com","token":"domain-token","state":"hosted","auto_renew":false,"private_whois":false,"expires_on":null,"created_at":"2014-12-06T15:56:55.573Z","updated_at":"2015-12-09T00:20:56.056Z"}}
		`)
	})

	accountID := "1010"
	domain, _, err := client.Domains.Get(accountID, "example.com")

	if err != nil {
		t.Errorf("Domains.Get() returned error: %v", err)
	}

	wantSingle := &Domain{
		ID:           1,
		AccountID:    1010,
		RegistrantID: 0,
		Name:         "example-alpha.com",
		UnicodeName:  "example-alpha.com",
		Token:        "domain-token",
		State:        "hosted",
		PrivateWhois: false,
		ExpiresOn:    "",
		CreatedAt:    "2014-12-06T15:56:55.573Z",
		UpdatedAt:    "2015-12-09T00:20:56.056Z"}

	if !reflect.DeepEqual(domain, wantSingle) {
		t.Fatalf("Domains.Get() returned %+v, want %+v", domain, wantSingle)
	}
}

func TestDomainsService_Delete(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeaders(t, r)
	})

	accountID := "1010"
	_, err := client.Domains.Delete(accountID, "example.com")

	if err != nil {
		t.Fatalf("Domains.Delete() returned error: %v", err)
	}
}
