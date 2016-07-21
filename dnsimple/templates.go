package dnsimple

import (
	"fmt"
)

// TemplatesService handles communication with the template related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/templates/
type TemplatesService struct {
	client *Client
}

// Template represents a Template in DNSimple.
type Template struct {
	ID          int    `json:"id,omitempty"`
	AccountID   int    `json:"account_id,omitempty"`
	Name        string `json:"name,omitempty"`
	ShortName   string `json:"short_name,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func templatePath(accountID string, templateID string) string {
	if templateID != "" {
		return fmt.Sprintf("/%v/templates/%v", accountID, templateID)
	}

	return fmt.Sprintf("/%v/templates", accountID)
}

// TemplateResponse represents a response from an API method that returns a Template struct.
type TemplateResponse struct {
	Response
	Data *Template `json:"data"`
}

// TemplatesResponse represents a response from an API method that returns a collection of Template struct.
type TemplatesResponse struct {
	Response
	Data []Template `json:"data"`
}

// ListTemplates list the templates for an account.
//
// See https://developer.dnsimple.com/v2/templates/#list
func (s *TemplatesService) ListTemplates(accountID string, options *ListOptions) (*TemplatesResponse, error) {
	path := versioned(templatePath(accountID, ""))
	templatesResponse := &TemplatesResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, templatesResponse)
	if err != nil {
		return templatesResponse, err
	}

	templatesResponse.HttpResponse = resp
	return templatesResponse, nil
}

// CreateTemplate creates a new template.
//
// See https://developer.dnsimple.com/v2/templates/#create
func (s *TemplatesService) CreateTemplate(accountID string, templateAttributes Template) (*TemplateResponse, error) {
	path := versioned(templatePath(accountID, ""))
	templateResponse := &TemplateResponse{}

	resp, err := s.client.post(path, templateAttributes, templateResponse)
	if err != nil {
		return nil, err
	}

	templateResponse.HttpResponse = resp
	return templateResponse, nil
}

// GetTemplate fetches a template.
//
// See https://developer.dnsimple.com/v2/templates/#get
func (s *TemplatesService) GetTemplate(accountID string, templateID string) (*TemplateResponse, error) {
	path := versioned(templatePath(accountID, templateID))
	templateResponse := &TemplateResponse{}

	resp, err := s.client.get(path, templateResponse)
	if err != nil {
		return nil, err
	}

	templateResponse.HttpResponse = resp
	return templateResponse, nil
}

// UpdateTemplate updates a template.
//
// See https://developer.dnsimple.com/v2/templates/#update
func (s *TemplatesService) UpdateTemplate(accountID string, templateID string, templateAttributes Template) (*TemplateResponse, error) {
	path := versioned(templatePath(accountID, templateID))
	templateResponse := &TemplateResponse{}

	resp, err := s.client.patch(path, templateAttributes, templateResponse)
	if err != nil {
		return nil, err
	}

	templateResponse.HttpResponse = resp
	return templateResponse, nil
}

// DeleteTemplate deletes a template.
//
// See https://developer.dnsimple.com/v2/templates/#delete
func (s *TemplatesService) DeleteTemplate(accountID string, templateID string) (*TemplateResponse, error) {
	path := versioned(templatePath(accountID, templateID))
	templateResponse := &TemplateResponse{}

	resp, err := s.client.delete(path, nil, nil)
	if err != nil {
		return nil, err
	}

	templateResponse.HttpResponse = resp
	return templateResponse, nil
}

// ApplyTemplate deletes a template.
//
// See https://developer.dnsimple.com/v2/templates/domains/#apply
func (s *TemplatesService) ApplyTemplate(accountID string, domainID string, templateID string) (*TemplateResponse, error) {
	path := versioned(fmt.Sprintf("%v/templates/%v", domainPath(accountID, domainID), templateID))
	templateResponse := &TemplateResponse{}

	resp, err := s.client.post(path, nil, nil)
	if err != nil {
		return nil, err
	}

	templateResponse.HttpResponse = resp
	return templateResponse, nil
}

// Template Records

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
