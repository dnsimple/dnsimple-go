package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDomainServices_domainServicesPath(t *testing.T) {
	if want, got := "/1010/domains/example.com/services", domainServicesPath("1010", "example.com", ""); want != got {
		t.Errorf("domainServicesPath(%v, %v, ) = %v, want %v", "1010", "example.com", got, want)
	}

	if want, got := "/1010/domains/example.com/services/1", domainServicesPath("1010", "example.com", "1"); want != got {
		t.Errorf("domainServicesPath(%v, %v, 1) = %v, want %v", "1010", "example.com", got, want)
	}
}

func TestDomainServicesService_AppliedServices(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/services", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/appliedServices/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	servicesResponse, err := client.DomainServices.AppliedServices("1010", "example.com", nil)
	if err != nil {
		t.Fatalf("DomainServices.AppliedServices() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 1}), servicesResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("DomainServices.AppliedServices() pagination expected to be %v, got %v", want, got)
	}

	services := servicesResponse.Data
	if want, got := 1, len(services); want != got {
		t.Errorf("DomainServices.AppliedServices() expected to return %v services, got %v", want, got)
	}

	if want, got := 1, services[0].ID; want != got {
		t.Fatalf("DomainServices.AppliedServices() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "wordpress", services[0].ShortName; want != got {
		t.Fatalf("DomainServices.AppliedServices() returned ShortName expected to be `%v`, got `%v`", want, got)
	}
}
