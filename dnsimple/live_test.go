package dnsimple

import (
	"os"
	"testing"
	"fmt"
)

var (
	dnsimpleLiveTest bool
	dnsimpleToken    string
	dnsimpleBaseURL  string
	dnsimpleClient   *Client
)

func init() {
	dnsimpleToken = os.Getenv("DNSIMPLE_TOKEN")
	dnsimpleBaseURL = os.Getenv("DOMAINR_BASE_URL")

	// Prevent peoeple from wiping out their entire production account by mistake
	if dnsimpleBaseURL == "" {
		dnsimpleBaseURL = "https://api.sandbox.dnsimple.com/"
	}

	if len(dnsimpleToken) > 0 {
		dnsimpleLiveTest = true
		dnsimpleClient = NewClient(NewOauthTokenCredentials(dnsimpleToken))
		dnsimpleClient.BaseURL = dnsimpleBaseURL
		dnsimpleClient.UserAgent = fmt.Sprintf("%v +livetest", dnsimpleClient.UserAgent)
	}
}

func TestLive_Whoami(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoamiResponse, err := dnsimpleClient.Auth.Whoami()
	if err != nil {
		t.Fatalf("Live whoami() returned error: %v", err)
	}
	whoami := whoamiResponse.Data

	fmt.Println(whoami.Account)
	fmt.Println(whoami.User)
}

func TestLive_Domains(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoamiResponse, err := dnsimpleClient.Auth.Whoami()
	if err != nil {
		t.Fatalf("Live whoami()/listDomains() returned error: %v", err)
	}

	whoami := whoamiResponse.Data
	accountID := whoami.Account.ID

	domainsResponse, err := dnsimpleClient.Domains.List(fmt.Sprintf("%v", accountID))
	if err != nil {
		t.Fatalf("Live listDomains() returned error: %v", err)
	}

	fmt.Println(domainsResponse.Data)
}
