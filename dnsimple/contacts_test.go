package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestContacts_contactPath(t *testing.T) {
	actual := contactPath("1", nil)
	expected := "/1/contacts"

	if actual != expected {
		t.Errorf("contactPath(\"1\", nil): actual %s, expected %s", actual, expected)
	}

	actual = contactPath("1", 1)
	expected = "/1/contacts/1"

	if actual != expected {
		t.Errorf("contactPath(\"1\", 1): actual %s, expected %s", actual, expected)
	}
}

func TestContactsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":[{"id":1,"account_id":1010,"label":"Default","first_name":"First","last_name":"User","job_title":"CEO","organization_name":"Awesome Company","email_address":"first@example.com","phone":"+18001234567","fax":"+18011234567","address1":"Italian Street, 10","address2":"","city":"Roma","state_province":"RM","postal_code":"00100","country":"IT","created_at":"2013-11-08T17:23:15.884Z","updated_at":"2015-01-08T21:30:50.228Z"},{"id":2,"account_id":1010,"label":"","first_name":"Second","last_name":"User","job_title":"","organization_name":"","email_address":"second@example.com","phone":"+18881234567","fax":"","address1":"French Street","address2":"c/o Someone","city":"Paris","state_province":"XY","postal_code":"00200","country":"FR","created_at":"2014-12-06T15:46:18.014Z","updated_at":"2014-12-06T15:46:18.014Z"}],"pagination":{"current_page":1,"per_page":30,"total_entries":2,"total_pages":1}}
		`)
	})

	accountID := "1010"

	contactsResponse, err := client.Contacts.List(accountID)
	if err != nil {
		t.Fatalf("Contacts.List() returned error: %v", err)
	}

	contacts := contactsResponse.Data()
	if want, got := 2, len(contacts); want != got {
		t.Errorf("Contacts.List() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 1, contacts[0].ID; want != got {
		t.Fatalf("Contacts.List() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Default", contacts[0].Label; want != got {
		t.Fatalf("Contacts.List() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestContactsService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"label": "Default"}
		testRequestJSON(t, r, want)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `
			{"data":{"id":1,"account_id":1010,"label":"Default","first_name":"First","last_name":"User","job_title":"CEO","organization_name":"Awesome Company","email_address":"first@example.com","phone":"+18001234567","fax":"+18011234567","address1":"Italian Street, 10","address2":"","city":"Roma","state_province":"RM","postal_code":"00100","country":"IT","created_at":"2016-01-19T20:50:26.066Z","updated_at":"2016-01-19T20:50:26.066Z"}}
		`)
	})

	accountID := "1010"
	contactAttributes := Contact{Label: "Default"}

	contactResponse, err := client.Contacts.Create(accountID, contactAttributes)
	if err != nil {
		t.Fatalf("Contacts.Create() returned error: %v", err)
	}

	contact := contactResponse.Data()
	if want, got := 1, contact.ID; want != got {
		t.Fatalf("Contacts.Create() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Default", contact.Label; want != got {
		t.Fatalf("Contacts.Create() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestContactsService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":{"id":1,"account_id":1010,"label":"Default","first_name":"First","last_name":"User","job_title":"CEO","organization_name":"Awesome Company","email_address":"first@example.com","phone":"+18001234567","fax":"+18011234567","address1":"Italian Street, 10","address2":"","city":"Roma","state_province":"RM","postal_code":"00100","country":"IT","created_at":"2016-01-19T20:50:26.066Z","updated_at":"2016-01-19T20:50:26.066Z"}}
		`)
	})

	accountID := "1010"
	contactID := 1

	contactResponse, err := client.Contacts.Get(accountID, contactID)
	if err != nil {
		t.Fatalf("Contacts.Get() returned error: %v", err)
	}

	contact := contactResponse.Data()
	wantSingle := &Contact{
		ID:            1,
		Label:         "Default",
		FirstName:     "First",
		LastName:      "User",
		JobTitle:      "CEO",
		Organization:  "Awesome Company",
		Email:         "first@example.com",
		Phone:         "+18001234567",
		Fax:           "+18011234567",
		Address1:      "Italian Street, 10",
		City:          "Roma",
		StateProvince: "RM",
		PostalCode:    "00100",
		Country:       "IT",
		CreatedAt:     "2016-01-19T20:50:26.066Z",
		UpdatedAt:     "2016-01-19T20:50:26.066Z"}

	if !reflect.DeepEqual(contact, wantSingle) {
		t.Fatalf("Contacts.Get() returned %+v, want %+v", contact, wantSingle)
	}
}

func TestContactsService_Update(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		want := map[string]interface{}{"label": "Default"}
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `
			{"data":{"id":1,"account_id":1010,"label":"Default","first_name":"First","last_name":"User","job_title":"CEO","organization_name":"Awesome Company","email_address":"first@example.com","phone":"+18001234567","fax":"+18011234567","address1":"Italian Street, 10","address2":"","city":"Roma","state_province":"RM","postal_code":"00100","country":"IT","created_at":"2016-01-19T20:50:26.066Z","updated_at":"2016-01-19T20:50:26.066Z"}}
		`)
	})

	contactAttributes := Contact{Label: "Default"}
	accountID := "1010"
	contactID := 1

	contactResponse, err := client.Contacts.Update(accountID, contactID, contactAttributes)
	if err != nil {
		t.Fatalf("Contacts.Update() returned error: %v", err)
	}

	contact := contactResponse.Data()
	if want, got := 1, contact.ID; want != got {
		t.Fatalf("Contacts.Update() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Default", contact.Label; want != got {
		t.Fatalf("Contacts.Update() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestContactsService_Delete(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(http.StatusNoContent)
	})

	accountID := "1010"
	contactID := 1

	_, err := client.Contacts.Delete(accountID, contactID)
	if err != nil {
		t.Fatalf("Contacts.Delete() returned error: %v", err)
	}
}
