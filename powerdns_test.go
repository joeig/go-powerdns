package powerdns_test

import (
	"github.com/joeig/go-powerdns"
	"reflect"
	"testing"
)

func TestNewClientHTTP(t *testing.T) {
	tmpl := &powerdns.PowerDNS{"http", "localhost", "8080", "localhost", "apipw"}
	p := powerdns.NewClient("http://localhost:8080", "localhost", "apipw")
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}

func TestNewClientHTTPS(t *testing.T) {
	tmpl := &powerdns.PowerDNS{"https", "localhost", "443", "localhost", "apipw"}
	p := powerdns.NewClient("https://localhost", "localhost", "apipw")
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}
