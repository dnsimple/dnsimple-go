package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatePath(t *testing.T) {
	assert.Equal(t, "/1010/templates", templatePath("1010", ""))
	assert.Equal(t, "/1010/templates/1", templatePath("1010", "1"))
}

func TestTemplatesService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listTemplates/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	templatesResponse, err := client.Templates.ListTemplates(context.Background(), "1010", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, templatesResponse.Pagination)
	templates := templatesResponse.Data
	assert.Len(t, templates, 2)
	assert.Equal(t, int64(1), templates[0].ID)
	assert.Equal(t, "Alpha", templates[0].Name)
}

func TestTemplatesService_List_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listTemplates/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Templates.ListTemplates(context.Background(), "1010", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestTemplatesService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createTemplate/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "Beta"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	templateAttributes := Template{Name: "Beta"}

	templateResponse, err := client.Templates.CreateTemplate(context.Background(), accountID, templateAttributes)

	assert.NoError(t, err)
	template := templateResponse.Data
	assert.Equal(t, int64(1), template.ID)
	assert.Equal(t, "Beta", template.Name)
}

func TestTemplatesService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getTemplate/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	templateResponse, err := client.Templates.GetTemplate(context.Background(), "1010", "1")

	assert.NoError(t, err)
	template := templateResponse.Data
	wantSingle := &Template{
		ID:          1,
		SID:         "alpha",
		AccountID:   1010,
		Name:        "Alpha",
		Description: "An alpha template.",
		CreatedAt:   "2016-03-22T11:08:58Z",
		UpdatedAt:   "2016-03-22T11:08:58Z",
	}
	assert.Equal(t, wantSingle, template)
}

func TestTemplatesService_UpdateTemplate(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/updateTemplate/success.http")

		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	templateAttributes := Template{Name: "Alpha"}
	templateResponse, err := client.Templates.UpdateTemplate(context.Background(), "1010", "1", templateAttributes)

	assert.NoError(t, err)
	template := templateResponse.Data
	wantSingle := &Template{
		ID:          1,
		SID:         "alpha",
		AccountID:   1010,
		Name:        "Alpha",
		Description: "An alpha template.",
		CreatedAt:   "2016-03-22T11:08:58Z",
		UpdatedAt:   "2016-03-22T11:08:58Z",
	}
	assert.Equal(t, wantSingle, template)
}

func TestTemplatesService_DeleteTemplate(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteTemplate/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Templates.DeleteTemplate(context.Background(), "1010", "1")

	assert.NoError(t, err)
}
