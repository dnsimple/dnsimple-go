package webhook

import (
	"testing"
)

func TestParse(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "personal-weppos-domain.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event, err := Parse([]byte(payload))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if want, got := "domain.create", event.Event(); want != got {
		t.Errorf("Parse event expected to be %v, got %v", want, got)
	}

	_, ok := event.(*DomainCreateEvent)
	if !ok {
		t.Fatalf("Parse returned error when typecasting: %v", err)
	}
}
