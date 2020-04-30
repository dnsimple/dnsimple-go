package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestTemplatesService_ApplyTemplate(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/templates/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/applyTemplate/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Templates.ApplyTemplate(context.Background(), "1010", "1", "example.com")
	if err != nil {
		t.Fatalf("Templates.ApplyTemplate() returned error: %v", err)
	}
}
