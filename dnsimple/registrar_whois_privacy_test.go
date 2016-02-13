package dnsimple

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestRegistrarService_GetWhoisPrivacy(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/whois_privacy", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getWhoisPrivacy/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	privacyResponse, err := client.Registrar.GetWhoisPrivacy(accountID, "example.com")
	if err != nil {
		t.Errorf("Registrar.GetWhoisPrivacy() returned error: %v", err)
	}

	privacy := privacyResponse.Data
	wantSingle := &WhoisPrivacy{
		ID:        153,
		DomainID:  5916,
		Enabled:   true,
		ExpiresOn: "2017-02-13",
		CreatedAt: "2016-02-13T14:34:50.135Z",
		UpdatedAt: "2016-02-13T14:34:52.571Z"}

	if !reflect.DeepEqual(privacy, wantSingle) {
		t.Fatalf("Registrar.GetWhoisPrivacy() returned %+v, want %+v", privacy, wantSingle)
	}
}
