package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	switch name {
	case "domain.create":
		return ParseDomainCreateEvent(payload)
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

// DomainCreateEvent represents an event sent when a domain is created.
type DomainCreateEvent struct {
	eventCore
	Data   *DomainCreateEvent `json:"data"`
	Domain *dnsimple.Domain   `json:"domain"`
}

// ParseDomainCreateEvent unpacks the data into a DomainCreateEvent.
func ParseDomainCreateEvent(payload []byte) (*DomainCreateEvent, error) {
	event := &DomainCreateEvent{}
	return event, event.parse(payload)
}

func (e *DomainCreateEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}
