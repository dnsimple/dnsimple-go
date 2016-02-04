package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestContacts_contactPath(t *testing.T) {
	actual := contactPath("1", nil)
	expected := "1/contacts"

	if actual != expected {
		t.Errorf("contactPath(\"1\", nil): actual %s, expected %s", actual, expected)
	}

	actual = contactPath("1", 1)
	expected = "1/contacts/1"

	if actual != expected {
		t.Errorf("contactPath(\"1\", 1): actual %s, expected %s", actual, expected)
	}
}

func TestContactsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data":[{"id":1,"label":"Default"},{"id":2,"label":"Simone"}]}`)
	})

	accountId := "1"
	contacts, _, err := client.Contacts.List(accountId)

	if err != nil {
		t.Errorf("Contacts.List returned error: %v", err)
	}

	if size, want := len(contacts), 2; size != want {
		t.Fatalf("Contacts.List returned %+v items, want %+v", size, want)
	}

	for i, item := range contacts {
		if kind, want := reflect.TypeOf(item).Name(), "Contact"; kind != want {
			t.Errorf("Contacts.List expected [%d] to be %s, got %s", i, want, kind)
		}
	}
}

func TestContactsService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/contacts", func(w http.ResponseWriter, r *http.Request) {
		want := map[string]interface{}{"label": "Default"}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprintf(w, `{"data":{"id":1, "label":"Default"}}`)
	})

	accountId := "1"
	contactId := 1
	contactValues := Contact{Label: "Default"}
	contact, _, err := client.Contacts.Create(accountId, contactValues)

	if err != nil {
		t.Errorf("Contacts.Create returned error: %v", err)
	}

	want := Contact{Id: contactId, Label: "Default"}
	if !reflect.DeepEqual(contact, want) {
		t.Fatalf("Contacts.Create returned %+v, want %+v", contact, want)
	}
}

func TestContactsService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data":{"id":1,"user_id":21,"label":"Default","first_name":"Simone","last_name":"Carletti","job_title":"Underwater Programmer","organization_name":"DNSimple","email_address":"simone.carletti@dnsimple.com","phone":"+1 111 4567890","fax":"+1 222 4567890","address1":"Awesome Street","address2":"c/o Someone","city":"Rome","state_province":"RM","postal_code":"00171","country":"IT"}}`)
	})

	accountId := "1"
	contactId := 1
	contact, _, err := client.Contacts.Get(accountId, contactId)

	if err != nil {
		t.Errorf("Contacts.Get returned error: %v", err)
	}

	want := Contact{
		Id:            contactId,
		Label:         "Default",
		FirstName:     "Simone",
		LastName:      "Carletti",
		JobTitle:      "Underwater Programmer",
		Organization:  "DNSimple",
		Email:         "simone.carletti@dnsimple.com",
		Phone:         "+1 111 4567890",
		Fax:           "+1 222 4567890",
		Address1:      "Awesome Street",
		Address2:      "c/o Someone",
		City:          "Rome",
		StateProvince: "RM",
		PostalCode:    "00171",
		Country:       "IT"}
	if !reflect.DeepEqual(contact, want) {
		t.Fatalf("Contacts.Get returned %+v, want %+v", contact, want)
	}
}

func TestContactsService_Update(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		want := map[string]interface{}{"label": "Default"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"data":{"id":1, "label": "Default"}}`)
	})

	contactValues := Contact{Label: "Default"}
	accountId := "1"
	contactId := 1
	contact, _, err := client.Contacts.Update(accountId, contactId, contactValues)

	if err != nil {
		t.Errorf("Contacts.Update returned error: %v", err)
	}

	want := Contact{Id: contactId, Label: "Default"}
	if !reflect.DeepEqual(contact, want) {
		t.Fatalf("Contacts.Update returned %+v, want %+v", contact, want)
	}
}

func TestContactsService_Delete(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	accountId := "1"
	contactId := 1
	_, err := client.Contacts.Delete(accountId, contactId)

	if err != nil {
		t.Errorf("Contacts.Delete returned error: %v", err)
	}
}
