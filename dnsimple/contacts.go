package dnsimple

import (
	"fmt"
)

// ContactsService handles communication with the contact related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/contacts/
type ContactsService struct {
	client *Client
}

type Contact struct {
	Id            int    `json:"id,omitempty"`
	Label         string `json:"label,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	JobTitle      string `json:"job_title,omitempty"`
	Organization  string `json:"organization_name,omitempty"`
	Email         string `json:"email_address,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Fax           string `json:"fax,omitempty"`
	Address1      string `json:"address1,omitempty"`
	Address2      string `json:"address2,omitempty"`
	City          string `json:"city,omitempty"`
	StateProvince string `json:"state_province,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Country       string `json:"country,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type contactsWrapper struct {
	Contacts []Contact `json:"data"`
}

type contactWrapper struct {
	Contact Contact `json:"data"`
}

// contactPath generates the resource path for given contact.
func contactPath(accountId string, contact interface{}) string {
	if contact != nil {
		return fmt.Sprintf("%v/contacts/%d", accountId, contact)
	}
	return fmt.Sprintf("%v/contacts", accountId)
}

// List the contacts.
//
// See https://developer.dnsimple.com/v2/contacts/#list
func (s *ContactsService) List(accountId string) ([]Contact, *Response, error) {
	path := contactPath(accountId, nil)
	data := contactsWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return []Contact{}, res, err
	}

	return data.Contacts, res, nil
}

// Create a new contact.
//
// See https://developer.dnsimple.com/v2/contacts/#create
func (s *ContactsService) Create(accountId string, contactAttributes Contact) (Contact, *Response, error) {
	path := contactPath(accountId, nil)
	data := contactWrapper{}

	res, err := s.client.post(path, contactAttributes, &data)
	if err != nil {
		return Contact{}, res, err
	}

	return data.Contact, res, nil
}

// Get a contact.
//
// See https://developer.dnsimple.com/v2/contacts/#get
func (s *ContactsService) Get(accountId string, contactId int) (Contact, *Response, error) {
	path := contactPath(accountId, contactId)
	data := contactWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return Contact{}, res, err
	}

	return data.Contact, res, nil
}

// Update a contact.
//
// See https://developer.dnsimple.com/v2/contacts/#update
func (s *ContactsService) Update(accountId string, contactId int, contactAttributes Contact) (Contact, *Response, error) {
	path := contactPath(accountId, contactId)
	data := contactWrapper{}

	res, err := s.client.patch(path, contactAttributes, &data)
	if err != nil {
		return Contact{}, res, err
	}

	return data.Contact, res, nil
}

// Delete a contact.
//
// See https://developer.dnsimple.com/v2/contacts/#delete
func (s *ContactsService) Delete(accountId string, contactId int) (*Response, error) {
	path := contactPath(accountId, contactId)

	return s.client.delete(path, nil)
}
