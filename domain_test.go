package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDomain_domainPath(t *testing.T) {
	var pathTests = []struct {
		input    interface{}
		expected string
	}{
		{nil, "domains"},
		{"example.com", "domains/example.com"},
		{42, "domains/42"},
		{Domain{Id: 64}, "domains/64"},
		{Record{DomainId: 23}, "domains/23"},
	}

	for _, pt := range pathTests {
		actual := domainPath(pt.input)
		if actual != pt.expected {
			t.Errorf("domainPath(%+v): expected %s, actual %s", pt.input, pt.expected)
		}
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
