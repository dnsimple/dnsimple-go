package webhook

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

var regexpUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func readHttpRequestFixture(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile("../../fixtures.http" + filename)
	if err != nil {
		t.Fatalf("Unable to read HTTP fixture: %v", err)
	}

	s := string(data[:])

	return s
}

func getHttpRequestFromFixture(t *testing.T, filename string) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(readHttpRequestFixture(t, filename))))
	if err != nil {
		t.Fatalf("Unable to create http.Request from fixture: %v", err)
	}
	return req
}

func getHttpRequestBodyFromFixture(t *testing.T, filename string) []byte {
	req := getHttpRequestFromFixture(t, filename)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Error reading fixture: %v", err)
	}
	return body
}

func TestParseGenericEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "generic", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "generic", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	dataPointer, ok := event.GetData().(*GenericEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	data := *dataPointer

	if want, got := "example.com", data["domain"].(map[string]interface{})["name"]; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_BillingSettingsUpdate(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.billing_settings_update/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.billing_settings_update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*AccountEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "hello@example.com", data.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.update/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*AccountEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "hello@example.com", data.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_UserInvitationAccept(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.user_invitation_accept/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.user_invitation_accept", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*AccountEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "xxxxx@xxxxx1.xxx", data.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_UserInvite(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.user_invite/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.user_invite", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*AccountEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "xxxxx@xxxxx2.xxx", data.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}
}

func TestParseCertificateEvent_Certificate_RemovePrivateKey(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/certificate.remove_private_key/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "certificate.remove_private_key", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*CertificateEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := int64(41203), data.Certificate.ID; want != got {
		t.Errorf("ParseEvent Certificate.ID expected to be %v, got %v", want, got)
	}
}

func TestParseContactEvent_Contact_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/contact.create/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ContactEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "Test", data.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Label expected to be %v, got %v", want, got)
	}
}

func TestParseContactEvent_Contact_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/contact.update/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ContactEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "Test", data.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Label expected to be %v, got %v", want, got)
	}
}

func TestParseContactEvent_Contact_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/contact.delete/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ContactEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "Test", data.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Label expected to be %v, got %v", want, got)
	}
}

func TestParseDNSSECEvent_DNSSEC_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/dnssec.create/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "dnssec.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	_, ok := event.GetData().(*DNSSECEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
}

func TestParseDNSSECEvent_DNSSEC_RotationStart(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_start/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "dnssec.rotation_start", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DNSSECEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "BD9D898E92D0F668E6BDBC5E79D52E5C3BAB12823A6EEE8C8B6DC633007DFABC", data.DelegationSignerRecord.Digest; want != got {
		t.Errorf("ParseEvent DelegationSignerRecord.Digest expected to be %v, got %v", want, got)
	}
}

func TestParseDNSSECEvent_DNSSEC_RotationComplete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_complete/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "dnssec.rotation_complete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DNSSECEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "EF1D343203E03F1C98120646971F7B96806B759B66622F0A224551DA1A1EFC9A", data.DelegationSignerRecord.Digest; want != got {
		t.Errorf("ParseEvent DelegationSignerRecord.Digest expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalDisable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_disable/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_disable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "foobarbaz.online", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalEnable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_enable/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_enable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "foobarbaz.online", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.create/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.zone", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.delete/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.zone", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Register(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.register/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.register", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example-20181109121341.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Renew(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.renew/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.renew", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if data.Auto != true {
		t.Errorf("ParseEvent auto expected to be %v", true)
	}
	if want, got := "example.test", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_DelegationChange(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.delegation_change/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.delegation_change", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "foo1bar2.cloud", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := (&dnsimple.Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com", "ns4.dnsimple.com"}), data.Delegation; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Delegation expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_RegistrantChange(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.registrant_change/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.registrant_change", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example-20181109121341.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := "Test", data.Registrant.Label; want != got {
		t.Errorf("ParseEvent Registrant.Label expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_ResolutionDisable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.resolution_disable/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.resolution_disable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.zone", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_ResolutionEnable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.resolution_enable/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.resolution_enable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.zone", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Transfer(t *testing.T) {
	//payload := `{"data": {"domain": {"id": 6637, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 24, "auto_renew": false, "created_at": "2016-03-24T21:03:49.392Z", "expires_on": null, "updated_at": "2016-03-24T21:03:49.392Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 409}}, "name": "domain.transfer:started", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "49901af0-569e-4acd-900f-6edf0ebc123c"}`
	payload := `{"data": {"domain": {"id": 6637, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 24, "auto_renew": false, "created_at": "2016-03-24T21:03:49.392Z", "expires_on": null, "updated_at": "2016-03-24T21:03:49.392Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 409}}, "name": "domain.transfer", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "49901af0-569e-4acd-900f-6edf0ebc123c"}`

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.transfer", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*DomainEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseEmailForwardEvent_EmailForward_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/email_forward.create/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "email_forward.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*EmailForwardEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example@example.zone", data.EmailForward.From; want != got {
		t.Errorf("ParseEvent EmailForward.From expected to be %v, got %v", want, got)
	}
}

func TestParseEmailForwardEvent_EmailForward_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/email_forward.delete/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "email_forward.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*EmailForwardEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := ".*@example.zone", data.EmailForward.From; want != got {
		t.Errorf("ParseEvent EmailForward.From expected to be %v, got %v", want, got)
	}
}

func TestParseEmailForwardEvent_EmailForward_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/email_forward.update/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "email_forward.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*EmailForwardEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := ".*@example.zone", data.EmailForward.From; want != got {
		t.Errorf("ParseEvent EmailForward.From expected to be %v, got %v", want, got)
	}
}

func TestParseWebhookEvent_Webhook_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/webhook.create/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "webhook.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*WebhookEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "https://xxxxxx-xxxxxxx-00000.herokuapp.com/xxxxxxxx", data.Webhook.URL; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}
}

func TestParseWebhookEvent_Webhook_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/webhook.delete/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "webhook.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*WebhookEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "https://xxxxxx-xxxxxxx-00000.herokuapp.com/xxxxxxxx", data.Webhook.URL; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Disable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/whois_privacy.disable/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.disable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "xxxxxxxx.email", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(39319), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Enable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/whois_privacy.enable/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.enable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "xxxxxxxx.email", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(39319), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Purchase(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/whois_privacy.purchase/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.purchase", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "xxxxxxxx.email", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(39319), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Renew(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2016-01-17T17:10:41.187Z", "expires_on": "2017-01-17", "updated_at": "2016-01-17T17:11:19.797Z", "unicode_name": "example.com", "private_whois": true, "registrant_id": 2}, "whois_privacy": {"id": 3, "enabled": true, "domain_id": 1, "created_at": "2016-01-17T17:10:50.713Z", "expires_on": "2017-01-17", "updated_at": "2016-03-20T16:45:57.409Z"}}, "name": "whois_privacy.renew", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "e3861a08-a771-4049-abc4-715a3f7b7d6f"}`

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.renew", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*WhoisPrivacyEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(3), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseZoneEvent_Zone_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone.create/example.http")

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ZoneEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.zone", data.Zone.Name; want != got {
		t.Errorf("ParseEvent Zone.Name expected to be %v, got %v", want, got)
	}
}

func TestParseZoneEvent_Zone_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone.delete/example.http")

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ZoneEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "example.zone", data.Zone.Name; want != got {
		t.Errorf("ParseEvent Zone.Name expected to be %v, got %v", want, got)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone_record.create/example.http")

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone_record.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ZoneRecordEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "", data.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent ZoneRecord.Name expected to be %v, got %v", want, got)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone_record.update/example.http")

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone_record.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ZoneRecordEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "www", data.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent ZoneRecord.Name expected to be %v, got %v", want, got)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone_record.delete/example.http")

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone_record.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*ZoneRecordEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := "www", data.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent ZoneRecord.Name expected to be %v, got %v", want, got)
	}
}
