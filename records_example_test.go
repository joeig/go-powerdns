package powerdns_test

import (
	"context"
	"log"

	"github.com/joeig/go-powerdns/v3"
)

func ExampleRecordsService_Add_basic() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.Records.Add(ctx, "example.com.", "www.example.com.", powerdns.RRTypeA, 1337, []string{"127.0.0.9"}); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleRecordsService_Add_mX() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.Records.Add(ctx, "example.com.", "www.example.com.", powerdns.RRTypeMX, 1337, []string{"10 mx1.example.com.", "20 mx2.example.com."}); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleRecordsService_Add_tXT() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.Records.Add(ctx, "example.com.", "www.example.com.", powerdns.RRTypeTXT, 1337, []string{"\"foo1\""}); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleRecordsService_Change() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.Records.Change(ctx, "example.com.", "www.example.com.", powerdns.RRTypeA, 42, []string{"127.0.0.10"}); err != nil {
		log.Fatalf("%v", err)
	}
}
