package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestZonesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/domains/example.com/zone", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `$ORIGIN example.com.\n$TTL 1h\nexample.com. 3600 IN SOA ns1.dnsimple.com. admin.dnsimple.com. 1418386497 86400 7200 604800 300\n`)
	})

	zone, _, err := client.Zones.Get("example.com")

	if err != nil {
		t.Errorf("Zones.Get returned error: %v", err)
	}

	want := `$ORIGIN example.com.\n$TTL 1h\nexample.com. 3600 IN SOA ns1.dnsimple.com. admin.dnsimple.com. 1418386497 86400 7200 604800 300\n`
	if !reflect.DeepEqual(zone, want) {
		t.Errorf("Zones.Get returned %+v, want %+v", zone, want)
	}
}
