package webhook

import (
	"regexp"
	"testing"
)

var regexpUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func TestParseGenericEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event := &GenericEvent{}
	err := ParseGenericEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseGenericEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Name; want != got {
		t.Errorf("ParseGenericEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseGenericEvent requestID expected to be an UUID, got %v", event.RequestID)
	}

	data := event.Data.(map[string]interface{})
	if want, got := "example.com", data["domain"].(map[string]interface{})["name"]; want != got {
		t.Errorf("ParseDomainCreateEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalEnable(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-13T12:33:22.723Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11}}, "name": "domain.auto_renewal_enable", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "91d47480-c2ce-411c-ac95-b5b54f346bff"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainAutoRenewalEnableEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_enable", event.Name; want != got {
		t.Errorf("ParseDomainAutoRenewalEnableEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseDomainAutoRenewalEnableEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainAutoRenewalEnableEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_AutoRenewalDisable(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-13T12:33:22.723Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11}}, "name": "domain.auto_renewal_disable", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "91d47480-c2ce-411c-ac95-b5b54f346bff"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainAutoRenewalDisableEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renewal_disable", event.Name; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseDomainAutoRenewalDisableEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Create(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainCreateEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Name; want != got {
		t.Errorf("ParseDomainCreateEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseDomainCreateEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainCreateEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_Delete(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "name": "domain.delete", "api_version": "v2", "request_identifier": "3e625f1c-3e8b-48fc-9326-9489f4b60e52"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainDeleteEvent returned error: %v", err)
	}

	if want, got := "domain.delete", event.Name; want != got {
		t.Errorf("ParseDomainDeleteEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseDomainDeleteEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainDeleteEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainEvent_Domain_TokenReset(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-07T23:26:16.368Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11549}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "name": "domain.token_reset", "api_version": "v2", "request_identifier": "33537afb-0e99-49ec-b69e-93ffcc3db763"}`

	event := &DomainEvent{}
	err := ParseDomainEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainTokenResetEvent returned error: %v", err)
	}

	if want, got := "domain.token_reset", event.Name; want != got {
		t.Errorf("ParseDomainTokenResetEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseDomainTokenResetEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainTokenResetEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseWebhookEvent_Webhook_Create(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 25, "url": "https://webhook.test"}}, "name": "webhook.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "d6362e1f-310b-4009-a29d-ce76c849d32c"}`

	event := &WebhookEvent{}
	err := ParseWebhookEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseWebhookCreateEvent returned error: %v", err)
	}

	if want, got := "webhook.create", event.Name; want != got {
		t.Errorf("ParseWebhookCreateEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseWebhookCreateEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	//if want, got := "https://webhook.test", event.Webhook.URL; want != got {
	//	t.Errorf("ParseWebhookCreateEvent Webhook.URL expected to be %v, got %v", want, got)
	//}
}

func TestParseWebhookEvent_Webhook_Delete(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 23, "url": "https://webhook.test"}}, "name": "webhook.delete", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "756bad5c-b432-43be-821a-2f4c4f285d19"}`

	event := &WebhookEvent{}
	err := ParseWebhookEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseWebhookDeleteEvent returned error: %v", err)
	}

	if want, got := "webhook.delete", event.Name; want != got {
		t.Errorf("ParseWebhookDeleteEvent name expected to be %v, got %v", want, got)
	}
	if !regexpUUID.MatchString(event.RequestID) {
		t.Errorf("ParseWebhookDeleteEvent requestID expected to be an UUID, got %v", event.RequestID)
	}
	//if want, got := "https://webhook.test", event.Webhook.URL; want != got {
	//	t.Errorf("ParseWebhookCreateEvent Webhook.URL expected to be %v, got %v", want, got)
	//}
}
