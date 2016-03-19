package webhook

import (
	"github.com/aetrion/dnsimple-go/dnsimple"
)

func switchEvent(name string, payload []byte) (Event, error) {
	var event Event

	switch name {
	case // account
		"account.update",                  // TODO
		"account.billing_settings_update", // TODO
		"account.payment_details_update",  // TODO
		"account.add_user",                // TODO
		"account.remove_user":             // TODO
		event = &AccountEvent{}
	case // contact
		"contact.create",
		"contact.update",
		"contact.delete":
		event = &ContactEvent{}
	case // domain
		"domain.auto_renewal_enable",
		"domain.auto_renewal_disable",
		"domain.create",
		"domain.delete",
		"domain.register",
		"domain.renew",
		"domain.delegation_change",  // TODO
		"domain.registrant_change",  // TODO
		"domain.resolution_enable",  // TODO
		"domain.resolution_disable", // TODO
		"domain.token_reset",
		"domain.transfer": // TODO
		event = &DomainEvent{}
	case // webhook
		"webhook.create",
		"webhook.delete":
		event = &WebhookEvent{}
	case // zone record
		"record.create",
		"record.update",
		"record.delete":
		event = &ZoneRecordEvent{}
	default:
		event = &GenericEvent{}
	}

	return event, event.parse(payload)
}

//
// GenericEvent
//

// GenericEvent represents a generic event, where the data is a simple map of strings.
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
// AccountEvent
//

// AccountEvent represents the base event sent for an account action.
type AccountEvent struct {
	Event_Header
	Data    *AccountEvent     `json:"data"`
	Account *dnsimple.Account `json:"account"`
}

// ParseAccountEvent unpacks the data into an AccountEvent.
func ParseAccountEvent(e *AccountEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *AccountEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// ContactEvent
//

// ContactEvent represents the base event sent for a contact action.
type ContactEvent struct {
	Event_Header
	Data    *ContactEvent     `json:"data"`
	Contact *dnsimple.Contact `json:"contact"`
}

// ParseContactEvent unpacks the data into a ContactEvent.
func ParseContactEvent(e *ContactEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *ContactEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// DomainEvent
//

// DomainEvent represents the base event sent for a domain action.
type DomainEvent struct {
	Event_Header
	Data   *DomainEvent     `json:"data"`
	Domain *dnsimple.Domain `json:"domain"`
}

// ParseDomainEvent unpacks the payload into a DomainEvent.
func ParseDomainEvent(e *DomainEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *DomainEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// ZoneRecordEvent
//

// ZoneRecordEvent represents the base event sent for a webhook action.
type ZoneRecordEvent struct {
	Event_Header
	Data       *ZoneRecordEvent     `json:"data"`
	ZoneRecord *dnsimple.ZoneRecord `json:"record"`
}

// ParseZoneRecordEvent unpacks the data into a ZoneRecordEvent.
func ParseZoneRecordEvent(e *ZoneRecordEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *ZoneRecordEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}

//
// WebhookEvent
//

// WebhookEvent represents the base event sent for a webhook action.
type WebhookEvent struct {
	Event_Header
	Data    *WebhookEvent     `json:"data"`
	Webhook *dnsimple.Webhook `json:"webhook"`
}

// ParseWebhookEvent unpacks the data into a WebhookEvent.
func ParseWebhookEvent(e *WebhookEvent, payload []byte) error {
	return e.parse(payload)
}

func (e *WebhookEvent) parse(payload []byte) error {
	e.payload, e.Data = payload, e
	return unmashalEvent(payload, e)
}
