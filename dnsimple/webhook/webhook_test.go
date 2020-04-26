package webhook

import (
	"testing"
)

func TestParseEvent_Generic(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "generic", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "generic", event.Name; want != got {
		t.Errorf("ParseEvent event Name expected to be %v, got %v", want, got)
	}

	eventAccount := event.Account
	if want, got := "User", eventAccount.Display; want != got {
		t.Errorf("ParseEvent event Account.Display expected to be %v, got %v", want, got)
	}

	_, ok := event.GetData().(*GenericEventData)
	if !ok {
		t.Fatalf("ParseEvent returned error when typecasting: %v", err)
	}
}

func TestParseEvent_Account(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.update/example.http")

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "account.update", event.Name; want != got {
		t.Errorf("ParseEvent event Name expected to be %v, got %v", want, got)
	}

	eventAccount := event.Account
	if want, got := "Personal2", eventAccount.Display; want != got {
		t.Errorf("ParseEvent event Account.Display expected to be %v, got %v", want, got)
	}

	_, ok := event.GetData().(*AccountEventData)
	if !ok {
		t.Fatalf("ParseEvent returned error when typecasting: %v", err)
	}
}

func TestParseEvent_Webhook(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 25, "url": "https://webhook.test"}}, "name": "webhook.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "d6362e1f-310b-4009-a29d-ce76c849d32c"}`

	event, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}

	if want, got := "webhook.create", event.Name; want != got {
		t.Errorf("ParseEvent event Name expected to be %v, got %v", want, got)
	}

	eventAccount := event.Account
	if want, got := "User", eventAccount.Display; want != got {
		t.Errorf("ParseEvent event Account.Display expected to be %v, got %v", want, got)
	}

	_, ok := event.GetData().(*WebhookEventData)
	if !ok {
		t.Fatalf("ParseEvent returned error when typecasting: %v", err)
	}
}
