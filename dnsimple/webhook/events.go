package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	switch name {
	case "domain.create":
		e := &DomainCreateEvent{}
		return e, ParseDomainCreateEvent(e, payload)
	case "domain.delete":
		e := &DomainDeleteEvent{}
		return e, ParseDomainDeleteEvent(e, payload)
	default:
		e := &GenericEvent{}
		return e, ParseGenericEvent(e, payload)
	}
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

// ParseDomainEvent unpacks the payload into a DomainEvent.
func ParseDomainEvent(e *DomainEvent, payload []byte) error { return e.parse(payload) }

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

type DomainCreateEvent struct{ DomainEvent }

func ParseDomainCreateEvent(e *DomainCreateEvent, p []byte) error { return e.DomainEvent.parse(p) }

type DomainDeleteEvent struct{ DomainEvent }

func ParseDomainDeleteEvent(e *DomainDeleteEvent, p []byte) error { return e.DomainEvent.parse(p) }
