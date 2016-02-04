package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestRecords_recordPath(t *testing.T) {
	actual := recordPath("1", "example.com", nil)
	expected := "1/zones/example.com/records"
	if actual != expected {
		t.Errorf("recordPath(\"1\", \"example.com\", nil): actual %s, expected %s", actual, expected)
	}

	actual = recordPath("1", "example.com", 2)
	expected = "1/zones/example.com/records/2"
	if actual != expected {
		t.Errorf("recordPath(\"1\", \"example.com\", 2): actual %s, expected %s", actual, expected)
	}

	actual = recordPath("1", 1, nil)
	expected = "1/zones/1/records"
	if actual != expected {
		t.Errorf("recordPath(\"1\", 1, nil): actual %s, expected %s", actual, expected)
	}

	actual = recordPath("1", 1, 2)
	expected = "1/zones/1/records/2"
	if actual != expected {
		t.Errorf("recordPath(\"1\", 1, 2): actual %s, expected %s", actual, expected)
	}
}

func TestDomainsService_ListRecords_all(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data":[{"id":1, "name":"foo.example.com"}]}`)
	})

	accountId := "1"
	records, _, err := client.Domains.ListRecords(accountId, "example.com", "", "")

	if err != nil {
		t.Errorf("Domains.ListRecords returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Fatalf("Domains.ListRecords returned %+v, want %+v", records, want)
	}
}

func TestDomainsService_ListRecords_subdomain(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "foo"})

		fmt.Fprint(w, `{"data":[{"id":1, "name":"foo.example.com"}]}`)
	})

	accountId := "1"
	records, _, err := client.Domains.ListRecords(accountId, "example.com", "foo", "")

	if err != nil {
		t.Errorf("Domains.ListRecords returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Fatalf("Domains.ListRecords returned %+v, want %+v", records, want)
	}
}

func TestDomainsService_ListRecords_type(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"name": "foo", "type": "CNAME"})

		fmt.Fprint(w, `{"data":[{"id":1, "name":"foo.example.com"}]}`)
	})

	accountId := "1"
	records, _, err := client.Domains.ListRecords(accountId, "example.com", "foo", "CNAME")

	if err != nil {
		t.Errorf("Domains.ListRecords returned error: %v", err)
	}

	want := []Record{{Id: 1, Name: "foo.example.com"}}
	if !reflect.DeepEqual(records, want) {
		t.Fatalf("Domains.ListRecords returned %+v, want %+v", records, want)
	}
}

func TestDomainsService_CreateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		want := map[string]interface{}{"name": "foo", "content": "192.168.0.10", "record_type": "A"}

		testMethod(t, r, "POST")
		testRequestJSON(t, r, want)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"data":{"id":2, "zone_id":1, "name":"foo"}}`)
	})

	accountId := "1"
	recordValues := Record{Name: "foo", Content: "192.168.0.10", Type: "A"}
	record, _, err := client.Domains.CreateRecord(accountId, "example.com", recordValues)

	if err != nil {
		t.Errorf("Domains.CreateRecord returned error: %v", err)
	}

	want := Record{Id: 2, ZoneId: 1, Name: "foo"}
	if !reflect.DeepEqual(record, want) {
		t.Fatalf("Domains.CreateRecord returned %+v, want %+v", record, want)
	}
}

func TestDomainsService_GetRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records/1539", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"data":{"id":2, "zone_id":1, "name":"foo"}}`)
	})

	accountId := "1"
	record, _, err := client.Domains.GetRecord(accountId, "example.com", 1539)

	if err != nil {
		t.Errorf("Domains.GetRecord returned error: %v", err)
	}

	want := Record{Id: 2, ZoneId: 1, Name: "foo"}
	if !reflect.DeepEqual(record, want) {
		t.Fatalf("Domains.GetRecord returned %+v, want %+v", record, want)
	}
}

func TestDomainsService_UpdateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		want := map[string]interface{}{"content": "192.168.0.10", "name": "bar"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"data":{"id":2, "zone_id":1, "name":"bar", "content": "192.168.0.10"}}`)
	})

	accountId := "1"
	recordValues := Record{Name: "bar", Content: "192.168.0.10", Type: "A"}
	record, _, err := client.Domains.UpdateRecord(accountId, "example.com", 2, recordValues)

	if err != nil {
		t.Errorf("Domains.UpdateRecord returned error: %v", err)
	}

	want := Record{Id: 2, ZoneId: 1, Name: "bar", Content: "192.168.0.10"}
	if !reflect.DeepEqual(record, want) {
		t.Fatalf("Domains.UpdateRecord returned %+v, want %+v", record, want)
	}
}

func TestDomainsService_DeleteRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		// fmt.Fprint(w, `{}`)
	})

	accountId := "1"
	_, err := client.Domains.DeleteRecord(accountId, "example.com", 2)

	if err != nil {
		t.Errorf("Domains.DeleteRecord returned error: %v", err)
	}
}

func TestDomainsService_DeleteRecord_failed(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"Invalid request"}`)
	})

	accountId := "1"
	_, err := client.Domains.DeleteRecord(accountId, "example.com", 2)
	if err == nil {
		t.Errorf("Domains.DeleteRecord expected error to be returned")
	}

	if match := "400 Invalid request"; !strings.Contains(err.Error(), match) {
		t.Errorf("Records.Delete returned %+v, should match %+v", err, match)
	}
}

func TestRecord_UpdateIP(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/zones/24/records/42", func(w http.ResponseWriter, r *http.Request) {
		want := map[string]interface{}{"name": "foo", "content": "192.168.0.1"}

		testMethod(t, r, "PUT")
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `{"data":{"id":24, "domain_id":42}}`)
	})

	accountId := "1"
	record := Record{Id: 42, ZoneId: 24, Name: "foo"}
	err := record.UpdateIP(client, "192.168.0.1", accountId)

	if err != nil {
		t.Errorf("UpdateIP returned error: %v", err)
	}
}
