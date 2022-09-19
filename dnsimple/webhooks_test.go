package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhooks_webhookPath(t *testing.T) {
	assert.Equal(t, "/1010/webhooks", webhookPath("1010", 0))
	assert.Equal(t, "/1010/webhooks/1", webhookPath("1010", 1))
}

func TestWebhooksService_ListWebhooks(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/webhooks", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listWebhooks/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	webhooksResponse, err := client.Webhooks.ListWebhooks(context.Background(), "1010", nil)

	assert.NoError(t, err)
	webhooks := webhooksResponse.Data
	assert.Len(t, webhooks, 2)
	assert.Equal(t, int64(1), webhooks[0].ID)
	assert.Equal(t, "https://webhook.test", webhooks[0].URL)
}

func TestWebhooksService_CreateWebhook(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/webhooks", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createWebhook/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"url": "https://webhook.test"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	webhookAttributes := Webhook{URL: "https://webhook.test"}

	webhookResponse, err := client.Webhooks.CreateWebhook(context.Background(), "1010", webhookAttributes)

	assert.NoError(t, err)
	webhook := webhookResponse.Data
	assert.Equal(t, int64(1), webhook.ID)
	assert.Equal(t, "https://webhook.test", webhook.URL)
}

func TestWebhooksService_GetWebhook(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/webhooks/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getWebhook/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	webhookResponse, err := client.Webhooks.GetWebhook(context.Background(), "1010", 1)

	assert.NoError(t, err)
	webhook := webhookResponse.Data
	wantSingle := &Webhook{
		ID:  1,
		URL: "https://webhook.test"}
	assert.Equal(t, wantSingle, webhook)
}

func TestWebhooksService_DeleteWebhook(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/webhooks/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteWebhook/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Webhooks.DeleteWebhook(context.Background(), "1010", 1)

	assert.NoError(t, err)
}
