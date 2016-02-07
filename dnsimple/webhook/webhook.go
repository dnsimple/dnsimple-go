// Package dnsimple/webhook provides the support for reading and parsing the events
// sent from DNSimple via webhook.
package webhook

import (
	"encoding/json"
)

type Action struct {
	Action string `json:"action"`
}

func ParseAction(data []byte) (action *Action, err error) {
	action = &Action{}
	err = json.Unmarshal(data, action)
	return
}

type Actor struct {
	ID     int    `json:"id"`
	Entity string `json:"entity"`
	Pretty string `json:"pretty"`
}

type Event interface {
	Payload() []byte
	parse([]byte) error
}

type eventCore struct {
	APIVersion string  `json:"api_version"`
	RequestID  string  `json:"request_identifier"`
	Actor      *Actor  `json:"actor"`
	Action     *Action `json:"action"`
	payload    []byte
}

func (e *eventCore) Payload() []byte {
	return e.payload
}

func Parse(data []byte) (Event, error) {
	action, err := ParseAction(data)
	if err != nil {
		return nil, err
	}

	switch action.Action {
	case "domain.create":
		return ParseDomainCreateEvent(data)
	default:
		return ParseGenericEvent(data)
	}

	return nil, nil
}

func jsonUnmarshalEvent(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
