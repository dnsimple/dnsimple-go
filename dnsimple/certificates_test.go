package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestCertificatesService_ListCertificates(t *testing.T) {
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

	if want, got := 1, certificates[0].ID; want != got {
		t.Fatalf("Certificates.ListCertificates() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "www.weppos.net", certificates[0].CommonName; want != got {
		t.Fatalf("Certificates.ListCertificates() returned CommonName expected to be `%v`, got `%v`", want, got)
	}
}

func TestCertificatesService_ListCertificates_WithOptions(t *testing.T) {
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

func TestCertificatesService_GetCertificate(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/certificates/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getCertificate/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	certificateResponse, err := client.Certificates.GetCertificate("1010", "example.com", 2)
	if err != nil {
		t.Errorf("Certificates.GetCertificate() returned error: %v", err)
	}

	certificate := certificateResponse.Data
	wantSingle := &Certificate{
		ID:                  1,
		DomainID:            10,
		CommonName:          "www.weppos.net",
		Years:               1,
		State:               "issued",
		AuthorityIdentifier: "letsencrypt",
		CreatedAt:           "2016-06-11T18:47:08.949Z",
		UpdatedAt:           "2016-06-11T18:47:37.546Z",
		ExpiresOn:           "2016-09-09",

		CertificateRequest: "-----BEGIN CERTIFICATE REQUEST-----\nMIICljCCAX4CAQAwGTEXMBUGA1UEAwwOd3d3LndlcHBvcy5uZXQwggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC3MJwx9ahBG3kAwRjQdRvYZqtovUaxY6jp\nhd09975gO+2eYPDbc1yhNftVJ4KBT0zdEqzX0CwIlxE1MsnZ2YOsC7IJO531hMBp\ndBxM4tSG07xPz70AVUi9rY6YCUoJHmxoFbclpHFbtXZocR393WyzUK8047uM2mlz\n03AZKcMdyfeuo2/9TcxpTSCkklGqwqS9wtTogckaDHJDoBunAkMioGfOSMe7Yi6E\nYRtG4yPJYsDaq2yPJWV8+i0PFR1Wi5RCnPt0YdQWstHuZrxABi45+XVkzKtz3TUc\nYxrvPBucVa6uzd953u8CixNFkiOefvb/dajsv1GIwH6/Cvc1ftz1AgMBAAGgODA2\nBgkqhkiG9w0BCQ4xKTAnMCUGA1UdEQQeMByCDnd3dy53ZXBwb3MubmV0ggp3ZXBw\nb3MubmV0MA0GCSqGSIb3DQEBCwUAA4IBAQCDnVBO9RdJX0eFeZzlv5c8yG8duhKP\nl0Vl+V88fJylb/cbNj9qFPkKTK0vTXmS2XUFBChKPtLucp8+Z754UswX+QCsdc7U\nTTSG0CkyilcSubdZUERGej1XfrVQhrokk7Fu0Jh3BdT6REP0SIDTpA8ku/aRQiAp\np+h19M37S7+w/DMGDAq2LSX8jOpJ1yIokRDyLZpmwyLxutC21DXMGoJ3xZeUFrUT\nqRNwzkn2dJzgTrPkzhaXalUBqv+nfXHqHaWljZa/O0NVCFrHCdTdd53/6EE2Yabv\nq5SFTkRCpaxrvM/7a8Tr4ixD1/VKD6rw3+WC00000000000000000000\n-----END CERTIFICATE REQUEST-----\n"}

	if !reflect.DeepEqual(certificate, wantSingle) {
		t.Fatalf("Certificates.GetCertificate() returned %+v, want %+v", certificate, wantSingle)
	}
}
