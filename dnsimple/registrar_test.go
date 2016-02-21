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

	registerRequest := &RegisterRequest{RegistrantID: 2}

	registrationResponse, err := client.Registrar.Register("1010", "example.com", registerRequest)
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

func TestRegistrarService_Transfer(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/transfer", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/transfer/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"registrant_id": float64(2), "auth_info": "x1y2z3"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	transferRequest := &TransferRequest{RegistrantID: 2, AuthInfo: "x1y2z3"}

	transferResponse, err := client.Registrar.Transfer("1010", "example.com", transferRequest)
	if err != nil {
		t.Fatalf("Registrar.Transfer() returned error: %v", err)
	}

	domain := transferResponse.Data
	if want, got := 1, domain.ID; want != got {
		t.Fatalf("Registrar.Transfer() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example.com", domain.Name; want != got {
		t.Fatalf("Registrar.Transfer() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestRegistrarService_TransferOut(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/transfer_out", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/transferOut/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.TransferOut("1010", "example.com")
	if err != nil {
		t.Fatalf("Registrar.TransferOut() returned error: %v", err)
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
