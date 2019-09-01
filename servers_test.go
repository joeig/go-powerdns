package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGetServers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				serversMock := []Server{
					{
						Type:       "Server",
						ID:         testVhost,
						DaemonType: "authoritative",
						Version:    "4.1.2",
						URL:        "/api/v1/servers/" + testVhost,
						ConfigURL:  "/api/v1/servers/" + testVhost + "/config{/config_setting}",
						ZonesURL:   "/api/v1/servers/" + testVhost + "/zones{/zone}",
					},
				}
				return httpmock.NewJsonResponse(200, serversMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	servers, err := p.GetServers()
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
					Type:       "Server",
					ID:         testVhost,
					DaemonType: "authoritative",
					Version:    "4.1.2",
					URL:        "/api/v1/servers/" + testVhost,
					ConfigURL:  "/api/v1/servers/" + testVhost + "/config{/config_setting}",
					ZonesURL:   "/api/v1/servers/" + testVhost + "/zones{/zone}",
				}
				return httpmock.NewJsonResponse(200, serverMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	server, err := p.GetServer()
	if err != nil {
		t.Errorf("%s", err)
	}
	if server.ID != testVhost {
		t.Error("Received no server")
	}
}
