package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRecord_recordPath(t *testing.T) {
	var pathTest = []struct {
		domainInput interface{}
		recordInput *Record
		expected    string
	}{
		{"example.com", nil, "domains/example.com/records"},
		{"example.com", &Record{Id: 42}, "domains/example.com/records/42"},
		{23, nil, "domains/23/records"},
		{23, &Record{Id: 42}, "domains/23/records/42"},
		{Domain{Id: 23}, nil, "domains/23/records"},
		{Domain{Id: 23}, &Record{Id: 42}, "domains/23/records/42"},
	}

	for _, pt := range pathTest {
		actual := recordPath(pt.domainInput, pt.recordInput)
		if actual != pt.expected {
			t.Errorf("recordPath(%+v, %+v): expected %s, actual %s", pt.domainInput, pt.recordInput, pt.expected, actual)
		}
	}
}

func TestRecord_Records_All(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"record":{"id":1, "name":"foo.example.com"}}]`)
	})

	records, err := client.Records("example.com", "", "")

	if err != nil {
		t.Errorf("Records returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Records returned %+v, want %+v", records, want)
	}
}

func TestRecord_Records_Subdomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "foo"})

		fmt.Fprint(w, `[{"record":{"id":1, "name":"foo.example.com"}}]`)
	})

	records, err := client.Records("example.com", "foo", "")

	if err != nil {
		t.Errorf("Records returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Records returned %+v, want %+v", records, want)
	}
}

func TestRecord_Records_RecordType(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "foo", "type": "CNAME"})

		fmt.Fprint(w, `[{"record":{"id":1, "name":"foo.example.com"}}]`)
	})

	records, err := client.Records("example.com", "foo", "CNAME")

	if err != nil {
		t.Errorf("Records returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Records returned %+v, want %+v", records, want)
	}
}

func TestRecord_CreateRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["record"] = map[string]interface{}{"name": "foo", "content": "192.168.0.10", "record_type": "A"}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		fmt.Fprintf(w, `{"record":{"id":42, "name":"foo"}}`)
	})

	recordValues := Record{Name: "foo", Content: "192.168.0.10", RecordType: "A"}
	record, err := client.CreateRecord("example.com", recordValues)

	if err != nil {
		t.Errorf("CreateRecord returned error: %v", err)
	}

	want := Record{Id: 42, Name: "foo"}
	if !reflect.DeepEqual(record, want) {
		t.Errorf("CreateRecord returned %+v, want %+v", record, want)
	}
}

func TestRecord_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/24/records/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	record := Record{Id: 42, DomainId: 24}
	err := record.Delete(client)

	if err != nil {
		t.Errorf("Delete returned error: %v", err)
	}
}

func TestRecord_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/24/records/42", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		// TODO: there's a problem when verifying adding prio and ttl integers. Why?
		want["record"] = map[string]interface{}{"content": "192.168.0.10", "name": "bar"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"record":{"id":24, "domain_id":42}}`)
	})

	record := Record{Id: 42, DomainId: 24, Name: "foo"}
	recordAttributes := Record{Name: "bar", Content: "192.168.0.10"}

	_, err := record.Update(client, recordAttributes)

	if err != nil {
		t.Errorf("Update returned error: %v", err)
	}
}

func TestRecord_UpdateIP(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/24/records/42", func(w http.ResponseWriter, r *http.Request) {
		want := make(map[string]interface{})
		want["record"] = map[string]interface{}{"name": "foo", "content": "192.168.0.1"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"record":{"id":24, "domain_id":42}}`)
	})

	record := Record{Id: 42, DomainId: 24, Name: "foo"}
	err := record.UpdateIP(client, "192.168.0.1")

	if err != nil {
		t.Errorf("UpdateIP returned error: %v", err)
	}
}
