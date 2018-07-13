package powerdns_test

import (
	"github.com/joeig/go-powerdns"
	"testing"
)

func TestGetZones(t *testing.T) {
	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	zones, err := p.GetZones()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(*zones) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestGetZone(t *testing.T) {
	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	zone, err := p.GetZone()
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID == "" {
		t.Error("Received no zone")
	}
}

func TestNotify(t *testing.T) {
	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	notifyResult, err := p.Notify()
	if err != nil {
		t.Errorf("%s", err)
	}
	if notifyResult.Result != "Notification queued" {
		t.Error("Notification was not queued successfully")
	}
}
