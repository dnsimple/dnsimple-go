package webhook

import (
	"github.com/dnsimple/dnsimple-go/v8/dnsimple"
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
		"account.user_invite",
		"account.user_remove":
		data = &AccountMembershipEventData{}
	case // sso
		"account.sso_user_add":
		data = &AccountSsoEventData{}
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
		"domain.delegation_change",
		"domain.delete",
		"domain.register",
		"domain.registrant_change",
		"domain.registrant_change:started",
		"domain.registrant_change:cancelled",
		"domain.renew",
		"domain.resolution_disable",
		"domain.resolution_enable",
		"domain.restore",
		"domain.transfer":
		data = &DomainEventData{}
	case // domain transfer lock
		"domain.transfer_lock_enable",
		"domain.transfer_lock_disable":
		data = &DomainTransferLockEventData{}
	case // email forward
		"email_forward.create",
		"email_forward.delete",
		"email_forward.update",
		"email_forward.activate",
		"email_forward.deactivate":
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
	case // subscription
		"subscription.migrate",
		"subscription.renew",
		"subscription.renew:failed",
		"subscription.subscribe",
		"subscription.unsubscribe":
		data = &SubscriptionEventData{}
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

// AccountMembershipEventData represents the data node for an Account event.
type AccountMembershipEventData struct {
	Account           *dnsimple.Account           `json:"account"`
	AccountInvitation *dnsimple.AccountInvitation `json:"account_invitation"`
	User              *dnsimple.User              `json:"user"`
}

func (d *AccountMembershipEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// SsoEvent
//

// AccountSsoEventData represents the data node for an single sign-on (SSO) event.
type AccountSsoEventData struct {
	Account                 *dnsimple.Account                 `json:"account"`
	User                    *dnsimple.User                    `json:"user"`
	AccountIdentityProvider *dnsimple.AccountIdentityProvider `json:"account_identity_provider"`
}

func (d *AccountSsoEventData) unmarshalEventData(payload []byte) error {
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
	Zone                   *dnsimple.Zone                   `json:"zone"`
	DnssecConfiguration    *dnsimple.Dnssec                 `json:"dnssec"`
	DelegationSignerRecord *dnsimple.DelegationSignerRecord `json:"delegation_signer_record"`
}

func (d *DNSSECEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// DomainEvent
//

// DomainEventData represents the data node for a Domain event.
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
// DomainTransferLockEvent
//

// DomainTransferLockEventData represents the data node for a DomainTransferLockEnable or DomainTransferLockDisable event.
type DomainTransferLockEventData struct {
	Domain *dnsimple.Domain `json:"domain"`
}

func (d *DomainTransferLockEventData) unmarshalEventData(payload []byte) error {
	return unmarshalEventData(payload, d)
}

//
// DomainRegistrantChangeEvent
//

// DomainRegistrantChangegEventData represents the data node for a DomainRegistrantChange event.
type DomainRegistrantChangeEventData struct {
	Domain     *dnsimple.Domain  `json:"domain"`
	Registrant *dnsimple.Contact `json:"registrant"`
}

func (d *DomainRegistrantChangeEventData) unmarshalEventData(payload []byte) error {
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
// SubscriptionEvent
//

// SubscriptionEventData represents the data node for a Subscription event.
type SubscriptionEventData struct {
	Subscription *dnsimple.Subscription `json:"subscription"`
}

func (d *SubscriptionEventData) unmarshalEventData(payload []byte) error {
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
