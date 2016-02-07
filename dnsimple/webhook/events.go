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
	default:
		event = &GenericEvent{}
	}

	return event, event.parse(payload)
}

//
// GenericEvent represents a generic event, where the data is a simple map of strings.
//
type GenericEvent struct {
	eventCore
	Data interface{} `json:"data"`
}

// ParseGenericEvent unpacks the data into a GenericEvent.
func ParseGenericEvent(e *GenericEvent, payload []byte) error { return e.parse(payload) }

func (e *GenericEvent) parse(payload []byte) error {
	e.payload = payload
	return unmashalEvent(payload, e)
}

//
// DomainEvent represents the base event sent for a domain action.
//
type DomainEvent struct {
	eventCore
	Data   *DomainEvent     `json:"data"`
	Domain *dnsimple.Domain `json:"domain"`
}
type DomainTokenResetEvent struct{ DomainEvent }
type DomainAutoRenewalEnableEvent struct{ DomainEvent }
type DomainCreateEvent struct{ DomainEvent }
type DomainDeleteEvent struct{ DomainEvent }

// ParseDomainEvent unpacks the payload into a DomainEvent.
func ParseDomainEvent(e *DomainEvent, payload []byte) error { return e.parse(payload) }

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

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}
