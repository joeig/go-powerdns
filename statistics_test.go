package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestListStatistics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIVhostURL()+"/statistics",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				statisticsMock := []Statistic{
					{
						Name:  String("corrupt-packets"),
						Type:  String("StatisticItem"),
						Value: String("1337"),
					},
					{
						Name:  String("deferred-cache-inserts"),
						Type:  String("StatisticItem"),
						Value: String("42"),
					},
					{
						Name:  String("deferred-cache-lookup"),
						Type:  String("StatisticItem"),
						Value: String("123"),
					},
					{
						Name:  String("deferred-packetcache-inserts"),
						Type:  String("StatisticItem"),
						Value: String("234"),
					},
					{
						Name:  String("deferred-packetcache-lookup"),
						Type:  String("StatisticItem"),
						Value: String("345"),
					},
					{
						Name:  String("dnsupdate-answers"),
						Type:  String("StatisticItem"),
						Value: String("456"),
					},
					{
						Name:  String("dnsupdate-changes"),
						Type:  String("StatisticItem"),
						Value: String("567"),
					},
				}
				return httpmock.NewJsonResponse(200, statisticsMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	statistics, err := p.Statistics.List()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(statistics) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}
