package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			TotalAmount:   "14.50",
			BalanceAmount: "0.00",
			Reference:     "1-2",
			State:         "collected",
			Items: []ChargeItem{
				{
					Description:      "Register bubble-registered.com",
					Amount:           "14.50",
					ProductId:        1,
					ProductType:      "domain-registration",
					ProductReference: "bubble-registered.com",
				},
			},
		},
		{
			InvoicedAt:    "2023-08-17T05:57:53Z",
			TotalAmount:   "14.50",
			BalanceAmount: "0.00",
			Reference:     "2-2",
			State:         "refunded",
			Items: []ChargeItem{
				{
					Description:      "Register example.com",
					Amount:           "14.50",
					ProductId:        2,
					ProductType:      "domain-registration",
					ProductReference: "example.com",
				},
			},
		},
		{
			InvoicedAt:    "2023-10-24T07:49:05Z",
			TotalAmount:   "1099999.99",
			BalanceAmount: "0.00",
			Reference:     "4-2",
			State:         "collected",
			Items: []ChargeItem{
				{
					Description:      "Test Line Item 1",
					Amount:           "99999.99",
					ProductId:        0,
					ProductType:      "manual",
					ProductReference: "",
				},
				{
					Description:      "Test Line Item 2",
					Amount:           "1000000.00",
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

	assert.Equal(t, err.(*ErrorResponse).Message, "Invalid date format must be ISO8601 (YYYY-MM-DD)")
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

	assert.Equal(t, err.(*ErrorResponse).Message, "Permission Denied. Required Scope: billing:*:read")
}

func TestChargeTotalAmountFloat(t *testing.T) {
	tests := []struct {
		name    string
		charge  Charge
		want    float64
		wantErr bool
	}{
		{
			name:    "empty total amount",
			charge:  Charge{TotalAmount: ""},
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "valid total amount",
			charge:  Charge{TotalAmount: "123.45"},
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "invalid total amount",
			charge:  Charge{TotalAmount: "abc"},
			want:    0.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.charge.TotalAmountFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("TotalAmountFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TotalAmountFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChargeBalanceAmountFloat(t *testing.T) {
	tests := []struct {
		name    string
		charge  Charge
		want    float64
		wantErr bool
	}{
		{
			name:    "empty total amount",
			charge:  Charge{BalanceAmount: ""},
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "valid total amount",
			charge:  Charge{BalanceAmount: "123.45"},
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "invalid total amount",
			charge:  Charge{BalanceAmount: "abc"},
			want:    0.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.charge.BalanceAmountFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceAmountFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BalanceAmountFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChargeItemAmountFloat(t *testing.T) {
	tests := []struct {
		name       string
		chargeItem ChargeItem
		want       float64
		wantErr    bool
	}{
		{
			name:       "empty total amount",
			chargeItem: ChargeItem{Amount: ""},
			want:       0.0,
			wantErr:    true,
		},
		{
			name:       "valid total amount",
			chargeItem: ChargeItem{Amount: "123.45"},
			want:       123.45,
			wantErr:    false,
		},
		{
			name:       "invalid total amount",
			chargeItem: ChargeItem{Amount: "abc"},
			want:       0.0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chargeItem.AmountFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("AmountFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AmountFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
