package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRecords_recordPath(t *testing.T) {
	actual := recordPath("1", "example.com", nil)
	expected := "/1/zones/example.com/records"
	if actual != expected {
		t.Errorf("recordPath(\"1\", \"example.com\", nil): actual %s, expected %s", actual, expected)
	}

	actual = recordPath("1", "example.com", 2)
	expected = "/1/zones/example.com/records/2"
	if actual != expected {
		t.Errorf("recordPath(\"1\", \"example.com\", 2): actual %s, expected %s", actual, expected)
	}

	actual = recordPath("1", 1, nil)
	expected = "/1/zones/1/records"
	if actual != expected {
		t.Errorf("recordPath(\"1\", 1, nil): actual %s, expected %s", actual, expected)
	}

	actual = recordPath("1", 1, 2)
	expected = "/1/zones/1/records/2"
	if actual != expected {
		t.Errorf("recordPath(\"1\", 1, 2): actual %s, expected %s", actual, expected)
	}
}

func TestDomainsService_ListRecords(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":[{"id":64779,"zone_id":"example.com","parent_id":null,"name":"","content":"ns1.dnsimple.com admin.dnsimple.com 1452184205 86400 7200 604800 300","ttl":3600,"priority":null,"type":"SOA","system_record":true,"created_at":"2016-01-07T16:30:05.379Z","updated_at":"2016-01-07T16:30:05.379Z"},{"id":64780,"zone_id":"example.com","parent_id":null,"name":"","content":"ns1.dnsimple.com","ttl":3600,"priority":null,"type":"NS","system_record":true,"created_at":"2016-01-07T16:30:05.422Z","updated_at":"2016-01-07T16:30:05.422Z"},{"id":64781,"zone_id":"example.com","parent_id":null,"name":"","content":"ns2.dnsimple.com","ttl":3600,"priority":null,"type":"NS","system_record":true,"created_at":"2016-01-07T16:30:05.433Z","updated_at":"2016-01-07T16:30:05.433Z"},{"id":64782,"zone_id":"example.com","parent_id":null,"name":"","content":"ns3.dnsimple.com","ttl":3600,"priority":null,"type":"NS","system_record":true,"created_at":"2016-01-07T16:30:05.445Z","updated_at":"2016-01-07T16:30:05.445Z"},{"id":64783,"zone_id":"example.com","parent_id":null,"name":"","content":"ns4.dnsimple.com","ttl":3600,"priority":null,"type":"NS","system_record":true,"created_at":"2016-01-07T16:30:05.457Z","updated_at":"2016-01-07T16:30:05.457Z"}],"pagination":{"current_page":1,"per_page":30,"total_entries":5,"total_pages":1}}
		`)
	})

	accountID := "1010"

	recordsResponse, err := client.Zones.ListRecords(accountID, "example.com")
	if err != nil {
		t.Fatalf("Zones.ListRecords() returned error: %v", err)
	}

	records := recordsResponse.Data
	if want, got := 5, len(records); want != got {
		t.Errorf("Zones.ListRecords() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 64779, records[0].ID; want != got {
		t.Fatalf("Zones.ListRecords() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "", records[0].Name; want != got {
		t.Fatalf("Zones.ListRecords() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_CreateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "foo", "content": "192.168.0.10", "type": "A"}
		testRequestJSON(t, r, want)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `
			{"data":{"id":64784,"zone_id":"example.com","parent_id":null,"name":"www","content":"127.0.0.1","ttl":600,"priority":null,"type":"A","system_record":false,"created_at":"2016-01-07T17:45:13.653Z","updated_at":"2016-01-07T17:45:13.653Z"}}
		`)
	})

	accountID := "1010"
	recordValues := Record{Name: "foo", Content: "192.168.0.10", Type: "A"}

	recordResponse, err := client.Zones.CreateRecord(accountID, "example.com", recordValues)
	if err != nil {
		t.Fatalf("Zones.CreateRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	if want, got := 64784, record.ID; want != got {
		t.Fatalf("Zones.CreateRecord() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "www", record.Name; want != got {
		t.Fatalf("Zones.CreateRecord() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_GetRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1539", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprintf(w, `
			{"data":{"id":64784,"zone_id":"example.com","parent_id":null,"name":"www","content":"127.0.0.1","ttl":600,"priority":null,"type":"A","system_record":false,"created_at":"2016-01-07T17:45:13.653Z","updated_at":"2016-01-07T17:45:13.653Z"}}
		`)
	})

	accountID := "1010"

	recordResponse, err := client.Zones.GetRecord(accountID, "example.com", 1539)
	if err != nil {
		t.Fatalf("Zones.GetRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	wantSingle := &Record{
		ID:           64784,
		ZoneID:       "example.com",
		ParentID:     0,
		Type:         "A",
		Name:         "www",
		Content:      "127.0.0.1",
		TTL:          600,
		Priority:     0,
		SystemRecord: false,
		CreatedAt:    "2016-01-07T17:45:13.653Z",
		UpdatedAt:    "2016-01-07T17:45:13.653Z"}

	if !reflect.DeepEqual(record, wantSingle) {
		t.Fatalf("Zones.GetRecord() returned %+v, want %+v", record, wantSingle)
	}
}

func TestDomainsService_UpdateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		want := map[string]interface{}{"content": "192.168.0.10", "name": "bar"}
		testRequestJSON(t, r, want)

		fmt.Fprint(w, `
			{"data":{"id":64784,"domain_id":5841,"parent_id":null,"name":"www","content":"127.0.0.1","ttl":600,"priority":null,"type":"A","system_record":false,"created_at":"2016-01-07T17:45:13.653Z","updated_at":"2016-01-07T17:54:46.722Z"}}
		`)
	})

	accountID := "1010"
	recordValues := Record{Name: "bar", Content: "192.168.0.10"}

	recordResponse, err := client.Zones.UpdateRecord(accountID, "example.com", 2, recordValues)
	if err != nil {
		t.Fatalf("Zones.UpdateRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	if want, got := 64784, record.ID; want != got {
		t.Fatalf("Zones.UpdateRecord() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "www", record.Name; want != got {
		t.Fatalf("Zones.UpdateRecord() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_DeleteRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeaders(t, r)
	})

	accountID := "1010"

	_, err := client.Zones.DeleteRecord(accountID, "example.com", 2)
	if err != nil {
		t.Fatalf("Zones.DeleteRecord() returned error: %v", err)
	}
}
