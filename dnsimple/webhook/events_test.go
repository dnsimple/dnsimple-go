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
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event := &GenericEvent{}
	err := ParseGenericEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data := event.Data.(map[string]interface{})
	if want, got := "example.com", data["domain"].(map[string]interface{})["name"]; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	if err != nil {
		t.Fatalf("Parse returned error when parsing: %v", err)
	}
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseAccountEvent_Account_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.update/example.http")

	event := &AccountEvent{}
	err := ParseAccountEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "hello@example.com", event.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*AccountEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseAccountEvent_Account_BillingSettingsUpdate(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.billing_settings_update/example.http")

	event := &AccountEvent{}
	err := ParseAccountEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.billing_settings_update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "hello@example.com", event.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*AccountEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseAccountEvent_Account_RemoveUser(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/account.remove_user/example.http")

	event := &AccountEvent{}
	err := ParseAccountEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.remove_user", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "xxxxxx@xxxxx.xxx", event.Account.Email; want != got {
		t.Errorf("ParseEvent Account.Email expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*AccountEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseContactEvent_Contact_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/contact.create/example.http")

	event := &ContactEvent{}
	err := ParseContactEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "Test", event.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ContactEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseContactEvent_Contact_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/contact.update/example.http")

	event := &ContactEvent{}
	err := ParseContactEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "Test", event.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ContactEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseContactEvent_Contact_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/contact.delete/example.http")

	event := &ContactEvent{}
	err := ParseContactEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "Test", event.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ContactEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDNSSECEvent_DNSSEC_RotationStart(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_start/example.http")

	event := &DNSSECEvent{}
	err := ParseDNSSECEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "dnssec.rotation_start", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "BD9D898E92D0F668E6BDBC5E79D52E5C3BAB12823A6EEE8C8B6DC633007DFABC", event.DelegationSignerRecord.Digest; want != got {
		t.Errorf("ParseEvent DelegationSignerRecord.Digest expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DNSSECEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDNSSECEvent_DNSSEC_RotationComplete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/dnssec.rotation_complete/example.http")

	event := &DNSSECEvent{}
	err := ParseDNSSECEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "dnssec.rotation_complete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "EF1D343203E03F1C98120646971F7B96806B759B66622F0A224551DA1A1EFC9A", event.DelegationSignerRecord.Digest; want != got {
		t.Errorf("ParseEvent DelegationSignerRecord.Digest expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DNSSECEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalDisable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_disable/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_disable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "foobarbaz.online", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalEnable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.auto_renewal_enable/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_enable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "foobarbaz.online", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.create/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.zone", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.delete/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.zone", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_Register(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-24T21:53:38.878Z", "expires_on": "2017-02-24", "updated_at": "2016-02-24T22:22:27.025Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 2}}, "name": "domain.register", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "8c92b76f-125d-43c0-8e72-b911e4bdbd96"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.register", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_Renew(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.renew/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.renew", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if event.Auto != true {
		t.Errorf("ParseEvent auto expected to be %v", true)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.test", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_DelegationChange(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.delegation_change/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)

	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.delegation_change", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "foo1bar2.cloud", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := (&dnsimple.Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com", "ns4.dnsimple.com"}), event.Delegation; !reflect.DeepEqual(want, got) {
		t.Errorf("ParseEvent Delegation expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_RegistrantChange(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-01-16T16:08:50.649Z", "expires_on": "2018-01-16", "updated_at": "2016-03-24T20:30:05.895Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 2}, "registrant": {"id": 2, "fax": "+39 339 1111111", "city": "Rome", "label": "Webhook", "phone": "+39 339 0000000", "country": "IT", "address1": "Some Street", "address2": "", "job_title": "Developer", "last_name": "Contact", "account_id": 1010, "created_at": "2016-02-13T13:11:29.388Z", "first_name": "Example", "updated_at": "2016-02-13T13:11:29.388Z", "postal_code": "12037", "email": "example@example.com", "state_province": "Italy", "organization_name": "Company"}}, "name": "domain.registrant_change", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "0391e4e2-7614-41bf-a7bd-7ba01232e434"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))

	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.registrant_change", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := "Webhook", event.Registrant.Label; want != got {
		t.Errorf("ParseEvent Registrant.Label expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_ResolutionDisable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.resolution_disable/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)

	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.resolution_disable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.zone", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_ResolutionEnable(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/domain.resolution_enable/example.http")

	event := &DomainEvent{}
	err := ParseDomainEvent(event, payload)

	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.resolution_enable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent RequestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.zone", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_Transfer(t *testing.T) {
	//payload := `{"data": {"domain": {"id": 6637, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 24, "auto_renew": false, "created_at": "2016-03-24T21:03:49.392Z", "expires_on": null, "updated_at": "2016-03-24T21:03:49.392Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 409}}, "name": "domain.transfer:started", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "49901af0-569e-4acd-900f-6edf0ebc123c"}`
	payload := `{"data": {"domain": {"id": 6637, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 24, "auto_renew": false, "created_at": "2016-03-24T21:03:49.392Z", "expires_on": null, "updated_at": "2016-03-24T21:03:49.392Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 409}}, "name": "domain.transfer", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "49901af0-569e-4acd-900f-6edf0ebc123c"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.transfer", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseEmailForwardEvent_EmailForward_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/email_forward.create/example.http")

	event := &EmailForwardEvent{}
	err := ParseEmailForwardEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "email_forward.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example@example.zone", event.EmailForward.From; want != got {
		t.Errorf("ParseEvent EmailForward.From expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*EmailForwardEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseEmailForwardEvent_EmailForward_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/email_forward.delete/example.http")

	event := &EmailForwardEvent{}
	err := ParseEmailForwardEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "email_forward.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := ".*@example.zone", event.EmailForward.From; want != got {
		t.Errorf("ParseEvent EmailForward.From expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*EmailForwardEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseEmailForwardEvent_EmailForward_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/email_forward.update/example.http")

	event := &EmailForwardEvent{}
	err := ParseEmailForwardEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "email_forward.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := ".*@example.zone", event.EmailForward.From; want != got {
		t.Errorf("ParseEvent EmailForward.From expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*EmailForwardEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWebhookEvent_Webhook_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/webhook.create/example.http")

	event := &WebhookEvent{}
	err := ParseWebhookEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "webhook.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "https://xxxxxx-xxxxxxx-00000.herokuapp.com/xxxxxxxx", event.Webhook.URL; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WebhookEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWebhookEvent_Webhook_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/webhook.delete/example.http")

	event := &WebhookEvent{}
	err := ParseWebhookEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "webhook.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "https://xxxxxx-xxxxxxx-00000.herokuapp.com/xxxxxxxx", event.Webhook.URL; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WebhookEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Disable(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2016-01-17T17:10:41.187Z", "expires_on": "2017-01-17", "updated_at": "2016-01-17T17:11:19.797Z", "unicode_name": "example.com", "private_whois": true, "registrant_id": 2}, "whois_privacy": {"id": 3, "enabled": true, "domain_id": 1, "created_at": "2016-01-17T17:10:50.713Z", "expires_on": "2017-01-17", "updated_at": "2016-03-20T16:45:57.409Z"}}, "name": "whois_privacy.disable", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "e3861a08-a771-4049-abc4-715a3f7b7d6f"}`

	event := &WhoisPrivacyEvent{}
	err := ParseWhoisPrivacyEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.disable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(3), event.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WhoisPrivacyEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Enable(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2016-01-17T17:10:41.187Z", "expires_on": "2017-01-17", "updated_at": "2016-01-17T17:11:19.797Z", "unicode_name": "example.com", "private_whois": true, "registrant_id": 2}, "whois_privacy": {"id": 3, "enabled": true, "domain_id": 1, "created_at": "2016-01-17T17:10:50.713Z", "expires_on": "2017-01-17", "updated_at": "2016-03-20T16:45:57.409Z"}}, "name": "whois_privacy.enable", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "e3861a08-a771-4049-abc4-715a3f7b7d6f"}`

	event := &WhoisPrivacyEvent{}
	err := ParseWhoisPrivacyEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.enable", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(3), event.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WhoisPrivacyEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Purchase(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2016-01-17T17:10:41.187Z", "expires_on": "2017-01-17", "updated_at": "2016-01-17T17:11:19.797Z", "unicode_name": "example.com", "private_whois": true, "registrant_id": 2}, "whois_privacy": {"id": 3, "enabled": true, "domain_id": 1, "created_at": "2016-01-17T17:10:50.713Z", "expires_on": "2017-01-17", "updated_at": "2016-03-20T16:45:57.409Z"}}, "name": "whois_privacy.purchase", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "e3861a08-a771-4049-abc4-715a3f7b7d6f"}`

	event := &WhoisPrivacyEvent{}
	err := ParseWhoisPrivacyEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.purchase", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(3), event.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WhoisPrivacyEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWhoisPrivacyEvent_WhoisPrivacy_Renew(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2016-01-17T17:10:41.187Z", "expires_on": "2017-01-17", "updated_at": "2016-01-17T17:11:19.797Z", "unicode_name": "example.com", "private_whois": true, "registrant_id": 2}, "whois_privacy": {"id": 3, "enabled": true, "domain_id": 1, "created_at": "2016-01-17T17:10:50.713Z", "expires_on": "2017-01-17", "updated_at": "2016-03-20T16:45:57.409Z"}}, "name": "whois_privacy.renew", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "e3861a08-a771-4049-abc4-715a3f7b7d6f"}`

	event := &WhoisPrivacyEvent{}
	err := ParseWhoisPrivacyEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "whois_privacy.renew", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}
	if want, got := int64(3), event.WhoisPrivacy.ID; want != got {
		t.Errorf("ParseEvent WhoisPrivacy.ID expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WhoisPrivacyEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneEvent_Zone_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone.create/example.http")

	event := &ZoneEvent{}
	err := ParseZoneEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.zone", event.Zone.Name; want != got {
		t.Errorf("ParseEvent Zone.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ZoneEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneEvent_Zone_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone.delete/example.http")

	event := &ZoneEvent{}
	err := ParseZoneEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.zone", event.Zone.Name; want != got {
		t.Errorf("ParseEvent Zone.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ZoneEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Create(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone_record.create/example.http")

	event := &ZoneRecordEvent{}
	err := ParseZoneRecordEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone_record.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "", event.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent ZoneRecord.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ZoneRecordEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Update(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone_record.update/example.http")

	event := &ZoneRecordEvent{}
	err := ParseZoneRecordEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone_record.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "www", event.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent ZoneRecord.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ZoneRecordEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Delete(t *testing.T) {
	payload := getHttpRequestBodyFromFixture(t, "/webhooks/zone_record.delete/example.http")

	event := &ZoneRecordEvent{}
	err := ParseZoneRecordEvent(event, payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "zone_record.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "www", event.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent ZoneRecord.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse(payload)
	_, ok := parsedEvent.(*ZoneRecordEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}
