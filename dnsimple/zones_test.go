package dnsimple

import (
	"io"
	"net/http"
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
