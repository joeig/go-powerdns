package powerdns_test

import (
	"context"
	"log"

	"github.com/joeig/go-powerdns/v3"
)

func ExampleZonesService_AddNative() {
	pdns := powerdns.NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
	ctx := context.Background()

	zone, err := pdns.Zones.AddNative(ctx, "example.com.", false, "", false, "", "", true, []string{"localhost."})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Zone: %v", zone)
}

func ExampleZonesService_Change() {
	pdns := powerdns.NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
	ctx := context.Background()
	zoneChangeSet := &powerdns.Zone{
		Account: powerdns.String("test"),
		DNSsec:  powerdns.Bool(true),
	}

	if err := pdns.Zones.Change(ctx, "example.com.", zoneChangeSet); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleZonesService_Get() {
	pdns := powerdns.NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
	ctx := context.Background()

	zone, err := pdns.Zones.Get(ctx, "example.com.")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Zone: %v", zone)
}
