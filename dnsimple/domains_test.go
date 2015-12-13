package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDomains_domainPath(t *testing.T) {
	actual := domainPath("1", nil)
	expected := "1/domains"

	if actual != expected {
		t.Errorf("domainPath(\"1\", nil): actual %s, expected %s", actual, expected)
	}

	actual = domainPath("1", "example.com")
	expected = "1/domains/example.com"

	if actual != expected {
		t.Errorf("domainPath(\"1\", \"example.com\", nil): actual %s, expected %s", actual, expected)
	}

	actual = domainPath("1", 1)
	expected = "1/domains/1"

	if actual != expected {
		t.Errorf("domainPath(\"1\", 1, nil): actual %s, expected %s", actual, expected)
	}
}

func TestDomainsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/1/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data":[{"id": 1, "name":"example.com"}]}`)
	})

	accountId := "1"
	domains, _, err := client.Domains.List(accountId)

	if err != nil {
		t.Errorf("Domains.List returned error: %v", err)
	}

	want := []Domain{{Id: 1, Name: "example.com"}}
	if !reflect.DeepEqual(domains, want) {
		t.Errorf("Domains.List returned %+v, want %+v", domains, want)
	}
}

func TestDomainsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/1/domains", func(w http.ResponseWriter, r *http.Request) {
		want := map[string]interface{}{"name": "example.com"}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"data":{"id":1, "name":"example.com"}}`)
	})

	accountId := "1"
	domainId := 1
	domainValues := Domain{Name: "example.com"}
	domain, _, err := client.Domains.Create(accountId, domainValues)

	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}

	want := Domain{Id: domainId, Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Fatalf("Domains.Create returned %+v, want %+v", domain, want)
	}
}

func TestDomainsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/1/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data": {"id":1, "name":"example.com"}}`)
	})

	accountId := "1"
	domain, _, err := client.Domains.Get(accountId, "example.com")

	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	want := Domain{Id: 1, Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Fatalf("Domains.Get returned %+v, want %+v", domain, want)
	}
}

func TestDomainsService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/1/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		// fmt.Fprint(w, `{}`)
	})

	accountId := "1"
	_, err := client.Domains.Delete(accountId, "example.com")

	if err != nil {
		t.Errorf("Domains.Delete returned error: %v", err)
	}
}
