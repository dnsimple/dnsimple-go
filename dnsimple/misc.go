package dnsimple

import ()

// MiscService handles communication with several miscellaneous
// methods of the DNSimple API.
//
// DNSimple API docs: https://developer.dnsimple.com/
type MiscService struct {
	client *Client
}

type whoamiWrapper struct {
	Whoami Whoami `json:"data"`
}

type Whoami struct {
	User    User    `json:"user,omitempty"`
	Account Account `json:"account,omitempty"`
}

// User gets the logged in user.
//
// DNSimple API docs: http://developer.dnsimple.com/v2/whoami/
func (s *MiscService) Whoami() (Whoami, *Response, error) {
	responseWrapper := whoamiWrapper{}

	res, err := s.client.get("whoami", &responseWrapper)
	if err != nil {
		return Whoami{}, res, err
	}

	return responseWrapper.Whoami, res, nil
}
