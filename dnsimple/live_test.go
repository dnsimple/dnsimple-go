package dnsimple

import (
	"fmt"
	"os"
	"testing"
	"time"
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
		dnsimpleBaseURL = "https://api.sandbox.dnsimple.com"
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

	whoamiResponse, err := dnsimpleClient.Identity.Whoami()
	if err != nil {
		t.Fatalf("Live Auth.Whoami() returned error: %v", err)
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

	whoami, err := Whoami(dnsimpleClient)
	if err != nil {
		t.Fatalf("Live Whoami() returned error: %v", err)
	}

	accountID := whoami.Account.ID

	domainsResponse, err := dnsimpleClient.Domains.ListDomains(fmt.Sprintf("%v", accountID), nil)
	if err != nil {
		t.Fatalf("Live Domains.List() returned error: %v", err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", domainsResponse.RateLimitRemaining(), domainsResponse.RateLimit(), domainsResponse.RateLimitReset())
	fmt.Printf("Domains: %+v\n", domainsResponse.Data)
}

func TestLive_Registration(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(dnsimpleClient)
	if err != nil {
		t.Fatalf("Live Whoami() returned error: %v", err)
	}

	accountID := whoami.Account.ID

	// TODO: fetch the registrant randomly
	registerRequest := &DomainRegisterRequest{RegistrantID: 2}
	registrationResponse, err := dnsimpleClient.Registrar.RegisterDomain(fmt.Sprintf("%v", accountID), fmt.Sprintf("example-%v.com", time.Now().Unix()), registerRequest)
	if err != nil {
		t.Fatalf("Live Registrar.Register() returned error: %v", err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", registrationResponse.RateLimitRemaining(), registrationResponse.RateLimit(), registrationResponse.RateLimitReset())
	fmt.Printf("Domain: %+v\n", registrationResponse.Data)
}

func TestLive_Webhooks(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	var err error
	var webhook *Webhook
	var webhookResponse *WebhookResponse
	var webhooksResponse *WebhooksResponse

	whoami, err := Whoami(dnsimpleClient)
	if err != nil {
		t.Fatalf("Live Auth.Whoami()/Domains.List() returned error: %v", err)
	}
	accountID := whoami.Account.ID

	webhooksResponse, err = dnsimpleClient.Webhooks.List(fmt.Sprintf("%v", accountID), nil)
	if err != nil {
		t.Fatalf("Live Webhooks.List() returned error: %v", err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", webhooksResponse.RateLimitRemaining(), webhooksResponse.RateLimit(), webhooksResponse.RateLimitReset())
	fmt.Printf("Webhooks: %+v\n", webhooksResponse.Data)

	webhookAttributes := Webhook{URL: "https://livetest.test"}
	webhookResponse, err = dnsimpleClient.Webhooks.Create(fmt.Sprintf("%v", accountID), webhookAttributes)
	if err != nil {
		t.Fatalf("Live Webhooks.Create() returned error: %v", err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", webhooksResponse.RateLimitRemaining(), webhooksResponse.RateLimit(), webhooksResponse.RateLimitReset())
	fmt.Printf("Webhook: %+v\n", webhookResponse.Data)
	webhook = webhookResponse.Data

	webhookResponse, err = dnsimpleClient.Webhooks.Delete(fmt.Sprintf("%v", accountID), webhook.ID)
	if err != nil {
		t.Fatalf("Live Webhooks.Delete(%v) returned error: %v", webhook.ID, err)
	}

	fmt.Printf("RateLimit: %v/%v until %v\n", webhooksResponse.RateLimitRemaining(), webhooksResponse.RateLimit(), webhooksResponse.RateLimitReset())
	webhook = webhookResponse.Data
}
