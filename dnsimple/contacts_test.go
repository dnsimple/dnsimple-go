package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactPath(t *testing.T) {
	assert.Equal(t, "/1010/contacts", contactPath("1010", 0))
	assert.Equal(t, "/1010/contacts/1", contactPath("1010", 1))
}

func TestContactsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listContacts/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	contactsResponse, err := client.Contacts.ListContacts(context.Background(), "1010", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, contactsResponse.Pagination)
	contacts := contactsResponse.Data
	assert.Len(t, contacts, 2)
	assert.Equal(t, int64(1), contacts[0].ID)
	assert.Equal(t, "Default", contacts[0].Label)
}

func TestContactsService_List_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listContacts/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Contacts.ListContacts(context.Background(), "1010", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestContactsService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createContact/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"label": "Default"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	contactAttributes := Contact{Label: "Default"}

	contactResponse, err := client.Contacts.CreateContact(context.Background(), accountID, contactAttributes)

	assert.NoError(t, err)
	contact := contactResponse.Data
	assert.Equal(t, int64(1), contact.ID)
	assert.Equal(t, "Default", contact.Label)
}

func TestContactsService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getContact/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	contactID := int64(1)

	contactResponse, err := client.Contacts.GetContact(context.Background(), accountID, contactID)

	assert.NoError(t, err)
	contact := contactResponse.Data
	wantSingle := &Contact{
		ID:            1,
		AccountID:     1010,
		Label:         "Default",
		FirstName:     "First",
		LastName:      "User",
		JobTitle:      "CEO",
		Organization:  "Awesome Company",
		Address1:      "Italian Street, 10",
		City:          "Roma",
		StateProvince: "RM",
		PostalCode:    "00100",
		Country:       "IT",
		Phone:         "+18001234567",
		Fax:           "+18011234567",
		Email:         "first@example.com",
		CreatedAt:     "2016-01-19T20:50:26Z",
		UpdatedAt:     "2016-01-19T20:50:26Z",
	}
	assert.Equal(t, wantSingle, contact)
}

func TestContactsService_Update(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateContact/success.http")

		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		want := map[string]interface{}{"label": "Default"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	contactAttributes := Contact{Label: "Default"}
	accountID := "1010"
	contactID := int64(1)

	contactResponse, err := client.Contacts.UpdateContact(context.Background(), accountID, contactID, contactAttributes)

	assert.NoError(t, err)
	contact := contactResponse.Data
	assert.Equal(t, int64(1), contact.ID)
	assert.Equal(t, "Default", contact.Label)
}

func TestContactsService_Delete(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteContact/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	contactID := int64(1)

	_, err := client.Contacts.DeleteContact(context.Background(), accountID, contactID)

	assert.NoError(t, err)
}
