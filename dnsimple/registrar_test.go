package dnsimple

import (
	"io"
	"net/http"
	"testing"
)

func TestRegistrarService_Register(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/registration", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/register/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"registrant_id": float64(2)}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	domainAttributes := Domain{RegistrantID: 2}

	registrationResponse, err := client.Registrar.Register("1010", "example.com", domainAttributes)
	if err != nil {
		t.Fatalf("Registrar.Register() returned error: %v", err)
	}

	domain := registrationResponse.Data
	if want, got := 1, domain.ID; want != got {
		t.Fatalf("Registrar.Register() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example.com", domain.Name; want != got {
		t.Fatalf("Registrar.Register() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_Renew(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/renewal", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/renew/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		//want := map[string]interface{}{}
		//testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	registrationResponse, err := client.Registrar.Renew("1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Registrar.Renew() returned error: %v", err)
	}

	domain := registrationResponse.Data
	if want, got := 1, domain.ID; want != got {
		t.Fatalf("Registrar.Renew() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example.com", domain.Name; want != got {
		t.Fatalf("Registrar.Renew() returned Name expected to be `%v`, got `%v`", want, got)
	}
}
