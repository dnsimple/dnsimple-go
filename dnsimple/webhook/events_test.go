package webhook

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

var regexpUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func readHTTPRequestFixture(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile("../../fixtures.http" + filename)
	if err != nil {
		t.Fatalf("Unable to read HTTP fixture: %v", err)
	}

	s := string(data[:])

	return s
}

func getHTTPRequestFromFixture(t *testing.T, filename string) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(readHTTPRequestFixture(t, filename))))
	if err != nil {
		t.Fatalf("Unable to create http.Request from fixture: %v", err)
	}
	return req
}

func getHTTPRequestBodyFromFixture(t *testing.T, filename string) []byte {
	req := getHTTPRequestFromFixture(t, filename)
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.billing_settings_update/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.update/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_invitation_accept/example.http")

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

	data, ok := event.GetData().(*AccountMembershipEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	if want, got := expectedAccount, *data.Account; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Account expected to be %v, got %v", want, got)
	}
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
	if want, got := expectedAccountInvitation, *data.AccountInvitation; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent AccountInvitation expected to be %v, got %v", want, got)
	}
	var expectedUser *dnsimple.User = nil
	if want, got := expectedUser, data.User; want != got {
		t.Errorf("ParseEvent User expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_UserInvitationRevoke(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_invitation_revoke/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.user_invitation_revoke", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*AccountMembershipEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	if want, got := expectedAccount, *data.Account; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}
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
	if want, got := expectedAccountInvitation, *data.AccountInvitation; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent AccountInvitation expected to be %v, got %v", want, got)
	}
	var expectedUser *dnsimple.User = nil
	if want, got := expectedUser, data.User; want != got {
		t.Errorf("ParseEvent User expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_UserInvite(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_invite/example.http")

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

	data, ok := event.GetData().(*AccountMembershipEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed: %v", ok)
	}
	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	if want, got := expectedAccount, *data.Account; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Account expected to be %v, got %v", want, got)
	}
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
	if want, got := expectedAccountInvitation, *data.AccountInvitation; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent AccountInvitation expected to be %v, got %v", want, got)
	}
	var expectedUser *dnsimple.User = nil
	if want, got := expectedUser, data.User; want != got {
		t.Errorf("ParseEvent User expected to be %v, got %v", want, got)
	}
}

func TestParseAccountEvent_Account_UserRemove(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.user_remove/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.user_remove", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*AccountMembershipEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed: %v", ok)
	}
	expectedAccount := dnsimple.Account{
		ID:             1111,
		Email:          "xxxxx@xxxxxx.xxx",
		CreatedAt:      "2012-03-16T16:02:54Z",
		UpdatedAt:      "2020-05-10T18:11:03Z",
		PlanIdentifier: "professional-v1-monthly",
	}
	if want, got := expectedAccount, *data.Account; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Account expected to be %v, got %v", want, got)
	}
	var expectedAccountInvitation *dnsimple.AccountInvitation = nil
	if want, got := expectedAccountInvitation, data.AccountInvitation; want != got {
		t.Errorf("ParseEvent AccountInvitation expected to be %v, got %v", want, got)
	}
	expectedUser := dnsimple.User{
		ID:    1120,
		Email: "xxxxxx@xxxxxx.xxx",
	}
	if want, got := expectedUser, *data.User; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent User expected to be %v, got %v", want, got)
	}
}

func TestParseCertificateEvent_Certificate_Issue(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/certificate.issue/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "certificate.issue", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data, ok := event.GetData().(*CertificateEventData)
	if !ok {
		t.Fatalf("ParseEvent type assertion failed")
	}
	if want, got := int64(86368), data.Certificate.ID; want != got {
		t.Errorf("ParseEvent Certificate.ID expected to be %v, got %v", want, got)
	}
}

func TestParseCertificateEvent_Certificate_RemovePrivateKey(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/certificate.remove_private_key/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/contact.create/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/contact.update/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/contact.delete/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.create/example.http")

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

func TestParseDNSSECEvent_DNSSEC_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.delete/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "dnssec.delete", event.Name; want != got {
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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_start/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_complete/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_disable/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalEnable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_enable/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.create/example.http")

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
	if want, got := "example-beta.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Delete(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.delete/example.http")

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
	if want, got := "example-delta.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Register(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.register/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Renew(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.renew/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_DelegationChange(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.delegation_change/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := (&dnsimple.Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com"}), data.Delegation; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Delegation expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_RegistrantChange(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.registrant_change/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := "new_contact", data.Registrant.Label; want != got {
		t.Errorf("ParseEvent Registrant.Label expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_ResolutionDisable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.resolution_disable/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_ResolutionEnable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/domain.resolution_enable/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/email_forward.create/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/email_forward.delete/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/email_forward.update/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/webhook.create/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/webhook.delete/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.disable/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(902), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Enable(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.enable/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(902), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Purchase(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.purchase/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(902), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Renew(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/whois_privacy.renew/example.http")

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
	if want, got := "example-alpha.com", data.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(902), data.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}
}

func TestParseZoneEvent_Zone_Create(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone.create/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone.delete/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone_record.create/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone_record.update/example.http")

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
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/zone_record.delete/example.http")

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
