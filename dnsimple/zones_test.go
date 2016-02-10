package dnsimple

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestZonesService_ListZones(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture("/listZones/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	zonesResponse, err := client.Zones.ListZones(accountID)
	if err != nil {
		t.Fatalf("Zones.ListZones() returned error: %v", err)
	}

	zones := zonesResponse.Data
	if want, got := 2, len(zones); want != got {
		t.Errorf("Zones.ListZones() expected to return %v zones, got %v", want, got)
	}

	if want, got := 1, zones[0].ID; want != got {
		t.Fatalf("Zones.ListZones() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", zones[0].Name; want != got {
		t.Fatalf("Zones.ListZones() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestZonesService_GetZone(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture("/getZone/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneResponse, err := client.Zones.GetZone(accountID, zoneName)
	if err != nil {
		t.Fatalf("Zones.GetZone() returned error: %v", err)
	}

	zone := zoneResponse.Data
	wantSingle := &Zone{
		ID:        1,
		AccountID: 1010,
		Name:      "example-alpha.com",
		Reverse:   false,
		CreatedAt: "2015-04-23T07:40:03.045Z",
		UpdatedAt: "2015-04-23T07:40:03.051Z"}

	if !reflect.DeepEqual(zone, wantSingle) {
		t.Fatalf("Zones.GetZone() returned %+v, want %+v", zone, wantSingle)
	}
}
