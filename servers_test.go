package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGetServers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				serversMock := []Server{
					{
						Type:       "Server",
						ID:         "localhost",
						DaemonType: "authoritative",
						Version:    "4.1.2",
						URL:        "/api/v1/servers/localhost",
						ConfigURL:  "/api/v1/servers/localhost/config{/config_setting}",
						ZonesURL:   "/api/v1/servers/localhost/zones{/zone}",
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
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				serverMock := Server{
					Type:       "Server",
					ID:         "localhost",
					DaemonType: "authoritative",
					Version:    "4.1.2",
					URL:        "/api/v1/servers/localhost",
					ConfigURL:  "/api/v1/servers/localhost/config{/config_setting}",
					ZonesURL:   "/api/v1/servers/localhost/zones{/zone}",
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
	if server.ID != "localhost" {
		t.Error("Received no server")
	}
}
