package webhook

import (
	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func switchEventData(event *Event) (EventDataContainer, error) {
	var data EventDataContainer

	switch event.Name {
	case // account
		"account.billing_settings_update",
		"account.update":
		data = &AccountEventData{}
	case // account_invitation
		"account.user_invitation_accept",
		"account.user_invitation_revoke",
		"account.user_invite":
		data = &AccountInvitationEventData{}
	case // certificate
		"certificate.issue",
		"certificate.remove_private_key":
		data = &CertificateEventData{}
	case // contact
		"contact.create",
		"contact.delete",
		"contact.update":
		data = &ContactEventData{}
	case // dnssec
		"dnssec.create",
		"dnssec.delete",
		"dnssec.rotation_complete",
		"dnssec.rotation_start":
		data = &DNSSECEventData{}
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
		data = &DomainEventData{}
	case // email forward
		"email_forward.create",
		"email_forward.delete",
		"email_forward.update":
		data = &EmailForwardEventData{}
	case // webhook
		"webhook.create",
		"webhook.delete":
		data = &WebhookEventData{}
	case // whois privacy
		"whois_privacy.disable",
		"whois_privacy.enable",
		"whois_privacy.purchase",
		"whois_privacy.renew":
		data = &WhoisPrivacyEventData{}
	case // zone
		"zone.create",
		"zone.delete":
		data = &ZoneEventData{}
	case // zone record
		"zone_record.create",
		"zone_record.delete",
		"zone_record.update":
		data = &ZoneRecordEventData{}
	default:
		data = &GenericEventData{}
	}

	err := data.unmarshalEventData(event.payload)
	return data, err
}

//
// GenericEvent
//

// GenericEventData represents the data node for a generic event, where the data is a simple map of strings.
type GenericEventData map[string]interface{}

func (d *GenericEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// AccountEvent
//

// AccountEventData represents the data node for an Account event.
type AccountEventData struct {
	Account *dnsimple.Account `json:"account"`
}

func (d *AccountEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// AccountInvitationEvent
//

// AccountInvitationEventData represents the data node for an Account event.
type AccountInvitationEventData struct {
	Account           *dnsimple.Account           `json:"account"`
	AccountInvitation *dnsimple.AccountInvitation `json:"account_invitation"`
}

func (d *AccountInvitationEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// CertificateEvent
//

// CertificateEventData represents the data node for a Certificate event.
type CertificateEventData struct {
	Certificate *dnsimple.Certificate `json:"certificate"`
}

func (d *CertificateEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// ContactEvent
//

// ContactEventData represents the data node for a Contact event.
type ContactEventData struct {
	Contact *dnsimple.Contact `json:"contact"`
}

func (d *ContactEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
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
	return unmarshalEventData(payload, d)
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
	return unmarshalEventData(payload, d)
}

//
// EmailForwardEvent
//

// EmailForwardEventData represents the data node for a EmailForward event.
type EmailForwardEventData struct {
	EmailForward *dnsimple.EmailForward `json:"email_forward"`
}

func (d *EmailForwardEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// WebhookEvent
//

// WebhookEventData represents the data node for a Webhook event.
type WebhookEventData struct {
	Webhook *dnsimple.Webhook `json:"webhook"`
}

func (d *WebhookEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
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
	return unmarshalEventData(payload, d)
}

//
// ZoneEvent
//

// ZoneEventData represents the data node for a Zone event.
type ZoneEventData struct {
	Zone *dnsimple.Zone `json:"zone"`
}

func (d *ZoneEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// ZoneRecordEvent
//

// ZoneRecordEventData represents the data node for a ZoneRecord event.
type ZoneRecordEventData struct {
	ZoneRecord *dnsimple.ZoneRecord `json:"zone_record"`
}

func (d *ZoneRecordEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}
