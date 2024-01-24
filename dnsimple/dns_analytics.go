package dnsimple

import (
	"context"
	"fmt"
)

// DnsAnalyticsService handles communication with the DNS Analytics related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/dns_analytics/
type DnsAnalyticsService struct {
	client *Client
}

// DnsAnalytics represents DNS Analytics data.
type DnsAnalytics struct {
	Volume int64
	Zone   string
	Date   string
}

type DnsAnalyticsQueryParameters struct {
	AccountId int64  `json:"account_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Sort      string `json:"sort"`
	Page      int64  `json:"page"`
	PerPage   int64  `json:"per_page"`
	Groupings string `json:"groupings"`
}

// DnsAnalyticsResponse represents a response from an API method that returns DnsAnalytics data.
type DnsAnalyticsResponse struct {
	Response
	Data    []DnsAnalytics
	Rows    [][]interface{}             `json:"data"`
	Headers []string                    `json:"headers"`
	Query   DnsAnalyticsQueryParameters `json:"query"`
}

func (r *DnsAnalyticsResponse) marshallData() {
	list := make([]DnsAnalytics, len(r.Rows))

	for i, row := range r.Rows {
		var dataEntry DnsAnalytics
		for j, header := range r.Headers {
			switch header {
			case "volume":
				dataEntry.Volume = int64(row[j].(float64))
			case "zone":
				dataEntry.Zone = row[j].(string)
			case "date":
				dataEntry.Date = row[j].(string)
			}
		}

		list[i] = dataEntry
	}
	r.Data = list
}

// Query gets DNS Analytics data for an account
//
// See https://developer.dnsimple.com/v2/dns_analytics/#query
func (s *DnsAnalyticsService) Query(ctx context.Context, accountID string, options *ListOptions) (*DnsAnalyticsResponse, error) {
	path := versioned(fmt.Sprintf("/%v/dns_analytics", accountID))
	dnsAnalyticsResponse := &DnsAnalyticsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(ctx, path, dnsAnalyticsResponse)
	if err != nil {
		return dnsAnalyticsResponse, err
	}

	dnsAnalyticsResponse.HTTPResponse = resp
	dnsAnalyticsResponse.marshallData()
	return dnsAnalyticsResponse, nil
}
