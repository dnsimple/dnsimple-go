package main

import (
	"github.com/weppos/go-dnsimple/dnsimple"
	"fmt"
	"log"
	"os"
)

// Simple script to dump a zone.
//
// Configuration:
//
// 		DNSIMPLE_API_TOKEN
// 		DNSIMPLE_EMAIL
//
// Usage:
//
//		$ go run examples/zone-dump.go "example.com"
//
func main() {
	dnsimpleToken := os.Getenv("DNSIMPLE_API_TOKEN")
	dnsimpleEmail := os.Getenv("DNSIMPLE_EMAIL")

	client := dnsimple.NewClient(dnsimpleToken, dnsimpleEmail)

	zone, _, error := client.Zones.Get(os.Args[1])
	if error != nil {
		log.Fatalln(error)
	}

	fmt.Println(zone)
}
