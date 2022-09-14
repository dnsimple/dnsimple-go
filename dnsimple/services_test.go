package dnsimple

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicePath(t *testing.T) {
	assert.Equal(t, "/services", servicePath(""))
	assert.Equal(t, "/services/name", servicePath("name"))
}

func TestServicesService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/services", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listServices/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	servicesResponse, err := client.Services.ListServices(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, &Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}, servicesResponse.Pagination)
	services := servicesResponse.Data
	assert.Len(t, services, 2)
	assert.Equal(t, int64(1), services[0].ID)
	assert.Equal(t, "Service 1", services[0].Name)
}

func TestServicesService_List_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/services", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/listServices/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	_, err := client.Services.ListServices(context.Background(), &ListOptions{Page: Int(2), PerPage: Int(20)})

	assert.NoError(t, err)
}

func TestServicesService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/services/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/api/getService/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	serviceID := "1"

	serviceResponse, err := client.Services.GetService(context.Background(), serviceID)

	assert.NoError(t, err)
	service := serviceResponse.Data
	wantSingle := &Service{
		ID:               1,
		SID:              "service1",
		Name:             "Service 1",
		Description:      "First service example.",
		SetupDescription: "",
		RequiresSetup:    true,
		DefaultSubdomain: "",
		CreatedAt:        "2014-02-14T19:15:19Z",
		UpdatedAt:        "2016-03-04T09:23:27Z",
		Settings: []ServiceSetting{
			{
				Name:        "username",
				Label:       "Service 1 Account Username",
				Append:      ".service1.com",
				Description: "Your Service 1 username is used to connect services to your account.",
				Example:     "username",
				Password:    false,
			},
		},
	}
	assert.Equal(t, wantSingle, service)
}
