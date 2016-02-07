package webhook

import (
	"testing"
	"fmt"
)

func TestParsePayload(t *testing.T) {
	p1 := `{"data": {"domain": {"id": 229375, "name": "personal-weppos-domain.com", "state": "hosted", "token": "Alp8OJ60i7vbhyi7MqCOhsrZTw00bFyw", "account_id": 981, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "personal-weppos-domain.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1120, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	p2 := `{"data": {"domain": {"id": 229375, "name": "personal-weppos-domain.com", "state": "hosted", "token": "Alp8OJ60i7vbhyi7MqCOhsrZTw00bFyw", "account_id": 981, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "personal-weppos-domain.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1120, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.delete", "api_version": "v2", "request_identifier": "3e625f1c-3e8b-48fc-9326-9489f4b60e52"}`

	var payload *Payload
	var err error

	payload, err = ParsePayload([]byte(p1))
	fmt.Println(payload)
	fmt.Println(err)

	payload, err = ParsePayload([]byte(p2))
	fmt.Println(payload)
	fmt.Println(err)

	e1 := Parse([]byte(p1))
	fmt.Println("Event")
	fmt.Println(e1)

	ee1 := e1.(*DomainCreateEvent)
	fmt.Println(*ee1)
	fmt.Println(ee1.Domain)
}


