package powerdns_test

import (
	"github.com/joeig/go-powerdns/v3"
)

func ExampleNewClient() {
	_ = powerdns.NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
}
