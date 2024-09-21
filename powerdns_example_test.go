package powerdns_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joeig/go-powerdns/v3"
	"log"
	"math/rand"
	"time"
)

func ExampleNewClient() {
	_ = powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
}

func Example() {
	domain := fmt.Sprintf("%d.example.com.", rand.Int())

	// Let's say
	// * PowerDNS Authoritative Server is listening on `http://localhost:80`,
	// * the virtual host is `localhost` and
	// * the API key is `apipw`.
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))

	// All API interactions support a Go context, which allow you to pass cancellation signals and deadlines.
	// If you don't need a context, `context.Background()` would be the right choice for the following examples.
	// If you want to learn more about how context helps you to build reliable APIs, see: https://go.dev/blog/context
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
