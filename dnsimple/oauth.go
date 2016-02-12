package dnsimple

// OauthService handles communication with the authorization related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/oauth/
type OauthService struct {
	client *Client
}

// AccessToken represents a DNSimple Oauth access token.
type AccessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	AccountID int    `json:"account_id"`
}

// ExchangeAuthorizationRequest represents a request to exchange
// an authorization code for an access token.
// RedirectURI is optional, all the other fields are mandatory.
type ExchangeAuthorizationRequest struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
}

// ExchangeAuthorizationForToken exchanges the short-lived authorization code for an access token
// you can use to authenticate your API calls.
func (s *OauthService) ExchangeAuthorizationForToken(authorization *ExchangeAuthorizationRequest) (*AccessToken, error) {
	accessToken := &AccessToken{}

	_, err := s.client.post("/oauth/access_token", authorization, accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
