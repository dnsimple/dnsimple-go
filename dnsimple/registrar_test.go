package dnsimple

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestRegistrarService_IsAvailable_available(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"name":"example.com", "status":"available"}`)
	})

	available, err := client.Registrar.IsAvailable("example.com")

	if err != nil {
		t.Errorf("Registrar.IsAvailable check returned %v", err)
	}

	if !available {
		t.Errorf("Registrar.IsAvailable returned false, want true")
	}
}

func TestRegistrarService_IsAvailable_unavailable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"name":"example.com", "status":"unavailable"}`)
	})

	available, err := client.Registrar.IsAvailable("example.com")

	if err != nil {
		t.Errorf("Domains.IsAvailable check returned %v", err)
	}

	if available {
		t.Errorf("Domains.IsAvailable returned true, want false")
	}
}

func TestRegistrarService_IsAvailable_failed400(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"Invalid request"}`)
	})

	_, err := client.Registrar.IsAvailable("example.com")

	if err == nil {
		t.Errorf("Domains.IsAvailable expected error to be returned")
	}

	if match := "400 Invalid request"; !strings.Contains(err.Error(), match) {
		t.Errorf("Domains.IsAvailable returned %+v, should match %+v", err, match)
	}
}

func TestRegistrarService_Renew(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domain_renewals", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["domain"] = map[string]interface{}{"name": "example.com", "renew_whois_privacy": true}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"domain":{"name":"example.com"}}`)
	})

	_, err := client.Registrar.Renew("example.com", true)

	if err != nil {
		t.Errorf("Domains.Renew returned %v", err)
	}
}
