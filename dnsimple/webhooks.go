package dnsimple

import (
	"fmt"
)

// WebhooksService handles communication with the webhook related
// methods of the DNSimple API.
//
// See PRIVATE
type WebhooksService struct {
	client *Client
}

// WebhookResponse represents a response from an API method that returns a Webhook struct.
type WebhookResponse struct {
	Response
	Data *Webhook `json:"data"`
}

// WebhookResponse represents a response from an API method that returns a collection of Webhook struct.
type WebhooksResponse struct {
	Response
	Data []Webhook `json:"data"`
}

// Webhook represents a DNSimple webhook.
type Webhook struct {
	ID  int    `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

// webhookPath generates the resource path for given webhook.
func webhookPath(accountID string, webhookID int) (path string) {
	path = fmt.Sprintf("/%v/webhooks", accountID)
	if webhookID != 0 {
		path = fmt.Sprintf("%v/%v", path, webhookID)
	}
	return
}

// List the webhooks.
//
// See PRIVATE
func (s *WebhooksService) List(accountID string) (*WebhooksResponse, error) {
	path := webhookPath(accountID, 0)
	webhooksResponse := &WebhooksResponse{}

	resp, err := s.client.get(path, webhooksResponse)
	if err != nil {
		return webhooksResponse, err
	}

	webhooksResponse.HttpResponse = resp
	return webhooksResponse, nil
}

// Create a new webhook.
//
// See PRIVATE
func (s *WebhooksService) Create(accountID string, webhookAttributes Webhook) (*WebhookResponse, error) {
	path := webhookPath(accountID, 0)
	webhookResponse := &WebhookResponse{}

	resp, err := s.client.post(path, webhookAttributes, webhookResponse)
	if err != nil {
		return nil, err
	}

	webhookResponse.HttpResponse = resp
	return webhookResponse, nil
}

// Get a webhook.
//
// See PRIVATE
func (s *WebhooksService) Get(accountID string, webhookID int) (*WebhookResponse, error) {
	path := webhookPath(accountID, webhookID)
	webhookResponse := &WebhookResponse{}

	resp, err := s.client.get(path, webhookResponse)
	if err != nil {
		return nil, err
	}

	webhookResponse.HttpResponse = resp
	return webhookResponse, nil
}
