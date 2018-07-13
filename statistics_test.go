package powerdns_test

import (
	"github.com/joeig/go-powerdns"
	"testing"
)

func TestGetStatistics(t *testing.T) {
	p := powerdns.NewClient("http://localhost:8080/", "localhost", "example.com", "apipw")
	statistics, err := p.GetStatistics()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(*statistics) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}
