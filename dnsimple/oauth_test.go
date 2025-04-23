package dnsimple

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOauthService_ExchangeAuthorizationForToken(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	code := "1234567890"
	clientID := "a1b2c3"
	clientSecret := "thisisasecret"

	mux.HandleFunc("/v2/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/oauthAccessToken/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"code": code, "client_id": clientID, "client_secret": clientSecret, "grant_type": "authorization_code"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	token, err := client.Oauth.ExchangeAuthorizationForToken(&ExchangeAuthorizationRequest{Code: code, ClientID: clientID, ClientSecret: clientSecret, GrantType: AuthorizationCodeGrant})

	assert.NoError(t, err)
	want := &AccessToken{Token: "zKQ7OLqF5N1gylcJweA9WodA000BUNJD", Type: "Bearer", AccountID: int64(1)}
	assert.Equal(t, want, token)
}

func TestOauthService_ExchangeAuthorizationForToken_Error(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/oauth/access_token", func(w http.ResponseWriter, _ *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/oauthAccessToken/error-invalid-request.http")

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Oauth.ExchangeAuthorizationForToken(&ExchangeAuthorizationRequest{Code: "1234567890", ClientID: "a1b2c3", ClientSecret: "thisisasecret", GrantType: "authorization_code"})
	if err == nil {
		t.Fatalf("Oauth.ExchangeAuthorizationForToken() expected to return an error")
	}

	var v *ExchangeAuthorizationError
	assert.ErrorAs(t, err, &v)
	assert.Equal(t, `Invalid "state": value doesn't match the "state" in the authorization request`, v.ErrorDescription)
}

func TestOauthService_AuthorizeURL(t *testing.T) {
	clientID := "a1b2c3"
	client.BaseURL = "https://api.host.test"

	assert.Equal(t, "https://host.test/oauth/authorize?client_id=a1b2c3&response_type=code", client.Oauth.AuthorizeURL(clientID, nil))
	assert.Equal(t, "https://host.test/oauth/authorize?client_id=a1b2c3&response_type=code&state=randomstate", client.Oauth.AuthorizeURL(clientID, &AuthorizationOptions{State: "randomstate"}))
}
