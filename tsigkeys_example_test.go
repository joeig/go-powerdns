package powerdns_test

import (
	"context"
	"log"

	"github.com/joeig/go-powerdns/v3"
)

var (
	exampleTSIGKey = powerdns.TSIGKey{
		Name:      powerdns.String("examplekey"),
		Algorithm: powerdns.String("hmac-sha256"),
		Key:       powerdns.String("ruTjBX2Jw/2BlE//5255fmKHaSRvLvp6p+YyDDAXThnBN/1Mz/VwMw+HQJVtkpDsAXvpPuNNZhucdKmhiOS4Tg=="),
	}
)

func ExampleTSIGKeyService_Create() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKeyHeader("apipw"))
	ctx := context.Background()

	_, err := pdns.TSIGKey.Create(ctx, *exampleTSIGKey.Name, *exampleTSIGKey.Algorithm, "")
	if err != nil {
		log.Fatalf("%v", err)
	}

}

func ExampleTSIGKeyService_List() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKeyHeader("apipw"))
	ctx := context.Background()

	if _, err := pdns.TSIGKey.List(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleTSIGKeyService_Get() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKeyHeader("apipw"))
	ctx := context.Background()

	if _, err := pdns.TSIGKey.Get(ctx, *exampleTSIGKey.ID); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleTSIGKeyService_Change() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKeyHeader("apipw"))
	ctx := context.Background()

	if _, err := pdns.TSIGKey.Change(ctx, *exampleTSIGKey.ID, exampleTSIGKey); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleTSIGKeyService_Delete() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKeyHeader("apipw"))
	ctx := context.Background()

	if err := pdns.TSIGKey.Delete(ctx, *exampleTSIGKey.ID); err != nil {
		log.Fatalf("%v", err)
	}
}
