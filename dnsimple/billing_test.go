package dnsimple

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func toDecimal(t *testing.T, s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		assert.Nilf(t, err, "toDecimal() error = %v", err)
	}

	return d
}

func TestBillingService_ListCharges_Success(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/billing/charges", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listCharges/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	response, err := client.Billing.ListCharges(context.Background(), "1010", ListChargesOptions{})

	assert.NoError(t, err)
	assert.Equal(t, response.Pagination, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 3})
	assert.Equal(t, response.Data, []Charge{
		{
			InvoicedAt:    "2023-08-17T05:53:36Z",
			TotalAmount:   toDecimal(t, "14.50"),
			BalanceAmount: toDecimal(t, "0.00"),
			Reference:     "1-2",
			State:         "collected",
			Items: []ChargeItem{
				{
					Description:      "Register bubble-registered.com",
					Amount:           toDecimal(t, "14.50"),
					ProductId:        1,
					ProductType:      "domain-registration",
					ProductReference: "bubble-registered.com",
				},
			},
		},
		{
			InvoicedAt:    "2023-08-17T05:57:53Z",
			TotalAmount:   toDecimal(t, "14.50"),
			BalanceAmount: toDecimal(t, "0.00"),
			Reference:     "2-2",
			State:         "refunded",
			Items: []ChargeItem{
				{
					Description:      "Register example.com",
					Amount:           toDecimal(t, "14.50"),
					ProductId:        2,
					ProductType:      "domain-registration",
					ProductReference: "example.com",
				},
			},
		},
		{
			InvoicedAt:    "2023-10-24T07:49:05Z",
			TotalAmount:   toDecimal(t, "1099999.99"),
			BalanceAmount: toDecimal(t, "0.00"),
			Reference:     "4-2",
			State:         "collected",
			Items: []ChargeItem{
				{
					Description:      "Test Line Item 1",
					Amount:           toDecimal(t, "99999.99"),
					ProductId:        0,
					ProductType:      "manual",
					ProductReference: "",
				},
				{
					Description:      "Test Line Item 2",
					Amount:           toDecimal(t, "1000000.00"),
					ProductId:        0,
					ProductType:      "manual",
					ProductReference: "",
				},
			},
		},
	})
}

func TestBillingService_ListCharges_Fail400BadFilter(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/billing/charges", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listCharges/fail-400-bad-filter.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Billing.ListCharges(context.Background(), "1010", ListChargesOptions{})

	assert.Equal(t, "Invalid date format must be ISO8601 (YYYY-MM-DD)", err.(*ErrorResponse).Message)
}

func TestBillingService_ListCharges_Fail403(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/billing/charges", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listCharges/fail-403.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Billing.ListCharges(context.Background(), "1010", ListChargesOptions{})

	assert.Equal(t, "Permission Denied. Required Scope: billing:*:read", err.(*ErrorResponse).Message)
}

func TestUnmarshalCharge(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    Charge
		wantErr bool
	}{
		{
			name:    "valid json",
			jsonStr: `{"total_amount": "123.45", "balance_amount": "67.89"}`,
			want: Charge{
				TotalAmount:   decimal.NewFromFloat(123.45),
				BalanceAmount: decimal.NewFromFloat(67.89),
			},
			wantErr: false,
		},
		{
			name:    "zero values",
			jsonStr: `{"total_amount": "0.00", "balance_amount": "0.00"}`,
			want: Charge{
				TotalAmount:   decimal.NewFromFloat(0.00),
				BalanceAmount: decimal.NewFromFloat(0.00),
			},
			wantErr: false,
		},
		{
			name:    "invalid amount value",
			jsonStr: `{"total_amount": "123.45", "balance_amount": "abc"}`,
			want: Charge{
				TotalAmount:   decimal.NewFromFloat(123.45),
				BalanceAmount: decimal.Decimal{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Charge
			err := json.Unmarshal([]byte(tt.jsonStr), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Truef(t, got.TotalAmount.Equal(tt.want.TotalAmount), "TotalAmount: got %v, want %v\nTesting: %s\n%s", got.TotalAmount, tt.want.TotalAmount, tt.name, tt.jsonStr)
			assert.Truef(t, got.BalanceAmount.Equal(tt.want.BalanceAmount), "BalanceAmount: got %v, want %v\nTesting: %s\n%s", got.BalanceAmount, tt.want.BalanceAmount, tt.name, tt.jsonStr)
		})
	}
}
