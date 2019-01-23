package powerdns

import (
	"reflect"
	"testing"
)

func TestNewClientHTTP(t *testing.T) {
	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	tmpl := &PowerDNS{"http", "localhost", "8080", "localhost", headers, nil}
	p := NewClient("http://localhost:8080", "localhost", headers, nil)
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}

func TestNewClientHTTPS(t *testing.T) {
	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	tmpl := &PowerDNS{"https", "localhost", "443", "localhost", headers, nil}
	p := NewClient("https://localhost", "localhost", headers, nil)
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
}
