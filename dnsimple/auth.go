package dnsimple

import ()

// AuthService handles communication with several authentication
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/
type AuthService struct {
	client *Client
}

type whoamiWrapper struct {
	Whoami Whoami `json:"data"`
}

type Whoami struct {
	User    User    `json:"user,omitempty"`
	Account Account `json:"account,omitempty"`
}

// Whoami gets the current authenticate context.
//
// See https://developer.dnsimple.com/v2/whoami
func (s *AuthService) Whoami() (Whoami, *Response, error) {
	responseWrapper := whoamiWrapper{}

	res, err := s.client.get("whoami", &responseWrapper)
	if err != nil {
		return Whoami{}, res, err
	}

	return responseWrapper.Whoami, res, nil
}
