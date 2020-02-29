package main

import (
	"encoding/json"
	"fmt"
	"github.com/joeig/go-powerdns/v2"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	domain := fmt.Sprintf("test-%d.com.", rand.Int())

	pdns := powerdns.NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)

	// Create a native zone
	zone, err := pdns.Zones.AddNative(domain, false, "", false, "", "", true, []string{"localhost."})
	if err != nil {
		log.Fatalf("%v", err)
	}

	o, _ := json.MarshalIndent(zone, "", "\t")
	fmt.Printf("zone: %s\n\n", o)

	// Add and change an A record
	if err := pdns.Records.Add(domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeA, 1337, []string{"127.0.0.9"}); err != nil {
		log.Fatalf("%v", err)
	}
	if err := pdns.Records.Change(domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeA, 42, []string{"127.0.0.10"}); err != nil {
		log.Fatalf("%v", err)
	}

	// Add a MX record with multiple values
	if err := pdns.Records.Add(domain, domain, powerdns.RRTypeMX, 1337, []string{"10 mx1.example.com.", "20 mx2.example.com."}); err != nil {
		log.Fatalf("%v", err)
	}

	// Add a TXT record
	if err := pdns.Records.Add(domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeTXT, 1337, []string{"\"foo1\""}); err != nil {
		log.Fatalf("%v", err)
	}

	// Change a zone
	zoneChangeSet := &powerdns.Zone{
		Account: powerdns.String("test"),
		DNSsec:  powerdns.Bool(true),
	}

	err = pdns.Zones.Change(domain, zoneChangeSet)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Retrieve zone attributes
	changedZone, err := pdns.Zones.Get(domain)
	if err != nil {
		log.Fatalf("%v", err)
	}

	o, _ = json.MarshalIndent(changedZone, "", "\t")
	fmt.Printf("changed zone: %s\n\n", o)
	fmt.Printf("Account is \"%s\" and DNSsec is %t\n\n", powerdns.StringValue(changedZone.Account), powerdns.BoolValue(changedZone.DNSsec))
}
