package dnsimple

import ()

// WebhooksService handles communication with the webhook related
// methods of the DNSimple API.
//
// See #
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
	URL string `json:"url,omitempty"`
}
