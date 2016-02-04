package dnsimple

import ()

// AuthService handles communication with several authentication
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/
type AuthService struct {
	client *Client
}

// WhoamiResponse represents a response from an API method that returns a Whoami struct.
type WhoamiResponse struct {
	Response
	Data *Whoami `json:"data"`
}

// Whoami represents an authenticated context
// that contains information about the current logged User and/or Account.
type Whoami struct {
	User    *User    `json:"user,omitempty"`
	Account *Account `json:"account,omitempty"`
}

// Whoami gets the current authenticate context.
//
// See https://developer.dnsimple.com/v2/whoami
func (s *AuthService) Whoami() (*WhoamiResponse, error) {
	whoamiResponse := &WhoamiResponse{}

	resp, err := s.client.get("/whoami", whoamiResponse)
	whoamiResponse.HttpResponse = resp
	if err != nil {
		return whoamiResponse, err
	}

	return whoamiResponse, nil
}
