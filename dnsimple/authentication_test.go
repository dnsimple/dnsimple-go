package dnsimple

import (
	"context"
	"testing"

	"golang.org/x/oauth2"
)

func TestBasicAuthHTTPClient(t *testing.T) {
	ctx := context.Background()
	x := "user"
	y := "pass"
	h := BasicAuthHTTPClient(ctx, x, y)
	rt := h.Transport

	ts, ok := rt.(*BasicAuthTransport)
	if !ok {
		t.Fatalf("Expected transport to be a dnsimple.BasicAuthTransport")
	}
	if ts.Username != x {
		t.Fatalf("Username mismathing, expected `%v`, got `%v`", x, ts.Username)
	}
	if ts.Password != y {
		t.Fatalf("Username mismathing, expected `%v`, got `%v`", y, ts.Password)
	}
}

func TestStaticTokenHTTPClient(t *testing.T) {
	ctx := context.Background()
	x := "123456"
	h := StaticTokenHTTPClient(ctx, x)
	rt := h.Transport

	ts, ok := rt.(*oauth2.Transport)
	if !ok {
		t.Fatalf("Expected transport to be a oauth2.Transport")
	}
	tk, _ := ts.Source.Token()
	if tk.AccessToken != x {
		t.Fatalf("Token mismathing, expected `%v`, got `%v`", x, tk.AccessToken)
	}
}