package dnsimple

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestRegistrarService_GetWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getWhoisPrivacy/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyResponse, err := client.Registrar.GetWhoisPrivacy(context.Background(), "1010", "example.com")
	if err != nil {
		t.Errorf("Registrar.GetWhoisPrivacy() returned error: %v", err)
	}

	privacy := privacyResponse.Data
	wantSingle := &WhoisPrivacy{
		ID:        1,
		DomainID:  2,
		Enabled:   true,
		ExpiresOn: "2017-02-13",
		CreatedAt: "2016-02-13T14:34:50Z",
		UpdatedAt: "2016-02-13T14:34:52Z"}

	if !reflect.DeepEqual(privacy, wantSingle) {
		t.Fatalf("Registrar.GetWhoisPrivacy() returned %+v, want %+v", privacy, wantSingle)
	}
}

func TestRegistrarService_EnableWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/enableWhoisPrivacy/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		//want := map[string]interface{}{}
		//testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyResponse, err := client.Registrar.EnableWhoisPrivacy(context.Background(), "1010", "example.com")
	if err != nil {
		t.Errorf("Registrar.EnableWhoisPrivacy() returned error: %v", err)
	}

	privacy := privacyResponse.Data
	if want, got := int64(1), privacy.ID; want != got {
		t.Fatalf("Registrar.EnableWhoisPrivacy() returned ID expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_DisableWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/disableWhoisPrivacy/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		//want := map[string]interface{}{}
		//testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyResponse, err := client.Registrar.DisableWhoisPrivacy(context.Background(), "1010", "example.com")
	if err != nil {
		t.Errorf("Registrar.DisableWhoisPrivacy() returned error: %v", err)
	}

	privacy := privacyResponse.Data
	if want, got := int64(1), privacy.ID; want != got {
		t.Fatalf("Registrar.DisableWhoisPrivacy() returned ID expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_RenewWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy/renewals", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/renewWhoisPrivacy/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	privacyRenewalResponse, err := client.Registrar.RenewWhoisPrivacy(context.Background(), "1010", "example.com")
	if err != nil {
		t.Errorf("Registrar.RenewWhoisPrivacy() returned error: %v", err)
	}

	privacyRenewal := privacyRenewalResponse.Data
	if want, got := int64(1), privacyRenewal.ID; want != got {
		t.Fatalf("Registrar.RenewWhoisPrivacy() returned ID expected to be `%v`, got `%v`", want, got)
	}
}
