package powerdns_test

import (
	"context"
	"log"

	"github.com/joeig/go-powerdns/v3"
)

func ExampleCryptokeysService_Create() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	cryptokey, err := pdns.Cryptokeys.Create(ctx, "example.com", powerdns.Cryptokey{
		KeyType:   powerdns.String("ksk"),
		Active:    powerdns.Bool(true),
		Published: powerdns.Bool(true),
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Cryptokey: %v", cryptokey)
}

func ExampleCryptokeysService_Change() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.Cryptokeys.Change(ctx, "example.com", 1, powerdns.Cryptokey{
		Active: powerdns.Bool(false),
	}); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleCryptokeysService_List() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if _, err := pdns.Cryptokeys.List(ctx, "example.com"); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleCryptokeysService_Get() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if _, err := pdns.Cryptokeys.Get(ctx, "example.com", 1); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleCryptokeysService_Delete() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.Cryptokeys.Delete(ctx, "example.com", 1); err != nil {
		log.Fatalf("%v", err)
	}
}
