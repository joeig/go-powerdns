package powerdns

import (
	"reflect"
	"testing"
)

func initialisePowerDNSTestClient() *PowerDNS {
	return NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
}

func TestNewClientHTTP(t *testing.T) {
	tmpl := &PowerDNS{"http", "localhost", "8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil}
	p := NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}

func TestNewClientHTTPS(t *testing.T) {
	tmpl := &PowerDNS{"https", "localhost", "443", "localhost", map[string]string{"X-API-Key": "apipw"}, nil}
	p := NewClient("https://localhost", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}
