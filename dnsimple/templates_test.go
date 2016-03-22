package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestTemplates_templatePath(t *testing.T) {
	if want, got := "/1010/templates", templatePath("1010", 0); want != got {
		t.Errorf("templatePath(%v,  ) = %v, want %v", "1010", got, want)
	}

	if want, got := "/1010/templates/1", templatePath("1010", 1); want != got {
		t.Errorf("templatePath(%v, 1) = %v, want %v", "1010", got, want)
	}
}

func TestTemplatesService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listTemplates/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	templatesResponse, err := client.Templates.ListTemplates("1010", nil)
	if err != nil {
		t.Fatalf("Templates.ListTemplates() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}), templatesResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Templates.ListTemplates() pagination expected to be %v, got %v", want, got)
	}

	templates := templatesResponse.Data
	if want, got := 2, len(templates); want != got {
		t.Errorf("Templates.ListTemplates() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 1, templates[0].ID; want != got {
		t.Fatalf("Templates.ListTemplates() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Alpha", templates[0].Name; want != got {
		t.Fatalf("Templates.ListTemplates() returned Name expected to be `%v`, got `%v`", want, got)
	}
}
