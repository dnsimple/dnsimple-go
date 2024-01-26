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
	assert.Equal(t, &Pagination{CurrentPage: 0, PerPage: 100, TotalPages: 1, TotalEntries: 93}, dnsAnalyticsResponse.Pagination)
	data := dnsAnalyticsResponse.Data
	assert.Len(t, data, 12)
	assert.Equal(t, int64(1200), data[0].Volume)
	assert.Equal(t, "2023-12-08", data[0].Date)
	assert.Equal(t, "bar.com", data[0].ZoneName)
}

func TestDnsAnalyticsService_Query_SupportsFiltering(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/dns_analytics", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/dnsAnalytics/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		expectedQueryParameters := url.Values{}
		expectedQueryParameters.Add("start_date", "2023-10-01")
		expectedQueryParameters.Add("end_date", "2023-11-01")
		testQuery(t, r, expectedQueryParameters)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, _ = client.DnsAnalytics.Query(context.Background(), "1", &DnsAnalyticsOptions{StartDate: String("2023-10-01"), EndDate: String("2023-11-01")})
}

func TestDnsAnalyticsService_Query_SupportsSorting(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/dns_analytics", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/dnsAnalytics/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		expectedQueryParameters := url.Values{}
		expectedQueryParameters.Add("sort", "date:desc,zone_name:asc")
		testQuery(t, r, expectedQueryParameters)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	options := DnsAnalyticsOptions{}
	options.Sort = String("date:desc,zone_name:asc")
	_, _ = client.DnsAnalytics.Query(context.Background(), "1", &options)
}

func TestDnsAnalyticsService_Query_SupportsPagination(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/dns_analytics", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/dnsAnalytics/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		expectedQueryParameters := url.Values{}
		expectedQueryParameters.Add("page", "33")
		expectedQueryParameters.Add("per_page", "200")
		testQuery(t, r, expectedQueryParameters)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	options := DnsAnalyticsOptions{}
	options.Page = Int(33)
	options.PerPage = Int(200)
	_, _ = client.DnsAnalytics.Query(context.Background(), "1", &options)
}

func TestDnsAnalyticsService_Query_SupportsGrouping(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/dns_analytics", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/dnsAnalytics/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		expectedQueryParameters := url.Values{}
		expectedQueryParameters.Add("groupings", "zone_name,date")
		testQuery(t, r, expectedQueryParameters)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, _ = client.DnsAnalytics.Query(context.Background(), "1", &DnsAnalyticsOptions{Groupings: String("zone_name,date")})
}
