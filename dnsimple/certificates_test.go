package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestCertificatesService_CertificatesList(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/certificates", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listCertificates/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	certificatesResponse, err := client.Certificates.ListCertificates("1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Certificates.ListCertificates() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 2}), certificatesResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Certificates.ListCertificates() pagination expected to be %v, got %v", want, got)
	}

	certificates := certificatesResponse.Data
	if want, got := 2, len(certificates); want != got {
		t.Errorf("Certificates.ListCertificates() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 22289, certificates[0].ID; want != got {
		t.Fatalf("Certificates.ListCertificates() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "www.weppos.net", certificates[0].CommonName; want != got {
		t.Fatalf("Certificates.ListCertificates() returned CommonName expected to be `%v`, got `%v`", want, got)
	}
}

func TestCertificatesService_CertificatesList_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/certificates", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listCertificates/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Certificates.ListCertificates("1010", "example.com", &ListOptions{Page: 2, PerPage: 20})
	if err != nil {
		t.Fatalf("Certificates.ListCertificates() returned error: %v", err)
	}
}
