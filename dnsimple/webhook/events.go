package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	var event Event

	switch name {
	case // contact
		"contact.create",
		"contact.update",
		"contact.delete":
		event = &ContactEvent{}
	case // domain
		"domain.auto_renewal_enable",
		"domain.auto_renewal_disable",
		"domain.create",
		"domain.delete",
		"domain.register",           // TODO
		"domain.renew",              // TODO
		"domain.delegation_change",  // TODO
		"domain.registrant_change",  // TODO
		"domain.resolution_enable",  // TODO
		"domain.resolution_disable", // TODO
		"domain.token_reset",
		"domain.transfer": // TODO
		event = &DomainEvent{}
	case // webhook
		"webhook.create",
		"webhook.delete":
		event = &WebhookEvent{}
	default:
		event = &GenericEvent{}
	}

	return event, event.parse(payload)
}

//
// GenericEvent represents a generic event, where the data is a simple map of strings.
//
type GenericEvent struct {
	Event_Header
	Data interface{} `json:"data"`
}

func (e *GenericEvent) parse(payload []byte) error {
	e.payload = payload
	return unmashalEvent(payload, e)
}

// ParseGenericEvent unpacks the data into a GenericEvent.
func ParseGenericEvent(e *GenericEvent, payload []byte) error {
	return e.parse(payload)
}

//
// ContactEvent represents the base event sent for a contact action.
//
type ContactEvent struct {
	Event_Header
	Data    *ContactEvent     `json:"data"`
	Contact *dnsimple.Contact `json:"contact"`
}

// ParseContactEvent unpacks the data into a ContactEvent.
func ParseContactEvent(e *ContactEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *ContactEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// DomainEvent represents the base event sent for a domain action.
//
type DomainEvent struct {
	Event_Header
	Data   *DomainEvent     `json:"data"`
	Domain *dnsimple.Domain `json:"domain"`
}

// ParseDomainEvent unpacks the payload into a DomainEvent.
func ParseDomainEvent(e *DomainEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// WebhookEvent represents the base event sent for a webhook action.
//
type WebhookEvent struct {
	Event_Header
	Data    *WebhookEvent     `json:"data"`
	Webhook *dnsimple.Webhook `json:"webhook"`
}

// ParseWebhookEvent unpacks the data into a WebhookEvent.
func ParseWebhookEvent(e *WebhookEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *WebhookEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}
