package webhook

import (
	"regexp"
	"testing"
)

var regexpUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

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

func TestParseDomainEvent_Domain_AutoRenewalEnable(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-13T12:33:22.723Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11}}, "name": "domain.auto_renewal_enable", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "91d47480-c2ce-411c-ac95-b5b54f346bff"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_enable", event.Name; want != got {
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

func TestParseDomainEvent_Domain_AutoRenewalDisable(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-13T12:33:22.723Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11}}, "name": "domain.auto_renewal_disable", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "91d47480-c2ce-411c-ac95-b5b54f346bff"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_disable", event.Name; want != got {
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

func TestParseDomainEvent_Domain_Create(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Name; want != got {
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

func TestParseDomainEvent_Domain_Delete(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.delete", "api_version": "v2", "request_identifier": "3e625f1c-3e8b-48fc-9326-9489f4b60e52"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.delete", event.Name; want != got {
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
	payload := `{"data": {"auto": true, "domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2014-04-01T08:37:15.729Z", "expires_on": "2017-04-01", "updated_at": "2016-03-04T07:40:02.334Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 2}}, "name": "domain.renew", "actor": {"id": "system", "entity": "dnsimple", "pretty": "support@dnsimple.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "9e8e47ef-f303-4455-b496-875f70ab5c00"}
`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
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
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseEvent Domain.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*DomainEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseDomainEvent_Domain_TokenReset(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-07T23:26:16.368Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11549}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.token_reset", "api_version": "v2", "request_identifier": "33537afb-0e99-49ec-b69e-93ffcc3db763"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "domain.token_reset", event.Name; want != got {
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

func TestParseContactEvent_Contact_Create(t *testing.T) {
	payload := `{"data": {"contact": {"id": 29032, "fax": "+39 339 1111111", "city": "Rome", "label": "Webhook", "phone": "+39 339 0000000", "country": "IT", "address1": "Some Street", "address2": "", "job_title": "Developer", "last_name": "Contact", "account_id": 981, "created_at": "2016-02-13T13:11:29.388Z", "first_name": "Example", "updated_at": "2016-02-13T13:11:29.388Z", "postal_code": "12037", "email_address": "example@example.com", "state_province": "Italy", "organization_name": "Company"}}, "name": "contact.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "3be0422c-8ca2-44d9-95d6-9f045b938781"}
`

	event := &ContactEvent{}
	err := ParseContactEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "Webhook", event.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ContactEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseContactEvent_Contact_Update(t *testing.T) {
	payload := `{"data": {"contact": {"id": 29032, "fax": "+39 339 1111111", "city": "Rome", "label": "Webhook", "phone": "+39 339 0000000", "country": "IT", "address1": "Some Street", "address2": "", "job_title": "Developer", "last_name": "Contact", "account_id": 981, "created_at": "2016-02-13T13:11:29.388Z", "first_name": "Example", "updated_at": "2016-02-13T13:11:29.388Z", "postal_code": "12037", "email_address": "example@example.com", "state_province": "Italy", "organization_name": "Company"}}, "name": "contact.update", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "3be0422c-8ca2-44d9-95d6-9f045b938781"}
`

	event := &ContactEvent{}
	err := ParseContactEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "Webhook", event.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ContactEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseContactEvent_Contact_Delete(t *testing.T) {
	payload := `{"data": {"contact": {"id": 29032, "fax": "+39 339 1111111", "city": "Rome", "label": "Webhook", "phone": "+39 339 0000000", "country": "IT", "address1": "Some Street", "address2": "", "job_title": "Developer", "last_name": "Contact", "account_id": 981, "created_at": "2016-02-13T13:11:29.388Z", "first_name": "Example", "updated_at": "2016-02-13T13:11:29.388Z", "postal_code": "12037", "email_address": "example@example.com", "state_province": "Italy", "organization_name": "Company"}}, "name": "contact.delete", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "3be0422c-8ca2-44d9-95d6-9f045b938781"}
`

	event := &ContactEvent{}
	err := ParseContactEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "contact.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "Webhook", event.Contact.Label; want != got {
		t.Errorf("ParseEvent Contact.Name expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ContactEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Create(t *testing.T) {
	payload := `{"data": {"record": {"id": 1, "ttl": 60, "name": "_frame", "type": "TXT", "content": "https://dnsimple.com/", "zone_id": "example.com", "priority": null, "parent_id": null, "created_at": "2016-02-22T21:06:48.957Z", "updated_at": "2016-02-22T21:23:22.503Z", "system_record": false}}, "name": "record.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "8f6cd405-2c87-453b-8b95-7a296982e4b8"}
`

	event := &ZoneRecordEvent{}
	err := ParseZoneRecordEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "record.create", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "_frame", event.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ZoneRecordEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Update(t *testing.T) {
	payload := `{"data": {"record": {"id": 1, "ttl": 60, "name": "_frame", "type": "TXT", "content": "https://dnsimple.com/", "zone_id": "example.com", "priority": null, "parent_id": null, "created_at": "2016-02-22T21:06:48.957Z", "updated_at": "2016-02-22T21:23:22.503Z", "system_record": false}}, "name": "record.update", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "8f6cd405-2c87-453b-8b95-7a296982e4b8"}
`

	event := &ZoneRecordEvent{}
	err := ParseZoneRecordEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "record.update", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "_frame", event.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ZoneRecordEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseZoneRecordEvent_ZoneRecord_Delete(t *testing.T) {
	payload := `{"data": {"record": {"id": 1, "ttl": 60, "name": "_frame", "type": "TXT", "content": "https://dnsimple.com/", "zone_id": "example.com", "priority": null, "parent_id": null, "created_at": "2016-02-22T21:06:48.957Z", "updated_at": "2016-02-22T21:23:22.503Z", "system_record": false}}, "name": "record.delete", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "8f6cd405-2c87-453b-8b95-7a296982e4b8"}
`

	event := &ZoneRecordEvent{}
	err := ParseZoneRecordEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "record.delete", event.Name; want != got {
		t.Errorf("ParseEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "_frame", event.ZoneRecord.Name; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*ZoneRecordEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWebhookEvent_Webhook_Create(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 25, "url": "https://webhook.test"}}, "name": "webhook.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "d6362e1f-310b-4009-a29d-ce76c849d32c"}`

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
	if want, got := "https://webhook.test", event.Webhook.URL; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WebhookEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}

func TestParseWebhookEvent_Webhook_Delete(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 23, "url": "https://webhook.test"}}, "name": "webhook.delete", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "756bad5c-b432-43be-821a-2f4c4f285d19"}`

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
	if want, got := "https://webhook.test", event.Webhook.URL; want != got {
		t.Errorf("ParseEvent Webhook.URL expected to be %v, got %v", want, got)
	}

	parsedEvent, err := Parse([]byte(payload))
	_, ok := parsedEvent.(*WebhookEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}
