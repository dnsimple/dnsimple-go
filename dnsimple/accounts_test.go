package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listAccounts/success-user.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountsResponse, err := client.Accounts.ListAccounts(context.Background(), nil)
	assert.NoError(t, err)

	accounts := accountsResponse.Data
	assert.Len(t, accounts, 2)
	assert.Equal(t, int64(123), accounts[0].ID)
	assert.Equal(t, "john@example.com", accounts[0].Email)
}
