package dnsimple

import (
	"encoding/base64"
)

const (
	httpHeaderAuthorization = "Authorization"
)

// Provides credentials that can be used for authenticating with DNSimple.
//
// See https://developer.dnsimple.com/v2/#authentication
type Credentials interface {
	// Returns the HTTP headers that should be set
	// to authenticate the HTTP Request.
	Headers() map[string]string
}

// HTTP basic authentication
type httpBasicCredentials struct {
	email    string
	password string
}

// NewHTTPBasicCredentials construct Credentials using HTTP Basic Auth.
func NewHTTPBasicCredentials(email, password string) Credentials {
	return &httpBasicCredentials{email, password}
}

func (c *httpBasicCredentials) Headers() map[string]string {
	return map[string]string{httpHeaderAuthorization: "Basic " + c.basicAuth(c.email, c.password)}
}

func (c *httpBasicCredentials) basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// OAuth token authentication
type oauthTokenCredentials struct {
	oauthToken string
}

// NewOauthTokenCredentials construct Credentials using the OAuth access token.
func NewOauthTokenCredentials(oauthToken string) Credentials {
	return &oauthTokenCredentials{oauthToken: oauthToken}
}

func (c *oauthTokenCredentials) Headers() map[string]string {
	return map[string]string{httpHeaderAuthorization: "Bearer " + c.oauthToken}
}
