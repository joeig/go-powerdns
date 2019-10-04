package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestListServers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				serversMock := []Server{
					{
						Type:       String("Server"),
						ID:         String(testVHost),
						DaemonType: String("authoritative"),
						Version:    String("4.1.2"),
						URL:        String("/api/v1/servers/" + testVHost),
						ConfigURL:  String("/api/v1/servers/" + testVHost + "/config{/config_setting}"),
						ZonesURL:   String("/api/v1/servers/" + testVHost + "/zones{/zone}"),
					},
				}
				return httpmock.NewJsonResponse(http.StatusOK, serversMock)
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	servers, err := p.Servers.List()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(servers) == 0 {
		t.Error("Received amount of servers is 0")
	}
}

func TestGetServer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIVhostURL(),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				serverMock := Server{
					Type:       String("Server"),
					ID:         String(testVHost),
					DaemonType: String("authoritative"),
					Version:    String("4.1.2"),
					URL:        String("/api/v1/servers/" + testVHost),
					ConfigURL:  String("/api/v1/servers/" + testVHost + "/config{/config_setting}"),
					ZonesURL:   String("/api/v1/servers/" + testVHost + "/zones{/zone}"),
				}
				return httpmock.NewJsonResponse(http.StatusOK, serverMock)
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	server, err := p.Servers.Get(testVHost)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *server.ID != testVHost {
		t.Error("Received no server")
	}
}
