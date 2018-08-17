package powerdns

import (
	"reflect"
	"testing"
)

func TestNewClientHTTP(t *testing.T) {
	tmpl := &PowerDNS{"http", "localhost", "8080", "localhost", "apipw"}
	p := NewClient("http://localhost:8080", "localhost", "apipw")
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}

func TestNewClientHTTPS(t *testing.T) {
	tmpl := &PowerDNS{"https", "localhost", "443", "localhost", "apipw"}
	p := NewClient("https://localhost", "localhost", "apipw")
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}
