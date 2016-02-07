package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	switch name {
	case "domain.create":
		return ParseDomainCreateEvent(payload)
	case "domain.delete":
		return ParseDomainDeleteEvent(payload)
	default:
		return ParseGenericEvent(payload)
	}
}

// GenericEvent represents a generic event, where the data is a simple map of strings.
type GenericEvent struct {
	eventCore
	Data interface{} `json:"data"`
}

// ParseGenericEvent unpacks the data into a GenericEvent.
func ParseGenericEvent(payload []byte) (*GenericEvent, error) {
	event := &GenericEvent{}
	return event, event.parse(payload)
}

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
func ParseDomainEvent(p []byte) (*DomainEvent, error) {
	e := &DomainEvent{}
	return e, e.parse(p)
}

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

type DomainCreateEvent struct{ DomainEvent }

func ParseDomainCreateEvent(p []byte) (*DomainCreateEvent, error) {
	e := &DomainCreateEvent{}
	return e, e.parse(p)
}

type DomainDeleteEvent struct{ DomainEvent }

func ParseDomainDeleteEvent(p []byte) (*DomainDeleteEvent, error) {
	e := &DomainDeleteEvent{}
	return e, e.parse(p)
}
