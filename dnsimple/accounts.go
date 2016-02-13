package dnsimple

// Account represents a DNSimple account.
type Account struct {
	ID    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}
