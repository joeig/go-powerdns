package powerdns_test

import (
	"context"
	"log"

	"github.com/joeig/go-powerdns/v3"
)

func ExampleSearchService_Data() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	results, err := pdns.Search.Data(ctx, "example*", 100, powerdns.SearchObjectTypeAll)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Search results: %v", results)
}

func ExampleSearchService_Data_zones() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	results, err := pdns.Search.Data(ctx, "example*", 100, powerdns.SearchObjectTypeZone)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Zone search results: %v", results)
}
