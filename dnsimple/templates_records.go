package dnsimple

import (
	"fmt"
)

// TemplateRecord represents a DNS record for a template in DNSimple.
type TemplateRecord struct {
	ID         int    `json:"id,omitempty"`
	TemplateID int    `json:"template_id,omitempty"`
	Name       string `json:"name"`
	Content    string `json:"content,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	Type       string `json:"type,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// TemplateRecordResponse represents a response from an API method that returns a TemplateRecord struct.
type TemplateRecordResponse struct {
	Response
	Data *TemplateRecord `json:"data"`
}

// TemplateRecordsResponse represents a response from an API method that returns a collection of TemplateRecord struct.
type TemplateRecordsResponse struct {
	Response
	Data []TemplateRecord `json:"data"`
}

func templateRecordPath(accountID string, templateID string, templateRecordID string) string {
	if templateRecordID != "" {
		return fmt.Sprintf("%v/records/%v", templatePath(accountID, templateID), templateRecordID)
	}

	return templatePath(accountID, templateID) + "/records"
}

// ListTemplateRecords list the templates for an account.
//
// See https://developer.dnsimple.com/v2/templates/records/#list
func (s *TemplatesService) ListTemplateRecords(accountID string, templateID string, options *ListOptions) (*TemplateRecordsResponse, error) {
	path := versioned(templateRecordPath(accountID, templateID, ""))
	templateRecordsResponse := &TemplateRecordsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, templateRecordsResponse)
	if err != nil {
		return templateRecordsResponse, err
	}

	templateRecordsResponse.HttpResponse = resp
	return templateRecordsResponse, nil
}

// CreateTemplateRecord creates a new template record.
//
// See https://developer.dnsimple.com/v2/templates/records/#create
func (s *TemplatesService) CreateTemplateRecord(accountID string, templateID string, templateRecordAttributes TemplateRecord) (*TemplateRecordResponse, error) {
	path := versioned(templateRecordPath(accountID, templateID, ""))
	templateRecordResponse := &TemplateRecordResponse{}

	resp, err := s.client.post(path, templateRecordAttributes, templateRecordResponse)
	if err != nil {
		return nil, err
	}

	templateRecordResponse.HttpResponse = resp
	return templateRecordResponse, nil
}

// GetTemplateRecord fetches a template record.
//
// See https://developer.dnsimple.com/v2/templates/records/#get
func (s *TemplatesService) GetTemplateRecord(accountID string, templateID string, templateRecordID string) (*TemplateRecordResponse, error) {
	path := versioned(templateRecordPath(accountID, templateID, templateRecordID))
	templateRecordResponse := &TemplateRecordResponse{}

	resp, err := s.client.get(path, templateRecordResponse)
	if err != nil {
		return nil, err
	}

	templateRecordResponse.HttpResponse = resp
	return templateRecordResponse, nil
}

// DeleteTemplateRecord deletes a template record.
//
// See https://developer.dnsimple.com/v2/templates/records/#delete
func (s *TemplatesService) DeleteTemplateRecord(accountID string, templateID string, templateRecordID string) (*TemplateRecordResponse, error) {
	path := versioned(templateRecordPath(accountID, templateID, templateRecordID))
	templateRecordResponse := &TemplateRecordResponse{}

	resp, err := s.client.delete(path, nil, nil)
	if err != nil {
		return nil, err
	}

	templateRecordResponse.HttpResponse = resp
	return templateRecordResponse, nil
}
