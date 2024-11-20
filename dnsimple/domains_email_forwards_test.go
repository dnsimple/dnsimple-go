package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var regexpEmail = regexp.MustCompile(`.+@.+`)

func TestEmailForwardPath(t *testing.T) {
	assert.Equal(t, "/1010/domains/example.com/email_forwards", emailForwardPath("1010", "example.com", 0))
	assert.Equal(t, "/1010/domains/example.com/email_forwards/2", emailForwardPath("1010", "example.com", 2))
}

func TestDomainsService_EmailForwardsList(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/email_forwards", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listEmailForwards/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	forwardsResponse, err := client.Domains.ListEmailForwards(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, forwardsResponse.Pagination)
	forwards := forwardsResponse.Data
	assert.Len(t, forwards, 2)
	assert.Equal(t, int64(17702), forwards[0].ID)
	assert.Regexp(t, regexpEmail, forwards[0].From)
}

func TestDomainsService_EmailForwardsList_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/email_forwards", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listEmailForwards/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListEmailForwards(context.Background(), "1010", "example.com", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestDomainsService_CreateEmailForward(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/email_forwards", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createEmailForward/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"from": "me"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	forwardAttributes := EmailForward{From: "me"}

	forwardResponse, err := client.Domains.CreateEmailForward(context.Background(), "1010", "example.com", forwardAttributes)
	assert.NoError(t, err)
	forward := forwardResponse.Data
	assert.Equal(t, int64(41872), forward.ID)
	assert.Regexp(t, regexpEmail, forward.From)
}

func TestDomainsService_GetEmailForward(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/email_forwards/41872", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getEmailForward/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	forwardResponse, err := client.Domains.GetEmailForward(context.Background(), "1010", "example.com", 41872)

	assert.NoError(t, err)
	forward := forwardResponse.Data
	wantSingle := &EmailForward{
		ID:               41872,
		DomainID:         235146,
		From:             "example@dnsimple.xyz",
		AliasName:        "",
		AliasEmail:       "example@dnsimple.xyz",
		To:               "example@example.com",
		DestinationEmail: "example@example.com",
		Active:           true,
		CreatedAt:        "2021-01-25T13:54:40Z",
		UpdatedAt:        "2021-01-25T13:54:40Z"}
	assert.Equal(t, wantSingle, forward)
}

func TestDomainsService_DeleteEmailForward(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/email_forwards/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteEmailForward/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.DeleteEmailForward(context.Background(), "1010", "example.com", 2)

	assert.NoError(t, err)
}
