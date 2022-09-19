package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrarService_CheckDomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/check", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkDomain/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	checkResponse, err := client.Registrar.CheckDomain(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	check := checkResponse.Data
	assert.Equal(t, "ruby.codes", check.Domain)
}

func TestRegistrarService_GetDomainPremiumPrice(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/premium_price", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainPremiumPrice/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	priceResponse, err := client.Registrar.GetDomainPremiumPrice(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	price := priceResponse.Data
	assert.Equal(t, "109.00", price.PremiumPrice)
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
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.GetDomainPremiumPrice(context.Background(), "1010", "example.com", &DomainPremiumPriceOptions{Action: "registration"})

	assert.NoError(t, err)
}

func TestRegistrarService_GetDomainPrices(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/bingo.pizza/prices", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainPrices/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	checkResponse, err := client.Registrar.GetDomainPrices(context.Background(), "1010", "bingo.pizza")

	assert.NoError(t, err)
	check := checkResponse.Data
	assert.Equal(t, "bingo.pizza", check.Domain)
	assert.True(t, check.Premium)
	assert.Equal(t, float64(20.0), check.RegistrationPrice)
	assert.Equal(t, float64(20.0), check.RenewalPrice)
	assert.Equal(t, float64(20.0), check.TransferPrice)
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
		_, _ = io.Copy(w, httpResponse.Body)
	})

	registerRequest := &RegisterDomainInput{RegistrantID: 2}

	registrationResponse, err := client.Registrar.RegisterDomain(context.Background(), "1010", "example.com", registerRequest)

	assert.NoError(t, err)
	registration := registrationResponse.Data
	assert.Equal(t, int64(1), registration.ID)
	assert.Equal(t, int64(999), registration.DomainID)
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
		_, _ = io.Copy(w, httpResponse.Body)
	})

	registerRequest := &RegisterDomainInput{RegistrantID: 2, ExtendedAttributes: map[string]string{"att1": "val1", "att2": "val2"}}

	_, err := client.Registrar.RegisterDomain(context.Background(), "1010", "example.com", registerRequest)

	assert.NoError(t, err)
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
		_, _ = io.Copy(w, httpResponse.Body)
	})

	transferRequest := &TransferDomainInput{RegistrantID: 2, AuthCode: "x1y2z3"}

	transferResponse, err := client.Registrar.TransferDomain(context.Background(), "1010", "example.com", transferRequest)

	assert.NoError(t, err)
	transfer := transferResponse.Data
	assert.Equal(t, int64(1), transfer.ID)
	assert.Equal(t, int64(999), transfer.DomainID)
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
		_, _ = io.Copy(w, httpResponse.Body)
	})

	transferRequest := &TransferDomainInput{RegistrantID: 2, AuthCode: "x1y2z3", ExtendedAttributes: map[string]string{"att1": "val1", "att2": "val2"}}

	_, err := client.Registrar.TransferDomain(context.Background(), "1010", "example.com", transferRequest)

	assert.NoError(t, err)
}

func TestRegistrarService_TransferDomainOut(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/authorize_transfer_out", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/authorizeDomainTransferOut/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.TransferDomainOut(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
}

func TestRegistrarService_GetDomainTransfer(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/transfers/361", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainTransfer/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	domainTransferResponse, err := client.Registrar.GetDomainTransfer(context.Background(), "1010", "example.com", 361)

	assert.NoError(t, err)
	domainTransfer := domainTransferResponse.Data
	wantSingle := &DomainTransfer{
		ID:                int64(361),
		DomainID:          int64(182245),
		RegistrantID:      int64(2715),
		State:             "cancelled",
		AutoRenew:         false,
		WhoisPrivacy:      false,
		StatusDescription: "Canceled by customer",
		CreatedAt:         "2020-06-05T18:08:00Z",
		UpdatedAt:         "2020-06-05T18:10:01Z"}
	assert.Equal(t, wantSingle, domainTransfer)
}

func TestRegistrarService_CancelDomainTransfer(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/transfers/361", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/cancelDomainTransfer/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	domainTransferResponse, err := client.Registrar.CancelDomainTransfer(context.Background(), "1010", "example.com", 361)

	assert.NoError(t, err)
	domainTransfer := domainTransferResponse.Data
	wantSingle := &DomainTransfer{
		ID:                int64(361),
		DomainID:          int64(182245),
		RegistrantID:      int64(2715),
		State:             "transferring",
		AutoRenew:         false,
		WhoisPrivacy:      false,
		StatusDescription: "",
		CreatedAt:         "2020-06-05T18:08:00Z",
		UpdatedAt:         "2020-06-05T18:08:04Z"}
	assert.Equal(t, wantSingle, domainTransfer)
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
		_, _ = io.Copy(w, httpResponse.Body)
	})

	renewalResponse, err := client.Registrar.RenewDomain(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	renewal := renewalResponse.Data
	assert.Equal(t, int64(1), renewal.ID)
	assert.Equal(t, int64(999), renewal.DomainID)
}
