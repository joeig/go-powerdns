package powerdns

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerServersMockResponder() {
	httpmock.RegisterResponder("GET", generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

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
		},
	)

	httpmock.RegisterResponder("GET", generateTestAPIVHostURL(),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

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
		},
	)
}

func registerCacheFlushMockResponder(testDomain string) {
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/cache/flush", generateTestAPIVHostURL()),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.URL.Query().Get("domain") != makeDomainCanonical(testDomain) {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Eneity"), nil
			}

			cacheFlushResultMock := CacheFlushResult{
				Count:  Uint32(1),
				Result: String("foo"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cacheFlushResultMock)
		},
	)
}

func TestListServers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerServersMockResponder()

	p := initialisePowerDNSTestClient()
	servers, err := p.Servers.List()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(servers) == 0 {
		t.Error("Received amount of servers is 0")
	}
}

func TestListServersError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Servers.List(); err == nil {
		t.Error("error is nil")
	}
}

func TestGetServer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerServersMockResponder()

	p := initialisePowerDNSTestClient()
	server, err := p.Servers.Get(testVHost)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *server.ID != testVHost {
		t.Error("Received no server")
	}
}

func TestGetServerError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Servers.Get(testVHost); err == nil {
		t.Error("error is nil")
	}
}

func TestCacheFlush(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerCacheFlushMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	cacheFlushResult, err := p.Servers.CacheFlush(testVHost, testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *cacheFlushResult.Count != 1 {
		t.Error("Received cache flush result is invalid")
	}
}

func TestCacheFlushResultError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Servers.CacheFlush(testVHost, testDomain); err == nil {
		t.Error("error is nil")
	}
}
