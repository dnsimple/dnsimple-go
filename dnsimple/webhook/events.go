package webhook

import (
	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func switchEventData(event *EventContainer) (EventDataContainer, error) {
	var dataContainer EventDataContainer

	switch event.Name {
	case // account
		"account.billing_settings_update",
		"account.update",
		"account.user_invitation_accept",
		"account.user_invite":
		dataContainer = &AccountEventData{}
	case // certificate
		"certificate.remove_private_key":
		dataContainer = &CertificateEventData{}
	case // contact
		"contact.create",
		"contact.delete",
		"contact.update":
		dataContainer = &ContactEventData{}
	case // dnssec
		"dnssec.rotation_complete",
		"dnssec.rotation_start":
		dataContainer = &DNSSECEventData{}
	case // domain
		"domain.auto_renewal_disable",
		"domain.auto_renewal_enable",
		"domain.create",
		"domain.delete",
		"domain.register",
		"domain.renew",
		"domain.delegation_change",
		"domain.registrant_change",
		"domain.resolution_disable",
		"domain.resolution_enable",
		"domain.transfer": // TODO
		dataContainer = &DomainEventData{}
	case // email forward
		"email_forward.create",
		"email_forward.delete",
		"email_forward.update":
		dataContainer = &EmailForwardEventData{}
	case // webhook
		"webhook.create",
		"webhook.delete":
		dataContainer = &WebhookEventData{}
	case // whois privacy
		"whois_privacy.disable",
		"whois_privacy.enable",
		"whois_privacy.purchase",
		"whois_privacy.renew": // TODO
		dataContainer = &WhoisPrivacyEventData{}
	case // zone
		"zone.create",
		"zone.delete":
		dataContainer = &ZoneEventData{}
	case // zone record
		"zone_record.create",
		"zone_record.delete",
		"zone_record.update":
		dataContainer = &ZoneRecordEventData{}
	default:
		dataContainer = &GenericEventData{}
	}

	err := dataContainer.unmarshalEventData(event.payload)
	return dataContainer, err
}

//
// GenericEvent
//

// GenericEventData represents the data node for a generic event, where the data is a simple map of strings.
type GenericEventData map[string]interface{}

func (d *GenericEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// AccountEvent
//

// AccountEventData represents the data node for an Account event.
type AccountEventData struct {
	Account *dnsimple.Account `json:"account"`
}

func (d *AccountEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// CertificateEvent
//

// CertificateEventData represents the data node for a Certificate event.
type CertificateEventData struct {
	Certificate *dnsimple.Certificate `json:"certificate"`
}

func (d *CertificateEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// ContactEvent
//

// ContactEventData represents the data node for a Contact event.
type ContactEventData struct {
	Contact *dnsimple.Contact `json:"contact"`
}

func (d *ContactEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// DNSSECEvent
//

// DNSSECEventData represents the data node for a DNSSEC event.
type DNSSECEventData struct {
	DelegationSignerRecord *dnsimple.DelegationSignerRecord `json:"delegation_signer_record"`
	//DNSSECConfig           *dnsimple.DNSSECConfig           `json:"dnssec"`
}

func (d *DNSSECEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// DomainEvent
//

// DomainEventData represents the data node for a Contact event.
type DomainEventData struct {
	Auto       bool                 `json:"auto"`
	Domain     *dnsimple.Domain     `json:"domain"`
	Registrant *dnsimple.Contact    `json:"registrant"`
	Delegation *dnsimple.Delegation `json:"name_servers"`
}

func (d *DomainEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// EmailForwardEvent
//

// EmailForwardEventData represents the data node for a EmailForward event.
type EmailForwardEventData struct {
	EmailForward *dnsimple.EmailForward `json:"email_forward"`
}

func (d *EmailForwardEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// WebhookEvent
//

// WebhookEventData represents the data node for a Webhook event.
type WebhookEventData struct {
	Webhook *dnsimple.Webhook `json:"webhook"`
}

func (d *WebhookEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// WhoisPrivacyEvent
//

// WhoisPrivacyEventData represents the data node for a WhoisPrivacy event.
type WhoisPrivacyEventData struct {
	Domain       *dnsimple.Domain       `json:"domain"`
	WhoisPrivacy *dnsimple.WhoisPrivacy `json:"whois_privacy"`
}

func (d *WhoisPrivacyEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// ZoneEvent
//

// ZoneEventData represents the data node for a Zone event.
type ZoneEventData struct {
	Zone *dnsimple.Zone `json:"zone"`
}

func (d *ZoneEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}

//
// ZoneRecordEvent
//

// ZoneRecordEventData represents the data node for a ZoneRecord event.
type ZoneRecordEventData struct {
	ZoneRecord *dnsimple.ZoneRecord `json:"zone_record"`
}

func (d *ZoneRecordEventData) unmarshalEventData(payload []byte) error {
	return unmashalEventData(payload, d)
}
