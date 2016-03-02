package dnsimple

import (
	"net/url"
	"strings"
)

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
	// Currently the only supported value is "authorization_code"
	GrantType string `json:"grant_type,omitempty"`
}

// ExchangeAuthorizationForToken exchanges the short-lived authorization code for an access token
// you can use to authenticate your API calls.
func (s *OauthService) ExchangeAuthorizationForToken(authorization *ExchangeAuthorizationRequest) (*AccessToken, error) {
	path := versioned("/oauth/access_token")
	accessToken := &AccessToken{}

	_, err := s.client.post(path, authorization, accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

// AuthorizationOptions represents the option you can use to generate an authorization URL.
type AuthorizationOptions struct {
	RedirectURI string `url:"redirect_uri,omitempty"`
	// Currently "state" is required by the DNSimple OAuth implementation,
	// so you must specify it.
	State string `url:"state,omitempty"`
}

// AuthorizeURL generates the URL to authorize an user for an application via the OAuth2 flow.
func (s *OauthService) AuthorizeURL(clientID string, options *AuthorizationOptions) string {
	uri, _ := url.Parse(strings.Replace(s.client.BaseURL, "api.", "", 1))
	uri.Path = "/oauth/authorize"
	query := uri.Query()
	query.Add("client_id", clientID)
	query.Add("response_type", "code")
	uri.RawQuery = query.Encode()

	path, _ := addURLQueryOptions(uri.String(), options)
	return path
}
