package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplates_templateRecordPath(t *testing.T) {
	assert.Equal(t, "/1010/templates/1/records", templateRecordPath("1010", "1", 0))
	assert.Equal(t, "/1010/templates/1/records/2", templateRecordPath("1010", "1", 2))
}

func TestTemplatesService_ListTemplateRecords(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listTemplateRecords/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	templatesRecordsResponse, err := client.Templates.ListTemplateRecords(context.Background(), "1010", "1", nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, templatesRecordsResponse.Pagination)
	templates := templatesRecordsResponse.Data
	assert.Len(t, templates, 2)
	assert.Equal(t, int64(296), templates[0].ID)
	assert.Equal(t, "192.168.1.1", templates[0].Content)
}

func TestTemplatesService_ListTemplateRecords_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listTemplateRecords/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Templates.ListTemplateRecords(context.Background(), "1010", "1", &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestTemplatesService_CreateTemplateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/createTemplateRecord/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "Beta"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	templateRecordAttributes := TemplateRecord{Name: "Beta"}

	templateRecordResponse, err := client.Templates.CreateTemplateRecord(context.Background(), "1010", "1", templateRecordAttributes)

	assert.NoError(t, err)
	templateRecord := templateRecordResponse.Data
	assert.Equal(t, int64(300), templateRecord.ID)
	assert.Equal(t, "mx.example.com", templateRecord.Content)
}

func TestTemplatesService_GetTemplateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1/records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getTemplateRecord/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	templateRecordResponse, err := client.Templates.GetTemplateRecord(context.Background(), "1010", "1", 2)

	assert.NoError(t, err)
	templateRecord := templateRecordResponse.Data
	wantSingle := &TemplateRecord{
		ID:         301,
		TemplateID: 268,
		Name:       "",
		Content:    "mx.example.com",
		TTL:        600,
		Priority:   10,
		Type:       "MX",
		CreatedAt:  "2016-05-03T08:03:26Z",
		UpdatedAt:  "2016-05-03T08:03:26Z"}
	assert.Equal(t, wantSingle, templateRecord)
}

func TestTemplatesService_DeleteTemplateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/templates/1/records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/deleteTemplateRecord/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Templates.DeleteTemplateRecord(context.Background(), "1010", "1", 2)

	assert.NoError(t, err)
}
