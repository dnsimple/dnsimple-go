package webhook

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/stretchr/testify/assert"
)

var regexpUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func readHTTPRequestFixture(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile("../../fixtures.http" + filename)
	assert.NoError(t, err)

	s := string(data[:])

	return s
}

func getHTTPRequestFromFixture(t *testing.T, filename string) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(readHTTPRequestFixture(t, filename))))
	assert.NoError(t, err)

	return req
}

func getHTTPRequestBodyFromFixture(t *testing.T, filename string) []byte {
	req := getHTTPRequestFromFixture(t, filename)
	body, err := ioutil.ReadAll(req.Body)
	assert.NoError(t, err)

	return body
}

func TestParseGenericEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "generic", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := ParseEvent([]byte(payload))

	assert.NoError(t, err)
	assert.Equal(t, "generic", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	dataPointer, ok := event.GetData().(*GenericEventData)
	assert.True(t, ok)
	data := *dataPointer
	assert.Equal(t, "example.com", data["domain"].(map[string]interface{})["name"])
}

func TestParseAccountEvent_Account_BillingSettingsUpdate(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.billing_settings_update/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.billing_settings_update", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*AccountEventData)
	assert.True(t, ok)
	assert.Equal(t, "hello@example.com", data.Account.Email)
}

func TestParseAccountEvent_Account_Update(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.update/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.update", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*AccountEventData)
	assert.True(t, ok)
	assert.Equal(t, "hello@example.com", data.Account.Email)
}

func TestParseAccountEvent_Account_UserInvitationAccept(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_invitation_accept/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.user_invitation_accept", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*AccountMembershipEventData)
	assert.True(t, ok)

	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	assert.Equal(t, expectedAccount, *data.Account)

	expectedAccountInvitation := dnsimple.AccountInvitation{
		ID:                   3523,
		Email:                "xxxxxx@xxxxxx.xxx",
		Token:                "eb5763dc-0f24-420b-b7f6-c7355c8b8309",
		AccountID:            1111,
		CreatedAt:            "2020-05-12T18:42:44Z",
		UpdatedAt:            "2020-05-12T18:43:44Z",
		InvitationSentAt:     "2020-05-12T18:42:44Z",
		InvitationAcceptedAt: "2020-05-12T18:43:44Z",
	}
	assert.Equal(t, expectedAccountInvitation, *data.AccountInvitation)

	assert.Nil(t, data.User)
}

func TestParseAccountEvent_Account_UserInvitationRevoke(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_invitation_revoke/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.user_invitation_revoke", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*AccountMembershipEventData)
	assert.True(t, ok)

	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	assert.Equal(t, expectedAccount, *data.Account)

	expectedAccountInvitation := dnsimple.AccountInvitation{
		ID:                   3522,
		Email:                "xxxxxx@xxxxxx.xxx",
		Token:                "be87d69b-a58a-43bd-9a21-aaf303829a60",
		AccountID:            1111,
		CreatedAt:            "2020-05-12T18:42:27Z",
		UpdatedAt:            "2020-05-12T18:42:27Z",
		InvitationSentAt:     "2020-05-12T18:42:27Z",
		InvitationAcceptedAt: "",
	}
	assert.Equal(t, expectedAccountInvitation, *data.AccountInvitation)

	assert.Nil(t, data.User)
}

func TestParseAccountEvent_Account_UserInvite(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_invite/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.user_invite", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*AccountMembershipEventData)
	assert.True(t, ok)

	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	assert.Equal(t, expectedAccount, *data.Account)

	expectedAccountInvitation := dnsimple.AccountInvitation{
		ID:                   3523,
		Email:                "xxxxxx@xxxxxx.xxx",
		Token:                "eb5763dc-0f24-420b-b7f6-c7355c8b8309",
		AccountID:            1111,
		CreatedAt:            "2020-05-12T18:42:44Z",
		UpdatedAt:            "2020-05-12T18:42:44Z",
		InvitationSentAt:     "2020-05-12T18:42:44Z",
		InvitationAcceptedAt: "",
	}
	assert.Equal(t, expectedAccountInvitation, *data.AccountInvitation)

	assert.Nil(t, data.User)
}

func TestParseAccountEvent_Account_UserRemove(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_remove/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.user_remove", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*AccountMembershipEventData)
	assert.True(t, ok)

	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	assert.Equal(t, expectedAccount, *data.Account)

	assert.Nil(t, data.AccountInvitation)

	expectedUser := dnsimple.User{
		ID:    1120,
		Email: "xxxxxx@xxxxxx.xxx",
	}
	assert.Equal(t, expectedUser, *data.User)
}

func TestParseCertificateEvent_Certificate_Issue(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/certificate.issue/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "certificate.issue", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*CertificateEventData)
	assert.True(t, ok)
	assert.Equal(t, int64(101967), data.Certificate.ID)
}

func TestParseCertificateEvent_Certificate_RemovePrivateKey(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/certificate.remove_private_key/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "certificate.remove_private_key", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*CertificateEventData)
	assert.True(t, ok)
	assert.Equal(t, int64(101972), data.Certificate.ID)
}

func TestParseContactEvent_Contact_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/contact.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "contact.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ContactEventData)
	assert.True(t, ok)
	assert.Equal(t, "Test", data.Contact.Label)
}

func TestParseContactEvent_Contact_Update(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/contact.update/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "contact.update", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ContactEventData)
	assert.True(t, ok)
	assert.Equal(t, "Test", data.Contact.Label)
}

func TestParseContactEvent_Contact_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/contact.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "contact.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ContactEventData)
	assert.True(t, ok)
	assert.Equal(t, "Test", data.Contact.Label)
}

func TestParseDNSSECEvent_DNSSEC_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "dnssec.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	_, ok := event.GetData().(*DNSSECEventData)
	assert.True(t, ok)
}

func TestParseDNSSECEvent_DNSSEC_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "dnssec.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	_, ok := event.GetData().(*DNSSECEventData)
	assert.True(t, ok)
}

func TestParseDNSSECEvent_DNSSEC_RotationStart(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_start/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "dnssec.rotation_start", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DNSSECEventData)
	assert.True(t, ok)
	assert.Equal(t, "BD9D898E92D0F668E6BDBC5E79D52E5C3BAB12823A6EEE8C8B6DC633007DFABC", data.DelegationSignerRecord.Digest)
}

func TestParseDNSSECEvent_DNSSEC_RotationComplete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_complete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "dnssec.rotation_complete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DNSSECEventData)
	assert.True(t, ok)
	assert.Equal(t, "EF1D343203E03F1C98120646971F7B96806B759B66622F0A224551DA1A1EFC9A", data.DelegationSignerRecord.Digest)
}

func TestParseDomainEvent_Domain_AutoRenewalDisable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_disable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.auto_renewal_disable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_AutoRenewalEnable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_enable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.auto_renewal_enable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-beta.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-delta.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_Register(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.register/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.register", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_Renew(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.renew/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.renew", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.True(t, data.Auto)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_DelegationChange(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.delegation_change/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.delegation_change", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
	assert.Equal(t, &dnsimple.Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com"}, data.Delegation)
}

func TestParseDomainEvent_Domain_RegistrantChange(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.registrant_change/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.registrant_change", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
	assert.Equal(t, "new_contact", data.Registrant.Label)
}

func TestParseDomainEvent_Domain_ResolutionDisable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.resolution_disable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.resolution_disable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_ResolutionEnable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.resolution_enable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.resolution_enable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_Transfer(t *testing.T) {
	payload := `{"data": {"domain": {"id": 6637, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 24, "auto_renew": false, "created_at": "2016-03-24T21:03:49.392Z", "expires_on": null, "updated_at": "2016-03-24T21:03:49.392Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 409}}, "name": "domain.transfer", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "49901af0-569e-4acd-900f-6edf0ebc123c"}`

	event, err := ParseEvent([]byte(payload))

	assert.NoError(t, err)
	assert.Equal(t, "domain.transfer", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainEventData)
	assert.True(t, ok)
	assert.Equal(t, "example.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_TransferLockEnable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.transfer_lock_enable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "domain.transfer_lock_enable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainTransferLockEventData)
	assert.True(t, ok)
	assert.Equal(t, "example.com", data.Domain.Name)
}

func TestParseDomainEvent_Domain_TransferLockDisable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.transfer_lock_disable/example.http")

	event, err := ParseEvent([]byte(payload))

	assert.NoError(t, err)
	assert.Equal(t, "domain.transfer_lock_disable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*DomainTransferLockEventData)
	assert.True(t, ok)
	assert.Equal(t, "example.com", data.Domain.Name)
}

func TestParseEmailForwardEvent_EmailForward_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/email_forward.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "email_forward.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*EmailForwardEventData)
	assert.True(t, ok)
	assert.Equal(t, "example@example.zone", data.EmailForward.From)
}

func TestParseEmailForwardEvent_EmailForward_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/email_forward.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "email_forward.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*EmailForwardEventData)
	assert.True(t, ok)
	assert.Equal(t, ".*@example.zone", data.EmailForward.From)
}

func TestParseEmailForwardEvent_EmailForward_Update(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/email_forward.update/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "email_forward.update", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*EmailForwardEventData)
	assert.True(t, ok)
	assert.Equal(t, ".*@example.zone", data.EmailForward.From)
}

func TestParseWebhookEvent_Webhook_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/webhook.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "webhook.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*WebhookEventData)
	assert.True(t, ok)
	assert.Equal(t, "https://xxxxxx-xxxxxxx-00000.herokuapp.com/xxxxxxxx", data.Webhook.URL)
}

func TestParseWebhookEvent_Webhook_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/webhook.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "webhook.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*WebhookEventData)
	assert.True(t, ok)
	assert.Equal(t, "https://xxxxxx-xxxxxxx-00000.herokuapp.com/xxxxxxxx", data.Webhook.URL)
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Disable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.disable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "whois_privacy.disable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
	assert.Equal(t, int64(902), data.WhoisPrivacy.ID)
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Enable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.enable/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "whois_privacy.enable", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
	assert.Equal(t, int64(902), data.WhoisPrivacy.ID)
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Purchase(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.purchase/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "whois_privacy.purchase", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
	assert.Equal(t, int64(902), data.WhoisPrivacy.ID)
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Renew(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.renew/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "whois_privacy.renew", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	assert.True(t, ok)
	assert.Equal(t, "example-alpha.com", data.Domain.Name)
	assert.Equal(t, int64(902), data.WhoisPrivacy.ID)
}

func TestParseZoneEvent_Zone_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "zone.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ZoneEventData)
	assert.True(t, ok)
	assert.Equal(t, "example.zone", data.Zone.Name)
}

func TestParseZoneEvent_Zone_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "zone.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ZoneEventData)
	assert.True(t, ok)
	assert.Equal(t, "example.zone", data.Zone.Name)
}

func TestParseZoneRecordEvent_ZoneRecord_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone_record.create/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "zone_record.create", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ZoneRecordEventData)
	assert.True(t, ok)
	assert.Equal(t, "", data.ZoneRecord.Name)
}

func TestParseZoneRecordEvent_ZoneRecord_Update(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone_record.update/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "zone_record.update", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ZoneRecordEventData)
	assert.True(t, ok)
	assert.Equal(t, "www", data.ZoneRecord.Name)
}

func TestParseZoneRecordEvent_ZoneRecord_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone_record.delete/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "zone_record.delete", event.Name)
	assert.Regexp(t, regexpUUID, event.RequestID)

	data, ok := event.GetData().(*ZoneRecordEventData)
	assert.True(t, ok)
	assert.Equal(t, "www", data.ZoneRecord.Name)
}
