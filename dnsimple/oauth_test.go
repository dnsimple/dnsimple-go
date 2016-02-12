package dnsimple

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestOauthService_ExchangeAuthorizationForToken(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	code := "1234567890"
	clientID := "a1b2c3"
	clientSecret := "thisisasecret"

	mux.HandleFunc("/v2/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/oauthAccessToken/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"code": code, "client_id": clientID, "client_secret": clientSecret}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	token, err := client.Oauth.ExchangeAuthorizationForToken(&ExchangeAuthorizationRequest{Code: code, ClientID: clientID, ClientSecret: clientSecret})
	if err != nil {
		t.Fatalf("Oauth.ExchangeAuthorizationForToken() returned error: %v", err)
	}

	want := &AccessToken{Token: "zKQ7OLqF5N1gylcJweA9WodA000BUNJD", Type: "Bearer", AccountID: 1}
	if !reflect.DeepEqual(token, want) {
		t.Errorf("Oauth.ExchangeAuthorizationForToken() returned %+v, want %+v", token, want)
	}
}

func TestOauthService_AuthorizeURL(t *testing.T) {
	clientID := "a1b2c3"
	client.BaseURL = "https://api.host.test"

	if want, got := "https://host.test?client_id=a1b2c3", client.Oauth.AuthorizeURL(clientID, nil); want != got {
		t.Errorf("AuthorizeURL = %v, want %v", got, want)
	}

	if want, got := "https://host.test?client_id=a1b2c3&state=randomstate", client.Oauth.AuthorizeURL(clientID, &AuthorizationOptions{State: "randomstate"}); want != got {
		t.Errorf("AuthorizeURL = %v, want %v", got, want)
	}
}
