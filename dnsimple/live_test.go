package dnsimple

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	dnsimpleLiveTest bool
	dnsimpleToken    string
	dnsimpleBaseURL  string
	dnsimpleClient   *Client
)

func init() {
	dnsimpleToken = os.Getenv("DNSIMPLE_TOKEN")
	dnsimpleBaseURL = os.Getenv("DNSIMPLE_BASE_URL")

	// Prevent people from wiping out their entire production account by mistake
	if dnsimpleBaseURL == "" {
		dnsimpleBaseURL = "https://api.sandbox.dnsimple.com"
	}

	if len(dnsimpleToken) > 0 {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: dnsimpleToken})
		tc := oauth2.NewClient(context.Background(), ts)

		dnsimpleLiveTest = true
		dnsimpleClient = NewClient(tc)
		dnsimpleClient.BaseURL = dnsimpleBaseURL
		dnsimpleClient.UserAgent = fmt.Sprintf("%v +livetest", dnsimpleClient.UserAgent)
	}
}

func TestLive_Whoami(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoamiResponse, err := dnsimpleClient.Identity.Whoami(context.Background())

	assert.NoError(t, err)

	fmt.Printf("RateLimit: %v/%v until %v\n", whoamiResponse.RateLimitRemaining(), whoamiResponse.RateLimit(), whoamiResponse.RateLimitReset())
	whoami := whoamiResponse.Data
	fmt.Printf("Account: %+v\n", whoami.Account)
	fmt.Printf("User: %+v\n", whoami.User)
}

func TestLive_Domains(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	if err != nil {
		t.Fatalf("Live Whoami() returned error: %v", err)
	}

	accountID := whoami.Account.ID

	domainsResponse, err := dnsimpleClient.Domains.ListDomains(context.Background(), fmt.Sprintf("%v", accountID), nil)
	assert.NoError(t, err)

	fmt.Printf("RateLimit: %v/%v until %v\n", domainsResponse.RateLimitRemaining(), domainsResponse.RateLimit(), domainsResponse.RateLimitReset())
	fmt.Printf("Domains: %+v\n", domainsResponse.Data)
}

func TestLive_Registration_ValidationError(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := whoami.Account.ID

	registerRequest := &RegisterDomainInput{RegistrantID: -1}
	_, err = dnsimpleClient.Registrar.RegisterDomain(context.Background(), fmt.Sprintf("%v", accountID), fmt.Sprintf("example-%v.com", time.Now().Unix()), registerRequest)

	assert.Equal(t, "Validation failed", err.(*ErrorResponse).Message)
}

func TestLive_Registration_ExtendedAttributesValidationError(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := whoami.Account.ID

	contactsResponse, err := dnsimpleClient.Contacts.ListContacts(context.Background(), fmt.Sprintf("%v", accountID), nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(contactsResponse.Data), 1, "At least one contact is required for this live test")
	registrantID := contactsResponse.Data[0].ID
	registerRequest := &RegisterDomainInput{RegistrantID: int(registrantID)}
	_, err = dnsimpleClient.Registrar.RegisterDomain(context.Background(), fmt.Sprintf("%v", accountID), fmt.Sprintf("example-%v.app", time.Now().Unix()), registerRequest)

	var got *ErrorResponse
	assert.ErrorAs(t, err, &got)
	assert.Contains(t, got.AttributeErrors["x-accept-ssl-requirement"], "it's required")
	assert.Equal(t, "Invalid extended attributes", got.Message)
}

func TestLive_Registration(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := whoami.Account.ID

	contactsResponse, err := dnsimpleClient.Contacts.ListContacts(context.Background(), fmt.Sprintf("%v", accountID), nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(contactsResponse.Data), 1, "At least one contact is required for this live test")
	registrantID := contactsResponse.Data[0].ID
	registerRequest := &RegisterDomainInput{RegistrantID: int(registrantID)}
	registrationResponse, err := dnsimpleClient.Registrar.RegisterDomain(context.Background(), fmt.Sprintf("%v", accountID), fmt.Sprintf("example-%v.com", time.Now().Unix()), registerRequest)
	assert.NoError(t, err)

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

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := whoami.Account.ID

	webhooksResponse, err = dnsimpleClient.Webhooks.ListWebhooks(context.Background(), fmt.Sprintf("%v", accountID), nil)
	assert.NoError(t, err)

	fmt.Printf("RateLimit: %v/%v until %v\n", webhooksResponse.RateLimitRemaining(), webhooksResponse.RateLimit(), webhooksResponse.RateLimitReset())
	fmt.Printf("Webhooks: %+v\n", webhooksResponse.Data)

	webhookAttributes := Webhook{URL: "https://livetest.test"}
	webhookResponse, err = dnsimpleClient.Webhooks.CreateWebhook(context.Background(), fmt.Sprintf("%v", accountID), webhookAttributes)
	assert.NoError(t, err)

	fmt.Printf("RateLimit: %v/%v until %v\n", webhooksResponse.RateLimitRemaining(), webhooksResponse.RateLimit(), webhooksResponse.RateLimitReset())
	fmt.Printf("Webhook: %+v\n", webhookResponse.Data)

	webhook = webhookResponse.Data
	webhookResponse, err = dnsimpleClient.Webhooks.DeleteWebhook(context.Background(), fmt.Sprintf("%v", accountID), webhook.ID)
	assert.NoError(t, err)

	fmt.Printf("RateLimit: %v/%v until %v\n", webhooksResponse.RateLimitRemaining(), webhooksResponse.RateLimit(), webhooksResponse.RateLimitReset())
	fmt.Printf("Webhook: %+v\n", webhookResponse.Data)
}

func TestLive_Zones(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := fmt.Sprintf("%v", whoami.Account.ID)

	domainResponse, err := dnsimpleClient.Domains.CreateDomain(context.Background(), accountID, Domain{Name: fmt.Sprintf("example-%v.test", time.Now().Unix())})
	assert.NoError(t, err)

	zoneName := domainResponse.Data.Name
	recordName := fmt.Sprintf("%v", time.Now().Unix())
	recordResponse, err := dnsimpleClient.Zones.CreateRecord(context.Background(), accountID, zoneName, ZoneRecordAttributes{Name: &recordName, Type: "TXT", Content: "Test"})
	assert.NoError(t, err)

	fmt.Printf("ZoneRecord: %+v\n", recordResponse.Data)
}

func TestLive_Error(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	_, err = dnsimpleClient.Registrar.RegisterDomain(context.Background(), fmt.Sprintf("%v", whoami.Account.ID), fmt.Sprintf("example-%v.test", time.Now().Unix()), &RegisterDomainInput{})

	e := err.(*ErrorResponse)
	fmt.Println(e.Message)
}
