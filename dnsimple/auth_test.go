package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAuthService_Whoami(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/whoami", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":{"user":null,"account":{"id":1,"email":"example-account@example.com"}}}
		`)
	})

	whoami, _, err := client.Auth.Whoami()

	if err != nil {
		t.Fatalf("Auth.Whoami() returned error: %v", err)
	}

	want := &Whoami{Account: &Account{Id: 1, Email: "example-account@example.com"}}
	if !reflect.DeepEqual(whoami, want) {
		t.Errorf("Auth.Whoami() returned %+v, want %+v", whoami, want)
	}
}
