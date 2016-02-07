package dnsimple

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWebhooks_webhookPath(t *testing.T) {
	if want, got := "/1010/webhooks", webhookPath("1010", 0); want != got {
		t.Errorf("webhookPath(%v,  ) = %v, want %v", "1010", got, want)
	}

	if want, got := "/1010/webhooks/1", webhookPath("1010", 1); want != got {
		t.Errorf("webhookPath(%v, 1) = %v, want %v", "1010", got, want)
	}
}

func TestWebhooksService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaders(t, r)

		fmt.Fprint(w, `
			{"data":[{"id":1,"url":"https://webhook.test"},{"id":2,"url":"http://callback.test"}]}
		`)
	})

	accountID := "1010"

	webhooksResponse, err := client.Webhooks.List(accountID)
	if err != nil {
		t.Fatalf("Webhooks.List() returned error: %v", err)
	}

	webhooks := webhooksResponse.Data
	if want, got := 2, len(webhooks); want != got {
		t.Errorf("Webhooks.List() expected to return %v webhooks, got %v", want, got)
	}

	if want, got := 1, webhooks[0].ID; want != got {
		t.Fatalf("Webhooks.List() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "https://webhook.test", webhooks[0].URL; want != got {
		t.Fatalf("Webhooks.List() returned URL expected to be `%v`, got `%v`", want, got)
	}
}
