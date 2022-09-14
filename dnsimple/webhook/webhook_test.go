package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEvent_Generic(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "generic", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := ParseEvent([]byte(payload))

	assert.NoError(t, err)
	assert.Equal(t, "generic", event.Name)
	eventAccount := event.Account
	assert.Equal(t, "User", eventAccount.Display)
	_, ok := event.GetData().(*GenericEventData)
	assert.True(t, ok)
}

func TestParseEvent_Account(t *testing.T) {
	payload := getHTTPRequestBodyFromFixture(t, "/webhooks/account.update/example.http")

	event, err := ParseEvent(payload)

	assert.NoError(t, err)
	assert.Equal(t, "account.update", event.Name)
	eventAccount := event.Account
	assert.Equal(t, "Personal2", eventAccount.Display)
	_, ok := event.GetData().(*AccountEventData)
	assert.True(t, ok)
}

func TestParseEvent_Webhook(t *testing.T) {
	payload := `{"data": {"webhook": {"id": 25, "url": "https://webhook.test"}}, "name": "webhook.create", "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1, "display": "User", "identifier": "user"}, "api_version": "v2", "request_identifier": "d6362e1f-310b-4009-a29d-ce76c849d32c"}`

	event, err := ParseEvent([]byte(payload))

	assert.NoError(t, err)
	assert.Equal(t, "webhook.create", event.Name)
	eventAccount := event.Account
	assert.Equal(t, "User", eventAccount.Display)
	_, ok := event.GetData().(*WebhookEventData)
	assert.True(t, ok)
}
