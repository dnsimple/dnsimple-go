package webhook

import (
	"testing"
)

func TestParse(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 23, "url": "https://test.host"}}, "name": "webhook.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": "1", "display": "Simone", "identifier": "example@example.com"}, "api_version": "v2", "request_identifier": "2f1cd735-0c02-4b1c-aa9d-20300520e62f"}`

	event, err := Parse([]byte(payload))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if want, got := "webhook.create", event.EventName(); want != got {
		t.Errorf("Parse event Name expected to be %v, got %v", want, got)
	}

	eventAccount := event.EventHeader().Account
	if want, got := "Simone", eventAccount.Display; want != got {
		t.Errorf("Parse event Account.Display expected to be %v, got %v", want, got)
	}

	_, ok := event.(*WebhookEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}
