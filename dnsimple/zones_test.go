package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZonesService_ListZones(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listZones/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	zonesResponse, err := client.Zones.ListZones(context.Background(), "1010", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, zonesResponse.Pagination)
	zones := zonesResponse.Data
	assert.Len(t, zones, 2)
	assert.Equal(t, int64(1), zones[0].ID)
	assert.Equal(t, "example-alpha.com", zones[0].Name)
}

func TestZonesService_ListZones_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listZones/success.http")

		testQuery(t, r, url.Values{
			"page":      []string{"2"},
			"per_page":  []string{"20"},
			"sort":      []string{"name,expiration:desc"},
			"name_like": []string{"example"},
		})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Zones.ListZones(context.Background(), "1010", &ZoneListOptions{String("example"), ListOptions{Page: Int(2), PerPage: Int(20), Sort: String("name,expiration:desc")}})

	assert.NoError(t, err)
}

func TestZonesService_GetZone(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getZone/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneResponse, err := client.Zones.GetZone(context.Background(), accountID, zoneName)

	assert.NoError(t, err)
	zone := zoneResponse.Data
	wantSingle := &Zone{
		ID:                1,
		AccountID:         1010,
		Name:              "example-alpha.com",
		Reverse:           false,
		Secondary:         false,
		LastTransferredAt: "",
		Active:            true,
		CreatedAt:         "2015-04-23T07:40:03Z",
		UpdatedAt:         "2015-04-23T07:40:03Z",
	}
	assert.Equal(t, wantSingle, zone)
}

func TestZonesService_GetZoneFile(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/file", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getZoneFile/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneFileResponse, err := client.Zones.GetZoneFile(context.Background(), accountID, zoneName)

	assert.NoError(t, err)
	zoneFile := zoneFileResponse.Data
	wantSingle := &ZoneFile{
		Zone: "$ORIGIN example.com.\n$TTL 1h\nexample.com. 3600 IN SOA ns1.dnsimple.com. admin.dnsimple.com. 1453132552 86400 7200 604800 300\nexample.com. 3600 IN NS ns1.dnsimple.com.\nexample.com. 3600 IN NS ns2.dnsimple.com.\nexample.com. 3600 IN NS ns3.dnsimple.com.\nexample.com. 3600 IN NS ns4.dnsimple.com.\n",
	}
	assert.Equal(t, wantSingle, zoneFile)
}

func TestZonesService_ActivateZoneDns(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/activation", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/activateZoneService/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneResponse, err := client.Zones.ActivateZoneDns(context.Background(), accountID, zoneName)

	assert.NoError(t, err)
	zone := zoneResponse.Data
	wantSingle := &Zone{
		ID:                1,
		AccountID:         1010,
		Name:              "example.com",
		Reverse:           false,
		Secondary:         false,
		LastTransferredAt: "",
		Active:            true,
		CreatedAt:         "2015-04-23T07:40:03Z",
		UpdatedAt:         "2015-04-23T07:40:03Z",
	}
	assert.Equal(t, wantSingle, zone)
}

func TestZonesService_DeactivateZoneDns(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/activation", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deactivateZoneService/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneResponse, err := client.Zones.DeactivateZoneDns(context.Background(), accountID, zoneName)

	assert.NoError(t, err)
	zone := zoneResponse.Data
	wantSingle := &Zone{
		ID:                1,
		AccountID:         1010,
		Name:              "example.com",
		Reverse:           false,
		Secondary:         false,
		LastTransferredAt: "",
		Active:            false,
		CreatedAt:         "2015-04-23T07:40:03Z",
		UpdatedAt:         "2015-04-23T07:40:03Z",
	}
	assert.Equal(t, wantSingle, zone)
}
