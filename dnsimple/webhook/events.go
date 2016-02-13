package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	var event Event

	switch name {
	case "domain.create":
		event = &DomainCreateEvent{}
	case "domain.delete":
		event = &DomainDeleteEvent{}
	case "domain.token_reset":
		event = &DomainTokenResetEvent{}
	case "domain.auto_renew_enable":
		event = &DomainAutoRenewalEnableEvent{}
	case "domain.auto_renew_disable":
		event = &DomainAutoRenewalDisableEvent{}
	case "webhook.create":
		event = &WebhookCreateEvent{}
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
// DomainEvent represents the base event sent for a domain action.
//
type DomainEvent struct {
	Event_Header
	Data   *DomainEvent     `json:"data"`
	Domain *dnsimple.Domain `json:"domain"`
}

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

// ParseDomainEvent unpacks the payload into a DomainEvent.
func ParseDomainEvent(e *DomainEvent, payload []byte) error {
	return e.parse(payload)
}

type DomainTokenResetEvent struct{ DomainEvent }
type DomainAutoRenewalEnableEvent struct{ DomainEvent }
type DomainAutoRenewalDisableEvent struct{ DomainEvent }
type DomainCreateEvent struct{ DomainEvent }
type DomainDeleteEvent struct{ DomainEvent }

func ParseDomainCreateEvent(e *DomainCreateEvent, p []byte) error {
	return e.DomainEvent.parse(p)
}
func ParseDomainDeleteEvent(e *DomainDeleteEvent, p []byte) error {
	return e.DomainEvent.parse(p)
}
func ParseDomainTokenResetEvent(e *DomainTokenResetEvent, p []byte) error {
	return e.DomainEvent.parse(p)
}
func ParseDomainAutoRenewalEnableEvent(e *DomainAutoRenewalEnableEvent, p []byte) error {
	return e.DomainEvent.parse(p)
}
func ParseDomainAutoRenewalDisableEvent(e *DomainAutoRenewalDisableEvent, p []byte) error {
	return e.DomainEvent.parse(p)
}

//
// Webhook represents a generic event, where the data is a simple map of strings.
//
type WebhookEvent struct {
	Event_Header
	Data *WebhookEvent `json:"data"`
}

// ParseWebhookEvent unpacks the data into a WebhookEvent.
func ParseWebhookEvent(e *GenericEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *WebhookEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

type WebhookCreateEvent struct{ WebhookEvent }

func ParseWebhookCreateEvent(e *WebhookCreateEvent, p []byte) error {
	return e.WebhookEvent.parse(p)
}
