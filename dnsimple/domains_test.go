package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
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
	}

	for _, pt := range pathTests {
		actual := domainPath(pt.input)
		if actual != pt.expected {
			t.Errorf("domainPath(%+v): expected %s, actual %s", pt.input, pt.expected)
		}
	}
}

func TestDomainsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"domain":{"id": 1, "name":"foo.com"}}]`)
	})

	domains, _, err := client.Domains.List()

	if err != nil {
		t.Errorf("Domains.List returned error: %v", err)
	}

	want := []Domain{{Id: 1, Name: "foo.com"}}
	if !reflect.DeepEqual(domains, want) {
		t.Errorf("Domains.List returned %+v, want %+v", domains, want)
	}
}

func TestDomainsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["domain"] = map[string]interface{}{"name": "example.com"}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprintf(w, `{"domain":{"id":1, "name":"example.com"}}`)
	})

	domainValues := Domain{Name: "example.com"}
	domain, _, err := client.Domains.Create(domainValues)

	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}

	want := Domain{Id: 1, Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Errorf("Domains.Create returned %+v, want %+v", domain, want)
	}
}

func TestDomainsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"domain": {"id":1, "name":"example.com"}}`)
	})

	domain, _, err := client.Domains.Get("example.com")

	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	want := Domain{Id: 1, Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Errorf("Domains.Get returned %+v, want %+v", domain, want)
	}
}

func TestDomainsService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		// fmt.Fprint(w, `{}`)
	})

	_, err := client.Domains.Delete("example.com")

	if err != nil {
		t.Errorf("Domains.Delete returned error: %v", err)
	}
}

func TestDomainsService_CheckAvailability_available(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"name":"example.com", "status":"available"}`)
	})

	available, _, err := client.Domains.CheckAvailability("example.com")

	if err != nil {
		t.Errorf("Domains.CheckAvailability check returned %v", err)
	}

	if !available {
		t.Errorf("Domains.CheckAvailability returned false, want true")
	}
}

func TestDomainsService_CheckAvailability_unavailable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"name":"example.com", "status":"unavailable"}`)
	})

	available, _, err := client.Domains.CheckAvailability("example.com")

	if err != nil {
		t.Errorf("Domains.CheckAvailability check returned %v", err)
	}

	if available {
		t.Errorf("Domains.CheckAvailability returned true, want false")
	}
}

func TestDomainsService_CheckAvailability_failed400(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"Invalid request"}`)
	})

	_, _, err := client.Domains.CheckAvailability("example.com")

	if err == nil {
		t.Errorf("Domains.CheckAvailability expected error to be returned")
	}

	if match := "400 Invalid request"; !strings.Contains(err.Error(), match) {
		t.Errorf("Domains.CheckAvailability returned %+v, should match %+v", err, match)
	}
}

func TestDomainsService_Renew(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domain_renewals", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["domain"] = map[string]interface{}{"name": "example.com", "renew_whois_privacy": true}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"domain":{"name":"example.com"}}`)
	})

	_, err := client.Domains.Renew("example.com", true)

	if err != nil {
		t.Errorf("Domains.Renew returned %v", err)
	}
}
