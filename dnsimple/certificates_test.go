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

func TestCertificatesService_GetCertificatePrivateKey(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/certificates/2/private_key", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getCertificatePrivateKey/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	certificateResponse, err := client.Certificates.GetCertificatePrivateKey("1010", "example.com", 2)
	if err != nil {
		t.Errorf("Certificates.GetPrivateKey() returned error: %v", err)
	}

	certificate := certificateResponse.Data
	wantSingle := &Certificate{
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAtzCcMfWoQRt5AMEY0HUb2GaraL1GsWOo6YXdPfe+YDvtnmDw\n23NcoTX7VSeCgU9M3RKs19AsCJcRNTLJ2dmDrAuyCTud9YTAaXQcTOLUhtO8T8+9\nAFVIva2OmAlKCR5saBW3JaRxW7V2aHEd/d1ss1CvNOO7jNppc9NwGSnDHcn3rqNv\n/U3MaU0gpJJRqsKkvcLU6IHJGgxyQ6AbpwJDIqBnzkjHu2IuhGEbRuMjyWLA2qts\njyVlfPotDxUdVouUQpz7dGHUFrLR7ma8QAYuOfl1ZMyrc901HGMa7zwbnFWurs3f\ned7vAosTRZIjnn72/3Wo7L9RiMB+vwr3NX7c9QIDAQABAoIBAEQx32OlzK34GTKT\nr7Yicmw7xEGofIGa1Q2h3Lut13whsxKLif5X0rrcyqRnoeibacS+qXXrJolIG4rP\nTl8/3wmUDQHs5J+6fJqFM+fXZUCP4AFiFzzhgsPBsVyd0KbWYYrZ0qU7s0ttoRe+\nTGjuHgIe3ip1QKNtx2Xr50YmytDydknmro79J5Gfrub1l2iA8SDm1eBrQ4SFaNQ2\nU709pHeSwX8pTihUX2Zy0ifpr0O1wYQjGLneMoG4rrNQJG/z6iUdhYczwwt1kDRQ\n4WkM2sovFOyxbBfoCQ3Gy/eem7OXfjNKUe47DAVLnPkKbqL/3Lo9FD7kcB8K87Ap\nr/vYrl0CgYEA413RAk7571w5dM+VftrdbFZ+Yi1OPhUshlPSehavro8kMGDEG5Ts\n74wEz2X3cfMxauMpMrBk/XnUCZ20AnWQClK73RB5fzPw5XNv473Tt/AFmt7eLOzl\nOcYrhpEHegtsD/ZaljlGtPqsjQAL9Ijhao03m1cGB1+uxI7FgacdckcCgYEAzkKP\n6xu9+WqOol73cnlYPS3sSZssyUF+eqWSzq2YJGRmfr1fbdtHqAS1ZbyC5fZVNZYV\nml1vfXi2LDcU0qS04JazurVyQr2rJZMTlCWVET1vhik7Y87wgCkLwKpbwamPDmlI\n9GY+fLNEa4yfAOOpvpTJpenUScxyKWH2cdYFOOMCgYBhrJnvffINC/d64Pp+BpP8\nyKN+lav5K6t3AWd4H2rVeJS5W7ijiLTIq8QdPNayUyE1o+S8695WrhGTF/aO3+ZD\nKQufikZHiQ7B43d7xL7BVBF0WK3lateGnEVyh7dIjMOdj92Wj4B6mv2pjQ2VvX/p\nAEWVLCtg24/+zL64VgxmXQKBgGosyXj1Zu2ldJcQ28AJxup3YVLilkNje4AXC2No\n6RCSvlAvm5gpcNGE2vvr9lX6YBKdl7FGt8WXBe/sysNEFfgmm45ZKOBCUn+dHk78\nqaeeQHKHdxMBy7utZWdgSqt+ZS299NgaacA3Z9kVIiSLDS4V2VeW7riujXXP/9TJ\nnxaRAoGBAMWXOfNVzfTyrKff6gvDWH+hqNICLyzvkEn2utNY9Q6WwqGuY9fvP/4Z\nXzc48AOBzUr8OeA4sHKJ79sJirOiWHNfD1swtvyVzsFZb6moiNwD3Ce/FzYCa3lQ\nU8blTH/uqpR2pSC6whzJ/lnSdqHUqhyp00000000000000000000\n-----END RSA PRIVATE KEY-----\n"}

	if !reflect.DeepEqual(certificate, wantSingle) {
		t.Fatalf("Certificates.GetPrivateKey() returned %+v, want %+v", certificate, wantSingle)
	}
}
