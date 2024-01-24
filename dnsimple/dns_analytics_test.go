package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsAnalyticsService_Query(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/dns_analytics", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/dnsAnalytics/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	dnsAnalyticsResponse, err := client.DnsAnalytics.Query(context.Background(), "1", nil)

	assert.NoError(t, err)
	//assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, dnsAnalyticsResponse.Pagination)
	data := dnsAnalyticsResponse.Data
	assert.Len(t, data, 12)
	assert.Equal(t, int64(1200), data[0].Volume)
	assert.Equal(t, "2023-12-08", data[0].Date)
	assert.Equal(t, "bar.com", data[0].Zone)
}
