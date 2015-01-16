package dnsimple

import (
	"fmt"
)

// ContactsService handles communication with the contact related
// methods of the DNSimple API.
//
// DNSimple API docs: http://developer.dnsimple.com/contacts/
type ContactsService struct {
	client *Client
}

type Contact struct {
	Id           int    `json:"id,omitempty"`
	UserId       int    `json:"user_id,omitempty"`
	Label        string `json:"label,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	JobTitle     string `json:"job_title,omitempty"`
	Organization string `json:"organization_name,omitempty"`
	Email        string `json:"email_address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Fax          string `json:"fax,omitempty"`
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city,omitempty"`
	Zip          string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type contactWrapper struct {
	Contact Contact `json:"contact"`
}

// contactPath generates the resource path for given contact.
func contactPath(contact interface{}) string {
	if contact != nil {
		return fmt.Sprintf("contacts/%d", contact)
	}
	return "contacts"
}

// List the contacts.
//
// DNSimple API docs: http://developer.dnsimple.com/contacts/#list
func (s *ContactsService) List() ([]Contact, *Response, error) {
	path := contactPath(nil)
	wrappedContacts := []contactWrapper{}

	res, err := s.client.get(path, &wrappedContacts)
	if err != nil {
		return []Contact{}, res, err
	}

	contacts := []Contact{}
	for _, contact := range wrappedContacts {
		contacts = append(contacts, contact.Contact)
	}

	return contacts, res, nil
}
