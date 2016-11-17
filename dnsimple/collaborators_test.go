package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestCollaboratorsService_ListCollaborators(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listCollaborators/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	collaboratorsResponse, err := client.Collaborators.ListCollaborators("1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Collaborators.ListCollaborators() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}), collaboratorsResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Collaborators.ListCollaborators() pagination expected to be %v, got %v", want, got)
	}

	collaborators := collaboratorsResponse.Data
	if want, got := 2, len(collaborators); want != got {
		t.Errorf("Collaborators.ListCollaborators() expected to return %v collaborators, got %v", want, got)
	}

	if want, got := 100, collaborators[0].ID; want != got {
		t.Fatalf("Collaborators.ListCollaborators() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example.com", collaborators[0].DomainName; want != got {
		t.Fatalf("Collaborators.ListCollaborators() returned DomainName expected to be `%v`, got `%v`", want, got)
	}
	if want, got := 999, collaborators[0].UserID; want != got {
		t.Fatalf("Collaborators.ListCollaborators() returned UserID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := false, collaborators[0].Invitation; want != got {
		t.Fatalf("Collaborators.ListCollaborators() returned Invitation expected to be `%v`, got `%v`", want, got)
	}
}

func TestCollaboratorsService_ListCollaborators_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listCollaborators/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Collaborators.ListCollaborators("1010", "example.com", &ListOptions{Page: 2, PerPage: 20})
	if err != nil {
		t.Fatalf("Collaborators.ListCollaborators() returned error: %v", err)
	}
}

func TestCollaboratorsService_AddCollaborator(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/addCollaborator/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"email": "existing-user@example.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainID := "example.com"
	collaboratorAttributes := CollaboratorAttributes{Email: "existing-user@example.com"}

	collaboratorResponse, err := client.Collaborators.AddCollaborator(accountID, domainID, collaboratorAttributes)
	if err != nil {
		t.Fatalf("Collaborators.AddCollaborator() returned error: %v", err)
	}

	collaborator := collaboratorResponse.Data
	if want, got := 100, collaborator.ID; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example.com", collaborator.DomainName; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned DomainName expected to be `%v`, got `%v`", want, got)
	}
	if want, got := false, collaborator.Invitation; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned Invitation expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "2016-10-07T08:53:41.643Z", collaborator.AcceptedAt; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned AcceptedAt expected to be `%v`, got `%v`", want, got)
	}
}

func TestCollaboratorsService_AddNonExistingCollaborator(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/addCollaborator/invite-success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"email": "invited-user@example.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainID := "example.com"
	collaboratorAttributes := CollaboratorAttributes{Email: "invited-user@example.com"}

	collaboratorResponse, err := client.Collaborators.AddCollaborator(accountID, domainID, collaboratorAttributes)
	if err != nil {
		t.Fatalf("Collaborators.AddCollaborator() returned error: %v", err)
	}

	collaborator := collaboratorResponse.Data
	if want, got := 101, collaborator.ID; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example.com", collaborator.DomainName; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned DomainName expected to be `%v`, got `%v`", want, got)
	}
	if want, got := true, collaborator.Invitation; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned Invitation expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "", collaborator.AcceptedAt; want != got {
		t.Fatalf("Collaborators.AddCollaborator() returned AcceptedAt expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_RemoveCollaborator(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/collaborators/100", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/removeCollaborator/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	domainID := "example.com"
	contactID := "100"

	_, err := client.Collaborators.RemoveCollaborator(accountID, domainID, contactID)
	if err != nil {
		t.Fatalf("Collaborators.RemoveCollaborator() returned error: %v", err)
	}
}
