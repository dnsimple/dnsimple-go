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

	var got *ErrorResponse
	assert.ErrorAs(t, err, &got)
	assert.Equal(t, "Validation failed", got.Message)
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
	assert.Equal(t, "Invalid extended attributes", got.Message)
	assert.Contains(t, got.AttributeErrors["x-accept-ssl-requirement"], "it's required")
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

func TestLive_BatchChangeZoneRecords(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := fmt.Sprintf("%v", whoami.Account.ID)

	// Create a test domain
	domainResponse, err := dnsimpleClient.Domains.CreateDomain(context.Background(), accountID, Domain{Name: fmt.Sprintf("batch-test-%v.test", time.Now().Unix())})
	assert.NoError(t, err)

	zoneName := domainResponse.Data.Name
	timestamp := fmt.Sprintf("%v", time.Now().Unix())

	// Create some initial records using batch API for testing updates and deletes
	setupRequest := BatchChangeZoneRecordsRequest{
		Creates: []ZoneRecordAttributes{
			{Name: String(fmt.Sprintf("update1-%v", timestamp)), Type: "A", Content: "1.2.3.4"},
			{Name: String(fmt.Sprintf("update2-%v", timestamp)), Type: "A", Content: "1.2.3.5"},
			{Name: String(fmt.Sprintf("delete1-%v", timestamp)), Type: "TXT", Content: "to-be-deleted"},
		},
	}

	setupResponse, err := dnsimpleClient.Zones.BatchChangeZoneRecords(context.Background(), accountID, zoneName, setupRequest)
	assert.NoError(t, err)
	assert.Len(t, setupResponse.Data.Creates, 3, "Expected 3 records to be created in setup")

	// Extract record IDs for later operations
	record1 := setupResponse.Data.Creates[0]
	record2 := setupResponse.Data.Creates[1]
	record3 := setupResponse.Data.Creates[2]

	// Perform batch operations
	batchRequest := BatchChangeZoneRecordsRequest{
		Creates: []ZoneRecordAttributes{
			{Type: "A", Content: "3.2.3.4", Name: String(fmt.Sprintf("create1-%v", timestamp))},
			{Type: "A", Content: "4.2.3.4", Name: String(fmt.Sprintf("create2-%v", timestamp))},
		},
		Updates: []ZoneRecordUpdateRequest{
			{ID: record1.ID, Content: "3.2.3.40", Name: String(fmt.Sprintf("updated1-%v", timestamp))},
			{ID: record2.ID, Content: "5.2.3.40", Name: String(fmt.Sprintf("updated2-%v", timestamp))},
		},
		Deletes: []ZoneRecordDeleteRequest{
			{ID: record3.ID},
		},
	}

	batchResponse, err := dnsimpleClient.Zones.BatchChangeZoneRecords(context.Background(), accountID, zoneName, batchRequest)
	assert.NoError(t, err)

	data := batchResponse.Data

	// Verify creates
	assert.Len(t, data.Creates, 2)
	fmt.Printf("Created records: %+v\n", data.Creates)

	// Verify updates
	assert.Len(t, data.Updates, 2)
	fmt.Printf("Updated records: %+v\n", data.Updates)

	// Verify deletes
	assert.Len(t, data.Deletes, 1)
	fmt.Printf("Deleted record IDs: %+v\n", data.Deletes)

	fmt.Printf("Batch operation completed successfully for zone: %v\n", zoneName)
}

func TestLive_BatchChangeZoneRecords_BatchError(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	accountID := fmt.Sprintf("%v", whoami.Account.ID)

	// Create a test domain
	domainResponse, err := dnsimpleClient.Domains.CreateDomain(context.Background(), accountID, Domain{Name: fmt.Sprintf("batch-error-test-%v.test", time.Now().Unix())})
	assert.NoError(t, err)

	zoneName := domainResponse.Data.Name

	// Test batch-specific error by trying to update non-existent records
	batchRequest := BatchChangeZoneRecordsRequest{
		Deletes: []ZoneRecordDeleteRequest{
			{ID: 88888888}, // Non-existent record ID
		},
	}

	_, err = dnsimpleClient.Zones.BatchChangeZoneRecords(context.Background(), accountID, zoneName, batchRequest)

	assert.Error(t, err)
	var errorResp *ErrorResponse
	if assert.ErrorAs(t, err, &errorResp) {
		assert.Equal(t, "Validation failed", errorResp.Message)
		assert.Contains(t, errorResp.AttributeErrors, "deletes[0]", "Expected delete error for non-existent record")
		assert.Equal(t, []string{"Record not found ID=88888888"}, errorResp.AttributeErrors["deletes[0]"])
		assert.Len(t, errorResp.AttributeErrors, 1, "Expected exactly 1 batch error")
	}
}

func TestLive_Error(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	whoami, err := Whoami(context.Background(), dnsimpleClient)
	assert.NoError(t, err)

	_, err = dnsimpleClient.Registrar.RegisterDomain(context.Background(), fmt.Sprintf("%v", whoami.Account.ID), fmt.Sprintf("example-%v.test", time.Now().Unix()), &RegisterDomainInput{})

	var got *ErrorResponse
	assert.ErrorAs(t, err, &got)

	fmt.Println(got.Message)
}
