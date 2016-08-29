package dnsimple

import (
	"io"
	"net/http"
	"testing"
)

func TestDomains_domainPushPath(t *testing.T) {
	actual := domainPushPath("1", "example.com", 0)
	expected := "/1/domains/example.com/pushes"

	if actual != expected {
		t.Errorf("domainPath(\"1\", \"example.com\", 0): actual %s, expected %s", actual, expected)
	}
}

func TestDomainsService_InitiatePush(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/pushes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/initiatePush/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"new_account_email": "admin@target-account.test"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	pushAttributes := DomainPushAttributes{NewAccountEmail: "admin@target-account.test"}

	pushResponse, err := client.Domains.InitiatePush("1010", "example.com", pushAttributes)
	if err != nil {
		t.Fatalf("Domains.InitiatePush() returned error: %v", err)
	}

	push := pushResponse.Data
	if want, got := 1, push.ID; want != got {
		t.Fatalf("Domains.InitiatePush() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := 2020, push.AccountID; want != got {
		t.Fatalf("Domains.InitiatePush() returned Account ID expected to be `%v`, got `%v`", want, got)
	}
}
