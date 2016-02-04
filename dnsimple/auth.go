package dnsimple

import (
)

// AuthService handles communication with several authentication
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/
type AuthService struct {
	client *Client
}

type WhoamiResponse struct {
	Response
	Data *Whoami `json:"data"`
}

//func (r *WhoamiResponse) Data() (*Whoami) {
//	return r.Data
//}

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
	whoamiResponse.HttpResponse = resp.HttpResponse
	if err != nil {
		return whoamiResponse, err
	}

	return whoamiResponse, nil
}
