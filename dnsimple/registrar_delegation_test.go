package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrarService_GetDomainDelegation(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/delegation", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getDomainDelegation/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	delegationResponse, err := client.Registrar.GetDomainDelegation(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
	delegation := delegationResponse.Data
	wantSingle := &Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com", "ns4.dnsimple.com"}
	assert.Equal(t, wantSingle, delegation)
}

func TestRegistrarService_ChangeDomainDelegation(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/delegation", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/changeDomainDelegation/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		want := []interface{}{"ns1.dnsimple.com", "ns2.dnsimple.com"}
		testRequestJSONArray(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	newDelegation := &Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com"}

	delegationResponse, err := client.Registrar.ChangeDomainDelegation(context.Background(), "1010", "example.com", newDelegation)

	assert.NoError(t, err)
	delegation := delegationResponse.Data
	wantSingle := &Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com", "ns4.dnsimple.com"}
	assert.Equal(t, wantSingle, delegation)
}

func TestRegistrarService_ChangeDomainDelegationToVanity(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/delegation/vanity", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/changeDomainDelegationToVanity/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		want := []interface{}{"ns1.example.com", "ns2.example.com"}
		testRequestJSONArray(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	newDelegation := &Delegation{"ns1.example.com", "ns2.example.com"}

	delegationResponse, err := client.Registrar.ChangeDomainDelegationToVanity(context.Background(), "1010", "example.com", newDelegation)

	assert.NoError(t, err)
	delegation := delegationResponse.Data[0].Name
	wantSingle := "ns1.example.com"
	assert.Equal(t, wantSingle, delegation)
}

func TestRegistrarService_ChangeDomainDelegationFromVanity(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/delegation/vanity", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/changeDomainDelegationFromVanity/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.ChangeDomainDelegationFromVanity(context.Background(), "1010", "example.com")

	assert.NoError(t, err)
}
