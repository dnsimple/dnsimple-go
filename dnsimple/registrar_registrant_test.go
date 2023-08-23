package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrarService_ListRegistrantChanges(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/registrant_changes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listRegistrantChanges/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.ListRegistrantChange(context.Background(), "1010", &RegistrantChangeListOptions{})

	assert.NoError(t, err)
	changes := res.Data
	assert.Equal(t, changes[0], RegistrantChange{
		Id:                  101,
		AccountId:           101,
		DomainId:            101,
		ContactId:           101,
		State:               "new",
		ExtendedAttributes:  map[string]string{},
		RegistryOwnerChange: true,
		IrtLockLiftedBy:     "",
		CreatedAt:           "2017-02-03T17:43:22Z",
		UpdatedAt:           "2017-02-03T17:43:22Z",
	})
}

func TestRegistrarService_CreateRegistrantChange(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/registrant_changes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createRegistrantChange/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.CreateRegistrantChange(context.Background(), "1010", &CreateRegistrantChangeInput{
		DomainId:           101,
		ContactId:          101,
		ExtendedAttributes: map[string]string{},
	})

	assert.NoError(t, err)
	change := res.Data
	assert.Equal(t, change, &RegistrantChange{
		Id:                  101,
		AccountId:           101,
		DomainId:            101,
		ContactId:           101,
		State:               "new",
		ExtendedAttributes:  map[string]string{},
		RegistryOwnerChange: true,
		IrtLockLiftedBy:     "",
		CreatedAt:           "2017-02-03T17:43:22Z",
		UpdatedAt:           "2017-02-03T17:43:22Z",
	})
}

func TestRegistrarService_CheckRegistrantChange(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/registrant_changes/check", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkRegistrantChange/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.CheckRegistrantChange(context.Background(), "1010", &CheckRegistrantChangeInput{
		DomainId:  101,
		ContactId: 101,
	})

	assert.NoError(t, err)
	change := res.Data
	assert.Equal(t, change, &RegistrantChangeCheck{
		DomainId:            101,
		ContactId:           101,
		ExtendedAttributes:  make([]ExtendedAttribute, 0),
		RegistryOwnerChange: true,
	})
}

func TestRegistrarService_GetRegistrantChange(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/registrant_changes/101", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getRegistrantChange/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	res, err := client.Registrar.GetRegistrantChange(context.Background(), "1010", 101)

	assert.NoError(t, err)
	change := res.Data
	assert.Equal(t, change, &RegistrantChange{
		Id:                  101,
		AccountId:           101,
		DomainId:            101,
		ContactId:           101,
		State:               "new",
		ExtendedAttributes:  map[string]string{},
		RegistryOwnerChange: true,
		IrtLockLiftedBy:     "",
		CreatedAt:           "2017-02-03T17:43:22Z",
		UpdatedAt:           "2017-02-03T17:43:22Z",
	})
}

func TestRegistrarService_DeleteRegistrantChange(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/registrant_changes/101", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteRegistrantChange/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Registrar.DeleteRegistrantChange(context.Background(), "1010", 101)

	assert.NoError(t, err)
}
