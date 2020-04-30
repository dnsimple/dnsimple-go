package dnsimple

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestZonesService_CheckZoneDistribution(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneDistribution/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneDistributionResponse, err := client.Zones.CheckZoneDistribution(context.Background(), accountID, zoneName)
	if err != nil {
		t.Fatalf("Zones.CheckZoneDistribution() returned error: %v", err)
	}

	zone := zoneDistributionResponse.Data
	wantSingle := &ZoneDistribution{
		Distributed: true,
	}

	if !reflect.DeepEqual(zone, wantSingle) {
		t.Fatalf("Zones.CheckZoneDistribution() returned %+v, want %+v", zone, wantSingle)
	}
}

func TestZonesService_CheckZoneDistributionFailure(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneDistribution/failure.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneDistributionResponse, err := client.Zones.CheckZoneDistribution(context.Background(), accountID, zoneName)
	if err != nil {
		t.Fatalf("Zones.CheckZoneDistribution() returned error: %v", err)
	}

	zone := zoneDistributionResponse.Data
	wantSingle := &ZoneDistribution{
		Distributed: false,
	}

	if !reflect.DeepEqual(zone, wantSingle) {
		t.Fatalf("Zones.CheckZoneDistribution() returned %+v, want %+v", zone, wantSingle)
	}
}

func TestZonesService_CheckZoneDistributionError(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneDistribution/error.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneDistributionResponse, err := client.Zones.CheckZoneDistribution(context.Background(), accountID, zoneName)
	if err == nil {
		t.Fatalf("Zones.CheckZoneDistribution() expected to return an error: %v", zoneDistributionResponse)
	}

	if zoneDistributionResponse != nil {
		t.Fatalf("Zones.CheckZoneDistribution() expected to return a nil response: %v", zoneDistributionResponse)
	}
}

func TestZonesService_CheckZoneRecordDistribution(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneRecordDistribution/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"
	recordID := int64(1)

	zoneDistributionResponse, err := client.Zones.CheckZoneRecordDistribution(context.Background(), accountID, zoneName, recordID)
	if err != nil {
		t.Fatalf("Zones.CheckZoneRecordDistribution() returned error: %v", err)
	}

	zone := zoneDistributionResponse.Data
	wantSingle := &ZoneDistribution{
		Distributed: true,
	}

	if !reflect.DeepEqual(zone, wantSingle) {
		t.Fatalf("Zones.CheckZoneRecordDistribution() returned %+v, want %+v", zone, wantSingle)
	}
}

func TestZonesService_CheckZoneRecordDistributionFailure(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneRecordDistribution/failure.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"
	recordID := int64(1)

	zoneDistributionResponse, err := client.Zones.CheckZoneRecordDistribution(context.Background(), accountID, zoneName, recordID)
	if err != nil {
		t.Fatalf("Zones.CheckZoneRecordDistribution() returned error: %v", err)
	}

	zone := zoneDistributionResponse.Data
	wantSingle := &ZoneDistribution{
		Distributed: false,
	}

	if !reflect.DeepEqual(zone, wantSingle) {
		t.Fatalf("Zones.CheckZoneRecordDistribution() returned %+v, want %+v", zone, wantSingle)
	}
}

func TestZonesService_CheckZoneRecordDistributionError(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneRecordDistribution/error.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"
	recordID := int64(1)

	zoneDistributionResponse, err := client.Zones.CheckZoneRecordDistribution(context.Background(), accountID, zoneName, recordID)
	if err == nil {
		t.Fatalf("Zones.CheckZoneRecordDistribution() expected to return an error: %v", zoneDistributionResponse)
	}

	if zoneDistributionResponse != nil {
		t.Fatalf("Zones.CheckZoneRecordDistribution() expected to return a nil response: %v", zoneDistributionResponse)
	}
}
