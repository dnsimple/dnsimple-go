package dnsimple

import (
	"context"
	"testing"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthHTTPClient(t *testing.T) {
	ctx := context.Background()
	x := "user"
	y := "pass"
	h := BasicAuthHTTPClient(ctx, x, y)
	rt := h.Transport

	ts, ok := rt.(*BasicAuthTransport)

	assert.True(t, ok)
	assert.Equal(t, x, ts.Username)
	assert.Equal(t, y, ts.Password)
}

func TestStaticTokenHTTPClient(t *testing.T) {
	ctx := context.Background()
	x := "123456"
	h := StaticTokenHTTPClient(ctx, x)
	rt := h.Transport

	ts, ok := rt.(*oauth2.Transport)

	assert.True(t, ok)
	tk, _ := ts.Source.Token()
	assert.Equal(t, x, tk.AccessToken)
}
