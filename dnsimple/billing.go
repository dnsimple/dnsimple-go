package dnsimple

import (
	"context"
	"fmt"
	"strconv"
)

type ListChargesOptions struct {
	// Only include results after the given date.
	StartDate string `url:"start_date,omitempty"`

	// Only include results before the given date.
	EndDate string `url:"end_date,omitempty"`

	// Sort results. Default sorting is by invoiced ascending.
	Sort string `url:"sort,omitempty"`
}

type ListChargesResponse struct {
	Response
	Data []Charge `json:"data"`
}

type Charge struct {
	InvoicedAt    string       `json:"invoiced_at,omitempty"`
	TotalAmount   string       `json:"total_amount,omitempty"`
	BalanceAmount string       `json:"balance_amount,omitempty"`
	Reference     string       `json:"reference,omitempty"`
	State         string       `json:"state,omitempty"`
	Items         []ChargeItem `json:"items,omitempty"`
}

type ChargeItem struct {
	Description      string `json:"description,omitempty"`
	Amount           string `json:"amount,omitempty"`
	ProductId        int64  `json:"product_id,omitempty"`
	ProductType      string `json:"product_type,omitempty"`
	ProductReference string `json:"product_reference,omitempty"`
}

type BillingService struct {
	client *Client
}

func (c *Charge) TotalAmountFloat() (float64, error) {
	if c.TotalAmount == "" {
		return 0.0, fmt.Errorf("total amount is empty")
	}

	totalAmount, err := strconv.ParseFloat(c.TotalAmount, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing total amount: %w", err)
	}

	return totalAmount, nil
}

func (c *Charge) BalanceAmountFloat() (float64, error) {
	if c.BalanceAmount == "" {
		return 0.0, fmt.Errorf("balance amount is empty")
	}

	balanceAmount, err := strconv.ParseFloat(c.BalanceAmount, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing balance amount: %w", err)
	}

	return balanceAmount, nil
}

func (c *ChargeItem) AmountFloat() (float64, error) {
	if c.Amount == "" {
		return 0.0, fmt.Errorf("total amount is empty")
	}

	amount, err := strconv.ParseFloat(c.Amount, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing total amount: %w", err)
	}

	return amount, nil
}

// Lists the billing charges for the account.
//
// See https://developer.dnsimple.com/v2/billing/#listCharges
func (s *BillingService) ListCharges(
	ctx context.Context,
	account string,
	options ListChargesOptions,
) (*ListChargesResponse, error) {
	res := &ListChargesResponse{}
	path := fmt.Sprintf("/v2/%v/billing/charges", account)

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	httpRes, err := s.client.get(
		ctx,
		path,
		res,
	)
	if err != nil {
		return nil, err
	}

	res.HTTPResponse = httpRes
	return res, nil
}
