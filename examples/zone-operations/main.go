package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/joeig/go-powerdns/v3"
)

func main() {
	domain := fmt.Sprintf("%d.example.com.", rand.Int())
	pdns := powerdns.NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
	ctx := context.Background()

	// Create a native zone
	zone, err := pdns.Zones.AddNative(ctx, domain, false, "", false, "", "", true, []string{"localhost."})
	if err != nil {
		log.Fatalf("%v", err)
	}

	o, _ := json.MarshalIndent(zone, "", "\t")
	log.Printf("Zone: %s\n\n", o)

	// Add and change an A record
	if err := pdns.Records.Add(ctx, domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeA, 1337, []string{"127.0.0.9"}); err != nil {
		log.Fatalf("%v", err)
	}
	if err := pdns.Records.Change(ctx, domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeA, 42, []string{"127.0.0.10"}); err != nil {
		log.Fatalf("%v", err)
	}

	// update the existing record with a comment
	comment := powerdns.Comment{
		Content:    powerdns.String("Example comment"),
		Account:    powerdns.String("example account"),
		ModifiedAt: powerdns.Uint64(uint64(time.Now().Unix())),
	}
	if err := pdns.Records.Change(ctx, domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeA, 42, []string{"127.0.0.10"}, powerdns.WithComments(comment)); err != nil {
		log.Fatalf("%v", err)
	}

	// Add a MX record with multiple values
	if err := pdns.Records.Add(ctx, domain, domain, powerdns.RRTypeMX, 1337, []string{"10 mx1.example.com.", "20 mx2.example.com."}); err != nil {
		log.Fatalf("%v", err)
	}

	// Add a TXT record
	if err := pdns.Records.Add(ctx, domain, fmt.Sprintf("www.%s", domain), powerdns.RRTypeTXT, 1337, []string{"\"foo1\""}); err != nil {
		log.Fatalf("%v", err)
	}

	// Create a TSIG Record
	exampleKey, err := pdns.TSIGKey.Create(ctx, "examplekey", "hmac-sha256", "")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Change a zone
	zoneChangeSet := &powerdns.Zone{
		Account:          powerdns.String("test"),
		DNSsec:           powerdns.Bool(true),
		MasterTSIGKeyIDs: []string{*exampleKey.ID},
	}

	if err := pdns.Zones.Change(ctx, domain, zoneChangeSet); err != nil {
		log.Fatalf("%v", err)
	}

	// Retrieve zone attributes
	changedZone, err := pdns.Zones.Get(ctx, domain)
	if err != nil {
		log.Fatalf("%v", err)
	}

	o, _ = json.MarshalIndent(changedZone, "", "\t")
	log.Printf("Changed zone: %q\n\n", o)
	log.Printf("Account is %q and DNSsec is %t\n\n", powerdns.StringValue(changedZone.Account), powerdns.BoolValue(changedZone.DNSsec))
}
