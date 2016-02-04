package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMiscService_Whoami(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/whoami", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{
			  "data": {
			    "user": null,
			    "account": {
			      "id": 24,
			      "email": "example-account@example.com"
			    }
			  }
			}
		`)
	})

	whoami, _, err := client.Misc.Whoami()

	if err != nil {
		t.Errorf("Misc.Whoami returned error: %v", err)
	}

	want := Whoami{Account: Account{Id: 24, Email: "example-account@example.com"}}
	if !reflect.DeepEqual(whoami, want) {
		t.Errorf("Misc.Whoami returned %+v, want %+v", whoami, want)
	}
}
