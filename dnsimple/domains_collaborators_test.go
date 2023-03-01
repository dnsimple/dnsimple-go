package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollaboratorsPath(t *testing.T) {
	t.Run("empty account id", func(t *testing.T) {
		path, err := collaboratorsPath("", "example.com")
		assert.Error(t, err)
		assert.Empty(t, path)
	})

	t.Run("empty domain identifier", func(t *testing.T) {
		path, err := collaboratorsPath("1010", "")
		assert.Error(t, err)
		assert.Empty(t, path)
	})

	t.Run("success", func(t *testing.T) {
		path, err := collaboratorsPath("1010", "example.com")
		assert.NoError(t, err)
		assert.Equal(t, "/1010/domains/example.com/collaborators", path)
	})
}

func TestCollaboratorPath(t *testing.T) {
	path, err := collaboratorPath("1010", "example.com", 2)
	assert.NoError(t, err)
	assert.Equal(t, "/1010/domains/example.com/collaborators/2", path)
}

func TestDomainsService_ListCollaborators(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listCollaborators/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	collaboratorsResponse, err := client.Domains.ListCollaborators(context.Background(), "1010", "example.com", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, collaboratorsResponse.Pagination)
	collaborators := collaboratorsResponse.Data
	assert.Len(t, collaborators, 2)
	assert.Equal(t, int64(100), collaborators[0].ID)
	assert.Equal(t, "example.com", collaborators[0].DomainName)
	assert.Equal(t, int64(999), collaborators[0].UserID)
	assert.False(t, collaborators[0].Invitation)
}

func TestDomainsService_ListCollaborators_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listCollaborators/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Domains.ListCollaborators(context.Background(), "1010", "example.com", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestDomainsService_AddCollaborator(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/addCollaborator/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"email": "existing-user@example.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainID := "example.com"
	collaboratorAttributes := CollaboratorAttributes{Email: "existing-user@example.com"}

	collaboratorResponse, err := client.Domains.AddCollaborator(context.Background(), accountID, domainID, collaboratorAttributes)

	assert.NoError(t, err)
	collaborator := collaboratorResponse.Data
	assert.Equal(t, int64(100), collaborator.ID)
	assert.Equal(t, "example.com", collaborator.DomainName)
	assert.False(t, collaborator.Invitation)
	assert.Equal(t, "2016-10-07T08:53:41Z", collaborator.AcceptedAt)
}

func TestDomainsService_AddNonExistingCollaborator(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/addCollaborator/invite-success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"email": "invited-user@example.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainID := "example.com"
	collaboratorAttributes := CollaboratorAttributes{Email: "invited-user@example.com"}

	collaboratorResponse, err := client.Domains.AddCollaborator(context.Background(), accountID, domainID, collaboratorAttributes)

	assert.NoError(t, err)
	collaborator := collaboratorResponse.Data
	assert.Equal(t, int64(101), collaborator.ID)
	assert.Equal(t, "example.com", collaborator.DomainName)
	assert.True(t, collaborator.Invitation)
	assert.Equal(t, "", collaborator.AcceptedAt)
}

func TestDomainsService_RemoveCollaborator(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators/100", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/removeCollaborator/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainID := "example.com"
	collaboratorID := int64(100)

	_, err := client.Domains.RemoveCollaborator(context.Background(), accountID, domainID, collaboratorID)

	assert.NoError(t, err)
}
