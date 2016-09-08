package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDomains_initiateDomainPushPath(t *testing.T) {
	actual := initiateDomainPushPath("1", "example.com")
	expected := "/1/domains/example.com/pushes"

	if actual != expected {
		t.Errorf("initiateDomainPushPathdomainPath(\"1\", \"example.com\", 0): actual %s, expected %s", actual, expected)
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

func TestDomainsService_DomainsPushesList(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listPushes/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	pushesResponse, err := client.Domains.ListPushes("2020", nil)
	if err != nil {
		t.Fatalf("Domains.ListPushes() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}), pushesResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Domains.ListPushes() pagination expected to be %v, got %v", want, got)
	}

	pushes := pushesResponse.Data
	if want, got := 2, len(pushes); want != got {
		t.Errorf("Domains.ListPushes() expected to return %v pushes, got %v", want, got)
	}

	if want, got := 1, pushes[0].ID; want != got {
		t.Fatalf("Domains.ListPushes() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := 2020, pushes[0].AccountID; want != got {
		t.Fatalf("Domains.ListPushes() returned Account ID expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_DomainsPushesList_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listPushes/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListPushes("2020", &ListOptions{Page: 2, PerPage: 20})
	if err != nil {
		t.Fatalf("Domains.ListPushes() returned error: %v", err)
	}
}

func TestDomainsService_AcceptPush(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/acceptPush/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"contact_id": "2"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	pushAttributes := DomainPushAttributes{ContactID: "2"}

	_, err := client.Domains.AcceptPush("2020", 1, pushAttributes)
	if err != nil {
		t.Fatalf("Domains.AcceptPush() returned error: %v", err)
	}
}

func TestDomainsService_RejectPush(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/2020/pushes/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/rejectPush/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.RejectPush("2020", 1)
	if err != nil {
		t.Fatalf("Domains.RejectPush() returned error: %v", err)
	}
}
