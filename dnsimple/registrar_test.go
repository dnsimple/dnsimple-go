package dnsimple

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRegistrarService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/registration", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "example.com", "registrant_id": float64(2)}
		testRequestJSON(t, r, want)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `
			{"data":{"id":1,"account_id":1010,"registrant_id":2,"name":"example.com","unicode_name":"example.com","token":"cc8h1jP8bDLw5rXycL16k8BivcGiT6K9","state":"registered","auto_renew":false,"private_whois":false,"expires_on":"2017-01-16","created_at":"2016-01-16T16:08:50.649Z","updated_at":"2016-01-16T16:09:01.161Z"}}
		`)
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
