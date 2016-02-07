package webhook

import (
	"testing"
)

func TestParseGenericEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 229375, "name": "example.com", "state": "hosted", "token": "Alp8OJ60i7vbhyi7MqCOhsrZTw00bFyw", "account_id": 981, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1120, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := ParseGenericEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseGenericEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Event; want != got {
		t.Errorf("ParseGenericEvent event expected to be %v, got %v", want, got)
	}
	if want, got := "096bfc29-2bf0-40c6-991b-f03b1f8521f1", event.RequestID; want != got {
		t.Errorf("ParseGenericEvent requestID expected to be %v, got %v", want, got)
	}

	data := event.Data.(map[string]interface{})
	if want, got := "example.com", data["domain"].(map[string]interface{})["name"]; want != got {
		t.Errorf("ParseDomainCreateEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainCreateEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 229375, "name": "example.com", "state": "hosted", "token": "Alp8OJ60i7vbhyi7MqCOhsrZTw00bFyw", "account_id": 981, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1120, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := ParseDomainCreateEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainCreateEvent returned error: %v", err)
	}

	if want, got := "domain.create", event.Event; want != got {
		t.Errorf("ParseDomainCreateEvent event expected to be %v, got %v", want, got)
	}
	if want, got := "096bfc29-2bf0-40c6-991b-f03b1f8521f1", event.RequestID; want != got {
		t.Errorf("ParseDomainCreateEvent requestID expected to be %v, got %v", want, got)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainCreateEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainDeleteEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 229375, "name": "example.com", "state": "hosted", "token": "Alp8OJ60i7vbhyi7MqCOhsrZTw00bFyw", "account_id": 981, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1120, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.delete", "api_version": "v2", "request_identifier": "3e625f1c-3e8b-48fc-9326-9489f4b60e52"}`

	event, err := ParseDomainDeleteEvent([]byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainDeleteEvent returned error: %v", err)
	}

	if want, got := "domain.delete", event.Event; want != got {
		t.Errorf("ParseDomainDeleteEvent event expected to be %v, got %v", want, got)
	}
	if want, got := "3e625f1c-3e8b-48fc-9326-9489f4b60e52", event.RequestID; want != got {
		t.Errorf("ParseDomainDeleteEvent requestID expected to be %v, got %v", want, got)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainDeleteEvent Domain.Name expected to be %v, got %v", want, got)
	}
}
