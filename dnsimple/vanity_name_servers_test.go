package dnsimple

import (
	"io"
	"net/http"
	// "net/url"
	// "reflect"
	"testing"
)

func TestVanityNameServers_vanityNameServerPath(t *testing.T) {
	if want, got := "/1010/vanity/example.com", vanityNameServerPath("1010", "example.com"); want != got {
		t.Errorf("vanity_name_serverPath(%v,  ) = %v, want %v", "1010", got, want)
	}
}

func TestVanityNameServersService_Enable(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/vanity/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/enableVanityNameServers/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	// newDelegation := &Delegation{"ns1.example.com", "ns2.example.com"}

	vanityNameServerResponse, err := client.VanityNameServers.Enable("1010", "example.com")
	if err != nil {
		t.Fatalf("VanityNameServers.Enable() returned error: %v", err)
	}

	delegation := vanityNameServerResponse.Data[0].Name
	wantSingle := "ns1.example.com"

	if delegation != wantSingle {
		t.Fatalf("VanityNameServers.Enable() returned %+v, want %+v", delegation, wantSingle)
	}
}
