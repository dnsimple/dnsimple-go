package dnsimple

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_Whoami(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/whoami", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/whoami/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	whoamiResponse, err := client.Identity.Whoami(context.Background())

	assert.NoError(t, err)
	whoami := whoamiResponse.Data
	assert.Nil(t, whoami.User)
	assert.NotNil(t, whoami.Account)
	account := whoami.Account
	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, "example-account@example.com", account.Email)
	assert.Equal(t, "teams-v1-monthly", account.PlanIdentifier)
}
