package dnsimple

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

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
