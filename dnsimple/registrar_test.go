package dnsimple

import (
	"io"
	"net/http"
	"testing"
)

func TestRegistrarService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/registration", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/register/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "example.com", "registrant_id": float64(2)}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainAttributes := Domain{Name: "example.com", RegistrantID: 2}

	registrationResponse, err := client.Registrar.Register(accountID, domainAttributes)
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
