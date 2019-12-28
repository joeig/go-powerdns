package powerdns

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestListStatistics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/statistics",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				statisticsMock := "[{\"name\": \"corrupt-packets\", \"type\": \"StatisticItem\", \"value\": \"0\"}, {\"name\": \"response-by-rcode\", \"type\": \"MapStatisticItem\", \"value\": [{\"name\": \"foo1\", \"value\": \"bar1\"}, {\"name\": \"foo2\", \"value\": \"bar2\"}]}, {\"name\": \"logmessages\", \"size\": \"10000\", \"type\": \"RingStatisticItem\", \"value\": [{\"name\": \"gmysql Connection successful. Connected to database 'powerdns' on 'mariadb'.\", \"value\": \"235\"}]}]"
				return httpmock.NewStringResponse(http.StatusOK, statisticsMock), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
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

func TestListStatisticsError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Statistics.List(); err == nil {
		t.Error("error is nil")
	}
}

func TestGetStatistics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/statistics",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") != testAPIKey {
				return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
			}

			if req.URL.Query().Get("statistic") != "corrupt-packets" {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			statisticsMock := "[{\"name\": \"corrupt-packets\", \"type\": \"StatisticItem\", \"value\": \"0\"}]"
			return httpmock.NewStringResponse(http.StatusOK, statisticsMock), nil
		},
	)

	p := initialisePowerDNSTestClient()
	statistics, err := p.Statistics.Get("corrupt-packets")
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(statistics) != 1 {
		t.Error("Received amount of statistics is not 1")
	}
}

func TestGetStatisticsError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Statistics.Get("corrupt-packets"); err == nil {
		t.Error("error is nil")
	}
}
