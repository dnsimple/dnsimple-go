package webhook

import (
	"testing"
)

func TestParseGenericEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event := &GenericEvent{}
	err := ParseGenericEvent(event, []byte(payload))
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
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	event := &DomainCreateEvent{}
	err := ParseDomainCreateEvent(event, []byte(payload))
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
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.delete", "api_version": "v2", "request_identifier": "3e625f1c-3e8b-48fc-9326-9489f4b60e52"}`

	event := &DomainDeleteEvent{}
	err := ParseDomainDeleteEvent(event, []byte(payload))
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

func TestParseDomainTokenResetEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-07T23:26:16.368Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11549}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.token_reset", "api_version": "v2", "request_identifier": "33537afb-0e99-49ec-b69e-93ffcc3db763"}`

	event := &DomainTokenResetEvent{}
	err := ParseDomainTokenResetEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainTokenResetEvent returned error: %v", err)
	}

	if want, got := "domain.token_reset", event.Event; want != got {
		t.Errorf("ParseDomainTokenResetEvent event expected to be %v, got %v", want, got)
	}
	if want, got := "33537afb-0e99-49ec-b69e-93ffcc3db763", event.RequestID; want != got {
		t.Errorf("ParseDomainTokenResetEvent requestID expected to be %v, got %v", want, got)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainTokenResetEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainAutoRenewalEnableEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": true, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-07T23:25:58.922Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11549}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.auto_renew_enable", "api_version": "v2", "request_identifier": "778a0c35-f9ed-4be9-a7a3-8c695f7872b6"}`

	event := &DomainAutoRenewalDisableEvent{}
	err := ParseDomainAutoRenewalDisableEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainAutoRenewalDisableEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renew_enable", event.Event; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent event expected to be %v, got %v", want, got)
	}
	if want, got := "778a0c35-f9ed-4be9-a7a3-8c695f7872b6", event.RequestID; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent requestID expected to be %v, got %v", want, got)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent Domain.Name expected to be %v, got %v", want, got)
	}
}

func TestParseDomainAutoRenewalDisableEvent(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2013-05-17T12:58:57.459Z", "expires_on": "2016-05-17", "updated_at": "2016-02-07T23:26:04.851Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 11549}}, "actor": {"id": 1111, "entity": "user", "pretty": "weppos@weppos.net"}, "action": "domain.auto_renew_disable", "api_version": "v2", "request_identifier": "394863e8-7669-4d92-98ab-372ce2f18dc1"}`

	event := &DomainAutoRenewalDisableEvent{}
	err := ParseDomainAutoRenewalDisableEvent(event, []byte(payload))
	if err != nil {
		t.Fatalf("ParseDomainAutoRenewalDisableEvent returned error: %v", err)
	}

	if want, got := "domain.auto_renew_disable", event.Event; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent event expected to be %v, got %v", want, got)
	}
	if want, got := "394863e8-7669-4d92-98ab-372ce2f18dc1", event.RequestID; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent requestID expected to be %v, got %v", want, got)
	}
	if want, got := "example.com", event.Domain.Name; want != got {
		t.Errorf("ParseDomainAutoRenewalDisableEvent Domain.Name expected to be %v, got %v", want, got)
	}
}
