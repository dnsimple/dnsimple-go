package dnsimple

import (
	"fmt"
	"os"
	"testing"
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

	fmt.Printf("RateLimit: %v/%v until %v\n", whoamiResponse.RateLimitRemaining(), whoamiResponse.RateLimit(), whoamiResponse.RateLimitReset())
	whoami := whoamiResponse.Data
	fmt.Printf("Account: %+v\n", whoami.Account)
	fmt.Printf("User: %+v\n", whoami.User)
}

func TestLive_Domains(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoamiResponse, err := dnsimpleClient.Auth.Whoami()
	if err != nil {
		t.Fatalf("Live whoami()/listDomains() returned error: %v", err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", whoamiResponse.RateLimitRemaining(), whoamiResponse.RateLimit(), whoamiResponse.RateLimitReset())
	whoami := whoamiResponse.Data
	accountID := whoami.Account.ID

	domainsResponse, err := dnsimpleClient.Domains.List(fmt.Sprintf("%v", accountID))
	if err != nil {
		t.Fatalf("Live listDomains() returned error: %v", err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", domainsResponse.RateLimitRemaining(), domainsResponse.RateLimit(), domainsResponse.RateLimitReset())
	fmt.Printf("Domains: %+v\n", domainsResponse.Data)
}
