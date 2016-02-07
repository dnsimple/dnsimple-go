package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

type GenericEvent struct {
	eventCore
	Data interface{} `json:"data"`
}

func (e *GenericEvent) parse(data []byte) error {
	e.payload = data
	return jsonUnmarshalEvent(data, e)
}

func ParseGenericEvent(data []byte) (*GenericEvent, error) {
	event := &GenericEvent{}
	return event, event.parse(data)
}

type DomainCreateEvent struct {
	eventCore
	Data   *DomainCreateEvent `json:"data"`
	Domain *dnsimple.Domain   `json:"domain"`
}

func (e *DomainCreateEvent) parse(data []byte) error {
	e.payload, e.Data = data, e
	return jsonUnmarshalEvent(data, e)
}

func ParseDomainCreateEvent(data []byte) (*DomainCreateEvent, error) {
	event := &DomainCreateEvent{}
	return event, event.parse(data)
}
