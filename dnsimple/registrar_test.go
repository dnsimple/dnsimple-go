package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestRegistrarService_CheckDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkDomain/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	checkResponse, err := client.Registrar.CheckDomain(context.Background(), "1010", "example.com")
	if err != nil {
		t.Fatalf("Registrar.CheckDomain() returned error: %v", err)
	}

	check := checkResponse.Data
	if want, got := "ruby.codes", check.Domain; want != got {
		t.Fatalf("Registrar.CheckDomain() returned Domain expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_GetDomainPremiumPrice(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/premium_price", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainPremiumPrice/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	priceResponse, err := client.Registrar.GetDomainPremiumPrice(context.Background(), "1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Registrar.GetDomainPremiumPrice() returned error: %v", err)
	}

	price := priceResponse.Data
	if want, got := "109.00", price.PremiumPrice; want != got {
		t.Fatalf("Registrar.GetDomainPremiumPrice() returned Domain expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_GetDomainPremiumPrice_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/premium_price", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainPremiumPrice/success.http")

		testQuery(t, r, url.Values{
			"action": []string{"registration"},
		})

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.GetDomainPremiumPrice(context.Background(), "1010", "example.com", &DomainPremiumPriceOptions{Action: "registration"})
	if err != nil {
		t.Fatalf("Registrar.GetDomainPremiumPrice() returned error: %v", err)
	}
}

func TestRegistrarService_RegisterDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/registrations", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/registerDomain/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"registrant_id": float64(2)}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	registerRequest := &DomainRegisterRequest{RegistrantID: 2}

	registrationResponse, err := client.Registrar.RegisterDomain(context.Background(), "1010", "example.com", registerRequest)
	if err != nil {
		t.Fatalf("Registrar.RegisterDomain() returned error: %v", err)
	}

	registration := registrationResponse.Data
	if want, got := 1, registration.ID; want != got {
		t.Fatalf("Registrar.RegisterDomain() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := 999, registration.DomainID; want != got {
		t.Fatalf("Registrar.RegisterDomain() returned DomainID expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_RegisterDomain_ExtendedAttributes(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/registrations", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/registerDomain/success.http")

		data, _ := getRequestJSON(r)

		if want, got := map[string]interface{}{"att1": "val1", "att2": "val2"}, data["extended_attributes"]; !reflect.DeepEqual(want, got) {
			t.Errorf("RegisterDomain() incorrect extended attributes payload, expected `%v`, got `%v`", want, got)
		}

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	registerRequest := &DomainRegisterRequest{RegistrantID: 2, ExtendedAttributes: map[string]string{"att1": "val1", "att2": "val2"}}

	if _, err := client.Registrar.RegisterDomain(context.Background(), "1010", "example.com", registerRequest); err != nil {
		t.Fatalf("Registrar.RegisterDomain() returned error: %v", err)
	}
}

func TestRegistrarService_TransferDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/transfers", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/transferDomain/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"registrant_id": float64(2), "auth_code": "x1y2z3"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	transferRequest := &DomainTransferRequest{RegistrantID: 2, AuthCode: "x1y2z3"}

	transferResponse, err := client.Registrar.TransferDomain(context.Background(), "1010", "example.com", transferRequest)
	if err != nil {
		t.Fatalf("Registrar.TransferDomain() returned error: %v", err)
	}

	transfer := transferResponse.Data
	if want, got := 1, transfer.ID; want != got {
		t.Fatalf("Registrar.TransferDomain() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := 999, transfer.DomainID; want != got {
		t.Fatalf("Registrar.TransferDomain() returned DomainID expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_TransferDomain_ExtendedAttributes(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/transfers", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/transferDomain/success.http")

		data, _ := getRequestJSON(r)

		if want, got := map[string]interface{}{"att1": "val1", "att2": "val2"}, data["extended_attributes"]; !reflect.DeepEqual(want, got) {
			t.Errorf("TransferDomain() incorrect extended attributes payload, expected `%v`, got `%v`", want, got)
		}

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	transferRequest := &DomainTransferRequest{RegistrantID: 2, AuthCode: "x1y2z3", ExtendedAttributes: map[string]string{"att1": "val1", "att2": "val2"}}

	if _, err := client.Registrar.TransferDomain(context.Background(), "1010", "example.com", transferRequest); err != nil {
		t.Fatalf("Registrar.TransferDomain() returned error: %v", err)
	}
}

func TestRegistrarService_TransferDomainOut(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/authorize_transfer_out", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/authorizeDomainTransferOut/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.TransferDomainOut(context.Background(), "1010", "example.com")
	if err != nil {
		t.Fatalf("Registrar.TransferOut() returned error: %v", err)
	}
}

func TestRegistrarService_RenewDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/renewals", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/renewDomain/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		//want := map[string]interface{}{}
		//testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	renewalResponse, err := client.Registrar.RenewDomain(context.Background(), "1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Registrar.RenewDomain() returned error: %v", err)
	}

	renewal := renewalResponse.Data
	if want, got := 1, renewal.ID; want != got {
		t.Fatalf("Registrar.RenewDomain() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := 999, renewal.DomainID; want != got {
		t.Fatalf("Registrar.RenewDomain() returned DomainID expected to be `%v`, got `%v`", want, got)
	}
}
