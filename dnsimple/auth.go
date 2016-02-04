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
	Data *WhoamiData `json:"data"`
}

// WhoamiData represents an authenticated context
// that contains information about the current logged User and/or Account.
type WhoamiData struct {
	User    *User    `json:"user,omitempty"`
	Account *Account `json:"account,omitempty"`
}

// Whoami gets the current authenticate context.
//
// See https://developer.dnsimple.com/v2/whoami
func (s *AuthService) Whoami() (*WhoamiResponse, error) {
	whoamiResponse := &WhoamiResponse{}

	resp, err := s.client.get("/whoami", whoamiResponse)
	if err != nil {
		return nil, err
	}

	whoamiResponse.HttpResponse = resp
	return whoamiResponse, nil
}

// Whoami is a state-less shortcut to client.Whoami()
// that returns only the relevant Data.
func Whoami(c *Client) (data *WhoamiData, err error) {
	resp, err := c.Auth.Whoami()
	if resp != nil {
		data = resp.Data
	}
	return
}
