package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGetStatistics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/statistics",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				statisticsMock := []Statistic{
					{
						Name:  "corrupt-packets",
						Type:  "StatisticItem",
						Value: "1337",
					},
					{
						Name:  "deferred-cache-inserts",
						Type:  "StatisticItem",
						Value: "42",
					},
					{
						Name:  "deferred-cache-lookup",
						Type:  "StatisticItem",
						Value: "123",
					},
					{
						Name:  "deferred-packetcache-inserts",
						Type:  "StatisticItem",
						Value: "234",
					},
					{
						Name:  "deferred-packetcache-lookup",
						Type:  "StatisticItem",
						Value: "345",
					},
					{
						Name:  "dnsupdate-answers",
						Type:  "StatisticItem",
						Value: "456",
					},
					{
						Name:  "dnsupdate-changes",
						Type:  "StatisticItem",
						Value: "567",
					},
				}
				return httpmock.NewJsonResponse(200, statisticsMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	statistics, err := p.GetStatistics()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(statistics) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}
