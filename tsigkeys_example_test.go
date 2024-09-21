package powerdns_test

import (
	"context"
	"log"

	"github.com/joeig/go-powerdns/v3"
)

var exampleTSIGKey = powerdns.TSIGKey{
	Name:      powerdns.String("examplekey"),
	Algorithm: powerdns.String("hmac-sha256"),
	Key:       powerdns.String("ruTjBX2Jw/2BlE//5255fmKHaSRvLvp6p+YyDDAXThnBN/1Mz/VwMw+HQJVtkpDsAXvpPuNNZhucdKmhiOS4Tg=="),
}

func ExampleTSIGKeysService_Create() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	_, err := pdns.TSIGKeys.Create(ctx, *exampleTSIGKey.Name, *exampleTSIGKey.Algorithm, "")
	if err != nil {
		log.Fatalf("%v", err)
	}

}

func ExampleTSIGKeysService_List() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if _, err := pdns.TSIGKeys.List(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleTSIGKeysService_Get() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if _, err := pdns.TSIGKeys.Get(ctx, *exampleTSIGKey.ID); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleTSIGKeysService_Change() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if _, err := pdns.TSIGKeys.Change(ctx, *exampleTSIGKey.ID, exampleTSIGKey); err != nil {
		log.Fatalf("%v", err)
	}
}

func ExampleTSIGKeysService_Delete() {
	pdns := powerdns.New("http://localhost:8080", "localhost", powerdns.WithAPIKey("apipw"))
	ctx := context.Background()

	if err := pdns.TSIGKeys.Delete(ctx, *exampleTSIGKey.ID); err != nil {
		log.Fatalf("%v", err)
	}
}
