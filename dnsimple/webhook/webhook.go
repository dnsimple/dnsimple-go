// Package dnsimple/webhook provides the support for reading and parsing the events
// sent from DNSimple via webhook.
package webhook

import (
	"encoding/json"

	"github.com/aetrion/dnsimple-go/dnsimple"
)

type Action struct {
	Action string `json:"action"`
}

type Actor struct {
	ID     int    `json:"id"`
	Entity string `json:"entity"`
	Pretty string `json:"pretty"`
}

type eventCore struct {
	APIVersion string `json:"api_version"`
	RequestID  string `json:"request_identifier"`
	Actor      Actor  `json:"actor"`
	Action     string `json:"action"`
	Payload   []byte             `json:"-"`
}

type Event interface {
	Parse([]byte) (error)
}

type DomainCreateEvent struct {
	eventCore
	RequestID string             `json:"request_identifier"`
	Domain    *dnsimple.Domain   `json:"domain"`
	Data      *DomainCreateEvent `json:"data"`
}

func ParseDomainCreateEvent(data []byte) (*DomainCreateEvent, error) {
	event := &DomainCreateEvent{}
	return event, event.Parse(data)
}

func (e *DomainCreateEvent) Parse(data []byte) (error) {
	e.Payload, e.Data = data, e
	return json.Unmarshal(data, e)
}

func Parse(data []byte) (Event, error) {
	action := &Action{}
	json.Unmarshal(data, &action)

	switch action.Action {
	case "domain.create":
		return ParseDomainCreateEvent(data)
	}

	return nil, nil
}
