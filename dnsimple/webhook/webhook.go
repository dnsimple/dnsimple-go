// Package dnsimple/webhook provides the support for reading and parsing the events
// sent from DNSimple via webhook.
package webhook

import(
	//"github.com/miekg/dns"
	"encoding/json"

	"github.com/aetrion/dnsimple-go/dnsimple"
	"fmt"
)

type Payload struct {
	APIVersion string      `json:"api_version"`
	RequestID  string      `json:"request_identifier"`
	Actor      Actor       `json:"actor"`
	Action     string      `json:"action"`
	What     interface{}      `json:"data"`
	Data     []byte      `json:"-"`
}

type Actor struct {
	ID     int    `json:"id"`
	Entity string `json:"entity"`
	Pretty string `json:"pretty"`
}

type Event interface {
	Payload() *Payload
}

type DomainCreateEvent struct {
	payload *Payload        `json:"-"`
	Domain *dnsimple.Domain `json:"domain"`
}

func (e *DomainCreateEvent) Payload() *Payload {
	return e.payload
}

func ParsePayload(data []byte) (*Payload, error) {
	payload := &Payload{Data:data}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func Parse(data []byte) (Event) {
	payload, _ := ParsePayload(data)
	var event Event

	switch payload.Action {
	case "domain.create":
		event = &DomainCreateEvent{payload: payload}
		v, _ := json.Marshal(payload.What)
		fmt.Println(v)
		json.Unmarshal(v, &event)
	}

	return event
}
