package powerdns_test

import (
	"github.com/joeig/go-powerdns"
	"testing"
)

func TestGetServers(t *testing.T) {
	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	servers, err := p.GetServers()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(*servers) == 0 {
		t.Error("Received amount of servers is 0")
	}
}

func TestGetServer(t *testing.T) {
	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	server, err := p.GetServer()
	if err != nil {
		t.Errorf("%s", err)
	}
	if server.ID == "" {
		t.Error("Received no server")
	}
}
