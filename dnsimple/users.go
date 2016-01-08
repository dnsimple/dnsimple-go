package dnsimple

import ()

// User represents a DNSimple user.
type User struct {
	Id    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}
