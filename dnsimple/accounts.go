package dnsimple

import ()

// Account represents a DNSimple account.
type Account struct {
	Id    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}
