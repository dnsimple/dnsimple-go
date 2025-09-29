// Package webhook provides the support for reading and parsing the events
// sent from DNSimple via webhook.
package webhook

import (
	"encoding/json"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple"
)

// Actor represents the entity that triggered the event. It can be either an user,
// a DNSimple support representative or the DNSimple system.
type Actor struct {
	ID     string `json:"id"`
	Entity string `json:"entity"`
	Pretty string `json:"pretty"`
}

// Account represents the account that this event is attached to.
type Account struct {
	dnsimple.Account

	// Display is a string that can be used as a display label
	// and it is sent in a webhook payload
	// It generally represent the name of the account.
	Display string `json:"display,omitempty"`

	// Identifier is a human-readable string identifier
	// and it is sent in a webhook payload
	// It generally represent the StringID or email of the account.
	Identifier string `json:"identifier,omitempty"`
}

// Event represents a webhook event.
type Event struct {
	APIVersion string   `json:"api_version"`
	RequestID  string   `json:"request_identifier"`
	Name       string   `json:"name"`
	Actor      *Actor   `json:"actor"`
	Account    *Account `json:"account"`
	data       EventDataContainer
	payload    []byte
}

// EventDataContainer defines the container for the event payload data.
type EventDataContainer interface {
	unmarshalEventData([]byte) error
}

// GetData returns the data container for the specific event type.
func (e *Event) GetData() EventDataContainer {
	return e.data
}

// GetPayload returns the data payload.
func (e *Event) GetPayload() []byte {
	return e.payload
}

// ParseEvent takes an event payload and attempts to deserialize the payload into an Event.
//
// The event data will be set with a data type that matches the event action in the payload.
// If no direct match is found, then a GenericEventData is assigned.
//
// The event data type is an EventContainerData interface. Therefore, you must perform typecasting
// to access any type-specific field.
func ParseEvent(payload []byte) (*Event, error) {
	e := &Event{payload: payload}

	if err := json.Unmarshal(payload, e); err != nil {
		return nil, err
	}

	data, err := switchEventData(e)
	if err != nil {
		return nil, err
	}

	e.data = data
	return e, nil
}

type eventDataStruct struct {
	Data interface{} `json:"data"`
}

func unmarshalEventData(data []byte, v interface{}) error {
	return json.Unmarshal(data, &eventDataStruct{Data: v})
}
