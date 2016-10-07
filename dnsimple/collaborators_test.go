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
