package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZonesService_CheckZoneDistribution(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneDistribution/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneDistributionResponse, err := client.Zones.CheckZoneDistribution(context.Background(), accountID, zoneName)

	assert.NoError(t, err)
	zone := zoneDistributionResponse.Data
	wantSingle := &ZoneDistribution{
		Distributed: true,
	}
	assert.Equal(t, wantSingle, zone)
}

func TestZonesService_CheckZoneDistributionFailure(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneDistribution/failure.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneDistributionResponse, err := client.Zones.CheckZoneDistribution(context.Background(), accountID, zoneName)

	assert.NoError(t, err)
	zone := zoneDistributionResponse.Data
	wantSingle := &ZoneDistribution{
		Distributed: false,
	}
	assert.Equal(t, wantSingle, zone)
}

func TestZonesService_CheckZoneDistributionError(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneDistribution/error.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"

	zoneDistributionResponse, err := client.Zones.CheckZoneDistribution(context.Background(), accountID, zoneName)

	assert.Error(t, err)
	assert.Nil(t, zoneDistributionResponse)
}

func TestZonesService_CheckZoneRecordDistribution(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneRecordDistribution/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"
	recordID := int64(1)

	zoneDistributionResponse, err := client.Zones.CheckZoneRecordDistribution(context.Background(), accountID, zoneName, recordID)

	assert.NoError(t, err)
	zone := zoneDistributionResponse.Data
	assert.Equal(t, &ZoneDistribution{Distributed: true}, zone)
}

func TestZonesService_CheckZoneRecordDistributionFailure(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneRecordDistribution/failure.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"
	recordID := int64(1)

	zoneDistributionResponse, err := client.Zones.CheckZoneRecordDistribution(context.Background(), accountID, zoneName, recordID)

	assert.NoError(t, err)
	zone := zoneDistributionResponse.Data
	assert.Equal(t, &ZoneDistribution{Distributed: false}, zone)
}

func TestZonesService_CheckZoneRecordDistributionError(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1/distribution", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/checkZoneRecordDistribution/error.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	zoneName := "example.com"
	recordID := int64(1)

	zoneDistributionResponse, err := client.Zones.CheckZoneRecordDistribution(context.Background(), accountID, zoneName, recordID)

	assert.Error(t, err)
	assert.Nil(t, zoneDistributionResponse)
}
