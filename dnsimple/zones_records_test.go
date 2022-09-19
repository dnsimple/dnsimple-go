package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneRecordPath(t *testing.T) {
	assert.Equal(t, "/1010/zones/example.com/records", zoneRecordPath("1010", "example.com", 0))
	assert.Equal(t, "/1010/zones/example.com/records/1", zoneRecordPath("1010", "example.com", 1))
}

func TestZonesService_ListRecords(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listZoneRecords/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordsResponse, err := client.Zones.ListRecords(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 5}, recordsResponse.Pagination)
	records := recordsResponse.Data
	assert.Len(t, records, 5)
	assert.Equal(t, int64(1), records[0].ID)
	assert.Equal(t, "", records[0].Name)
	assert.Equal(t, []string{"global"}, records[0].Regions)
}

func TestZonesService_ListRecords_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listZoneRecords/success.http")

		testQuery(t, r, url.Values{
			"page":      []string{"2"},
			"per_page":  []string{"20"},
			"sort":      []string{"name,expiration:desc"},
			"name":      []string{"example"},
			"name_like": []string{"www"},
			"type":      []string{"A"},
		})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Zones.ListRecords(context.Background(), "1010", "example.com", &ZoneRecordListOptions{String("example"), String("www"), String("A"), ListOptions{Page: Int(2), PerPage: Int(20), Sort: String("name,expiration:desc")}})

	assert.NoError(t, err)
}

func TestZonesService_ListRecords_WithOptionsSomeBlank(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listZoneRecords/success.http")

		testQuery(t, r, url.Values{
			"page": []string{"2"},
			"sort": []string{"name,expiration:desc"},
			"name": []string{"example"},
			"type": []string{"A"},
		})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Zones.ListRecords(context.Background(), "1010", "example.com", &ZoneRecordListOptions{Name: String("example"), Type: String("A"), ListOptions: ListOptions{Page: Int(2), Sort: String("name,expiration:desc")}})

	assert.NoError(t, err)
}

func TestZonesService_CreateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createZoneRecord/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "foo", "content": "mxa.example.com", "type": "MX"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	recordValues := ZoneRecordAttributes{Name: String("foo"), Content: "mxa.example.com", Type: "MX"}

	recordResponse, err := client.Zones.CreateRecord(context.Background(), accountID, "example.com", recordValues)

	assert.NoError(t, err)
	record := recordResponse.Data
	assert.Equal(t, int64(1), record.ID)
	assert.Equal(t, "www", record.Name)
	assert.Equal(t, "A", record.Type)
	assert.Equal(t, []string{"global"}, record.Regions)
}

func TestZonesService_CreateRecord_BlankName(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createZoneRecord/created-apex.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "", "content": "127.0.0.1", "type": "A"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues := ZoneRecordAttributes{Name: String(""), Content: "127.0.0.1", Type: "A"}

	recordResponse, err := client.Zones.CreateRecord(context.Background(), "1010", "example.com", recordValues)

	assert.NoError(t, err)
	record := recordResponse.Data
	assert.Equal(t, "", record.Name)
	assert.Equal(t, []string{"global"}, record.Regions)
}

func TestZonesService_CreateRecord_Regions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	var recordValues ZoneRecordAttributes

	mux.HandleFunc("/v2/1/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createZoneRecord/created.http")

		want := map[string]interface{}{"name": "foo"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues = ZoneRecordAttributes{Name: String("foo"), Regions: []string{}}

	_, err := client.Zones.CreateRecord(context.Background(), "1", "example.com", recordValues)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/2/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createZoneRecord/created.http")

		want := map[string]interface{}{"name": "foo", "regions": []interface{}{"global"}}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues = ZoneRecordAttributes{Name: String("foo"), Regions: []string{"global"}}

	_, err = client.Zones.CreateRecord(context.Background(), "2", "example.com", recordValues)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/3/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createZoneRecord/created.http")

		want := map[string]interface{}{"name": "foo", "regions": []interface{}{"global"}}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues = ZoneRecordAttributes{Name: String("foo"), Regions: []string{"global"}}

	_, err = client.Zones.CreateRecord(context.Background(), "2", "example.com", recordValues)
	assert.NoError(t, err)
}

func TestZonesService_GetRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1539", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getZoneRecord/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	recordResponse, err := client.Zones.GetRecord(context.Background(), accountID, "example.com", 1539)

	assert.NoError(t, err)
	record := recordResponse.Data
	wantSingle := &ZoneRecord{
		ID:           5,
		ZoneID:       "example.com",
		ParentID:     0,
		Type:         "MX",
		Name:         "",
		Content:      "mxa.example.com",
		TTL:          600,
		Priority:     10,
		SystemRecord: false,
		Regions:      []string{"SV1", "IAD"},
		CreatedAt:    "2016-10-05T09:51:35Z",
		UpdatedAt:    "2016-10-05T09:51:35Z"}
	assert.Equal(t, wantSingle, record)
}

func TestZonesService_UpdateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/5", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateZoneRecord/success.http")

		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "foo", "content": "127.0.0.1"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	recordValues := ZoneRecordAttributes{Name: String("foo"), Content: "127.0.0.1"}

	recordResponse, err := client.Zones.UpdateRecord(context.Background(), accountID, "example.com", 5, recordValues)

	assert.NoError(t, err)
	record := recordResponse.Data
	assert.Equal(t, int64(5), record.ID)
	assert.Equal(t, "mxb.example.com", record.Content)
}

func TestZonesService_UpdateRecord_NameNotProvided(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/5", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateZoneRecord/success.http")

		want := map[string]interface{}{"content": "127.0.0.1"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues := ZoneRecordAttributes{Content: "127.0.0.1"}

	_, err := client.Zones.UpdateRecord(context.Background(), "1010", "example.com", 5, recordValues)

	assert.NoError(t, err)
}

func TestZonesService_UpdateRecord_Regions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	var recordValues ZoneRecordAttributes

	mux.HandleFunc("/v2/1/zones/example.com/records/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateZoneRecord/success.http")

		want := map[string]interface{}{"name": "foo"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues = ZoneRecordAttributes{Name: String("foo"), Regions: []string{}}

	_, err := client.Zones.UpdateRecord(context.Background(), "1", "example.com", 1, recordValues)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/2/zones/example.com/records/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateZoneRecord/success.http")

		want := map[string]interface{}{"name": "foo", "regions": []interface{}{"global"}}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues = ZoneRecordAttributes{Name: String("foo"), Regions: []string{"global"}}

	_, err = client.Zones.UpdateRecord(context.Background(), "2", "example.com", 1, recordValues)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/3/zones/example.com/records/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateZoneRecord/success.http")

		want := map[string]interface{}{"name": "foo", "regions": []interface{}{"global"}}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	recordValues = ZoneRecordAttributes{Name: String("foo"), Regions: []string{"global"}}

	_, err = client.Zones.UpdateRecord(context.Background(), "2", "example.com", 1, recordValues)
	assert.NoError(t, err)
}

func TestZonesService_DeleteRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteZoneRecord/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	_, err := client.Zones.DeleteRecord(context.Background(), accountID, "example.com", 2)

	assert.NoError(t, err)
}
