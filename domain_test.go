package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDomain_domainPath(t *testing.T) {
	var path string

	path = domainPath(nil)
	if path != "domains" {
		t.Errorf("domainPath returned %+v, want domains", path)
	}

	path = domainPath("example.com")
	if path != "domains/example.com" {
		t.Errorf("domainPath returned %+v, want domains/example.com", path)
	}

	path = domainPath(42)
	if path != "domains/42" {
		t.Errorf("domainPath returned %+v, want domains/42", path)
	}

	path = domainPath(Domain{Id: 64})
	if path != "domains/64" {
		t.Errorf("domainPath returned %+v, want domains/64", path)
	}

	path = domainPath(Record{DomainId: 23})
	if path != "domains/23" {
		t.Errorf("domainPath returned %+v, want domains/23", path)
	}
}

func TestDomain_Domains(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"domain":{"id": 1, "name":"foo.com"}}]`)
	})

	domains, err := client.Domains()

	if err != nil {
		t.Errorf("Domains returned error: %v", err)
	}

	want := []Domain{{Id: 1, Name: "foo.com"}}
	if !reflect.DeepEqual(domains, want) {
		t.Errorf("Domains returned %+v, want %+v", domains, want)
	}
}

func TestDomain_Domain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"domain": {"id":1, "name":"example.com"}}`)
	})

	domain, err := client.Domain("example.com")

	if err != nil {
		t.Errorf("Domain returned error: %v", err)
	}

	want := Domain{Id: 1, Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Errorf("Domains returned %+v, want %+v", domain, want)
	}
}

func TestDomain_DomainAvailable_available(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"name":"example.com", "status":"available"}`)
	})

	available, err := client.DomainAvailable("example.com")

	if err != nil {
		t.Errorf("Availability check returned %v", err)
	}

	if !available {
		t.Errorf("Availability returned false, want true")
	}
}

func TestDomain_DomainAavailable_unavailable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"name":"example.com", "status":"unavailable"}`)
	})

	available, err := client.DomainAvailable("example.com")

	if err != nil {
		t.Errorf("Availability check returned %v", err)
	}

	if available {
		t.Errorf("Availability returned true, want false")
	}
}

func TestDomain_SetAutoRenew_enable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/example.com/auto_renewal", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	err := client.SetAutorenew("example.com", true)

	if err != nil {
		t.Errorf("Autorenew (enable) returned %v", err)
	}
}

func TestDomain_SetAutoRenew_disable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/example.com/auto_renewal", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.SetAutorenew("example.com", false)

	if err != nil {
		t.Errorf("Autorenew (disable) returned %v", err)
	}
}

func TestDomain_Renew(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domain_renewals", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["domain"] = map[string]interface{}{"name": "example.com", "renew_whois_privacy": true}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"domain":{"name":"example.com"}}`)
	})

	err := client.Renew("example.com", true)

	if err != nil {
		t.Errorf("Renew returned %v", err)
	}
}
